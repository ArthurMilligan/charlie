package storage

import (
	"context"
	"fmt"

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

func (c *MessagesStorage) GetLastMessages(ctx context.Context, id uint32) ([]models.DBMessage, error) {
	query := `
		SELECT id, direction, content
		FROM (
			SELECT id, direction, content
			FROM messages
			WHERE id < $1
			ORDER BY id DESC
			LIMIT 20
		) AS recent
		ORDER BY id ASC;
	`
	rows, err := c.db.Query(ctx, query, id)

	if err != nil {
		fmt.Print(err)

		return nil, err
	}

	defer rows.Close()

	messages := make([]models.DBMessage, 0)

	for rows.Next() {
		var message models.DBMessage

		err := rows.Scan(
			&message.ID,
			&message.Direction,
			&message.Content,
		)
		if err != nil {
			fmt.Print(err)

			return nil, err
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	return messages, nil
}
