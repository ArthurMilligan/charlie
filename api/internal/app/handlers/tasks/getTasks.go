package notes

import (
	"encoding/json"
	"net/http"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type getTasksRequest struct {
	Status *models.TaskStatus `json:"status"`
}

func GetTasks(n interfaces.TasksService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req getTasksRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		tasks, err := n.GetTasks(r.Context(), req.Status)
		if err != nil {
			http.Error(w, "invalid args", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"tasks": tasks,
		})
	}
}
