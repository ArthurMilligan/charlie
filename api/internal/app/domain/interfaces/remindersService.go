package interfaces

import (
	"context"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type RemindersService interface {
	AddRemind(
		ctx context.Context,
		taskID *uint32,
		title string,
		description *string,
		recurrenceType *models.RecurrenceType,
		timezone string,
		chatID uint32,
		localTime time.Time,
		startDate time.Time,
		endDate *time.Time,
		nextRemindAt time.Time,
		status *models.RemindStatus,
	) (remindID uint32, err error)
	GetDueReminders(ctx context.Context) ([]models.DBRemind, error)
	UpdateReminderStatusByID(ctx context.Context, id uint32, status models.RemindStatus) (bool, error)
}
