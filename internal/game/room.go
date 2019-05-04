package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
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
	wGroup *sync.WaitGroup

	doneM *sync.RWMutex
	_done bool

	playersM *sync.RWMutex
	_Players *OnlinePlayers `json:"players,omitempty"`

	observersM *sync.RWMutex
	_Observers *Connections `json:"observers,omitempty"`

	historyM *sync.RWMutex
	_History []*PlayerAction `json:"history,omitempty"`

	messagesM *sync.Mutex
	_Messages []*models.Message `json:"messages"`

	killedM *sync.RWMutex
	_killed int //amount of killed users

	ID     string `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`

	lobby *Lobby
	Field *Field `json:"field,omitempty"`

	Date       time.Time `json:"date,omitempty"`
	chanFinish chan struct{}

	Settings *models.RoomSettings `json:"settings"`
}

// NewRoom return new instance of room
func NewRoom(rs *models.RoomSettings, id string, lobby *Lobby) (*Room, error) {
	fmt.Println("NewRoom rs = ", *rs)
	if !rs.AreCorrect() {
		return nil, re.ErrorInvalidRoomSettings()
	}
	field := NewField(rs)
	room := &Room{
		wGroup: &sync.WaitGroup{},

		doneM: &sync.RWMutex{},
		_done: false,

		playersM: &sync.RWMutex{},
		_Players: newOnlinePlayers(rs.Players, *field),

		observersM: &sync.RWMutex{},
		_Observers: NewConnections(rs.Observers),

		historyM: &sync.RWMutex{},
		_History: make([]*PlayerAction, 0),

		messagesM: &sync.Mutex{},
		_Messages: make([]*models.Message, 0),

		killedM: &sync.RWMutex{},
		_killed: 0,

		ID:     id,
		Name:   rs.Name,
		Status: StatusPeopleFinding,

		lobby: lobby,
		Field: field,

		Date:       time.Now(),
		chanFinish: make(chan struct{}),

		Settings: rs,
	}
	return room, nil
}

func (room *Room) debug() {
	if room == nil {
		fmt.Println("cant debug nil room")
		return
	}
	fmt.Println("Room id    :", room.ID)
	fmt.Println("Room name  :", room.Name)
	fmt.Println("Room status:", room.Status)
	fmt.Println("Room date  :", room.Date)
	fmt.Println("Room killed:", room.killed())
	players := room.players()
	if len(players) == 0 {
		fmt.Println("cant debug nil players")
		return
	}
	for _, player := range players {
		fmt.Println("Player", player.ID)
		fmt.Println("Player points 	:", player.Points)
		fmt.Println("Player Finished:", player.Finished)
	}
	if room.Field == nil {
		fmt.Println("cant debug nil field")
		return
	}
	fmt.Println("Field width		:", room.Field.Width)
	fmt.Println("Field height 	:", room.Field.Height)
	fmt.Println("Field cellsleft:", room.Field.CellsLeft)
	fmt.Println("Field mines		:", room.Field.Mines)
	if room.Field.History == nil {
		fmt.Println("no field history")
	} else {
		for _, cell := range room.Field.History {
			fmt.Printf("Cell(%d,%d) with value %d", cell.X, cell.Y, cell.Value)
			fmt.Println("Cell Owner	:", cell.PlayerID)
			fmt.Println("Cell Time  :", cell.Time)
		}
	}
	history := room.history()
	if history == nil {
		fmt.Println("no action history")
	} else {
		for i, action := range history {
			fmt.Println("action", i)
			fmt.Println("action value  :", action.Action)
			fmt.Println("action Owner	:", action.Player)
			fmt.Println("action Time  :", action.Time)
		}
	}

}

// SameAs compare  one room with another
func (room *Room) SameAs(another *Room) bool {
	return room.Field.SameAs(another.Field)
}

/* Examples of json

message
{"message":{"text":"hello"}}


room search
{"send":{"RoomSettings":{"name":"my best room","id":"create","width":12,"height":12,"players":2,"observers":10,"prepare":10, "play":100, "mines":5}},"get":null}

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
	Send    *RoomSend       `json:"send"`
	Message *models.Message `json:"message"`
	Get     *RoomGet        `json:"get"`
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
