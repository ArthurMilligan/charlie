package messages

import (
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	messagesHandlers "github.com/ArthurMilligan/charlie/api/internal/app/handlers/messages"
	"github.com/go-chi/chi/v5"
)

func NewAuthRouter(r chi.Router, m interfaces.MessagesService) {
	r.Post("/add", messagesHandlers.AddMessage(m))
}
