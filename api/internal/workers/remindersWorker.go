package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	gptInterfaces "github.com/ArthurMilligan/charlie/api/internal/ai/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
	"github.com/go-telegram/bot"
)

func RunRemindersWorker(ctx context.Context, r interfaces.RemindersService, telegramBot *bot.Bot, ai gptInterfaces.GPTClient) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			reminders, err := r.GetDueReminders(ctx)
			if err != nil {
				log.Printf("process due reminders: %v", err)
				continue
			}

			if len(reminders) == 0 {
				continue
			}

			for _, reminder := range reminders {
				message, err := ai.FormatReminder(ctx, reminder)
				if err != nil {
					fmt.Print(err)
					continue
				}

				_, err = telegramBot.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: reminder.ChatID,
					Text:   message,
				})
				if err != nil {
					log.Printf("send reminder %d: %v", reminder.ID, err)
					continue
				}
				r.UpdateReminderStatusByID(ctx, reminder.ID, models.RemindStatusCompleted)
			}

			fmt.Print(reminders)
		}
	}
}
