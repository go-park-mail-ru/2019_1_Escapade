package engine

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/constants"
)

// ConnectionAction is a bundle of Connection and action
type ConnectionAction struct {
	conn   *Connection
	action int
}

// Room consist of players and observers, field and history
type Room struct {
	info             RoomInformationI
	field            FieldProxyI
	sync             synced.SyncI
	api              RoomRequestsI
	lobby            LobbyProxyI
	models           RModelsI
	sender           RSendI
	client           RClientI
	people           PeopleI
	events           EventsI
	record           ActionRecorderI
	metrics          *RoomMetrics
	messages         MessagesI
	garbageCollector GarbageCollectorI
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
	}
	return true
}

type RoomArgs struct {
	c        *config.Room
	lobby    *Lobby
	rs       *models.RoomSettings
	id       string
	DBchatID int32
	DBRoomID int32
	Field    *Field
}

// NewRoom return new instance of room
func NewRoom(c *config.Room, lobby *Lobby,
	game *models.Game, id string) (*Room, error) {
	if err := constants.Check(game.Settings); err != nil {
		return nil, err
	}

	var room = &Room{}
	dbchatID, dbRoomID, err := room.registerInDB(lobby, game, id)
	if err != nil {
		return room, err
	}
	var ra = &RoomArgs{
		c:        c,
		lobby:    lobby,
		rs:       game.Settings,
		id:       id,
		DBchatID: dbchatID,
		DBRoomID: dbRoomID,
	}
	room.configureAndStart(ra)

	return room, err
}
func (room *Room) GetSync() synced.SyncI {
	return room.sync
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
func (room *Room) configureAndStart(ra *RoomArgs) {
	var (
		components = &RoomBuilder{}
	)
	ra.Field = NewField(ra.rs, &ra.c.Field)
	components.Build(room, ra)
	go room.events.Run(room)
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
