package interfaces

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type NotesService interface {
	AddNote(ctx context.Context, content string, sourceMessageID uint32) (noteID uint32, err error)
	GetNotes(ctx context.Context) ([]models.DBNote, error)
}
