package engine

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	action_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/action"
)

// RoomJSON is a wrapper for sending Room by JSON
//easyjson:json
type RoomJSON struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`

	Players   OnlinePlayersJSON       `json:"players"`
	Observers ConnectionsJSON         `json:"observers,omitempty"`
	History   []*action_.PlayerAction `json:"history,omitempty"`
	Messages  []*models.Message       `json:"messages"`

	Field    FieldJSON            `json:"field,omitempty"`
	Date     time.Time            `json:"date,omitempty"`
	Settings *models.RoomSettings `json:"settings"`
}

// MarshalJSON - overriding the standard method json.Marshal
func (room *Room) MarshalJSON() ([]byte, error) {
	return room.models.JSON().MarshalJSON()
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (room *Room) UnmarshalJSON(b []byte) error {
	temp := &RoomJSON{}

	if err := temp.UnmarshalJSON(b); err != nil {
		return err
	}
	room.info.setName(temp.Name)
	room.events.configure(temp.Status, temp.Date)
	room.record.setHistory(temp.History)
	room.messages.setMessages(temp.Messages)
	room.info.setSettings(temp.Settings)

	return nil
}
