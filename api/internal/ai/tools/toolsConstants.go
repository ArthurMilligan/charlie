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
	AddTaskSuccess              = "Задача успешно создана"
	AddReminderSuccess          = "Напоминание успешно добавлено"
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
	{
		OfFunction: &responses.FunctionToolParam{
			Name: "add_task",
			Description: openai.String(
				"Создать задачу по явной просьбе пользователя. " +
					"Не выдумывать даты, время и длительность.",
			),
			Strict: openai.Bool(true),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title": map[string]any{
						"type":        "string",
						"description": "Краткое название задачи.",
					},
					"status": map[string]any{
						"type": []string{"string", "null"},
						"enum": []any{
							"todo",
							"in_progress",
							"done",
							"cancelled",
							nil,
						},
						"description": "Статус задачи. Если пользователь не указал его — null.",
					},
					"description": map[string]any{
						"type":        []string{"string", "null"},
						"description": "Дополнительное описание задачи или null.",
					},
					"planned_on": map[string]any{
						"type": []string{"string", "null"},
						"description": "Дата, на которую запланирована задача, " +
							"в формате RFC3339. Если дата неизвестна — null.",
					},
					"planned_time": map[string]any{
						"type": []string{"string", "null"},
						"description": "Запланированное время начала в формате RFC3339. " +
							"Если время неизвестно — null. " +
							"Передавать только вместе с planned_on.",
					},
					"duration_minutes": map[string]any{
						"type": []string{"integer", "null"},
						"description": "Планируемая длительность в минутах. " +
							"Если не указана — null.",
					},
					"due_on": map[string]any{
						"type": []string{"string", "null"},
						"description": "Дата крайнего срока в формате RFC3339, " +
							"если известно только число. Иначе — null. " +
							"Не заполнять одновременно с due_at.",
					},
					"due_at": map[string]any{
						"type": []string{"string", "null"},
						"description": "Точный момент крайнего срока в формате RFC3339. " +
							"Если точное время неизвестно — null. " +
							"Не заполнять одновременно с due_on.",
					},
				},
				"required": []string{
					"title",
					"status",
					"description",
					"planned_on",
					"planned_time",
					"duration_minutes",
					"due_on",
					"due_at",
				},
				"additionalProperties": false,
			},
		},
	},
	{
		OfFunction: &responses.FunctionToolParam{
			Name: "get_tasks",
			Description: openai.String(
				"Получить список сохранённых задач. " +
					"Если пользователь указал статус, отфильтровать задачи по нему.",
			),
			Strict: openai.Bool(true),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"status": map[string]any{
						"type": []string{"string", "null"},
						"enum": []any{
							"todo",
							"in_progress",
							"done",
							"cancelled",
							nil,
						},
						"description": "Статус задач для фильтрации. " +
							"Если статус не указан — null, чтобы получить все задачи.",
					},
				},
				"required":             []string{"status"},
				"additionalProperties": false,
			},
		},
	},
	{
		OfFunction: &responses.FunctionToolParam{
			Name: "add_reminder",
			Description: openai.String(
				"Создать напоминание по явной просьбе пользователя. " +
					"Если дата или время неоднозначны, сначала уточни их. " +
					"Не выдумывай идентификатор задачи.",
			),
			Strict: openai.Bool(true),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"task_id": map[string]any{
						"type": []string{"integer", "null"},
						"description": "ID существующей связанной задачи. " +
							"Если задача неизвестна — null.",
					},
					"title": map[string]any{
						"type":        "string",
						"description": "Краткое содержание напоминания.",
					},
					"description": map[string]any{
						"type": []string{"string", "null"},
						"description": "Дополнительные детали напоминания. " +
							"Если деталей нет — null.",
					},
					"recurrence_type": map[string]any{
						"type": []string{"string", "null"},
						"enum": []any{
							"once",
							"daily",
							"weekly",
							"monthly",
							"yearly",
							nil,
						},
						"description": "Тип повторения. " +
							"Для разового напоминания — once.",
					},
					"timezone": map[string]any{
						"type": "string",
						"description": "Часовой пояс IANA, например Europe/Moscow. " +
							"Используй часовой пояс из контекста приложения.",
					},
					"local_time": map[string]any{
						"type": "string",
						"description": "Локальные дата и время отправки в формате RFC3339, " +
							"например 2026-11-20T16:00:00+03:00.",
					},
					"start_date": map[string]any{
						"type": "string",
						"description": "Дата начала расписания в формате RFC3339, " +
							"например 2026-11-20T00:00:00+03:00.",
					},
					"end_date": map[string]any{
						"type": []string{"string", "null"},
						"description": "Дата окончания повторения в формате RFC3339. " +
							"Если дата окончания не указана — null.",
					},
					"next_remind_at": map[string]any{
						"type": "string",
						"description": "Ближайший момент отправки напоминания " +
							"в формате RFC3339 с часовым поясом.",
					},
				},
				"required": []string{
					"task_id",
					"title",
					"description",
					"recurrence_type",
					"timezone",
					"local_time",
					"start_date",
					"end_date",
					"next_remind_at",
				},
				"additionalProperties": false,
			},
		},
	},
}
