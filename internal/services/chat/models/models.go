package models

import (
	"time"

	//
	_ "github.com/lib/pq"
)

type UserRepository interface {
	Get() *User
	Set(*User)
}

type User struct {
	ID     int32
	Name   string
	Photo  string
	Status int32
}

type MessageRepository interface {
	Get() *Message
	Set(*Message)
}

type Message struct {
	ID     int32
	Answer *Message
	Text   string
	From   *User
	To     *User
	ChatID int32
	Time   time.Time
	Edited bool
}

type Chat struct {
	ID       int32
	Type     int32
	TypeID   int32
	Messages []*Messages
}

type Messages struct {
	Messages    []*Message
	BlockSize   int32
	BlockAmount int32
	BlockNumber int32
}

type UserInGroup struct {
	User *User
	Chat *Chat
}

type ChatWithUsers struct {
	Type   int32
	TypeID int32
	Users  []*User
}
