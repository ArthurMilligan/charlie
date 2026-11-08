package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/http/messages"
	"github.com/go-chi/chi/v5"
)

func NewHttpServer(m interfaces.MessagesService) nethttp.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	r.Route("/messages", func(r chi.Router) {
		messages.NewAuthRouter(r, m)
	})

	return r
}
