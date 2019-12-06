package constants

import (
	"strings"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
)

// RoomConfiguration - the limits of the characteristics of the room
//easyjson:json
type roomConfiguration struct {
	Set              bool
	NameMin          int32 `json:"nameMin"`
	NameMax          int32 `json:"nameMax"`
	TimeToPrepareMin int32 `json:"timeToPrepareMin"`
	TimeToPrepareMax int32 `json:"timeToPrepareMax"`
	TimeToPlayMin    int32 `json:"timeToPlayMin"`
	TimeToPlayMax    int32 `json:"timeToPlayMax"`
	PlayersMin       int32 `json:"playersMin"`
	PlayersMax       int32 `json:"playersMax"`
	ObserversMax     int32 `json:"observersMax"`
}

// ROOM - singleton of room constants
var ROOM = roomConfiguration{}

// InitRoom initializes ROOM
func InitRoom(rep RepositoryI, path string) error {
	room, err := rep.getRoom(path)
	if err != nil {
		return err
	}
	ROOM = room
	ROOM.Set = true

	return nil
}

// Check room's characteristics are valid
func Check(rs *models.RoomSettings) error {
	if ROOM.Set && FIELD.Set {
		rs.Name = strings.Trim(rs.Name, " ")
		namelen := int32(len(rs.Name))
		if namelen < ROOM.NameMin || namelen > ROOM.NameMax {
			return ErrorRoomName(rs)
		}
		if rs.Width < FIELD.WidthMin || rs.Width > FIELD.WidthMax {
			return ErrorFieldWidth(rs)
		}
		if rs.Height < FIELD.HeightMin || rs.Height > FIELD.HeightMax {
			return ErrorFieldHeight(rs)
		}
		s := rs.Width * rs.Height
		if rs.Players < ROOM.PlayersMin || rs.Players > ROOM.PlayersMax ||
			rs.Players > s {
			return ErrorPlayers(rs)
		}
		if rs.Observers < 0 || rs.Observers > ROOM.ObserversMax {
			return ErrorObservers(rs)
		}
		p := rs.TimeToPrepare
		if p < ROOM.TimeToPrepareMin || p > ROOM.TimeToPrepareMax {
			return ErrorTimeToPrepare(rs)
		}
		p = rs.TimeToPlay
		if p < ROOM.TimeToPlayMin || p > ROOM.TimeToPlayMax {
			return ErrorTimeToPlay(rs)
		}
		max := s - rs.Players
		if rs.Mines < 0 || rs.Mines > max {
			return ErrorMines(rs.Mines, max)
		}
	} else {
		return ErrorConstantsNotSet()
	}
	return nil
}
