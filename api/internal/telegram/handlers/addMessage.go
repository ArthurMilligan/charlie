package handlers

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func GetMessageHandler(m interfaces.MessagesService) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		messageID := uint32(update.Message.ID)
		_, err := m.AddMessage(ctx, update.Message.Text, "tasks", "telegram", &messageID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Ошибка(",
			})
			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   update.Message.Text,
		})
	}
}
