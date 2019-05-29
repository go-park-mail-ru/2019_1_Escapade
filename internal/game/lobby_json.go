package game

import (
	"encoding/json"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// LobbyJSON is a wrapper for sending Lobby by JSON
type LobbyJSON struct {
	AllRooms  Rooms             `json:"allRooms"`
	FreeRooms Rooms             `json:"freeRooms"`
	Waiting   Connections       `json:"waiting"`
	Playing   Connections       `json:"playing"`
	Messages  []*models.Message `json:"messages"`
}

// JSON convert Lobby to LobbyJSON
func (lobby *Lobby) JSON() LobbyJSON {
	return LobbyJSON{
		AllRooms:  *lobby._allRooms,
		FreeRooms: *lobby._freeRooms,
		Waiting:   *lobby.Waiting,
		Playing:   *lobby.Playing,
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
	lobby.Waiting = &temp.Waiting
	lobby.Playing = &temp.Playing
	lobby._messages = temp.Messages

	return nil
}
