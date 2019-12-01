package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
)

// LobbyJSON is a wrapper for sending Lobby by JSON
//easyjson:json
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
		Messages:  lobby.Messages(),
	}
}

// MarshalJSON - overriding the standard method json.Marshal
func (lobby Lobby) MarshalJSON() ([]byte, error) {
	return lobby.JSON().MarshalJSON()
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (lobby *Lobby) UnmarshalJSON(b []byte) error {
	temp := &LobbyJSON{}

	if err := temp.UnmarshalJSON(b); err != nil {
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
