package notes

import (
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	notesHandlers "github.com/ArthurMilligan/charlie/api/internal/app/handlers/notes"
	"github.com/go-chi/chi/v5"
)

func NewNotesRouter(r chi.Router, n interfaces.NotesService) {
	r.Post("/add", notesHandlers.AddNote(n))
	r.Post("/get", notesHandlers.GetNote(n))
}
