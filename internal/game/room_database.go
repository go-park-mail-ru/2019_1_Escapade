package game

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

// Save save room information to database
func (room *Room) Save() (err error) {
	if room.done() {
		return re.ErrorRoomDone()
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	players := room.Players.RPlayers()
	game := models.Game{
		RoomID:        room.ID,
		Name:          room.Name,
		Status:        room.Status,
		Players:       len(players),
		TimeToPrepare: room.Settings.TimeToPrepare,
		TimeToPlay:    room.Settings.TimeToPlay,
		Date:          room.Date,
	}

	winners := room.Winners()
	gamers := make([]models.Gamer, 0)
	for id, player := range players {
		gamer := models.Gamer{
			ID:        player.ID,
			Score:     player.Points,
			Explosion: player.Died,
			Won:       room.Winner(winners, id),
		}
		gamers = append(gamers, gamer)
	}

	field := models.Field{
		Width:     room.Field.Width,
		Height:    room.Field.Height,
		CellsLeft: room.Field.CellsLeft,
		Difficult: 0,
		Mines:     room.Field.Mines,
	}

	cells := make([]models.Cell, 0)
	for _, cellHistory := range room.Field.History {
		cell := models.Cell{
			PlayerID: cellHistory.PlayerID,
			X:        cellHistory.X,
			Y:        cellHistory.Y,
			Value:    cellHistory.Value,
			Date:     cellHistory.Time,
		}
		cells = append(cells, cell)
	}

	history := room.history()
	actions := make([]models.Action, 0)
	for _, actionHistory := range history {
		action := models.Action{
			PlayerID: actionHistory.Player,
			ActionID: actionHistory.Action,
			Date:     actionHistory.Time,
		}
		actions = append(actions, action)
	}

	gameInformation := models.GameInformation{
		Game:    game,
		Gamers:  gamers,
		Field:   field,
		Actions: actions,
		Cells:   cells,
	}

	if err = room.lobby.db.SaveGame(gameInformation); err != nil {
		fmt.Println("err. Cant save.", err.Error())
	}

	return
}

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
	if info, err = lobby.db.GetGame(id); err != nil {
		return
	}

	// settings
	settings := &models.RoomSettings{
		ID:            info.Game.RoomID,
		Name:          info.Game.Name,
		Width:         info.Field.Width,
		Height:        info.Field.Height,
		Players:       info.Game.Players,
		Observers:     1,
		TimeToPrepare: info.Game.TimeToPrepare,
		TimeToPlay:    info.Game.TimeToPlay,
	}

	if room, err = NewRoom(settings, id, lobby); err != nil {
		return
	}

	// main info
	room.ID = info.Game.RoomID
	room.Name = info.Game.Name
	room.Status = info.Game.Status
	room.setKilled(info.Game.Players)
	room.Date = info.Game.Date

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
	room.Field.CellsLeft = info.Field.CellsLeft
	room.Field.Mines = info.Field.Mines

	// cells
	room.Field.History = make([]Cell, 0)
	for _, cellDB := range info.Cells {
		cell := Cell{
			X:        cellDB.X,
			Y:        cellDB.Y,
			Value:    cellDB.Value,
			PlayerID: cellDB.PlayerID,
			Time:     cellDB.Date,
		}
		room.Field.History = append(room.Field.History, cell)
	}

	// players
	room.Players = newOnlinePlayers(info.Game.Players, *room.Field)
	for i, gamer := range info.Gamers {
		room.Players.SetPlayer(i, Player{
			ID:       gamer.ID,
			Points:   gamer.Score,
			Died:     gamer.Explosion,
			Finished: true,
		})
	}

	room._messages, err = room.lobby.db.LoadMessages(true, info.Game.RoomID)

	room.Status = StatusHistory

	return
}
