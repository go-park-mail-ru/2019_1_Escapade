package game

import (
	"encoding/json"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// LobbyJSON is a wrapper for sending Lobby by JSON
type LobbyJSON struct {
	AllRooms  Rooms             `json:"allRooms"`
	FreeRooms Rooms             `json:"freeRooms"`
	Waiting   ConnectionsJSON   `json:"waiting"`
	Playing   ConnectionsJSON   `json:"playing"`
	Messages  []*models.Message `json:"messages"`
}

// JSON convert Lobby to LobbyJSON
func (lobby *Lobby) JSON() LobbyJSON {
	return LobbyJSON{
		AllRooms:  *lobby._allRooms,
		FreeRooms: *lobby._freeRooms,
		Waiting:   lobby.Waiting.JSON(),
		Playing:   lobby.Playing.JSON(),
		Messages:  lobby._messages,
	}
}

// MarshalJSON - overriding the standard method json.Marshal
func (lobby Lobby) MarshalJSON() ([]byte, error) {
	return json.Marshal(lobby.JSON())
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (lobby *Lobby) UnmarshalJSON(b []byte) error {
	temp := &LobbyJSON{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	lobby._allRooms = &temp.AllRooms
	lobby._freeRooms = &temp.FreeRooms
	lobby.Waiting._get = temp.Waiting.Get
	lobby.Waiting._capacity = temp.Waiting.Capacity
	lobby.Playing._get = temp.Playing.Get
	lobby.Playing._capacity = temp.Playing.Capacity
	lobby._messages = temp.Messages

	return nil
}
