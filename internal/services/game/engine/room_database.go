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
}

// Save save room information to database
func (room *RoomModelsConverter) Save(wg *sync.WaitGroup) error {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	if room.r.done() {
		return re.ErrorRoomDone()
	}
	room.r.wGroup.Add(1)
	defer func() {
		room.r.wGroup.Done()
	}()

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

	if err := room.r.lobby.db().Save(gameInformation); err != nil {
		utils.Debug(false, "err. Cant save.", err.Error())
		room.r.lobby.AddNotSavedGame(&gameInformation)
	}

	return nil
}

func (room *RoomModelsConverter) toModelGame() models.Game {
	return models.Game{
		ID:              room.r.dbRoomID,
		Settings:        room.r.Settings,
		RecruitmentTime: room.r.recruitmentTime(),
		PlayingTime:     room.r.playingTime(),
		ChatID:          room.r.dbChatID,
		Status:          int32(room.r.Status()),
		Date:            room.r.Date(),
	}
}

func (room *RoomModelsConverter) toModelGamers() []models.Gamer {
	winners := room.r.Winners()
	players := room.r.Players.m.RPlayers()
	gamers := make([]models.Gamer, 0)
	for id, player := range players {
		gamer := models.Gamer{
			ID:        player.ID,
			Score:     player.Points,
			Explosion: player.Died,
			Won:       room.r.Winner(winners, id),
		}
		gamers = append(gamers, gamer)
	}
	return gamers
}

func (room *RoomModelsConverter) toModelField() models.Field {
	return models.Field{
		Width:     room.r.Field.Width,
		Height:    room.r.Field.Height,
		CellsLeft: room.r.Field._cellsLeft,
		Difficult: 0,
		Mines:     room.r.Field.Mines,
	}
}

func (room *RoomModelsConverter) toModelCells() []models.Cell {
	cells := make([]models.Cell, 0)
	for _, cellHistory := range room.r.Field.History() {
		cell := models.Cell{
			PlayerID: cellHistory.PlayerID,
			X:        cellHistory.X,
			Y:        cellHistory.Y,
			Value:    cellHistory.Value,
			Date:     cellHistory.Time,
		}
		cells = append(cells, cell)
	}
	return cells
}

func (room *RoomModelsConverter) toModelActions() []models.Action {
	history := room.r.history()
	actions := make([]models.Action, 0)
	for _, actionHistory := range history {
		action := models.Action{
			PlayerID: actionHistory.Player,
			ActionID: actionHistory.Action,
			Date:     actionHistory.Time,
		}
		actions = append(actions, action)
	}
	return actions
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
			Players: room.r.Players.m.RPlayers(),
			Cells:   cells,
			Winners: room.r.Winners(),
			Timer:   timer,
		},
	}
}

func (room *RoomModelsConverter) responseRoomStatus(
	status int) models.Response {
	var leftTime int32
	since := int32(time.Since(room.r.Date()).Seconds())
	if status == StatusFlagPlacing {
		leftTime = room.r.Settings.TimeToPrepare - since
	} else if status == StatusRunning {
		leftTime = room.r.Settings.TimeToPlay - since
	}
	return models.Response{
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
	conn *Connection, isPlayer bool) models.Response {
	var flag Flag
	if room.r.Settings.Deathmatch {
		index := conn.Index()
		if index >= 0 {
			flag = room.r.Players.m.Flag(index)
		}
	} else {
		flag = Flag{Cell: *NewCell(-1, -1, 0, 0)}
	}

	//leftTime := room.Settings.TimeToPlay + room.Settings.TimeToPrepare - int(time.Since(room.Date).Seconds())

	return models.Response{
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
	room.setStatus(int(info.Game.Status))
	room.setKilled(info.Game.Settings.Players)
	room.setDate(info.Game.Date)

	// actions
	for _, actionDB := range info.Actions {
		action := &PlayerAction{
			Player: actionDB.PlayerID,
			Action: actionDB.ActionID,
			Time:   actionDB.Date,
		}
		room.appendAction(action)
	}

	// field
	room.Field.Width = info.Field.Width
	room.Field.Height = info.Field.Height
	room.Field.setCellsLeft(info.Field.CellsLeft)
	room.Field.Mines = info.Field.Mines

	// cells
	room.Field.setHistory(make([]*Cell, 0))
	for _, cellDB := range info.Cells {
		cell := &Cell{
			X:        cellDB.X,
			Y:        cellDB.Y,
			Value:    cellDB.Value,
			PlayerID: cellDB.PlayerID,
			Time:     cellDB.Date,
		}
		room.Field.setToHistory(cell)
	}

	// players
	room.Players = newOnlinePlayers(info.Game.Settings.Players, *room.Field)
	for i, gamer := range info.Gamers {
		room.Players.m.SetPlayer(i, Player{
			ID:       gamer.ID,
			Points:   gamer.Score,
			Died:     gamer.Explosion,
			Finished: true,
		})
	}

	_, room._messages, err = GetChatIDAndMessages(lobby.location(),
		pChat.ChatType_ROOM, room.dbChatID, room.lobby.SetImage)

	//room._messages, err = room.lobby.db.LoadMessages(true, info.Game.RoomID)

	room.setStatus(StatusHistory)

	return
}
