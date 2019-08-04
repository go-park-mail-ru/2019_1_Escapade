package game

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/constants"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// Game status
const (
	StatusRecruitment = 0
	StatusAborted     = 1 // in case of error
	StatusFlagPlacing = 2
	StatusRunning     = 3
	StatusFinished    = 4
	StatusHistory     = 5
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

	dbRoomID int32
	dbChatID int32

	//playersM *sync.RWMutex
	Players *OnlinePlayers

	//observersM *sync.RWMutex
	Observers *Connections

	historyM *sync.RWMutex
	_history []*PlayerAction

	messagesM *sync.Mutex
	_messages []*models.Message

	killedM *sync.RWMutex
	_killed int32 //amount of killed users

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

	recruitmentTimeM *sync.RWMutex
	_recruitmentTime time.Duration

	playingTimeM *sync.RWMutex
	_playingTime time.Duration

	play    *time.Timer
	prepare *time.Timer

	chanStatus     chan int
	chanConnection chan *ConnectionAction

	Settings *models.RoomSettings
}

// CharacteristicsCheck check room's characteristics are valid
func CharacteristicsCheck(rs *models.RoomSettings) bool {
	if constants.ROOM.Set {
		namelen := int32(len(rs.Name))
		if namelen < constants.ROOM.NameMin ||
			namelen > constants.ROOM.NameMax {
			utils.Debug(false, "Name length is invalid:",
				namelen, ". Need more then", constants.ROOM.NameMin,
				" and less then", constants.ROOM.NameMax)
			return false
		}
		if rs.Width < constants.FIELD.WidthMin ||
			rs.Width > constants.FIELD.WidthMax {
			utils.Debug(false, "Width is invalid:", rs.Width)
			return false
		}
		if rs.Height < constants.FIELD.HeightMin ||
			rs.Height > constants.FIELD.HeightMax {
			utils.Debug(false, "Height is invalid:", rs.Height)
			return false
		}
		if rs.Players < constants.ROOM.PlayersMin ||
			rs.Players > constants.ROOM.PlayersMax {
			utils.Debug(false, "Amount of players is invalid:", rs.Players)
			return false
		}
		if rs.Observers > constants.ROOM.ObserversMax {
			utils.Debug(false, "Amount of observers is invalid:", rs.Observers)
			return false
		}
		if rs.TimeToPrepare < constants.ROOM.TimeToPrepareMin ||
			rs.TimeToPrepare > constants.ROOM.TimeToPrepareMax {
			utils.Debug(false, "Time to prepare is invalid:", rs.TimeToPrepare)
			return false
		}
		if rs.TimeToPlay < constants.ROOM.TimeToPlayMin ||
			rs.TimeToPlay > constants.ROOM.TimeToPlayMax {
			utils.Debug(false, "Time to play is invalid:", rs.TimeToPlay)
			return false
		}
	} else {
		panic(3)
	}
	return true
}

// NewRoom return new instance of room
func NewRoom(config *config.FieldConfig, lobby *Lobby,
	game *models.Game, id string) (*Room, error) {
	if !CharacteristicsCheck(game.Settings) || !game.Settings.FieldCheck() {
		return nil, re.ErrorInvalidRoomSettings()
	}
	var (
		room           = &Room{}
		roomID, chatID int32
		err            error
	)
	game.Settings.ID = id
	if game.ID == 0 {
		game.Date = time.Now()
		// we create chat here, not when all people will be find, because
		// with this chat people can message while battle is finding players
		roomID, chatID, err = lobby.db().CreateGame(game)
		if err != nil {
			return nil, err
		}
	} else {
		roomID = game.ID
		chatID = game.ChatID
	}

	room.dbRoomID = roomID
	room.dbChatID = chatID
	room.Init(config, lobby, game.Settings, id)
	return room, err
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

	room.recruitmentTimeM = &sync.RWMutex{}

	room.playingTimeM = &sync.RWMutex{}

	room.dateM = &sync.RWMutex{}
	room._date = time.Now()

	room.Observers = NewConnections(room.Settings.Observers)

	room.chanStatus = make(chan int)
	room.chanConnection = make(chan *ConnectionAction)

	room.setHistory(make([]*PlayerAction, 0))
	room.setMessages(make([]*models.Message, 0))
	room.setKilled(0)
	room.setID(utils.RandomString(16))
	room.setStatus(StatusRecruitment)

	room.Field = field

	room.setDate(time.Now().In(room.lobby.location()))

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
//easyjson:json
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
//easyjson:json
type RoomSend struct {
	Cell     *Cell            `json:"cell,omitempty"`
	Action   *int             `json:"action,omitempty"`
	Messages *models.Messages `json:"messages,omitempty"`
}

// RoomGet is struct of flags, that client can get from room
//easyjson:json
type RoomGet struct {
	Players   bool `json:"players"`
	Observers bool `json:"observers"`
	Field     bool `json:"field"`
	History   bool `json:"history"`
}
