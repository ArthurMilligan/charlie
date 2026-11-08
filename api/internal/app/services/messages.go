package services

import (
	"context"
	"errors"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

var (
	ErrInvalidArgs = errors.New("invalid args")
)

type Messages struct {
	MessagesStorage interfaces.MessagesStorage
}

func NewMessagesService(MessagesStorage interfaces.MessagesStorage) *Messages {
	return &Messages{MessagesStorage: MessagesStorage}
}

func (m *Messages) AddMessage(ctx context.Context, content string, direction string, source string, externalMessageID *uint32) (messageID uint32, err error) {
	if content == "" || source == "" || direction == "" {
		return 0, ErrInvalidArgs
	}

	if externalMessageID != nil && *externalMessageID == 0 {
		return 0, ErrInvalidArgs
	}

	var message = models.Message{
		Content:           content,
		Direction:         direction,
		Source:            source,
		ExternalMessageID: externalMessageID,
	}

	messageID, err = m.MessagesStorage.AddMessage(ctx, message)

	if err != nil {
		return 0, err
	}

	return messageID, nil
}
