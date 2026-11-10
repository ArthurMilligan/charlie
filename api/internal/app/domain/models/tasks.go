package models

import "time"

type Task struct {
	Title           string
	Description     *string
	SourceMessageID uint32
	Status          *string
	PlanedOn        *time.Time
	PlannedTime     *time.Time
	DurationMinutes *int
	DueOn           *time.Time
	DueAt           *time.Time
	CompletedAt     *time.Time
}

type DBTask struct {
	ID              uint32
	Title           string
	Description     *string
	SourceMessageID uint32
	Status          string     //in_progress, done, cancelled,todo
	PlanedOn        *time.Time //День, на который пользователь запланировал работу.
	PlannedTime     *time.Time //Время начала работы в запланированный день.
	DurationMinutes *int       //Сколько минут пользователь планирует потратить.
	DueOn           *time.Time //Дедлайн с точностью до дня
	DueAt           *time.Time //Дедлайн с точностью до конкретного момента.
	CompletedAt     *time.Time //Момент, когда задачу отметили выполненной.
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type TaskStatus string

const (
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
	StatusCancelled  TaskStatus = "cancelled"
	StatusToDo       TaskStatus = "todo"
)

type TaskFilter struct {
	Status *TaskStatus
}
