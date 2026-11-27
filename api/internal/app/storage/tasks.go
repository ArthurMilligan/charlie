package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
	"github.com/jackc/pgx/v5"
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

func (c *TasksStorage) UpdateTaskByID(ctx context.Context, taskID uint32, task models.TaskPatch) (uint32, error) {
	const query = `
   	UPDATE tasks
	SET
		title            = COALESCE(@title, title),
		status           = COALESCE(@status, status),
		description      = CASE WHEN @set_description      THEN @description      ELSE description END,
		planned_on       = CASE WHEN @set_planned_on       THEN @planned_on       ELSE planned_on END,
		planned_time     = CASE WHEN @set_planned_time     THEN @planned_time     ELSE planned_time END,
		duration_minutes = CASE WHEN @set_duration_minutes THEN @duration_minutes ELSE duration_minutes END,
		due_on           = CASE WHEN @set_due_on           THEN @due_on           ELSE due_on END,
		due_at           = CASE WHEN @set_due_at           THEN @due_at           ELSE due_at END,
		completed_at = CASE
			   	WHEN @status = 'done' THEN COALESCE(@completed_at::timestamptz, now())
				WHEN @status IS NULL THEN completed_at
				ELSE NULL
		END,
		updated_at = now ()
	WHERE id = @id
	RETURNING id;
`

	args := pgx.StrictNamedArgs{
		"id":                   taskID,
		"title":                task.Title,
		"status":               task.Status,
		"set_description":      task.Description.Set,
		"description":          task.Description.Value,
		"planned_on":           task.PlannedOn.Value,
		"set_planned_on":       task.PlannedOn.Set,
		"planned_time":         task.PlannedTime.Value,
		"set_planned_time":     task.PlannedTime.Set,
		"duration_minutes":     task.DurationMinutes.Value,
		"set_duration_minutes": task.DurationMinutes.Set,
		"due_on":               task.DueOn.Value,
		"set_due_on":           task.DueOn.Set,
		"due_at":               task.DueAt.Value,
		"set_due_at":           task.DueAt.Set,
		"completed_at":         task.CompletedAt,
	}

	var updatedID uint32
	err := c.db.QueryRow(ctx, query, args).Scan(&updatedID)

	if err != nil {
		return 0, err
	}

	return updatedID, nil
}
