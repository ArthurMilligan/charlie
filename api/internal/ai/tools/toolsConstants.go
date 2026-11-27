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
	AddMemoryLoreSuccess        = "Воспоминание о прошлом Charlie успешно добавлено"
	AddUserSectionSuccess       = "Секция успешно добавлено"
	AddUserFactSuccess          = "Факт успешно добавлено"
	UpdateUserFactSuccess       = "Факт успешно обновлен"
)

var Tools = []responses.ToolUnionParam{
	responses.ToolParamOfWebSearch(responses.WebSearchToolTypeWebSearch),
	{
		OfFunction: &responses.FunctionToolParam{
			Name: "get_user_profile_sections",
			Description: openai.String(
				"Получить каталог разделов профиля пользователя: ID, ключ, название и описание каждого раздела. " +
					"Вызывай, когда для ответа нужны сведения о пользователе или когда нужно выбрать раздел для нового факта. " +
					"Если актуальный каталог уже есть в контексте, используй его. " +
					"Для чтения и записи фактов передавай числовой ID выбранного раздела в section_id.",
			),
			Strict: openai.Bool(true),
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
			Name: "get_user_profile_facts",
			Description: openai.String(
				"Прочитать сохранённые факты из одного раздела профиля. " +
					"Выбери раздел по его названию и описанию в каталоге и передай его числовой ID. " +
					"Вызывай, чтобы учесть предпочтения пользователя в ответе или проверить, есть ли факт перед записью. " +
					"Если подходящих сведений нет, продолжи с доступным контекстом. " +
					"Задай короткий уточняющий вопрос, если без этих сведений нельзя дать полезный ответ.",
			),
			Strict: openai.Bool(true),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"section_id": map[string]any{
						"type":        "integer",
						"minimum":     1,
						"description": "Числовой ID раздела из каталога профиля.",
					},
				},
				"required":             []string{"section_id"},
				"additionalProperties": false,
			},
		},
	},
	{
		OfFunction: &responses.FunctionToolParam{
			Name: "update_profile_fact",
			Description: openai.String(
				"Обновить существующий факт в профиле пользователя. " +
					"Вызывай, если пользователь в текущей реплике явно исправил " +
					"сохранённый факт или сообщил, что его значение изменилось. " +
					"Найди раздел через get_user_profile_sections и прочитай его факты " +
					"через get_user_profile_facts. Используй ID найденного раздела " +
					"и точный FactKey существующей записи. " +
					"Передай set_value: true и полное новое значение. " +
					"Если подходящей записи нет, для нового факта используй add_profile_fact.",
			),
			Strict: openai.Bool(false),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"section_id": map[string]any{
						"type":        "integer",
						"minimum":     1,
						"description": "Числовой ID раздела, в котором уже хранится факт.",
					},
					"fact_key": map[string]any{
						"type":        "string",
						"description": "Точный FactKey из результата get_user_profile_facts.",
					},
					"value": map[string]any{
						"description": "Полное новое значение факта в JSON. " +
							"Передавай число, строку, булево значение, массив или объект. " +
							"Массивы и объекты передавай без дополнительных кавычек.",
					},
					"set_value": map[string]any{
						"type":        "boolean",
						"description": "Передай true, чтобы заменить сохранённое значение.",
					},
				},
				"required":             []string{"section_id", "fact_key", "value", "set_value"},
				"additionalProperties": false,
			},
		},
	},
	{
		OfFunction: &responses.FunctionToolParam{
			Name: "add_profile_section",
			Description: openai.String(
				"Создать раздел для устойчивых сведений, которые пользователь сообщает о себе. " +
					"Сначала проверь каталог. Если подходящий раздел уже есть, используй его. " +
					"Создавай раздел для темы, в которой может быть несколько связанных фактов, например «Питомцы». " +
					"После успешного создания снова вызови get_user_profile_sections, чтобы получить числовой ID нового раздела.",
			),
			Strict: openai.Bool(true),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"section_key": map[string]any{
						"type":        "string",
						"pattern":     "^[a-z][a-z0-9_]*$",
						"description": "Стабильный ключ раздела латиницей в нижнем регистре, например pets или food_preferences.",
					},
					"title": map[string]any{
						"type":        "string",
						"description": "Короткое название раздела для пользователя, например «Питомцы».",
					},
					"description": map[string]any{
						"type":        "string",
						"description": "Какие факты относятся к разделу и при каких запросах их полезно читать.",
					},
				},
				"required":             []string{"section_key", "title", "description"},
				"additionalProperties": false,
			},
		},
	},
	{
		OfFunction: &responses.FunctionToolParam{
			Name: "add_profile_fact",
			Description: openai.String(
				"Сохранить новый устойчивый факт, который пользователь явно сообщил о себе в текущей реплике. " +
					"Выбери раздел из каталога и прочитай его факты перед записью. " +
					"Используй стабильный ключ, обозначающий предмет факта. " +
					"Если такой факт уже сохранён, пропусти запись. " +
					"Если пользователь поправил существующий факт, используй инструмент обновления, если он доступен. " +
					"Не создавай для исправления второй ключ. " +
					"Если пользователь явно попросил запомнить факт, подтверждай сохранение после успешного вызова.",
			),
			Strict: openai.Bool(false),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"section_id": map[string]any{
						"type":        "integer",
						"minimum":     1,
						"description": "Числовой ID существующего раздела из каталога профиля.",
					},
					"fact_key": map[string]any{
						"type":        "string",
						"description": "Стабильный ключ факта внутри раздела, например preferred_language, height_cm или cuisine:georgian.",
					},
					"value": map[string]any{
						"description": `Само JSON-значение факта. Например, строка "likes", число 180, массив ["jazz","rock"] или объект {"name":"Энни","species":"собака"}. Массивы и объекты передавай без дополнительных кавычек.`,
					},
				},
				"required":             []string{"section_id", "fact_key", "value"},
				"additionalProperties": false,
			},
		},
	},
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
				"Создать задачу по просьбе пользователя. " +
					"Различай запланированное выполнение и крайний срок. " +
					"Если пользователь указал только дату, сохрани её и оставь время выполнения равным null. " +
					"Не выдумывай даты, время или длительность.",
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
						"type":        []string{"string", "null"},
						"enum":        []any{"todo", "in_progress", "cancelled", nil},
						"description": "Статус новой задачи. Если не указан, передай null.",
					},
					"description": map[string]any{
						"type":        []string{"string", "null"},
						"description": "Дополнительные детали задачи. Если их нет, передай null.",
					},
					"planned_on": map[string]any{
						"type":   []string{"string", "null"},
						"format": "date-time",
						"description": "Запланированная дата выполнения. Передавай полночь этой даты " +
							"со смещением часового пояса, например 2026-11-26T00:00:00+03:00. " +
							"Если дата неизвестна, передай null.",
					},
					"planned_time": map[string]any{
						"type":   []string{"string", "null"},
						"format": "date-time",
						"description": "Дата из planned_on и указанное пользователем время начала, " +
							"например 2026-11-26T16:30:00+03:00. Заполняй только вместе с planned_on. " +
							"Если время начала неизвестно, передай null.",
					},
					"duration_minutes": map[string]any{
						"type":        []string{"integer", "null"},
						"minimum":     1,
						"description": "Длительность задачи в минутах. Если не указана, передай null.",
					},
					"due_on": map[string]any{
						"type":   []string{"string", "null"},
						"format": "date-time",
						"description": "Крайний срок, если указана только дата. Передавай полночь " +
							"этой даты, например 2026-11-26T00:00:00+03:00. " +
							"Если срока нет или заполнен due_at, передай null.",
					},
					"due_at": map[string]any{
						"type":   []string{"string", "null"},
						"format": "date-time",
						"description": "Точный крайний срок с датой, временем и смещением часового пояса, " +
							"например 2026-11-26T18:00:00+03:00. " +
							"Если точное время срока неизвестно, передай null. " +
							"При заполнении due_at оставь due_on равным null.",
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
			Name: "update_task",
			Description: openai.String(
				"Обновить существующую задачу по task_id. Если ID неизвестен, сначала найди задачу через get_tasks. " +
					"Если подходят несколько задач, уточни, какую изменить. " +
					"Меняй только поля, о которых попросил пользователь. " +
					"Для поля с флагом set_* передай false и значение null, чтобы сохранить его. " +
					"Передай true и новое значение, чтобы обновить поле, или true и null, чтобы очистить его.",
			),
			Strict: openai.Bool(true),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"task_id": map[string]any{
						"type":        "integer",
						"minimum":     1,
						"description": "ID существующей задачи.",
					},
					"title": map[string]any{
						"type":        []string{"string", "null"},
						"description": "Новое непустое название. Передай null, чтобы сохранить текущее.",
					},
					"status": map[string]any{
						"type":        []string{"string", "null"},
						"enum":        []any{"todo", "in_progress", "done", "cancelled", nil},
						"description": "Новый статус. Передай null, чтобы сохранить текущий.",
					},
					"completed_at": map[string]any{
						"type":   []string{"string", "null"},
						"format": "date-time",
						"description": "Фактическое время завершения в RFC3339 со смещением. " +
							"Заполняй только при status: done и явно указанном времени. " +
							"При status: done и completed_at: null сервер поставит текущее время.",
					},
					"set_description": map[string]any{"type": "boolean"},
					"description": map[string]any{
						"type":        []string{"string", "null"},
						"description": "Новое описание задачи.",
					},
					"set_planned_on": map[string]any{"type": "boolean"},
					"planned_on": map[string]any{
						"type":   []string{"string", "null"},
						"format": "date-time",
						"description": "Запланированная дата. Передавай полночь со смещением, " +
							"например 2026-11-26T00:00:00+03:00. При очистке даты также очисти planned_time.",
					},
					"set_planned_time": map[string]any{"type": "boolean"},
					"planned_time": map[string]any{
						"type":   []string{"string", "null"},
						"format": "date-time",
						"description": "Время начала в запланированный день. У задачи должна быть planned_on. " +
							"При необходимости задай дату в том же вызове.",
					},
					"set_duration_minutes": map[string]any{"type": "boolean"},
					"duration_minutes": map[string]any{
						"type":        []string{"integer", "null"},
						"minimum":     1,
						"description": "Длительность задачи в минутах.",
					},
					"set_due_on": map[string]any{"type": "boolean"},
					"due_on": map[string]any{
						"type":   []string{"string", "null"},
						"format": "date-time",
						"description": "Крайний срок с точностью до дня. Передавай полночь со смещением. " +
							"При переходе с точного срока очисти due_at в том же вызове.",
					},
					"set_due_at": map[string]any{"type": "boolean"},
					"due_at": map[string]any{
						"type":   []string{"string", "null"},
						"format": "date-time",
						"description": "Точный крайний срок в RFC3339 со смещением. " +
							"При переходе со срока на день очисти due_on в том же вызове.",
					},
				},
				"required": []string{
					"task_id", "title", "status", "completed_at",
					"set_description", "description",
					"set_planned_on", "planned_on",
					"set_planned_time", "planned_time",
					"set_duration_minutes", "duration_minutes",
					"set_due_on", "due_on",
					"set_due_at", "due_at",
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
	{
		OfFunction: &responses.FunctionToolParam{
			Name: "add_lore_memory",
			Description: openai.String(
				"Сохранить факт, событие или деталь из истории Charlie. " +
					"Вызывай, когда пользователь явно просит добавить что-то в лор Charlie. " +
					"Подтверждай сохранение после успешного результата инструмента.",
			),
			Strict: openai.Bool(true),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"content": map[string]any{
						"type": "string",
						"description": "Самодостаточная запись о Charlie. " +
							"Сохрани имена, даты и детали, которые сообщил пользователь.",
					},
				},
				"required":             []string{"content"},
				"additionalProperties": false,
			},
		},
	},
	{
		OfFunction: &responses.FunctionToolParam{
			Name:        "call_to_memory",
			Description: openai.String(ToolDescriptionCallToMemory),
			Strict:      openai.Bool(true),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"content": map[string]any{
						"type": "string",
						"description": "Короткая поисковая фраза с темой и ключевыми словами. " +
							"Сохраняй точные имена и названия, например «Redis причины переноса в V2».",
					},
				},
				"required":             []string{"content"},
				"additionalProperties": false,
			},
		},
	},
	{
		OfFunction: &responses.FunctionToolParam{
			Name: "get_memory_context",
			Description: openai.String(
				"Прочитать целиком запись из памяти по ID одного из её фрагментов. " +
					"Вызывай, когда пользователь просит полный текст, подробности или контекст " +
					"и результат call_to_memory указывает на нужный источник. " +
					"Читай полученные фрагменты в порядке ChunkIndex.",
			),
			Strict: openai.Bool(true),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"memory_id": map[string]any{
						"type": "integer",
						"description": "Значение поля ID подходящего фрагмента " +
							"из результата call_to_memory. SourceID здесь не подходит.",
					},
				},
				"required":             []string{"memory_id"},
				"additionalProperties": false,
			},
		},
	},
}

func taskUpdateField(description string, valueSchema map[string]any) map[string]any {
	return map[string]any{
		"type": []string{"object", "null"},
		"description": description + " Передай null, чтобы сохранить текущее значение. " +
			"Передай объект с value: null, чтобы очистить поле. " +
			"Передай объект со значением в value, чтобы обновить поле.",
		"properties": map[string]any{
			"value": valueSchema,
		},
		"required":             []string{"value"},
		"additionalProperties": false,
	}
}
