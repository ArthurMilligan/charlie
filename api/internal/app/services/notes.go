package services

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type Notes struct {
	NotesStorage interfaces.NotesStorage
}

func NewNotesService(NotesStorage interfaces.NotesStorage) *Notes {
	return &Notes{NotesStorage: NotesStorage}
}

func (m *Notes) AddNote(ctx context.Context, content string, sourceMessageID uint32) (noteID uint32, err error) {
	if content == "" || sourceMessageID == 0 {
		return 0, ErrInvalidArgs
	}

	var note = models.Note{
		Content:         content,
		SourceMessageID: sourceMessageID,
	}

	noteID, err = m.NotesStorage.AddNote(ctx, note)

	if err != nil {
		return 0, err
	}

	return noteID, nil
}

func (m *Notes) GetNotes(ctx context.Context) ([]models.DBNote, error) {
	notes, err := m.NotesStorage.GetNotes(ctx)

	if err != nil {
		return nil, err
	}

	return notes, nil
}
