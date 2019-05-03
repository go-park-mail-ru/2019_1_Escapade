package game

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// Save save room information to database
func (room *Room) Save() (err error) {
	game := models.Game{
		RoomID:        room.ID,
		Name:          room.Name,
		Status:        room.Status,
		Players:       len(room.Players.Players),
		TimeToPrepare: room.Settings.TimeToPrepare,
		TimeToPlay:    room.Settings.TimeToPlay,
		Date:          room.Date,
	}

	idWin := room.Winner()
	gamers := make([]models.Gamer, 0)
	for id, player := range room.Players.Players {
		gamer := models.Gamer{
			ID:        player.ID,
			Score:     player.Points,
			Explosion: player.Died,
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

// Load load room information from database
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
	room.killed = info.Game.Players
	room.Date = info.Game.Date

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

	room.Field.History = make([]Cell, 0)
	for _, cellHistory := range info.Cells {
		cell := Cell{
			PlayerID: cellHistory.PlayerID,
			X:        cellHistory.X,
			Y:        cellHistory.Y,
			Value:    cellHistory.Value,
			Time:     cellHistory.Date,
		}
		room.Field.History = append(room.Field.History, cell)
	}

	// actions
	room.History = make([]*PlayerAction, 0)
	for _, actionHistory := range info.Actions {
		action := &PlayerAction{
			Player: actionHistory.PlayerID,
			Action: actionHistory.ActionID,
			Time:   actionHistory.Date,
		}
		room.History = append(room.History, action)
	}

	// players
	room.Players = newOnlinePlayers(info.Game.Players, *room.Field)
	for i, gamer := range info.Gamers {
		room.Players.Players[i] = Player{
			ID:       gamer.ID,
			Points:   gamer.Score,
			Died:     gamer.Explosion,
			Finished: true,
		}
	}

	room.Messages, err = room.lobby.db.LoadMessages(true, info.Game.RoomID)

	return
}
