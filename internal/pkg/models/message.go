package models

import (
	sql "database/sql"
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

// NullTime overriding the time type.Time to be able to retrieve time
//  from the database, even if the corresponding field is nil
type NullTime time.Time

// Scan allow to fill in a field of the type NullTime from the database
func (t *NullTime) Scan(v interface{}) error {
	if v == nil {
		*t = NullTime(time.Now())
		return nil
	}
	vt, err := time.Parse("2006-01-02 15:04:05 +300 MSK", v.(time.Time).String())
	if err != nil {
		return err
	}
	*t = NullTime(vt)
	return nil
}

/*
MessageUserSQL - wrapper to retrieve user data from the Database, given that
 the User may not exist
*/
type MessageUserSQL struct {
	ID     sql.NullInt64  `json:"-"`
	Name   sql.NullString `json:"-"`
	Photo  sql.NullString `json:"-"`
	Status sql.NullInt64  `json:"-"`
}

/*
MessageSQL - wrapper to retrieve Message data from the Database, given that
 the Message may not exist
*/
type MessageSQL struct {
	ID     sql.NullInt64   `json:"-"`
	Answer *MessageSQL     `json:"-"`
	Text   sql.NullString  `json:"-"`
	From   *MessageUserSQL `json:"-"`
	To     *MessageUserSQL `json:"-"`
	ChatID sql.NullInt64   `json:"-"`
	Time   time.Time       `json:"-"`
	Edited sql.NullBool    `json:"-"`
}

// Messages slice of the messages
//easyjson:json
type Messages struct {
	Messages []*Message `json:"Messages"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
	Capacity int        `json:"capacity"`
}
