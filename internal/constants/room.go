package constants

import (
	"io/ioutil"
)

// RoomConfiguration - the limits of the characteristics of the room
//easyjson:json
type roomConfiguration struct {
	Set              bool
	NameMin          int `json:"nameMinLength"`
	NameMax          int `json:"nameMaxLength"`
	TimeToPrepareMin int `json:"timeToPrepareMin"`
	TimeToPrepareMax int `json:"timeToPrepareMax"`
	TimeToPlayMin    int `json:"timeToPlayMin"`
	TimeToPlayMax    int `json:"timeToPlayMax"`
	PlayersMin       int `json:"playersMin"`
	PlayersMax       int `json:"playersMax"`
	ObserversMax     int `json:"observersMax"`
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
