package services

import (
	"context"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type Reminders struct {
	RemindersStorage interfaces.RemindersStorage
}

func NewRemindersService(RemindersStorage interfaces.RemindersStorage) *Reminders {
	return &Reminders{RemindersStorage: RemindersStorage}
}

func (r *Reminders) AddRemind(
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
) (remindID uint32, err error) {
	var remind = models.Remind{
		TaskId:         taskID,
		Title:          title,
		Description:    description,
		RecurrenceType: recurrenceType,
		Timezone:       timezone,
		ChatID:         chatID,
		LocalTime:      localTime,
		StartDate:      startDate,
		EndDate:        endDate,
		NextRemindAt:   nextRemindAt,
		Status:         status,
	}

	remindID, err = r.RemindersStorage.AddRemind(ctx, remind)

	if err != nil {
		return 0, err
	}

	return remindID, nil
}

func (m *Reminders) GetDueReminders(ctx context.Context) ([]models.DBRemind, error) {

	reminders, err := m.RemindersStorage.GetDueReminders(ctx)

	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func (m *Reminders) UpdateReminderStatusByID(ctx context.Context, id uint32, status models.RemindStatus) (bool, error) {
	success, err := m.RemindersStorage.UpdateReminderStatusByID(ctx, id, status)

	if err != nil {
		return false, err
	}

	return success, nil
}
