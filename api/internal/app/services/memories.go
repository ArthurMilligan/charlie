package services

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/pkoukk/tiktoken-go"
)

type Memories struct {
	MemoriesStorage interfaces.MemoriesStorage
	ai              *openai.Client
	countTokens     func(string) int
}

const embeddingDimensions = 1536

func NewMemoriesService(memoriesStorage interfaces.MemoriesStorage, GPTToken string) (*Memories, error) {
	openAIClient := openai.NewClient(option.WithAPIKey(GPTToken))

	encoding, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return nil, fmt.Errorf("load memory tokenizer: %w", err)
	}

	return &Memories{
		MemoriesStorage: memoriesStorage,
		ai:              &openAIClient,
		countTokens: func(text string) int {
			return len(encoding.EncodeOrdinary(text))
		},
	}, nil
}

func (m *Memories) AddMemories(
	ctx context.Context,
	chatID int64,
	sourceType models.MemorySourceType,
	sourceID uint32,
	text string,
) ([]uint32, error) {
	chunks, err := splitMemoryText(text, 500, m.countTokens)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("memory text is empty")
	}
	if err != nil {
		return nil, err
	}

	memories := make([]models.Memory, 0, len(chunks))

	for i, chunk := range chunks {
		embedding, err := embedText(ctx, m.ai, chunk)
		if err != nil {
			return nil, fmt.Errorf("embed chunk %d: %w", i, err)
		}

		memories = append(memories, models.Memory{
			ChatID:     chatID,
			SourceType: sourceType,
			SourceID:   sourceID,
			ChunkIndex: uint32(i),
			Content:    chunk,
			Embedding:  embedding,
		})
	}

	return m.MemoriesStorage.AddMemories(ctx, memories)

}

func (m *Memories) FindMemories(ctx context.Context, text string, chatID uint32) ([]models.DBMemory, error) {
	embedding, err := embedText(ctx, m.ai, text)
	if err != nil {
		return nil, err
	}

	memories, err := m.MemoriesStorage.FindMemories(ctx, text, embedding, chatID)

	if err != nil {
		return nil, err
	}

	return memories, nil
}

func (m *Memories) GetMemoryBySourceId(ctx context.Context, sourceID uint32, sourceType models.MemorySourceType) ([]models.DBMemory, error) {
	memories, err := m.MemoriesStorage.GetMemoryBySourceId(ctx, sourceID, sourceType)

	if err != nil {
		return nil, err
	}

	return memories, nil
}

func (m *Memories) GetMemoryContext(ctx context.Context, memoryID uint32) ([]models.DBMemory, error) {
	memory, err := m.MemoriesStorage.GetMemoryById(ctx, memoryID)

	if err != nil {
		return nil, err
	}

	memories, err := m.MemoriesStorage.GetMemoryBySourceId(ctx, memory.SourceID, memory.SourceType)

	if err != nil {
		return nil, err
	}

	return memories, nil
}

func splitMemoryText(
	text string,
	maxTokens int,
	countTokens func(string) int,
) ([]string, error) {
	if maxTokens <= 0 {
		return nil, fmt.Errorf("maxTokens must be positive")
	}

	text = strings.ReplaceAll(text, "\r\n", "\n")
	var chunks []string
	current := ""

	for _, block := range strings.Split(text, "\n\n") {
		paragraph := strings.TrimSpace(block)
		if paragraph == "" {
			continue
		}

		if current != "" {
			combined := current + "\n\n" + paragraph
			if countTokens(combined) <= maxTokens {
				current = combined
				continue
			}

			chunks = append(chunks, current)
			current = ""
		}

		for countTokens(paragraph) > maxTokens {
			runes := []rune(paragraph)
			cut := fittingPrefix(runes, maxTokens, countTokens)
			if cut == 0 {
				return nil, fmt.Errorf("token limit is too small for this text")
			}

			// Ищем границу в последней половине подходящего префикса.
			wordCut := cut
			for i := cut - 1; i > cut/2; i-- {
				if !unicode.IsSpace(runes[i]) {
					continue
				}
				if wordCut == cut {
					wordCut = i
				}
				if strings.ContainsRune(".!?…", runes[i-1]) {
					wordCut = i
					break
				}
			}

			if countTokens(string(runes[:wordCut])) <= maxTokens {
				cut = wordCut
			}

			chunks = append(chunks, strings.TrimSpace(string(runes[:cut])))
			paragraph = strings.TrimSpace(string(runes[cut:]))
		}

		current = paragraph
	}

	if current != "" {
		chunks = append(chunks, current)
	}
	return chunks, nil
}

// Подбирает префикс, который помещается в лимит, по границам рун.
func fittingPrefix(
	runes []rune,
	maxTokens int,
	countTokens func(string) int,
) int {
	low, high := 0, len(runes)
	for low < high {
		mid := low + (high-low+1)/2
		if countTokens(string(runes[:mid])) <= maxTokens {
			low = mid
		} else {
			high = mid - 1
		}
	}
	return low
}

func embedText(ctx context.Context, ai *openai.Client, text string) ([]float64, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("embedding input is empty")
	}

	resp, err := ai.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model:      openai.EmbeddingModelTextEmbedding3Small,
		Dimensions: openai.Int(embeddingDimensions),
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create embedding: %w", err)
	}

	if len(resp.Data) != 1 || len(resp.Data[0].Embedding) != embeddingDimensions {
		return nil, fmt.Errorf("unexpected embedding response")
	}

	return resp.Data[0].Embedding, nil
}
