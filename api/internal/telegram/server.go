package telegram

import (
	gptInterfaces "github.com/ArthurMilligan/charlie/api/internal/ai/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/telegram/handlers"
	"github.com/go-telegram/bot"
)

func GetTelegramBot(m interfaces.MessagesService, n interfaces.NotesService, token string, g gptInterfaces.GPTClient) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.GetMessageHandler(m, n, g)),
	}

	b, err := bot.New(token, opts...)
	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.
		return nil, err
	}

	return b, nil
}
