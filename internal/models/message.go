package models

import "time"

// Message is the message struct
type Message struct {
	User *UserPublicInfo `json:"user"`
	Text string          `json:"text"`
	Time time.Time       `json:"time"`
}
