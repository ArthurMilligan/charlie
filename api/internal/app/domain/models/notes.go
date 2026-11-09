package models

import "time"

type Note struct {
	Content         string
	SourceMessageID uint32
}

type DBNote struct {
	ID              uint32
	Content         string
	SourceMessageID uint32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
