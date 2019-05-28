package game

import (
	"time"

	"encoding/json"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// Room consist of players and observers, field and history
type RoomJSON struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`

	Players   *OnlinePlayers    `json:"players"`
	Observers *Connections      `json:"observers,omitempty"`
	History   []*PlayerAction   `json:"history,omitempty"`
	Messages  []*models.Message `json:"messages"`

	Field    *Field               `json:"field,omitempty"`
	Date     time.Time            `json:"date,omitempty"`
	Settings *models.RoomSettings `json:"settings"`
}

func (room *Room) JSON() RoomJSON {
	return RoomJSON{
		ID:        room.ID,
		Name:      room.Name,
		Status:    room.Status,
		Players:   room.Players,
		Observers: room.Observers,
		History:   room.history(),
		Messages:  room._Messages,
		Field:     room.Field,
		Date:      room.Date,
		Settings:  room.Settings,
	}
}

func (room *Room) MarshalJSON() ([]byte, error) {
	return json.Marshal(room.JSON())
}

func (room *Room) UnmarshalJSON(b []byte) error {
	temp := &RoomJSON{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	room.Name = temp.Name
	room.Status = temp.Status
	room.Players = temp.Players
	room.Observers = temp.Observers
	room._History = temp.History
	room._Messages = temp.Messages
	room.Date = temp.Date
	room.Settings = temp.Settings

	return nil
}
