package reminders

import (
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	remindersHandlers "github.com/ArthurMilligan/charlie/api/internal/app/handlers/reminders"
	"github.com/go-chi/chi/v5"
)

func NewRemindersRouter(r chi.Router, n interfaces.RemindersService) {
	r.Post("/add", remindersHandlers.AddRemind(n))
}
