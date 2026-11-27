package interfaces

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type MessagesStorage interface {
	AddMessage(ctx context.Context, message models.Message) (uint32, error)
	GetLastMessages(ctx context.Context, id uint32) ([]models.DBMessage, error)
}
