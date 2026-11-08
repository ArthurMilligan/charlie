package storage

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessagesStorage struct {
	db *pgxpool.Pool
}

func NewMessagesStorage(db *pgxpool.Pool) *MessagesStorage {
	return &MessagesStorage{db: db}
}

func (c *MessagesStorage) AddMessage(ctx context.Context, message models.Message) (uint32, error) {
	query := `
	INSERT INTO messages (content, direction, source, external_message_id)
	VALUES ($1, $2, $3, $4) 
	RETURNING id
	`
	var messageID uint32

	err := c.db.QueryRow(ctx, query, message.Content, message.Direction, message.Source, message.ExternalMessageID).Scan(&messageID)

	if err != nil {
		return 0, err
	}

	return messageID, nil
}
