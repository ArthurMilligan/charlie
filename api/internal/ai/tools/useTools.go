package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

func UseTools(
	ctx context.Context,
	n interfaces.NotesService,
	t interfaces.TasksService,
	r interfaces.RemindersService,
	m interfaces.MemoriesService,
	p interfaces.ProfilesService,
	toolName string,
	args string,
	messageID uint32,
	chatID int64,
) (string, error) {
	switch toolName {
	case "add_profile_section":
		fmt.Println("use add_profile_section")
		var msg struct {
			// UserID int64 `json:"user_id"`
			SectionKey  string `json:"section_key"`
			Title       string `json:"title"`
			Description string `json:"description"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		_, err := p.AddUserProfileSection(ctx, 1, msg.SectionKey, msg.Title, msg.Description)
		if err != nil {
			return UnknownToolError, err
		}
		return AddUserSectionSuccess, nil

	case "add_profile_fact":
		fmt.Println("use add_profile_fact")
		var msg struct {
			// UserID int64 `json:"user_id"`
			SectionID int64           `json:"section_id"`
			FactKey   string          `json:"fact_key"`
			Value     json.RawMessage `json:"value"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		err := p.AddUserProfileFact(ctx, 1, msg.SectionID, msg.FactKey, msg.Value, int64(messageID))
		if err != nil {
			fmt.Println(err)
			return UnknownToolError, err
		}
		return AddUserFactSuccess, nil

	case "update_profile_fact":
		fmt.Println("use update_profile_fact")
		var msg struct {
			// UserID int64 `json:"user_id"`
			SectionID int64            `json:"section_id"`
			FactKey   string           `json:"fact_key"`
			Value     *json.RawMessage `json:"value"`
			SetValue  bool             `json:"set_value"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		err := p.UpdateUserProfileFact(ctx, 1, msg.SectionID, msg.FactKey, msg.Value, msg.SetValue)
		if err != nil {
			fmt.Println(err)
			return UnknownToolError, err
		}
		return UpdateUserFactSuccess, nil

	case "get_user_profile_sections":
		fmt.Println("use get_user_profile_sections")

		var msg struct {
			// UserID int64 `json:"user_id"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		sections, err := p.GetUserProfileSections(ctx, 1)

		if err != nil {
			return UnknownToolError, err
		}

		sectionsJSON, err := json.Marshal(sections)

		if err != nil {
			return UnknownToolError, err
		}

		return string(sectionsJSON), nil

	case "get_user_profile_facts":
		fmt.Println("use get_user_profile_facts")

		var msg struct {
			// UserID int64 `json:"user_id"`
			SectionID int64 `json:"section_id"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		facts, err := p.GetUserProfileFactsBySection(ctx, 1, msg.SectionID)

		if err != nil {
			return UnknownToolError, err
		}

		factsJSON, err := json.Marshal(facts)

		if err != nil {
			return UnknownToolError, err
		}

		return string(factsJSON), nil

	case "add_note":
		fmt.Println("use add_note")

		var msg struct {
			Content string `json:"content"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		_, err := n.AddNote(ctx, msg.Content, messageID)
		if err != nil {
			return UnknownToolError, err
		}
		return AddNoteSuccess, nil
	case "get_memory_context":
		fmt.Println("use get_memory_context")

		var msg struct {
			MemoryID uint32 `json:"memory_id"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		memories, err := m.GetMemoryContext(ctx, msg.MemoryID)

		if err != nil {
			return UnknownToolError, err
		}

		memoriesJSON, err := json.Marshal(memories)

		if err != nil {
			return UnknownToolError, err
		}

		return string(memoriesJSON), nil

	case "add_lore_memory":
		fmt.Println("use add_lore_memory")
		var msg struct {
			Content string `json:"content"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		_, err := m.AddMemories(ctx, 554822003, models.MemorySourceTypeLore, messageID, msg.Content)
		if err != nil {
			return UnknownToolError, err
		}
		return AddMemoryLoreSuccess, nil
	case "call_to_memory":
		var msg struct {
			Content string `json:"content"`
		}
		fmt.Println("use call_to_memory")

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		memories, err := m.FindMemories(ctx, msg.Content, 554822003)

		if err != nil {
			return UnknownToolError, err
		}

		memoriesJSON, err := json.Marshal(memories)

		if err != nil {
			return UnknownToolError, err
		}

		return string(memoriesJSON), nil

	case "get_notes":
		fmt.Println("use get_notes")

		notes, err := n.GetNotes(ctx)
		if err != nil {
			return UnknownToolError, err
		}

		notesJSON, err := json.Marshal(notes)
		if err != nil {
			return UnknownToolError, err
		}
		return string(notesJSON), nil

	case "add_reminder":
		fmt.Println("use add_reminder")

		var msg struct {
			TaskID         *uint32                `json:"task_id"`
			Title          string                 `json:"title"`
			Description    *string                `json:"description"`
			RecurrenceType *models.RecurrenceType `json:"recurrence_type"`
			Timezone       string                 `json:"timezone"`
			LocalTime      time.Time              `json:"local_time"`
			StartDate      time.Time              `json:"start_date"`
			EndDate        *time.Time             `json:"end_date"`
			NextRemindAt   time.Time              `json:"next_remind_at"`
			Status         *models.RemindStatus   `json:"status"`
		}
		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		_, err := r.AddRemind(ctx, msg.TaskID, msg.Title, msg.Description, msg.RecurrenceType, msg.Timezone, uint32(chatID), msg.LocalTime, msg.StartDate, msg.EndDate, msg.NextRemindAt, msg.Status)
		if err != nil {
			return UnknownToolError, err
		}
		return AddReminderSuccess, nil

	case "add_task":
		fmt.Println("use add_task")

		var msg struct {
			Title           string     `json:"title"`
			Description     *string    `json:"description"`
			Status          *string    `json:"status"`
			PlanedOn        *time.Time `json:"planned_on"`
			PlannedTime     *time.Time `json:"planned_time"`
			DurationMinutes *int       `json:"duration_minutes"`
			DueOn           *time.Time `json:"due_on"`
			DueAt           *time.Time `json:"due_at"`
			CompletedAt     *time.Time `json:"completed_at"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		taskID, err := t.AddTask(ctx, msg.Title, msg.Description, messageID, msg.Status, msg.PlanedOn, msg.PlannedTime, msg.DurationMinutes, msg.DueOn, msg.DueAt, msg.CompletedAt)
		if err != nil {
			return UnknownToolError, err
		}
		return fmt.Sprintf(`{"success":true,"task_id":%d}`, taskID), nil
	case "get_tasks":
		fmt.Println("use get_tasks")

		var msg struct {
			Status *models.TaskStatus `json:"status"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}
		tasks, err := t.GetTasks(ctx, msg.Status)
		if err != nil {
			return UnknownToolError, err
		}

		tasksJSON, err := json.Marshal(tasks)
		if err != nil {
			return UnknownToolError, err
		}

		return string(tasksJSON), nil
	case "update_task":
		fmt.Println("use update_task")

		var msg struct {
			TaskID             uint32             `json:"task_id"`
			Title              *string            `json:"title"`
			Status             *models.TaskStatus `json:"status"`
			CompletedAt        *time.Time         `json:"completed_at"`
			SetDescription     bool               `json:"set_description"`
			Description        *string            `json:"description"`
			SetPlannedOn       bool               `json:"set_planned_on"`
			PlannedOn          *time.Time         `json:"planned_on"`
			SetPlannedTime     bool               `json:"set_planned_time"`
			PlannedTime        *time.Time         `json:"planned_time"`
			SetDurationMinutes bool               `json:"set_duration_minutes"`
			DurationMinutes    *int               `json:"duration_minutes"`
			SetDueOn           bool               `json:"set_due_on"`
			DueOn              *time.Time         `json:"due_on"`
			SetDueAt           bool               `json:"set_due_at"`
			DueAt              *time.Time         `json:"due_at"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("parse update_task arguments: %w", err)
			return UnknownToolError, fmt.Errorf("parse update_task arguments: %w", err)
		}
		if msg.TaskID == 0 {
			fmt.Println("task_id must be positive")

			return UnknownToolError, fmt.Errorf("task_id must be positive")
		}

		patch := models.TaskPatch{
			Title:       msg.Title,
			Status:      msg.Status,
			CompletedAt: msg.CompletedAt,
			Description: models.PatchField[string]{
				Set: msg.SetDescription, Value: msg.Description,
			},
			PlannedOn: models.PatchField[time.Time]{
				Set: msg.SetPlannedOn, Value: msg.PlannedOn,
			},
			PlannedTime: models.PatchField[time.Time]{
				Set: msg.SetPlannedTime, Value: msg.PlannedTime,
			},
			DurationMinutes: models.PatchField[int]{
				Set: msg.SetDurationMinutes, Value: msg.DurationMinutes,
			},
			DueOn: models.PatchField[time.Time]{
				Set: msg.SetDueOn, Value: msg.DueOn,
			},
			DueAt: models.PatchField[time.Time]{
				Set: msg.SetDueAt, Value: msg.DueAt,
			},
		}

		updatedID, err := t.UpdateTaskByID(ctx, msg.TaskID, patch)
		if err != nil {
			fmt.Println("do update_task arguments: %w", err)

			return UnknownToolError, err
		}

		return fmt.Sprintf(`{"success":true,"task_id":%d}`, updatedID), nil
	default:
		fmt.Println("Неизвестная команда")
		return UnknownToolError, nil
	}
}
