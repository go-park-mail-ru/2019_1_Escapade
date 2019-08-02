package constants

import (
	"io/ioutil"
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
func InitRoom(path string) error {
	var (
		data []byte
		err  error
	)

	if data, err = ioutil.ReadFile(path); err != nil {
		return err
	}

	var tmp roomConfiguration
	if err = tmp.UnmarshalJSON(data); err != nil {
		return err
	}
	tmp.Set = true
	ROOM = tmp

	return nil
}
