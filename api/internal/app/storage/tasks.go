package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TasksStorage struct {
	db *pgxpool.Pool
}

func NewTasksStorage(db *pgxpool.Pool) *TasksStorage {
	return &TasksStorage{db: db}
}

func (c *TasksStorage) AddTask(ctx context.Context, task models.Task) (uint32, error) {
	query := `
	INSERT INTO
	tasks (
		title,
		description,
		source_message_id,
		status,
		planned_on,
		planned_time,
		duration_minutes,
		due_on,
		due_at,
		completed_at
	)
	VALUES ($1, $2, $3, COALESCE($4, 'todo'), $5, $6, $7, $8, $9, $10) 
	RETURNING id
	`
	var taskID uint32

	var plannedTime pgtype.Time

	if task.PlannedTime != nil {
		t := task.PlannedTime

		sinceMidnight := time.Duration(t.Hour())*time.Hour +
			time.Duration(t.Minute())*time.Minute +
			time.Duration(t.Second())*time.Second +
			time.Duration(t.Nanosecond())

		plannedTime = pgtype.Time{
			Microseconds: sinceMidnight.Microseconds(),
			Valid:        true,
		}
	}
	err := c.db.QueryRow(ctx, query, task.Title, task.Description, task.SourceMessageID, task.Status, task.PlanedOn, plannedTime, task.DurationMinutes, task.DueOn, task.DueAt, task.CompletedAt).Scan(&taskID)

	if err != nil {
		return 0, err
	}

	return taskID, nil
}

func (c *TasksStorage) GetTasks(ctx context.Context, filter models.TaskFilter) ([]models.DBTask, error) {
	query := `
	SELECT id,
        title,
        description,
        source_message_id,
        status,
        planned_on,
        planned_time,
        duration_minutes,
        due_on,
        due_at,
        completed_at,
        created_at,
        updated_at
	FROM tasks
	WHERE ($1::text IS NULL OR status = $1)
	ORDER BY created_at ASC
	`

	rows, err := c.db.Query(ctx, query, filter.Status)

	if err != nil {
		fmt.Print(err)

		return nil, err
	}

	defer rows.Close()

	tasks := make([]models.DBTask, 0)

	for rows.Next() {
		var task models.DBTask
		var plannedTime pgtype.Time

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.SourceMessageID,
			&task.Status,
			&task.PlanedOn,
			&plannedTime,
			&task.DurationMinutes,
			&task.DueOn,
			&task.DueAt,
			&task.CompletedAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			fmt.Print(err)

			return nil, err
		}

		if plannedTime.Valid {
			duration := time.Duration(plannedTime.Microseconds) * time.Microsecond

			// Базовая дата нужна только потому, что time.Time
			// не умеет хранить время суток без даты.
			base := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
			value := base.Add(duration)

			task.PlannedTime = &value
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	return tasks, nil
}
