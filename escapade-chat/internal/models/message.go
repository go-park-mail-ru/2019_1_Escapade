package models

import "time"

type Message struct {
	Type    string          `json:"type"`
	User    *UserPublicInfo `json:"user"`
	Message string          `json:"message"`
	Time    time.Time       `json:"time"`
}

type Messages struct {
	Type     string     `json:"type"`
	Messages []*Message `json:"messages"`
	Capacity int        `json:"capacity"`
}
