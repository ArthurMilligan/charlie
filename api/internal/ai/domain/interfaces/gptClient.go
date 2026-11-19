package interfaces

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type GPTClient interface {
	GetAnswer(ctx context.Context, text string, messageID uint32, chatID int64) (string, error)
	FormatReminder(
		ctx context.Context,
		reminder models.DBRemind,
	) (string, error)
}
