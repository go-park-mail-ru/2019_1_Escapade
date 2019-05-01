package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"fmt"
	"time"
)

// Game status
const (
	StatusPeopleFinding = 0
	StatusAborted       = 1 // in case of error
	StatusFlagPlacing   = 2
	StatusRunning       = 3
	StatusFinished      = 4
)

// Room consist of players and observers, field and history
type Room struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`

	Players   *OnlinePlayers `json:"players,omitempty"`
	Observers *Connections   `json:"observers,omitempty"`

	History []*PlayerAction `json:"history,omitempty"`

	lobby *Lobby
	Field *Field `json:"field,omitempty"`

	Date       time.Time `json:"date,omitempty"`
	chanFinish chan struct{}

	// for save game room
	settings *models.RoomSettings

	killed int //amount of killed users
}

// NewRoom return new instance of room
func NewRoom(rs *models.RoomSettings, id string, lobby *Lobby) *Room {
	fmt.Println("NewRoom rs = ", *rs)
	room := &Room{
		ID:        id,
		Name:      rs.Name,
		Status:    StatusPeopleFinding,
		Players:   newOnlinePlayers(rs.Players),
		Observers: NewConnections(rs.Observers),

		History: make([]*PlayerAction, 0),

		lobby: lobby,
		Field: NewField(rs),

		Date:       time.Now(),
		chanFinish: make(chan struct{}),
		settings:   rs,
		killed:     0,
	}
	return room
}

// SameAs compare  one room with another
func (room *Room) SameAs(another *Room) bool {
	return room.Field.SameAs(another.Field)
}

/* Examples of json

room search
{"send":{"RoomSettings":{"name":"my best room","id":"create","width":12,"height":12,"players":3,"observers":10,"prepare":10, "play":100, "mines":5}},"get":null}

send cell
{"send":{"cell":{"x":2,"y":1,"value":0,"PlayerID":0}, "action":null},"get":null}

send action(all actions are in action.go). Server iswaiting only one of these:
ActionStop 5
ActionContinue 6
ActionGiveUp 13
ActionBackToLobby 14

give up
{"send":{"cell":null, "action":13,"get":null}}

back to lobby
{"send":{"cell":null, "action":14,"get":null}}

get lobby all info
{"send":null,"get":{"allRooms":true,"freeRooms":true,"waiting":true,"playing":true}}

	Players   bool `json:"players"`
	Observers bool `json:"observers"`
	Field     bool `json:"field"`
	History   bool `json:"history"`
{"send":null,"get":{"players":true,"observers":true,"field":true,"history":true}}

*/

// RoomRequest is request from client to room
type RoomRequest struct {
	Send *RoomSend `json:"send"`
	Get  *RoomGet  `json:"get"`
}

// IsGet check if client want get information
func (rr *RoomRequest) IsGet() bool {
	return rr.Get != nil
}

// IsSend check if client want send information
func (rr *RoomRequest) IsSend() bool {
	return rr.Send != nil
}

// RoomSend is struct of information, that client can send to room
type RoomSend struct {
	Cell   *Cell `json:"cell,omitempty"`
	Action *int  `json:"action,omitempty"`
}

// RoomGet is struct of flags, that client can get from room
type RoomGet struct {
	Players   bool `json:"players"`
	Observers bool `json:"observers"`
	Field     bool `json:"field"`
	History   bool `json:"history"`
}
