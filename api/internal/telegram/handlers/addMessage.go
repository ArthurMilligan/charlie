package handlers

import (
	"context"
	"log"
	"strings"
	"unicode"

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

		_, typingErr := b.SendChatAction(ctx, &bot.SendChatActionParams{
			ChatID: msg.Chat.ID,
			Action: models.ChatActionTyping,
		})

		if typingErr != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Ошибка(",
			})
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

		answer, err := g.GetAnswer(ctx, update.Message.Text, savedMessageID, msg.From.ID)

		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Ошибка(",
			})
			return
		}

		for i, part := range splitTelegramText(answer) {
			_, err = m.AddMessage(ctx, part, "output", "telegram", &telegramMessageID)

			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Ошибка(",
				})
				return
			}

			_, err := b.SendRichMessage(ctx, &bot.SendRichMessageParams{
				ChatID: msg.Chat.ID,
				RichMessage: models.InputRichMessage{
					Markdown: part,
				},
			})

			if err != nil {
				log.Printf("telegram: failed to send part %d of reply to message %d (%T)",
					i+1, telegramMessageID, err)
				return
			}
		}
	}
}

func splitTelegramText(text string) []string {
	const limit = 3900

	remaining := []rune(strings.TrimSpace(text))
	var parts []string

	for len(remaining) > 0 {
		end, size := 0, 0

		for end < len(remaining) {
			width := 1
			if remaining[end] > 0xFFFF {
				width = 2 // Запас для эмодзи и других символов вне BMP.
			}
			if size+width > limit {
				break
			}
			size += width
			end++
		}

		if end < len(remaining) {
			cut := end
			for i := end; i > end/2; i-- {
				if remaining[i-1] == '\n' && remaining[i] == '\n' {
					cut = i
					break
				}
				if cut == end && unicode.IsSpace(remaining[i]) {
					cut = i
				}
			}
			end = cut
		}

		part := strings.TrimSpace(string(remaining[:end]))
		if part != "" {
			parts = append(parts, part)
		}
		remaining = remaining[end:]
	}

	return parts
}
