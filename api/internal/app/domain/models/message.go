package models

import "time"

type Message struct {
	Content           string
	Direction         string
	Source            string
	ExternalMessageID *uint32
}

type DBMessage struct {
	ID                uint32
	Content           string
	Direction         string
	Source            string
	ExternalMessageID *uint32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
