package game

import (
	"fmt"
	"sync"

	pChat "github.com/go-park-mail-ru/2019_1_Escapade/chat/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

// Save save room information to database
func (room *Room) Save(wg *sync.WaitGroup) (err error) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	if room.done() {
		return re.ErrorRoomDone()
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	players := room.Players.RPlayers()

	// made in NewRoom
	//room.Settings.ID = room.ID()

	game := models.Game{
		ID:              room.dbRoomID,
		Settings:        room.Settings,
		RecruitmentTime: room.recruitmentTime(),
		PlayingTime:     room.playingTime(),
		ChatID:          room.dbChatID,
		Status:          int32(room.Status()),
		Date:            room.Date(),
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
		CellsLeft: room.Field._cellsLeft,
		Difficult: 0,
		Mines:     room.Field.Mines,
	}

	cells := make([]models.Cell, 0)
	for _, cellHistory := range room.Field.History() {
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

	if err = room.lobby.db().SaveGame(gameInformation); err != nil {
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
	if info, err = lobby.db().GetGame(id); err != nil {
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
		room.Players.SetPlayer(i, Player{
			ID:       gamer.ID,
			Points:   gamer.Score,
			Died:     gamer.Explosion,
			Finished: true,
		})
	}

	_, room._messages, err = GetChatIDAndMessages(lobby.location(), pChat.ChatType_ROOM, room.dbChatID)

	//room._messages, err = room.lobby.db.LoadMessages(true, info.Game.RoomID)

	room.setStatus(StatusHistory)

	return
}
