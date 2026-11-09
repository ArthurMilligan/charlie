package notes

import (
	"encoding/json"
	"net/http"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
)

type addNoteRequest struct {
	Content         string `json:"content"`
	SourceMessageID uint32 `json:"source_message_id"`
}

func AddNote(n interfaces.NotesService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addNoteRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		noteID, err := n.AddNote(r.Context(), req.Content, req.SourceMessageID)
		if err != nil {
			http.Error(w, "invalid args", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"note_id": noteID,
		})
	}
}
