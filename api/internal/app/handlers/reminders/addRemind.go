package reminders

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type addRemindRequest struct {
	TaskID         *uint32                `json:"task_id"`
	Title          string                 `json:"title"`
	Description    *string                `json:"description"`
	RecurrenceType *models.RecurrenceType `json:"recurrence_type"`
	Timezone       string                 `json:"timezone"`
	ChatID         uint32                 `json:"chat_id"`
	LocalTime      time.Time              `json:"local_time"`
	StartDate      time.Time              `json:"start_date"`
	EndDate        *time.Time             `json:"end_date"`
	NextRemindAt   time.Time              `json:"next_remind_at"`
	Status         *models.RemindStatus   `json:"status"`
}

func AddRemind(s interfaces.RemindersService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addRemindRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		remindID, err := s.AddRemind(r.Context(), req.TaskID, req.Title, req.Description, req.RecurrenceType, req.Timezone, req.ChatID, req.LocalTime, req.StartDate, req.EndDate, req.NextRemindAt, req.Status)
		if err != nil {
			http.Error(w, "invalid args", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"remind_id": remindID,
		})
	}
}
