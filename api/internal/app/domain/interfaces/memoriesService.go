package interfaces

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type MemoriesService interface {
	AddMemories(
		ctx context.Context,
		chatID int64,
		sourceType models.MemorySourceType,
		sourceID uint32,
		text string,
	) ([]uint32, error)
	FindMemories(ctx context.Context, text string, chatID uint32) ([]models.DBMemory, error)
	GetMemoryBySourceId(ctx context.Context, sourceID uint32, sourceType models.MemorySourceType) ([]models.DBMemory, error)
	GetMemoryContext(ctx context.Context, memoryID uint32) ([]models.DBMemory, error)
}
