package ai

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/ai/tools"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type GPTClient struct {
	client       *openai.Client
	notesService interfaces.NotesService
	tasksService interfaces.TasksService
}

func NewGPTClient(token string, n interfaces.NotesService, t interfaces.TasksService) *GPTClient {
	client := openai.NewClient(
		option.WithAPIKey(token), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	return &GPTClient{client: &client, notesService: n, tasksService: t}
}

func (c *GPTClient) GetAnswer(ctx context.Context, text string, messageID uint32) (string, error) {
	resp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Instructions: openai.String(tools.AgentInstruction),
		Input:        responses.ResponseNewParamsInputUnion{OfString: openai.String(text)},
		Model:        openai.ChatModelGPT5_2,
		Tools:        tools.Tools,
	})

	if err != nil {
		return "", err
	}

	var toolOutputs []responses.ResponseInputItemUnionParam

	for _, item := range resp.Output {
		switch item.Type {
		case "function_call":
			call := item.AsFunctionCall()
			resp, toolErr := tools.UseTools(
				ctx,
				c.notesService,
				c.tasksService,
				call.Name,
				call.Arguments,
				messageID,
			)
			if toolErr != nil {
				return "", toolErr
			}

			toolOutputs = append(toolOutputs, responses.ResponseInputItemUnionParam{
				OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
					CallID: openai.String(call.CallID),
					Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
						OfString: openai.String(resp),
					},
				},
			})

		}
	}

	if len(toolOutputs) == 0 {
		return resp.OutputText(), nil
	}

	finalResp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model:              openai.ChatModelGPT5_2,
		PreviousResponseID: openai.String(resp.ID),
		Instructions:       openai.String(tools.AgentResponseInstruction),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: toolOutputs,
		},
	})
	if err != nil {
		return "", err
	}

	return finalResp.OutputText(), nil
}
