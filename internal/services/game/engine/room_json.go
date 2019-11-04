package engine

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// RoomJSON is a wrapper for sending Room by JSON
//easyjson:json
type RoomJSON struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`

	Players   OnlinePlayersJSON `json:"players"`
	Observers ConnectionsJSON   `json:"observers,omitempty"`
	History   []*PlayerAction   `json:"history,omitempty"`
	Messages  []*models.Message `json:"messages"`

	Field    FieldJSON            `json:"field,omitempty"`
	Date     time.Time            `json:"date,omitempty"`
	Settings *models.RoomSettings `json:"settings"`
}

// JSON convert Room to RoomJSON
func (room *Room) JSON() RoomJSON {
	return RoomJSON{
		ID:        room.ID(),
		Name:      room.Name(),
		Status:    room.events.Status(),
		Players:   room.people.Players.JSON(),
		Observers: room.people.Observers.JSON(),
		History:   room.connEvents.notify.history(),
		Messages:  room.messages.Messages(),
		Field:     room.field.JSON(),
		Date:      room.events.Date(),
		Settings:  room.Settings,
	}
}

// MarshalJSON - overriding the standard method json.Marshal
func (room *Room) MarshalJSON() ([]byte, error) {
	return room.JSON().MarshalJSON()
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (room *Room) UnmarshalJSON(b []byte) error {
	temp := &RoomJSON{}

	if err := temp.UnmarshalJSON(b); err != nil {
		return err
	}
	room.setName(temp.Name)
	room.events.setStatus(temp.Status)
	room.connEvents.notify._history = temp.History
	room.messages.setMessages(temp.Messages)
	room.events.setDate(temp.Date)
	room.Settings = temp.Settings

	return nil
}
