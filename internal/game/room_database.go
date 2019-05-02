package game

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

func (room *Room) Save() (err error) {
	game := models.Game{
		RoomID:        room.ID,
		Name:          room.Name,
		Status:        room.Status,
		Players:       len(room.Players.Players),
		TimeToPrepare: room.settings.TimeToPrepare,
		TimeToPlay:    room.settings.TimeToPlay,
		Date:          room.Date,
	}

	idWin := room.Winner()
	gamers := make([]models.Gamer, 0)
	for id, player := range room.Players.Players {
		gamer := models.Gamer{
			ID:        player.ID,
			Score:     player.Points,
			Explosion: player.Finished,
			Won:       idWin == id,
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

	actions := make([]models.Action, 0)
	for _, actionHistory := range room.History {
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

	fmt.Println("gameInformation:", gameInformation)

	if err = room.lobby.db.SaveGame(gameInformation); err != nil {
		fmt.Println("err. Cant save.", err.Error())
	}

	var room1 *Room
	if room1, err = lobby.Load(room.ID); err != nil {
		fmt.Println("err. Cant load.", err.Error())
	} else {
		room1.debug()
	}

	return
}

func (lobby *Lobby) Load(id string) (room *Room, err error) {

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
		TimeToPrepare: info.Game.TimeToPrepare,
		TimeToPlay:    info.Game.TimeToPlay,
	}

	room = NewRoom(settings, id, lobby)
	room.settings = settings

	// main info
	room.ID = info.Game.RoomID
	room.Name = info.Game.Name
	room.Status = info.Game.Status
	room.killed = info.Game.Players
	room.Date = info.Game.Date

	// players
	room.Players = newOnlinePlayers(info.Game.Players)
	for i, gamer := range info.Gamers {
		room.Players.Players[i] = Player{
			ID:       gamer.ID,
			Points:   gamer.Score,
			Finished: gamer.Explosion,
		}
	}

	// actions
	room.History = make([]*PlayerAction, 0)
	for _, actionDB := range info.Actions {
		action := &PlayerAction{
			Player: actionDB.PlayerID,
			Action: actionDB.ActionID,
			Time:   actionDB.Date,
		}
		room.History = append(room.History, action)
	}

	// field
	room.Field = &Field{
		Width:     info.Field.Width,
		Height:    info.Field.Height,
		CellsLeft: info.Field.CellsLeft,
		Mines:     info.Field.Mines,
	}

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

	actions := make([]models.Action, 0)
	for _, actionHistory := range room.History {
		action := models.Action{
			PlayerID: actionHistory.Player,
			ActionID: actionHistory.Action,
			Date:     actionHistory.Time,
		}
		actions = append(actions, action)
	}
	return
}
