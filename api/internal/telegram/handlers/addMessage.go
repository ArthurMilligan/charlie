package handlers

import (
	"context"

	gptInterfaces "github.com/ArthurMilligan/charlie/api/internal/ai/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func GetMessageHandler(m interfaces.MessagesService, n interfaces.NotesService, g gptInterfaces.GPTClient) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		telegramMessageID := uint32(update.Message.ID)
		savedMessageID, err := m.AddMessage(ctx, update.Message.Text, "input", "telegram", &telegramMessageID)

		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Ошибка(",
			})
			return
		}

		answer, err := g.GetAnswer(ctx, update.Message.Text, savedMessageID)

		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Ошибка(",
			})
			return
		}
		_, err = m.AddMessage(ctx, answer, "output", "telegram", &telegramMessageID)

		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Ошибка(",
			})
			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   answer,
		})
	}
}
