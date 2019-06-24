package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// Game status
const (
	StatusPeopleFinding = 0
	StatusAborted       = 1 // in case of error
	StatusFlagPlacing   = 2
	StatusRunning       = 3
	StatusFinished      = 4
	StatusHistory       = 5
)

// ConnectionAction is a bundle of Connection and action
type ConnectionAction struct {
	conn   *Connection
	action int
}

// Room consist of players and observers, field and history
type Room struct {
	wGroup *sync.WaitGroup

	doneM *sync.RWMutex
	_done bool

	//playersM *sync.RWMutex
	Players *OnlinePlayers

	//observersM *sync.RWMutex
	Observers *Connections

	historyM *sync.RWMutex
	_history []*PlayerAction

	messagesM *sync.Mutex
	_messages []*models.Message

	killedM *sync.RWMutex
	_killed int //amount of killed users

	idM *sync.RWMutex
	_id string

	nameM *sync.RWMutex
	_name string

	statusM *sync.RWMutex
	_status int

	nextM *sync.RWMutex
	_next *Room

	lobby *Lobby
	Field *Field

	dateM *sync.RWMutex
	_date time.Time

	chanFinish chan struct{}

	chanStatus     chan int
	chanConnection chan *ConnectionAction

	play    *time.Timer
	prepare *time.Timer

	Settings *models.RoomSettings
}

// NewRoom return new instance of room
func NewRoom(config *config.FieldConfig, lobby *Lobby, rs *models.RoomSettings, id string) (*Room, error) {
	if !rs.AreCorrect() {
		return nil, re.ErrorInvalidRoomSettings()
	}
	var room = &Room{}

	room.Init(config, lobby, rs, id)
	return room, nil
}

// Init init instance of room
func (room *Room) Init(config *config.FieldConfig, lobby *Lobby,
	rs *models.RoomSettings, id string) {

	room.wGroup = &sync.WaitGroup{}

	field := NewField(rs, config)
	room._done = false
	room.doneM = &sync.RWMutex{}

	// cant use Restart cause need to
	room.Players = newOnlinePlayers(rs.Players, *field)
	room.historyM = &sync.RWMutex{}
	room.messagesM = &sync.Mutex{}
	room.killedM = &sync.RWMutex{}

	room.nameM = &sync.RWMutex{}
	room._name = rs.Name

	room.statusM = &sync.RWMutex{}

	room.lobby = lobby

	room.Settings = rs

	room.idM = &sync.RWMutex{}
	room._id = id

	room.nextM = &sync.RWMutex{}
	room._next = nil

	room.dateM = &sync.RWMutex{}

	room.Observers = NewConnections(room.Settings.Observers)

	room.chanFinish = make(chan struct{})
	room.chanStatus = make(chan int)
	room.chanConnection = make(chan *ConnectionAction)

	room.setHistory(make([]*PlayerAction, 0))
	room.setMessages(make([]*models.Message, 0))
	room.setKilled(0)
	room.setID(utils.RandomString(16))
	room.setStatus(StatusPeopleFinding)

	room.Field = field

	loc, _ := time.LoadLocation(room.lobby.config.Location)

	room.setDate(time.Now().In(loc))

	go room.runRoom()

	return
}

// Restart fill in the room fields with the original values

func (room *Room) Restart(conn *Connection) {

	if room.Next() == nil || room.Next().done() {
		pa := *room.addAction(conn.ID(), ActionRestart)
		room.sendAction(pa, room.All)
		next, err := room.lobby.CreateAndAddToRoom(room.Settings, conn)
		if err != nil {
			panic("next")
		}
		room.setNext(next)
	}
	room.processActionBackToLobby(conn)
	room.Next().Enter(conn)
	return
}

// debug print all room fields
func (room *Room) debug() {
	if room == nil {
		fmt.Println("cant debug nil room")
		return
	}
	fmt.Println("Room id    :", room._id)
	fmt.Println("Room name  :", room._name)
	fmt.Println("Room status:", room._status)
	fmt.Println("Room date  :", room.Date)
	fmt.Println("Room killed:", room.killed())
	players := room.Players.RPlayers()
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

// Empty check room has no people
func (room *Room) Empty() bool {
	if room.done() {
		return true
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	return room.Players.Connections.len()+room.Observers.len() == 0
}

// IsActive check if game is started and results not known
func (room *Room) IsActive() bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()
	return room.Status() == StatusFlagPlacing || room.Status() == StatusRunning
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
	Cell     *Cell            `json:"cell,omitempty"`
	Action   *int             `json:"action,omitempty"`
	Messages *models.Messages `json:"messages,omitempty"`
}

// RoomGet is struct of flags, that client can get from room
type RoomGet struct {
	Players   bool `json:"players"`
	Observers bool `json:"observers"`
	Field     bool `json:"field"`
	History   bool `json:"history"`
}
