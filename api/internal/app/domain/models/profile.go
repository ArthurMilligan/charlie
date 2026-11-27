package models

import (
	"encoding/json"
	"time"
)

type ProfileSection struct {
	UserID      int64
	SectionKey  string
	Title       string
	Description string
}

type DBProfileSection struct {
	ID          int64
	UserID      int64
	SectionKey  string
	Title       string
	Description string
}

type ProfileFact struct {
	UserID          int64
	SectionID       int64
	FactKey         string
	Value           json.RawMessage
	SourceMessageID int64
	UpdatedAt       *time.Time
}

type ProfileFactPatch struct {
	UserID    int64
	SectionID int64
	FactKey   string
	SetValue  bool
	Value     *json.RawMessage
}
