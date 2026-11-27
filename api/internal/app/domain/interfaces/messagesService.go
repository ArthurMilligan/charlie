package interfaces

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type MessagesService interface {
	AddMessage(ctx context.Context, content string, direction string, source string, externalMessageID *uint32) (messageID uint32, err error)
	GetLastMessages(ctx context.Context, id uint32) ([]models.DBMessage, error)
}
