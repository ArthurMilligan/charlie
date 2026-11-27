package models

import "time"

type Memory struct {
	SourceType MemorySourceType
	SourceID   uint32
	ChunkIndex uint32
	Content    string
	Embedding  []float64
	ChatID     int64
}

type DBMemory struct {
	ID         uint32
	SourceType MemorySourceType
	SourceID   uint32
	ChunkIndex uint32
	Content    string
	Embedding  []float64
	ChatID     int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MemorySourceType string

const (
	MemorySourceTypeNote MemorySourceType = "note"
	MemorySourceTypeChat MemorySourceType = "chat"
	MemorySourceTypeLore MemorySourceType = "lore"
)
