package models

import "time"

// Cell type
const (
	StatusLobby = iota
	StatusPlayer
	StatusObserver
)

// Message is the message struct
type Message struct {
	User   *UserPublicInfo `json:"user"`
	Text   string          `json:"text"`
	Time   time.Time       `json:"time"`
	Status int             `json:"status"`
}
