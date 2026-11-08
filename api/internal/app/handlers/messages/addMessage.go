package messages

import (
	"encoding/json"
	"net/http"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
)

type addMessageRequest struct {
	Content           string  `json:"content"`
	Direction         string  `json:"direction"`
	Source            string  `json:"source"`
	ExternalMessageID *uint32 `json:"external_message_id"`
}

func AddMessage(m interfaces.MessagesService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addMessageRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		messageID, err := m.AddMessage(r.Context(), req.Content, req.Direction, req.Source, req.ExternalMessageID)
		if err != nil {
			http.Error(w, "invalid args", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"message_id": messageID,
		})
	}
}
