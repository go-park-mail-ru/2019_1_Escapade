package engine

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"

	pChat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/database"
)

// RModelsI turns Game-type structures into models that can be sent
// to the client or to databases
// room model interface - adapter pattern
type RModelsI interface {
	Save(wg *sync.WaitGroup) error
	JSON() RoomJSON

	responseRoomGameOver(timer bool, cells []Cell) *models.Response
	responseRoomStatus(status int) *models.Response
	responseRoom(conn *Connection, isPlayer bool) *models.Response

	fromModelPlayerAction(action models.Action) *PlayerAction
	toModelPlayerAction(action *PlayerAction) models.Action
}

// RoomModels impelements RModelsI
type RoomModels struct {
	s  synced.SyncI
	i  RoomInformationI
	l  LobbyProxyI
	e  EventsI
	m  MessagesI
	p  PeopleI
	re ActionRecorderI
	f  FieldProxyI
}

func (room *RoomModels) Init(builder RBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildInformation(&room.i)
	builder.BuildLobby(&room.l)
	builder.BuildEvents(&room.e)
	builder.BuildMessages(&room.m)
	builder.BuildPeople(&room.p)
	builder.BuildRecorder(&room.re)
	builder.BuildField(&room.f)
}

// Save save room information to database
func (room *RoomModels) Save(wg *sync.WaitGroup) error {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	var err error
	room.s.Do(func() {
		// made in NewRoom
		//room.Settings.ID = room.ID()

		game := room.toModelGame()
		gamers := room.toModelGamers()
		field := room.toModelField()
		cells := room.toModelCells()
		actions := room.toModelActions()

		gameInformation := models.GameInformation{
			Game:    game,
			Gamers:  gamers,
			Field:   field,
			Actions: actions,
			Cells:   cells,
		}
		err = room.l.SaveGame(gameInformation)
	})

	return err
}

func (room *RoomModels) toModelGame() models.Game {
	return models.Game{
		ID:              room.i.RoomID(),
		Settings:        room.i.Settings(),
		RecruitmentTime: room.e.recruitmentTime(),
		PlayingTime:     room.e.playingTime(),
		ChatID:          room.m.ChatID(),
		Status:          int32(room.e.Status()),
		Date:            room.e.Date(),
	}
}

func (room *RoomModels) toModelGamer(index int, player Player) models.Gamer {
	return models.Gamer{
		ID:        player.ID,
		Score:     player.Points,
		Explosion: player.Died,
		Won:       room.p.IsWinner(index),
	}
}

func (room *RoomModels) fromModelPlayerAction(actionDB models.Action) *PlayerAction {
	return &PlayerAction{
		Player: actionDB.PlayerID,
		Action: actionDB.ActionID,
		Time:   actionDB.Date,
	}
}

func (room *RoomModels) toModelPlayerAction(action *PlayerAction) models.Action {
	return models.Action{
		PlayerID: action.Player,
		ActionID: action.Action,
		Date:     action.Time,
	}
}

func (room *RoomModels) toModelGamers() []models.Gamer {
	gamers := make([]models.Gamer, 0)
	room.p.players().ForEach(room.getGamers(gamers))
	return gamers
}

func (room *RoomModels) getGamers(gamers []models.Gamer) func(int, Player) {
	return func(index int, player Player) {
		gamers = append(gamers, room.toModelGamer(index, player))
	}
}

func (room *RoomModels) toModelField() models.Field {
	return room.f.Field().Model()
}

func (room *RoomModels) toModelCells() []models.Cell {
	return room.f.ModelCells()
}

func (room *RoomModels) toModelActions() []models.Action {
	return room.re.ModelActions()
}

// JSON convert Room to RoomJSON
func (room *RoomModels) JSON() RoomJSON {
	return RoomJSON{
		ID:        room.i.ID(),
		Name:      room.i.Name(),
		Status:    room.e.Status(),
		Players:   room.p.players().JSON(),
		Observers: room.p.observers().JSON(),
		History:   room.re.history(),
		Messages:  room.m.Messages(),
		Field:     room.f.Field().JSON(),
		Date:      room.e.Date(),
		Settings:  room.i.Settings(),
	}
}

////////// sender models //////////

func (room *RoomModels) responseRoomGameOver(timer bool,
	cells []Cell) *models.Response {
	return &models.Response{
		Type: "RoomGameOver",
		Value: struct {
			Players []Player `json:"players"`
			Cells   []Cell   `json:"cells"`
			Winners []int    `json:"winners"`
			Timer   bool     `json:"timer"`
		}{
			Players: room.p.PlayersSlice(),
			Cells:   cells,
			Winners: room.p.Winners(),
			Timer:   timer,
		},
	}
}

func (room *RoomModels) responseRoomStatus(
	status int) *models.Response {
	var leftTime int32
	since := int32(time.Since(room.e.Date()).Seconds())
	if status == StatusFlagPlacing {
		leftTime = room.i.Settings().TimeToPrepare - since
	} else if status == StatusRunning {
		leftTime = room.i.Settings().TimeToPlay - since
	}
	return &models.Response{
		Type: "RoomStatus",
		Value: struct {
			ID     string `json:"id"`
			Status int    `json:"status"`
			Time   int32  `json:"time"`
		}{
			ID:     room.i.ID(),
			Status: status,
			Time:   leftTime,
		},
	}
}

func (room *RoomModels) responseRoom(
	conn *Connection, isPlayer bool) *models.Response {
	var flag Flag
	if room.i.Settings().Deathmatch {
		index := conn.Index()
		if index >= 0 {
			flag = room.p.Flag(index)
		}
	} else {
		flag = Flag{Cell: *NewCell(-1, -1, 0, 0)}
	}

	//leftTime := room.Settings.TimeToPlay + room.Settings.TimeToPrepare - int(time.Since(room.Date).Seconds())

	return &models.Response{
		Type: "Room",
		Value: struct {
			Room RoomJSON              `json:"room"`
			You  models.UserPublicInfo `json:"you"`
			Flag Flag                  `json:"flag,omitempty"`
			//Time     int                   `json:"time"`
			IsPlayer bool `json:"isPlayer"`
		}{
			Room: room.JSON(),
			You:  *conn.User,
			Flag: flag,
			//Time:     leftTime,
			IsPlayer: isPlayer,
		},
	}
}

///////////////////////////////////

// Load load room information from database
func (lobby *Lobby) Load(id string) (*Room, error) {
	var (
		room *Room
		err  error
	)
	lobby.s.Do(func() {
		var info models.GameInformation
		if info, err = lobby.db().FetchOneGame(id); err != nil {
			return
		}

		room, err = NewRoom(lobby.rconfig(), lobby, &info.Game, id)
		if err != nil {
			return
		}

		room.events.configure(StatusHistory, info.Game.Date)
		room.record.configure(info.Actions)
		room.field.Configure(info)
		room.people.configure(info)

		_, messages, err := GetChatIDAndMessages(lobby.ChatService, lobby.location(),
			pChat.RoomType, room.messages.ChatID(), lobby.SetImage)

		if err == nil {
			room.messages.setMessages(messages)
		}
	})
	return room, err

	//room._messages, err = room.lobby.db.LoadMessages(true, info.Game.RoomID)
}
