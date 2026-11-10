package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/http/messages"
	"github.com/ArthurMilligan/charlie/api/internal/app/http/notes"
	"github.com/ArthurMilligan/charlie/api/internal/app/http/tasks"
	"github.com/go-chi/chi/v5"
)

func NewHttpServer(m interfaces.MessagesService, n interfaces.NotesService, t interfaces.TasksService) nethttp.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	r.Route("/messages", func(r chi.Router) {
		messages.NewMessagesRouter(r, m)
	})

	r.Route("/notes", func(r chi.Router) {
		notes.NewNotesRouter(r, n)
	})

	r.Route("/tasks", func(r chi.Router) {
		tasks.NewTasksRouter(r, t)
	})

	return r
}
