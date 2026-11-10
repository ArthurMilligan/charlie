package interfaces

import (
	"context"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type TasksService interface {
	AddTask(ctx context.Context, title string, description *string, sourceMessageID uint32, status *string, plannedOn *time.Time, PlannedTime *time.Time, durationMinutes *int, dueOn *time.Time, dueAt *time.Time, completedAt *time.Time) (taskID uint32, err error)
	GetTasks(ctx context.Context, status *models.TaskStatus) ([]models.DBTask, error)
}
