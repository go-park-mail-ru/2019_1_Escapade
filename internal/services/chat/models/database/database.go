package database

import (
	"database/sql"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/models"
)

/*
User - wrapper to retrieve user data from the Database, given that
 the User may not exist
*/
type User struct {
	ID     sql.NullInt32
	Name   sql.NullString
	Photo  sql.NullString
	Status sql.NullInt32
}

func (user *User) Get() *models.User {
	if user == nil {
		return nil
	}
	if user.ID.Valid &&
		user.Name.Valid &&
		user.Photo.Valid &&
		user.Status.Valid {
		return &models.User{
			ID:     user.ID.Int32,
			Name:   user.Name.String,
			Photo:  user.Photo.String,
			Status: user.Status.Int32,
		}
	}
	return nil
}

func (user *User) Set(m *models.User) {
	if m == nil {
		user.ID.Scan(nil)
		user.Name.Scan(nil)
		user.Photo.Scan(nil)
		user.Status.Scan(nil)
	} else {
		user.ID.Int32 = m.ID
		user.Name.String = m.Name
		user.Photo.String = m.Photo
		user.Status.Int32 = m.Status
	}
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
Message - wrapper to retrieve Message data from the Database, given that
 the Message may not exist
*/
type Message struct {
	ID     sql.NullInt32
	Answer *Message
	Text   sql.NullString
	From   *User
	To     *User
	ChatID sql.NullInt32
	Time   NullTime
	Edited sql.NullBool
}

func (message *Message) Get() *models.Message {
	if message == nil {
		return nil
	}
	if message.ID.Valid &&
		message.Text.Valid &&
		message.ChatID.Valid &&
		message.Edited.Valid {
		return &models.Message{
			ID:     message.ID.Int32,
			Answer: message.Answer.Get(),
			Text:   message.Text.String,
			From:   message.From.Get(),
			To:     message.To.Get(),
			ChatID: message.ChatID.Int32,
			Time:   time.Time(message.Time),
			Edited: message.Edited.Bool,
		}
	}
	return nil
}

func (message *Message) Set(m *models.Message) {
	if m == nil {
		message.ID.Scan(nil)
		message.Answer = nil
		message.Text.Scan(nil)
		message.From = nil
		message.To = nil
		message.ChatID.Scan(nil)
		message.Edited.Scan(nil)
	} else {
		message.ID.Int32 = m.ID
		if m.Answer != nil {
			message.Answer = new(Message)
			message.Answer.Set(m.Answer)
		}
		message.Text.String = m.Text
		if m.From != nil {
			message.From = new(User)
			message.From.Set(m.From)
		}
		if m.To != nil {
			message.To = new(User)
			message.To.Set(m.To)
		}
		message.ChatID.Int32 = m.ChatID
		message.Time = NullTime(m.Time)
		message.Edited.Bool = m.Edited
	}
}
