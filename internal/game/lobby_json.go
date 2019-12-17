package game

import (
	"encoding/json"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// LobbyJSON is a wrapper for sending Lobby by JSON
type LobbyJSON struct {
	AllRooms  RoomsJSON         `json:"allRooms"`
	FreeRooms RoomsJSON         `json:"freeRooms"`
	Waiting   ConnectionsJSON   `json:"waiting"`
	Playing   ConnectionsJSON   `json:"playing"`
	Messages  []*models.Message `json:"messages"`
}

// JSON convert Lobby to LobbyJSON
func (lobby *Lobby) JSON() LobbyJSON {
	return LobbyJSON{
		AllRooms:  lobby.allRooms.JSON(),
		FreeRooms: lobby.freeRooms.JSON(),
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
	lobby.allRooms.Set(temp.AllRooms.Get)
	lobby.allRooms.SetCapacity(temp.AllRooms.Capacity)

	lobby.freeRooms.Set(temp.FreeRooms.Get)
	lobby.freeRooms.SetCapacity(temp.FreeRooms.Capacity)

	lobby.Waiting.Set(temp.Waiting.Get)
	lobby.Waiting.SetCapacity(temp.Waiting.Capacity)

	lobby.Playing.Set(temp.Playing.Get)
	lobby.Playing.SetCapacity(temp.Playing.Capacity)

	lobby._messages = temp.Messages

	return nil
}
