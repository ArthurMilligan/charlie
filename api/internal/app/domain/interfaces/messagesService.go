package interfaces

import "context"

type MessagesService interface {
	AddMessage(ctx context.Context, content string, direction string, source string, externalMessageID *uint32) (messageID uint32, err error)
}
