package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

func NewHttpServer() nethttp.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	return r
}
