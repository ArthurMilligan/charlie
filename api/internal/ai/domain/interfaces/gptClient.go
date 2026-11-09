package interfaces

import "context"

type GPTClient interface {
	GetAnswer(ctx context.Context, text string, messageID uint32) (string, error)
}
