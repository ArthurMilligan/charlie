package services

import (
	"context"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type Tasks struct {
	TasksStorage interfaces.TasksStorage
}

func NewTasksService(TasksStorage interfaces.TasksStorage) *Tasks {
	return &Tasks{TasksStorage: TasksStorage}
}

func (m *Tasks) AddTask(ctx context.Context, title string, description *string, sourceMessageID uint32, status *string, plannedOn *time.Time, PlannedTime *time.Time, durationMinutes *int, dueOn *time.Time, dueAt *time.Time, completedAt *time.Time) (taskID uint32, err error) {
	var task = models.Task{
		Title:           title,
		Description:     description,
		SourceMessageID: sourceMessageID,
		Status:          status,
		PlanedOn:        plannedOn,
		PlannedTime:     PlannedTime,
		DurationMinutes: durationMinutes,
		DueOn:           dueOn,
		DueAt:           dueAt,
		CompletedAt:     completedAt,
	}

	taskID, err = m.TasksStorage.AddTask(ctx, task)

	if err != nil {
		return 0, err
	}

	return taskID, nil
}

func (m *Tasks) GetTasks(ctx context.Context, status *models.TaskStatus) ([]models.DBTask, error) {
	var filter = models.TaskFilter{
		Status: status,
	}

	tasks, err := m.TasksStorage.GetTasks(ctx, filter)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
