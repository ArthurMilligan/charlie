package notes

import (
	"encoding/json"
	"net/http"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
)

func GetNote(n interfaces.NotesService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		notes, err := n.GetNotes(r.Context())
		if err != nil {
			http.Error(w, "invalid args", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"notes": notes,
		})
	}
}
