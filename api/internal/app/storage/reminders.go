package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RemindersStorage struct {
	db *pgxpool.Pool
}

func NewRemindersStorage(db *pgxpool.Pool) *RemindersStorage {
	return &RemindersStorage{db: db}
}

func (r *RemindersStorage) AddRemind(ctx context.Context, remind models.Remind) (uint32, error) {
	query := `
	INSERT INTO
		reminders (
			task_id,
			title,
			description,
			recurrence_type,
			timezone,
			chat_id,
			local_time,
			start_date,
			end_date,
			next_remind_at,
			status
		)
	VALUES ($1, $2, $3, COALESCE($4, 'once'), $5, $6, $7, $8, $9, $10, COALESCE($11, 'active')) 
	RETURNING id
	`

	var remindID uint32

	var localTime pgtype.Time

	t := remind.LocalTime

	sinceMidnight := time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second +
		time.Duration(t.Nanosecond())

	localTime = pgtype.Time{
		Microseconds: sinceMidnight.Microseconds(),
		Valid:        true,
	}

	err := r.db.QueryRow(ctx, query, remind.TaskId, remind.Title, remind.Description, remind.RecurrenceType, remind.Timezone, remind.ChatID, localTime, remind.StartDate, remind.EndDate, remind.NextRemindAt, remind.Status).Scan(&remindID)

	if err != nil {
		return 0, err
	}

	return remindID, nil
}

func (r *RemindersStorage) GetDueReminders(ctx context.Context) ([]models.DBRemind, error) {
	query := `
	SELECT
		id,
		task_id,
		title,
		description,
		recurrence_type,
		timezone,
		chat_id,
		local_time,
		start_date,
		end_date,
		next_remind_at,
		status,
		created_at,
		updated_at
	FROM
		reminders
	WHERE
		status = 'active'
		AND next_remind_at <= now ()
	`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		fmt.Print(err)

		return nil, err
	}

	defer rows.Close()

	reminders := make([]models.DBRemind, 0)

	for rows.Next() {
		var reminder models.DBRemind
		var localTime pgtype.Time

		err := rows.Scan(
			&reminder.ID,
			&reminder.TaskId,
			&reminder.Title,
			&reminder.Description,
			&reminder.RecurrenceType,
			&reminder.Timezone,
			&reminder.ChatID,
			&localTime,
			&reminder.StartDate,
			&reminder.EndDate,
			&reminder.NextRemindAt,
			&reminder.Status,
			&reminder.CreatedAt,
			&reminder.UpdatedAt,
		)
		if err != nil {
			fmt.Print(err)

			return nil, err
		}

		if localTime.Valid {
			duration := time.Duration(localTime.Microseconds) * time.Microsecond

			// Базовая дата нужна только потому, что time.Time
			// не умеет хранить время суток без даты.
			base := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
			value := base.Add(duration)

			reminder.LocalTime = value
		}

		reminders = append(reminders, reminder)
	}

	if err := rows.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	return reminders, nil
}

func (r *RemindersStorage) UpdateReminderStatusByID(ctx context.Context, id uint32, status models.RemindStatus) (bool, error) {
	query := `
		UPDATE reminders
		SET status = $1
		WHERE id = $2
	`

	tag, err := r.db.Exec(ctx, query, status, id)

	if err != nil {
		return false, err
	}

	if tag.RowsAffected() == 0 {
		return false, err
	}

	return true, nil
}
