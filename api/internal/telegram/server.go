package telegram

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/telegram/handlers"
	"github.com/go-telegram/bot"
)

func NewTelegramConnection(m interfaces.MessagesService, token string) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.GetMessageHandler(m)),
	}

	b, err := bot.New(token, opts...)
	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.
		panic(err)
	}

	b.Start(context.TODO())
}
