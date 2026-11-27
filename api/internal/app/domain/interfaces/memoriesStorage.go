package interfaces

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type MemoriesStorage interface {
	AddMemories(ctx context.Context, memories []models.Memory) ([]uint32, error)
	FindMemories(ctx context.Context, text string, embedding []float64, chatID uint32) ([]models.DBMemory, error)
	GetMemoryBySourceId(ctx context.Context, sourceID uint32, sourceType models.MemorySourceType) ([]models.DBMemory, error)
	GetMemoryById(ctx context.Context, id uint32) (models.DBMemory, error)
}
