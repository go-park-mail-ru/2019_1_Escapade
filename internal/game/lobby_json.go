package game

import (
	"encoding/json"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// Room consist of players and observers, field and history
type LobbyJSON struct {
	AllRooms  Rooms             `json:"allRooms"`
	FreeRooms Rooms             `json:"freeRooms"`
	Waiting   Connections       `json:"waiting"`
	Playing   Connections       `json:"playing"`
	Messages  []*models.Message `json:"messages"`
}

func (lobby *Lobby) JSON() LobbyJSON {
	return LobbyJSON{
		AllRooms:  *lobby._AllRooms,
		FreeRooms: *lobby._FreeRooms,
		Waiting:   *lobby._Waiting,
		Playing:   *lobby._Playing,
		Messages:  lobby._Messages,
	}
}

func (lobby *Lobby) MarshalJSON() ([]byte, error) {
	return json.Marshal(lobby.JSON())
}

func (lobby *Lobby) UnmarshalJSON(b []byte) error {
	temp := &LobbyJSON{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	lobby._AllRooms = &temp.AllRooms
	lobby._FreeRooms = &temp.FreeRooms
	lobby._Waiting = &temp.Waiting
	lobby._Playing = &temp.Playing
	lobby._Messages = temp.Messages

	return nil
}
