package game

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"fmt"
	"time"
)

// Enter handle user joining as player or observer
func (room *Room) Enter(conn *Connection) bool {
	if room.done() {
		utils.Debug(true, "room is done")
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	var done bool

	// if room is searching new players
	if room.Status() == StatusPeopleFinding {
		if room.addConnection(conn, true, false) {
			utils.Debug(false, "You will be player!")
			done = true
		}
	} else if room.addConnection(conn, false, false) {
		utils.Debug(false, "You will be observer!")
		done = true
	}
	return done
}

// Free clear all resources. Call it when no
//  observers and players inside
func (room *Room) Free() {

	if room.done() {
		return
	}

	//room.setDone()
	room.wGroup.Wait()
	if room.done() {
		return
	}
	room.setDone()

	fmt.Println("room.setDone()")

	room.chanStatus <- StatusAborted

	room.setStatus(StatusFinished)
	go room.historyFree()
	go room.messagesFree()
	go room.Players.Free()
	go room.Observers.Free()
	go room.Field.Free()

	close(room.chanFinish)
	close(room.chanStatus)
	close(room.chanConnection)
}

// Close drives away players out of the room, free resources
// and inform lobby, that rooms closes
func (room *Room) Close() bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	utils.Debug(false, "Can close?", room.lobby.config.CanClose)
	if !room.lobby.config.CanClose {
		return false
	}
	utils.Debug(false, "We closed room :С")
	room.lobby.CloseRoom(room)
	utils.Debug(false, "Prepare to free!")
	go room.Free()
	utils.Debug(false, "We did it")
	return true
}

// LeaveAll make every room connection to leave
func (room *Room) LeaveAll() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	playersIterator := NewConnectionsIterator(room.Players.Connections)
	for playersIterator.Next() {
		player := playersIterator.Value()
		go room.Leave(player, true)
	}

	observersIterator := NewConnectionsIterator(room.Observers)
	for observersIterator.Next() {
		observer := observersIterator.Value()
		go room.Leave(observer, false)
	}
}

// LeavePlayer handle player going back to lobby
func (room *Room) Leave(conn *Connection, isPlayer bool) bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	var done bool
	if isPlayer {
		//done = room.Players.Connections.Remove(conn)
		utils.Debug(false, "isPlayer delete")
	} else {
		done = room.Observers.Remove(conn)
		utils.Debug(false, "notPlayer delete")
	}
	utils.Debug(false, "done", done)
	if !done {
		utils.Debug(false, "not found", conn.ID())
		return false
	}

	if room.Status() == StatusPeopleFinding {
		room.lobby.greet(conn)
		room.lobby.sendRoomUpdate(room, All)
	} else if isPlayer {
		room.GiveUp(conn)
	}

	utils.Debug(false, "letsCheckIfNil", room.Players.Connections.len(), room.Observers.len())
	if room.Empty() {
		room.Close()
	}
	return true
}

// LeaveMeta update metainformation about user leaving room
func (room *Room) LeaveMeta(conn *Connection, action int) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	if room.lobby.config.Metrics {
		metrics.Players.WithLabelValues(room.ID(), conn.User.Name).Dec()
	}

	pa := *room.addAction(conn.ID(), action)
	if room.Empty() {
		if room.lobby.config.Metrics {
			metrics.Rooms.Dec()
		}
	} else {
		room.sendAction(pa, room.AllExceptThat(conn))
	}

	return
}

// applyAction applies the effects of opening a cell
func (room *Room) applyAction(conn *Connection, cell *Cell) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	index := conn.Index()

	switch {
	case cell.Value < CellMine:
		room.Players.IncreasePlayerPoints(index, 1000*float64(cell.Value)/float64(room.Settings.Width*room.Settings.Height))
	case cell.Value == CellMine:
		room.Players.IncreasePlayerPoints(index, float64(-1000))
		room.Kill(conn, ActionExplode)
	case cell.Value > CellIncrement:
		room.FlagFound(*conn, cell)
	}
}

// OpenCell open cell
func (room *Room) OpenCell(conn *Connection, cell *Cell) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	// if user try set open cell before game launch
	if room.Status() != StatusRunning {
		return
	}

	// if wrong cell
	if !room.Field.IsInside(cell) {
		return
	}

	// if user died
	if !room.isAlive(conn) {
		return
	}
	//index := conn.Index()

	// set who try open cell(for history)
	cell.PlayerID = conn.ID()
	cells := room.Field.OpenCell(cell)
	if len(cells) == 1 {
		newCell := cells[0]
		room.applyAction(conn, &newCell)
	} else {
		for _, foundCell := range cells {
			room.applyAction(conn, &foundCell)
		}
	}

	if len(cells) > 0 {
		room.sendPlayerPoints(room.Players.Player(conn.Index()), room.All)
		go room.sendNewCells(room.All, cells...)
	}
	if room.Field.IsCleared() {
		room.chanStatus <- StatusFinished
		//go room.FinishGame(false)
	}
	return
}

// CellHandle processes the Cell came from the user
func (room *Room) CellHandle(conn *Connection, cell *Cell) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	utils.Debug(false, "cellHandle")
	status := room.Status()
	if status == StatusFlagPlacing {
		room.SetFlag(conn, cell)
	} else if status == StatusRunning {
		room.OpenCell(conn, cell)
	}
	return
}

// HandleRequest processes the equest came from the user
func (room *Room) HandleRequest(conn *Connection, rr *RoomRequest) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	if room == nil {
		return
	}

	utils.Debug(false, "room handle conn", room.Status, conn.Index(), conn.Disconnected())
	found, _ := room.Search(conn)
	if found != nil {
		utils.Debug(false, "room handle found", room.isAlive(found), found.Index(), found.Disconnected())
	}
	if rr.IsGet() {
		//go room.greet(conn)
	} else if rr.IsSend() {
		switch {
		case rr.Send.Messages != nil:
			Messages(conn, rr.Send.Messages, room.Messages())
		case rr.Send.Cell != nil:
			if room.isAlive(conn) {
				go room.CellHandle(conn, rr.Send.Cell)
			}
		case rr.Send.Action != nil:
			utils.Debug(false, "rr.Send.Action != nil ")
			room.chanConnection <- &ConnectionAction{
				conn:   conn,
				action: *rr.Send.Action,
			}
		}
	} else if rr.Message != nil {
		if conn.Index() < 0 {
			rr.Message.Status = models.StatusObserver
		} else {
			rr.Message.Status = models.StatusPlayer
		}
		Message(room.lobby, conn, rr.Message, room.appendMessage,
			room.setMessage, room.removeMessage, room.findMessage,
			room.send, room.All, true, room.ID())
	}
}

// StartFlagPlacing prepare field, players and observers
func (room *Room) StartFlagPlacing() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	room.FillField()

	room.setStatus(StatusFlagPlacing)

	playersIterator := NewConnectionsIterator(room.Players.Connections)
	for playersIterator.Next() {
		player := playersIterator.Value()
		room.greet(player, true)
		room.lobby.waiterToPlayer(player, room)
	}

	observersIterator := NewConnectionsIterator(room.Observers)
	for playersIterator.Next() {
		observer := observersIterator.Value()
		room.greet(observer, true)
		room.lobby.waiterToPlayer(observer, room)
	}
	room.Players.Init(room.Field)

	room.lobby.RoomStart(room)
	go room.runGame()

	loc, _ := time.LoadLocation(room.lobby.config.Location)
	room.setDate(time.Now().In(loc))
	go room.sendStatus(room.All, nil)
	go room.sendField(room.All)
	go room.sendMessage("Battle will be start soon! Set your flag!", room.All)
}

// StartGame start game
func (room *Room) StartGame() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()
	utils.Debug(false, "StartGame")

	room.FillField()

	open := float64(room.Settings.Mines) / float64(room.Settings.Width*room.Settings.Height) * float64(100)
	utils.Debug(false, "opennn", open, room.Settings.Width*room.Settings.Height)

	cells := room.Field.OpenZero() //room.Field.OpenSave(int(open))
	go room.sendNewCells(room.All, cells...)
	room.setStatus(StatusRunning)
	loc, _ := time.LoadLocation(room.lobby.config.Location)
	room.setDate(time.Now().In(loc))
	go room.sendStatus(room.All, nil)
	go room.sendMessage("Battle began! Destroy your enemy!", room.All)
}

// FinishGame finish game
func (room *Room) FinishGame(timer bool) {
	if room == nil {
		utils.Debug(true, "room nil")
	}
	if room.done() {
		utils.Debug(true, "room.done()!")
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	if room.Status() == StatusFinished {
		return
	}

	if !timer && room.Status() != StatusPeopleFinding {
		room.chanFinish <- struct{}{}
	}
	room.setStatus(StatusFinished)

	// save Group
	saveAndSendGroup := &sync.WaitGroup{}

	cells := make([]Cell, 0)
	room.Field.OpenEverything(&cells)

	saveAndSendGroup.Add(4)
	go room.sendStatus(room.All, saveAndSendGroup)
	go room.sendGameOver(timer, room.All, cells, saveAndSendGroup)
	go room.Save(saveAndSendGroup)
	go room.Players.Finish(saveAndSendGroup)

	saveAndSendGroup.Wait()
	/*
		playersIterator := NewConnectionsIterator(room.Players.Connections)
		for playersIterator.Next() {
			player := playersIterator.Value()
			player.SetIndex(-1)
			room.Observers.Add(player)
		}
		room.Players = newOnlinePlayers(room.Settings.Players, *room.Field)
	*/
	room.lobby.roomFinish(room)
}
