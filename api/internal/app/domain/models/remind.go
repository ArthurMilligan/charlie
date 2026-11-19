package models

import "time"

type Remind struct {
	TaskId         *uint32
	Title          string
	Description    *string
	RecurrenceType *RecurrenceType
	Timezone       string
	ChatID         uint32
	LocalTime      time.Time
	StartDate      time.Time
	EndDate        *time.Time
	NextRemindAt   time.Time
	Status         *RemindStatus
}

type DBRemind struct {
	ID             uint32
	TaskId         *uint32
	Title          string
	Description    *string
	RecurrenceType *RecurrenceType
	Timezone       string
	ChatID         uint32
	LocalTime      time.Time
	StartDate      time.Time
	EndDate        *time.Time
	NextRemindAt   time.Time
	Status         *RemindStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RecurrenceType string

const (
	TypeOnce    RecurrenceType = "once"
	TypeDaily   RecurrenceType = "daily"
	TypeWeekly  RecurrenceType = "weekly"
	TypeMonthly RecurrenceType = "monthly"
	TypeYearly  RecurrenceType = "yearly"
)

type RemindStatus string

const (
	RemindStatusActive    RemindStatus = "active"
	RemindStatusPaused    RemindStatus = "paused"
	RemindStatusCompleted RemindStatus = "completed"
	RemindStatusCancelled RemindStatus = "cancelled"
)
