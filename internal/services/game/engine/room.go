package engine

import (
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
	info             *RoomInformation
	field            *RoomField
	sync             *RoomSync
	api              *RoomAPI
	lobby            *RoomLobbyCommunication
	models           *RoomModelsConverter
	sender           *RoomSender
	connEvents       *RoomConnectionEvents
	people           *RoomPeople
	events           *RoomEvents
	record           *RoomRecorder
	metrics          *RoomMetrics
	messages         *RoomMessages
	garbageCollector RoomGarbageCollectorI
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
func NewRoom(config *config.Field, lobby *Lobby,
	game *models.Game, id string) (*Room, error) {
	if !CharacteristicsCheck(game.Settings) || !game.Settings.FieldCheck() {
		return nil, re.ErrorInvalidRoomSettings()
	}

	var room = &Room{}
	chatID, dbRoomID, err := room.registerInDB(lobby, game, id)
	if err != nil {
		return room, err
	}
	room.configureAndStart(config, lobby, game.Settings, id, chatID, dbRoomID)

	return room, err
}

func (room *Room) registerInDB(lobby *Lobby, game *models.Game, id string) (int32, int32, error) {
	var (
		dbRoomID, chatID int32
		err              error
	)
	game.Settings.ID = id
	if game.ID == 0 {
		game.Date = time.Now()
		// we create chat here, not when all people will be find, because
		// with this chat people can message while battle is finding players
		dbRoomID, chatID, err = lobby.db().Create(game)
		if err != nil {
			return dbRoomID, chatID, err
		}
	} else {
		dbRoomID = game.ID
		chatID = game.ChatID
	}
	return dbRoomID, chatID, err
}

// Init init instance of room
func (room *Room) configureAndStart(config *config.Field, lobby *Lobby,
	rs *models.RoomSettings, id string, chatID, roomID int32) {
	room.init()
	room.configureDependencies(config, lobby, rs, id, chatID, roomID)
	go room.events.Run()
}

func (room *Room) init() {
	room.sync = &RoomSync{}
	room.info = &RoomInformation{}
	room.api = &RoomAPI{}
	room.lobby = &RoomLobbyCommunication{}
	room.field = &RoomField{}
	room.models = &RoomModelsConverter{}
	room.sender = &RoomSender{}
	room.people = &RoomPeople{}
	room.connEvents = &RoomConnectionEvents{}
	room.events = &RoomEvents{}
	room.metrics = &RoomMetrics{}
	room.record = &RoomRecorder{}
	room.messages = &RoomMessages{}
	room.garbageCollector = &RoomGarbageCollector{}
}

func (room *Room) configureDependencies(config *config.Field, lobby *Lobby,
	rs *models.RoomSettings, id string, chatID, roomID int32) {
	field := NewField(rs, config)
	room.sync.Init(room)
	room.info.Init(rs, id, roomID)
	room.api.Init(room.sync, room.messages, room.connEvents, room.sender,
		room.events, room.info)
	room.lobby.Init(room, room.sync, room.info, lobby)
	room.field.Init(room.sync, room.record, room.sender, room.events,
		room.people, field, rs.Deathmatch)
	room.models.Init(room.sync, room.info, room.lobby, room.events,
		room.messages, room.people, room.record, room.field)
	room.sender.Init(room.sync, room.events, room.people, room.connEvents,
		room.info, room.models)
	room.people.Init(room.sync, room.connEvents, room.events, room.info,
		room.lobby, room.field, room.record, rs.Players, rs.Observers)
	room.connEvents.Init(room.sync, room.lobby, room.record, room.sender,
		room.info, room.events, room.people)
	room.events.Init(room.sync, room.info, room.lobby, room.people, room.field,
		room.garbageCollector, room.models, room.metrics, room.messages,
		room.record, room.sender)
	room.metrics.Init(room, room.sync, room.events, room.field, room.info)
	room.record.Init(room.sync, room.info, room.lobby, room.people,
		room.field, room.sender)
	room.messages.Init(room.sync, room.info, room.lobby, room.sender, chatID)

	// в конфиг
	t := Timeouts{
		timeoutPeopleFinding:   2.,
		timeoutRunningPlayer:   60.,
		timeoutRunningObserver: 5.,
		timeoutFinished:        20.,
	}
	room.garbageCollector.Init(room.sync, room.events, room.people,
		room.connEvents, t)
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

/*
func (room *Room) runHistory(conn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	//players := *room.Players
	actions := room.history()
	cells := room.Field.History
	actionsSize := len(actions)
	cellsSize := len(cells)
	actionsI := 0
	cellsI := 0
	actionTime := time.Now()
	cellTime := time.Now()
	if actionsSize > 0 {
		actionTime = actions[0].Time
	}
	if cellsSize > 0 {
		cellTime = cells[0].Time
	}

	// offline, err := room.lobby.createRoom(room.Settings)
	// if err != nil {
	// 	panic("offline doesnt work")
	// }
	// room.Leave(conn, ActionBackToLobby)
	// offline.Enter(conn)
	// offline.MakeObserver(conn, true)

	// offline.StartFlagPlacing()

	ticker := time.NewTicker(time.Millisecond * 10)
	defer func() {
		room.Status = StatusHistory
		ticker.Stop()
	}()

	for {
		select {
		case value := <-ticker.C:
			if (actionsSize + cellsSize) == (actionsI + cellsI) {
				return
			}
			for actionsI < actionsSize && value > time.Since() {

			}

		}
	}
}*/
