package models

import (
	"time"
)

// Status, who sent message
const (
	StatusLobby = iota
	StatusPlayer
	StatusObserver
)

// Action associated with the message
const (
	Write = iota
	Update
	Delete
	StartWrite
	FinishWrite
)

// Message is the message struct
//easyjson:json
type Message struct {
	ID     int32           `json:"id"`
	User   *UserPublicInfo `json:"user"`
	Text   string          `json:"text"`
	Time   time.Time       `json:"time"`
	Status int32           `json:"status"`
	Action int32           `json:"action"`
	Edited bool            `json:"edited"`
}

// Messages slice of the messages
//easyjson:json
type Messages struct {
	Messages []*Message `json:"Messages"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
	Capacity int        `json:"capacity"`
}
