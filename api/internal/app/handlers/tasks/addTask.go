package notes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
)

type addTaskRequest struct {
	SourceMessageID uint32     `json:"source_message_id"`
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

func AddTask(n interfaces.TasksService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addTaskRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		taskID, err := n.AddTask(r.Context(), req.Title, req.Description, req.SourceMessageID, req.Status, req.PlanedOn, req.PlannedTime, req.DurationMinutes, req.DueOn, req.DueAt, req.CompletedAt)
		if err != nil {
			fmt.Print(err)
			http.Error(w, "invalid args", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"task_id": taskID,
		})
	}
}
