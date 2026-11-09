package tools

import (
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

const (
	ToolAddNote                 = "add_note"
	ToolGetNotes                = "get_notes"
	InternalErrorMessageOnTools = "Internal Error"
	UnknownToolError            = "Error. Unknown tool"
	AddNoteSuccess              = "Заметка успешно сохранена"
)

var Tools = []responses.ToolUnionParam{
	{
		OfFunction: &responses.FunctionToolParam{
			Name:        "add_note",
			Description: openai.String("Сохранить заметку по просьбе пользователя."),
			Strict:      openai.Bool(true),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"content": map[string]any{
						"type":        "string",
						"description": "Текст заметки для сохранения.",
					},
				},
				"required":             []string{"content"},
				"additionalProperties": false,
			},
		},
	},
	{
		OfFunction: &responses.FunctionToolParam{
			Name:        "get_notes",
			Description: openai.String("Получить сохранённые заметки по просьбе пользователя."),
			Strict:      openai.Bool(true),
			Parameters: map[string]any{
				"type":                 "object",
				"properties":           map[string]any{},
				"required":             []string{},
				"additionalProperties": false,
			},
		},
	},
}
