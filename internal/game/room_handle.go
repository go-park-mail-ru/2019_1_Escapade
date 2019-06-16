package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"fmt"
	"time"
)

// Enter handle user joining as player or observer
func (room *Room) Enter(conn *Connection) bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	var done bool

	// if room is searching new players
	if room.Status() == StatusPeopleFinding {
		conn.debug("You will be player!")
		if room.addConnection(conn, true, false) {
			done = true
		}
	} else if room.addConnection(conn, false, false) {
		conn.debug("You will be observer!")
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

	fmt.Println("Can close?", room.lobby.canCloseRooms)
	if !room.lobby.canCloseRooms {
		return false
	}
	fmt.Println("We closed room :С")
	room.lobby.CloseRoom(room)
	//room.LeaveAll()
	fmt.Println("Prepare to free!")
	go room.Free()
	fmt.Println("We did it")
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
		go room.LeavePlayer(player)
	}

	observersIterator := NewConnectionsIterator(room.Observers)
	for observersIterator.Next() {
		observer := observersIterator.Value()
		go room.LeaveObserver(observer)
	}
	/*
		players := room.Players.Connections.RGet()
		for _, conn := range players {
			go room.LeavePlayer(conn)
		}
		observers := room.Observers.RGet()
		for _, conn := range observers {
			go room.LeaveObserver(conn)
		}*/
}

// Empty check room has no people
func (room *Room) Empty() bool {
	if room.done() {
		return true
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	return room.Players.Connections.len()+room.Observers.len() == 0
}

// LeavePlayer handle player going back to lobby
func (room *Room) LeavePlayer(conn *Connection) bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	done := room.Players.Connections.Remove(conn)
	if done {
		if room.Status() == StatusPeopleFinding {
			room.lobby.greet(conn)
			room.lobby.sendRoomUpdate(room, All)
		} else {
			room.GiveUp(conn)
		}
	}
	//fmt.Println("LeavePlayer", room.Empty(), len(room.Observers.RGet()))
	if room.Empty() {
		fmt.Println("room.Close()")
		room.Close()
	}
	return done
}

// LeaveObserver handle observer leaves room
func (room *Room) LeaveObserver(conn *Connection) bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	done := room.Observers.Remove(conn)
	if done {
		if room.Status() == StatusPeopleFinding {
			room.lobby.greet(conn)
			room.lobby.sendRoomUpdate(room, All)
		}
	}
	if room.Empty() {
		room.Close()
	}
	return done
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

	//conn.setDisconnected()
	if room.lobby.Metrics() {
		metrics.Players.WithLabelValues(room.ID(), conn.User.Name).Dec()
	}

	//done = room.RemoveFromGame(conn, action == ActionDisconnect)

	pa := *room.addAction(conn.ID(), action)
	if room.Empty() {
		if room.lobby.Metrics() {
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

	fmt.Println("cellHandle")
	status := room.Status()
	if status == StatusFlagPlacing {
		room.SetFlag(conn, cell)
	} else if status == StatusRunning {
		room.OpenCell(conn, cell)
	}
	return
}

// IsActive check if game is started and results not known
func (room *Room) IsActive() bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()
	return room.Status() == StatusFlagPlacing || room.Status() == StatusRunning
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

	fmt.Println("room handle conn", room.Status, conn.Index(), conn.Disconnected())
	found, _ := room.Search(conn)
	if found != nil {
		fmt.Println("room handle found", room.isAlive(found), found.Index(), found.Disconnected())
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
			fmt.Println("rr.Send.Action != nil ")
			room.chanConnection <- ConnectionAction{
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

	obsrversIterator := NewConnectionsIterator(room.Observers)
	for playersIterator.Next() {
		observer := obsrversIterator.Value()
		room.greet(observer, true)
		room.lobby.waiterToPlayer(observer, room)
	}
	/*
		players := room.Players.Connections.RGet()
		for _, conn := range players {
			room.greet(conn, true)
			room.lobby.waiterToPlayer(conn, room)
		}
		observers := room.Observers.RGet()
		for _, conn := range observers {
			room.greet(conn, false)
			room.lobby.waiterToPlayer(conn, room)
		}
	*/
	room.Players.Init(room.Field)

	room.lobby.RoomStart(room)
	go room.runGame()

	loc, _ := time.LoadLocation("Europe/Moscow")
	room.setDate(time.Now().In(loc))
	room.sendStatus(room.All)
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
	fmt.Println("StartGame")

	room.FillField()

	open := float64(room.Settings.Mines) / float64(room.Settings.Width*room.Settings.Height) * float64(100)
	fmt.Println("opennn", open, room.Settings.Width*room.Settings.Height)

	cells := room.Field.OpenSave(int(open))
	room.sendNewCells(room.All, cells...)
	room.setStatus(StatusRunning)
	loc, _ := time.LoadLocation("Europe/Moscow")
	room.setDate(time.Now().In(loc))
	room.sendStatus(room.All)
	room.sendMessage("Battle began! Destroy your enemy!", room.All)
}

// FinishGame finish game
func (room *Room) FinishGame(timer bool) {
	if room.done() {
		fmt.Println("room.done()!")
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	if room.Status() == StatusFinished {
		return
	}
	if !timer {
		room.chanFinish <- struct{}{}
	}

	room.setStatus(StatusFinished)

	room.sendStatus(room.All)
	room.sendGameOver(timer, room.All)
	room.Save()
	room.Players.Finish()

	playersIterator := NewConnectionsIterator(room.Players.Connections)
	for playersIterator.Next() {
		player := playersIterator.Value()
		room.Observers.Add(player /*, false*/)
	}
	/*
		playersConns := room.Players.Connections.RGet()
		for _, conn := range playersConns {
			room.Observers.Add(conn)
		}*/
	room.Players.RefreshConnections()
	room.lobby.roomFinish(room)
}
