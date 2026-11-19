package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

func UseTools(ctx context.Context, n interfaces.NotesService, t interfaces.TasksService, r interfaces.RemindersService, toolName string, args string, messageID uint32, chatID int64) (string, error) {
	switch toolName {
	case "add_note":
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
	case "get_notes":
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
	default:
		fmt.Println("Неизвестная команда")
		return UnknownToolError, nil
	}
}
