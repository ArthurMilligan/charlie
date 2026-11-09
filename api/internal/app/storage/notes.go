package storage

import (
	"context"
	"fmt"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotesStorage struct {
	db *pgxpool.Pool
}

func NewNotesStorage(db *pgxpool.Pool) *NotesStorage {
	return &NotesStorage{db: db}
}

func (c *NotesStorage) AddNote(ctx context.Context, note models.Note) (uint32, error) {
	query := `
	INSERT INTO notes (content, source_message_id)
	VALUES ($1, $2) 
	RETURNING id
	`
	var noteID uint32

	err := c.db.QueryRow(ctx, query, note.Content, note.SourceMessageID).Scan(&noteID)

	if err != nil {
		return 0, err
	}

	return noteID, nil
}

func (c *NotesStorage) GetNotes(ctx context.Context) ([]models.DBNote, error) {
	query := `
	SELECT
		id,
		content,
		source_message_id,
		created_at,
		updated_at
	FROM
		notes
	ORDER BY
		created_at ASC
	`

	rows, err := c.db.Query(ctx, query)

	if err != nil {
		fmt.Print(err)

		return nil, err
	}

	defer rows.Close()

	notes := make([]models.DBNote, 0)

	for rows.Next() {
		var note models.DBNote

		err := rows.Scan(
			&note.ID,
			&note.Content,
			&note.SourceMessageID,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			fmt.Print(err)

			return nil, err
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	return notes, nil
}
