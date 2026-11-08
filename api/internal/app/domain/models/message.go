package models

type Message struct {
	Content           string
	Direction         string
	Source            string
	ExternalMessageID *uint32
}
