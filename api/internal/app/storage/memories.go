package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemoriesStorage struct {
	db *pgxpool.Pool
}

func NewMemoriesStorage(db *pgxpool.Pool) *MemoriesStorage {
	return &MemoriesStorage{db: db}
}

func (c *MemoriesStorage) AddMemories(ctx context.Context, memories []models.Memory) ([]uint32, error) {
	if len(memories) == 0 {
		return nil, nil
	}
	query := `
	INSERT INTO memories (source_type, source_id, chunk_index, content, embedding, chat_id)
	VALUES ($1, $2, $3, $4, $5::text::vector, $6) 
	RETURNING id
	`

	ids := make([]uint32, 0, len(memories))

	err := pgx.BeginFunc(ctx, c.db, func(tx pgx.Tx) error {
		for i, memory := range memories {
			vectorJSON, err := json.Marshal(memory.Embedding)
			if err != nil {
				return fmt.Errorf("encode embedding for memory %d: %w", i, err)
			}

			var id uint32
			err = tx.QueryRow(ctx, query,
				memory.SourceType,
				memory.SourceID,
				memory.ChunkIndex,
				memory.Content,
				string(vectorJSON),
				memory.ChatID,
			).Scan(&id)
			if err != nil {
				return fmt.Errorf("insert memory %d: %w", i, err)
			}

			ids = append(ids, id)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (c *MemoriesStorage) FindMemories(ctx context.Context, text string, embedding []float64, chatID uint32) ([]models.DBMemory, error) {
	query := `
	WITH search_query AS (
    SELECT websearch_to_tsquery('russian', $3) AS q
	),
	by_vector AS (
		SELECT id, embedding <=> ($2::text::vector) AS distance
		FROM memories
		WHERE chat_id = $1
		ORDER BY distance
		LIMIT 20
	),
	by_text AS (
		SELECT id, ts_rank_cd(to_tsvector('russian', content), q) AS relevance
		FROM memories
		CROSS JOIN search_query
		WHERE chat_id = $1
		AND to_tsvector('russian', content) @@ q
		ORDER BY relevance DESC, id
		LIMIT 20
	),
	scores AS (
		SELECT id, 1.0 / (60 + row_number() OVER (ORDER BY distance, id)) AS score
		FROM by_vector
		UNION ALL
		SELECT id, 1.0 / (60 + row_number() OVER (ORDER BY relevance DESC, id))
		FROM by_text
	),
	ranked AS (
		SELECT id, sum(score) AS score
		FROM scores
		GROUP BY id
	)
	SELECT m.id, m.source_type, m.source_id, m.chunk_index, m.content
	FROM ranked r
	JOIN memories m ON m.id = r.id
	ORDER BY r.score DESC, m.id
	LIMIT 10;
	`

	vectorJSON, err := json.Marshal(embedding)
	if err != nil {
		return nil, err
	}

	rows, err := c.db.Query(ctx, query, chatID, string(vectorJSON), text)
	if err != nil {
		fmt.Print(err)

		return nil, err
	}

	defer rows.Close()

	memories := make([]models.DBMemory, 0)

	for rows.Next() {
		var memory models.DBMemory

		err := rows.Scan(
			&memory.ID,
			&memory.SourceType,
			&memory.SourceID,
			&memory.ChunkIndex,
			&memory.Content,
		)
		if err != nil {
			fmt.Print(err)

			return nil, err
		}

		memories = append(memories, memory)
	}

	if err := rows.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	return memories, nil
}

func (m *MemoriesStorage) GetMemoryById(
	ctx context.Context,
	id uint32,
) (models.DBMemory, error) {
	const query = `
        SELECT id, source_type, source_id
        FROM memories
        WHERE id = $1
    `

	var memory models.DBMemory
	err := m.db.QueryRow(ctx, query, id).Scan(
		&memory.ID,
		&memory.SourceType,
		&memory.SourceID,
	)
	if err != nil {
		return models.DBMemory{}, err
	}

	return memory, nil
}

func (m *MemoriesStorage) GetMemoryBySourceId(ctx context.Context, sourceID uint32, sourceType models.MemorySourceType) ([]models.DBMemory, error) {
	query := `
		SELECT
			chunk_index,
			content
		FROM
			memories
		WHERE source_id = $1 AND source_type = $2
		ORDER BY
			chunk_index ASC
	`

	rows, err := m.db.Query(ctx, query, sourceID, sourceType)

	if err != nil {
		fmt.Print(err)

		return nil, err
	}

	defer rows.Close()

	memories := make([]models.DBMemory, 0)

	for rows.Next() {
		var memory models.DBMemory

		err := rows.Scan(
			&memory.ChunkIndex,
			&memory.Content,
		)
		if err != nil {
			fmt.Print(err)

			return nil, err
		}

		memories = append(memories, memory)
	}

	if err := rows.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	return memories, nil
}
