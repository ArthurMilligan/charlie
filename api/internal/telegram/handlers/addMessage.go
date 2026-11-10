package handlers

import (
	"context"
	"fmt"

	gptInterfaces "github.com/ArthurMilligan/charlie/api/internal/ai/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func GetMessageHandler(m interfaces.MessagesService, n interfaces.NotesService, g gptInterfaces.GPTClient) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		telegramMessageID := uint32(update.Message.ID)

		msg := update.Message

		if msg.Chat.Type != "private" {
			return
		}

		if msg.From == nil || (msg.From.ID != 554822003 && msg.From.ID != 1821406104) {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: msg.Chat.ID,
				Text:   "Прошу простить, но я исполняю лишь поручения моего хозяина. Вашу просьбу вынужден отклонить.",
			})
			if err != nil {
				// Зафиксировать ошибку отправки без токенов.
			}
			return
		}

		savedMessageID, err := m.AddMessage(ctx, update.Message.Text, "input", "telegram", &telegramMessageID)

		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Ошибка(",
			})
			return
		}

		answer, err := g.GetAnswer(ctx, update.Message.Text, savedMessageID)

		fmt.Print(err)

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
