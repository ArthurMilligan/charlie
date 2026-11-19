package interfaces

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type RemindersStorage interface {
	AddRemind(ctx context.Context, remind models.Remind) (uint32, error)
	GetDueReminders(ctx context.Context) ([]models.DBRemind, error)
	UpdateReminderStatusByID(ctx context.Context, id uint32, status models.RemindStatus) (bool, error)
}
