package interfaces

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type TasksStorage interface {
	AddTask(ctx context.Context, task models.Task) (uint32, error)
	GetTasks(ctx context.Context, filter models.TaskFilter) ([]models.DBTask, error)
}
