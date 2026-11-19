package ai

import (
	"context"
	"encoding/json"

	"github.com/ArthurMilligan/charlie/api/internal/ai/tools"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type GPTClient struct {
	client           *openai.Client
	notesService     interfaces.NotesService
	tasksService     interfaces.TasksService
	remindersService interfaces.RemindersService
}

func NewGPTClient(token string, n interfaces.NotesService, t interfaces.TasksService, r interfaces.RemindersService) *GPTClient {
	client := openai.NewClient(
		option.WithAPIKey(token), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	return &GPTClient{client: &client, notesService: n, tasksService: t, remindersService: r}
}

func (c *GPTClient) GetAnswer(ctx context.Context, text string, messageID uint32, chatID int64) (string, error) {
	resp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Instructions: openai.String(tools.AgentInstruction),
		Input:        responses.ResponseNewParamsInputUnion{OfString: openai.String(text)},
		Model:        openai.ChatModelGPT5_2,
		Tools:        tools.Tools,
	})

	if err != nil {
		return "", err
	}

	finalResponse, err := handleResponse(ctx, c, resp, messageID, chatID)

	if err != nil {
		return "", err
	}

	return finalResponse, nil
}

func (c *GPTClient) FormatReminder(
	ctx context.Context,
	reminder models.DBRemind,
) (string, error) {
	// Только данные, необходимые для текста напоминания.
	input := struct {
		Title       string  `json:"title"`
		Description *string `json:"description,omitempty"`
	}{
		Title:       reminder.Title,
		Description: reminder.Description,
	}

	jsonReminder, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Instructions: openai.String(tools.AgentFormatReminderInstruction),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(string(jsonReminder)),
		},
		Model: openai.ChatModelGPT5_2,
	})
	if err != nil {
		return "", err
	}

	return resp.OutputText(), nil
}

func handleResponse(ctx context.Context, c *GPTClient, resp *responses.Response, messageID uint32, chatID int64) (string, error) {
	var toolOutputs []responses.ResponseInputItemUnionParam

	for _, item := range resp.Output {
		switch item.Type {
		case "function_call":
			call := item.AsFunctionCall()
			resp, toolErr := tools.UseTools(
				ctx,
				c.notesService,
				c.tasksService,
				c.remindersService,
				call.Name,
				call.Arguments,
				messageID,
				chatID,
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

	resp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model:              openai.ChatModelGPT5_2,
		PreviousResponseID: openai.String(resp.ID),
		Instructions:       openai.String(tools.AgentInstruction),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: toolOutputs,
		},
		Tools: tools.Tools,
	})

	if err != nil {
		return "", err
	}

	return handleResponse(ctx, c, resp, messageID, chatID)
}
