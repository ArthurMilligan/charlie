package interfaces

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type MessagesStorage interface {
	AddMessage(ctx context.Context, message models.Message) (uint32, error)
}
