package models

import "time"

const (
	StatusLobby = iota
	StatusPlayer
	StatusObserver
)

const (
	Write = iota
	Update
	Delete
	StartWrite
	FinishWrite
)

// Message is the message struct
type Message struct {
	ID     int             `json:"id"`
	User   *UserPublicInfo `json:"user"`
	Text   string          `json:"text"`
	Time   time.Time       `json:"time"`
	Status int             `json:"status"`
	Action int             `json:"action"`
}

type Messages struct {
	Messages []*Message `json:"Messages"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
	Capacity int        `json:"capacity"`
}
