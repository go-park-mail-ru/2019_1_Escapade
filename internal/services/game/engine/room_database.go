package engine

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	pChat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

type RoomModelsConverter struct {
	r *Room
	s SyncI
}

func (room *RoomModelsConverter) Init(r *Room, s SyncI) {
	room.r = r
	room.s = s
}

// Save save room information to database
func (room *RoomModelsConverter) Save(wg *sync.WaitGroup) error {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	var err error
	room.s.do(func() {
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

		if err = room.r.lobby.db().Save(gameInformation); err != nil {
			utils.Debug(false, "err. Cant save.", err.Error())
			room.r.lobby.AddNotSavedGame(&gameInformation)
		}
	})

	return err
}

func (room *RoomModelsConverter) toModelGame() models.Game {
	return models.Game{
		ID:              room.r.dbRoomID,
		Settings:        room.r.Settings,
		RecruitmentTime: room.r.events.recruitmentTime(),
		PlayingTime:     room.r.events.playingTime(),
		ChatID:          room.r.messages.dbChatID,
		Status:          int32(room.r.events.Status()),
		Date:            room.r.events.Date(),
	}
}

func (room *RoomModelsConverter) toModelGamer(index int, player Player) models.Gamer {
	return models.Gamer{
		ID:        player.ID,
		Score:     player.Points,
		Explosion: player.Died,
		Won:       room.r.people.IsWinner(index),
	}
}

func (room *RoomModelsConverter) toModelGamers() []models.Gamer {
	gamers := make([]models.Gamer, 0)
	room.r.people.Players.ForEach(room.getGamers(gamers))
	return gamers
}

func (room *RoomModelsConverter) getGamers(gamers []models.Gamer) func(int, Player) {
	return func(index int, player Player) {
		gamers = append(gamers, room.toModelGamer(index, player))
	}
}

func (room *RoomModelsConverter) toModelField() models.Field {
	return room.r.field.Model()
}

func (room *RoomModelsConverter) toModelCells() []models.Cell {
	return room.r.field.ModelCells()
}

func (room *RoomModelsConverter) toModelActions() []models.Action {
	return room.r.record.ModelActions()
}

////////// sender models //////////

func (room *RoomModelsConverter) responseRoomGameOver(timer bool,
	cells []Cell) *models.Response {
	return &models.Response{
		Type: "RoomGameOver",
		Value: struct {
			Players []Player `json:"players"`
			Cells   []Cell   `json:"cells"`
			Winners []int    `json:"winners"`
			Timer   bool     `json:"timer"`
		}{
			Players: room.r.people.PlayersSlice(),
			Cells:   cells,
			Winners: room.r.people.Winners(),
			Timer:   timer,
		},
	}
}

func (room *RoomModelsConverter) responseRoomStatus(
	status int) *models.Response {
	var leftTime int32
	since := int32(time.Since(room.r.events.Date()).Seconds())
	if status == StatusFlagPlacing {
		leftTime = room.r.Settings.TimeToPrepare - since
	} else if status == StatusRunning {
		leftTime = room.r.Settings.TimeToPlay - since
	}
	return &models.Response{
		Type: "RoomStatus",
		Value: struct {
			ID     string `json:"id"`
			Status int    `json:"status"`
			Time   int32  `json:"time"`
		}{
			ID:     room.r.ID(),
			Status: status,
			Time:   leftTime,
		},
	}
}

func (room *RoomModelsConverter) responseRoom(
	conn *Connection, isPlayer bool) *models.Response {
	var flag Flag
	if room.r.Settings.Deathmatch {
		index := conn.Index()
		if index >= 0 {
			flag = room.r.people.Flag(index)
		}
	} else {
		flag = Flag{Cell: *NewCell(-1, -1, 0, 0)}
	}

	//leftTime := room.Settings.TimeToPlay + room.Settings.TimeToPrepare - int(time.Since(room.Date).Seconds())

	return &models.Response{
		Type: "Room",
		Value: struct {
			Room *Room                 `json:"room"`
			You  models.UserPublicInfo `json:"you"`
			Flag Flag                  `json:"flag,omitempty"`
			//Time     int                   `json:"time"`
			IsPlayer bool `json:"isPlayer"`
		}{
			Room: room.r,
			You:  *conn.User,
			Flag: flag,
			//Time:     leftTime,
			IsPlayer: isPlayer,
		},
	}
}

///////////////////////////////////

// Load load room information from database
func (lobby *Lobby) Load(id string) (room *Room, err error) {
	if lobby.done() {
		return nil, re.ErrorLobbyDone()
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
	}()

	var info models.GameInformation
	if info, err = lobby.db().FetchOneGame(id); err != nil {
		return
	}

	if room, err = NewRoom(lobby.config().Field, lobby, &info.Game, id); err != nil {
		return
	}

	// main info
	room.events.setStatus(int(info.Game.Status))
	room.people.setKilled(info.Game.Settings.Players)
	room.events.setDate(info.Game.Date)

	// actions
	for _, actionDB := range info.Actions {
		action := &PlayerAction{
			Player: actionDB.PlayerID,
			Action: actionDB.ActionID,
			Time:   actionDB.Date,
		}
		room.record.appendAction(action)
	}

	// field
	room.field.Field.Width = info.Field.Width
	room.field.Field.Height = info.Field.Height
	room.field.Field.setCellsLeft(info.Field.CellsLeft)
	room.field.Field.Mines = info.Field.Mines

	// cells
	room.field.Field.setHistory(make([]*Cell, 0))
	for _, cellDB := range info.Cells {
		cell := &Cell{
			X:        cellDB.X,
			Y:        cellDB.Y,
			Value:    cellDB.Value,
			PlayerID: cellDB.PlayerID,
			Time:     cellDB.Date,
		}
		room.field.Field.setToHistory(cell)
	}

	// players
	room.people.Players = newOnlinePlayers(info.Game.Settings.Players)
	for i, gamer := range info.Gamers {
		room.people.Players.m.SetPlayer(i, Player{
			ID:       gamer.ID,
			Points:   gamer.Score,
			Died:     gamer.Explosion,
			Finished: true,
		})
	}

	_, messages, err := GetChatIDAndMessages(lobby.location(),
		pChat.ChatType_ROOM, room.messages.dbChatID, room.lobby.SetImage)

	if err == nil {
		room.messages.setMessages(messages)
	}

	//room._messages, err = room.lobby.db.LoadMessages(true, info.Game.RoomID)

	room.events.setStatus(StatusHistory)

	return
}
