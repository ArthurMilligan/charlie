package interfaces

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type NotesStorage interface {
	AddNote(ctx context.Context, message models.Note) (uint32, error)
	GetNotes(ctx context.Context) ([]models.DBNote, error)
}
