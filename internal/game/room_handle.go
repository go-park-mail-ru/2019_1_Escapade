package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

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
	if room.Status == StatusPeopleFinding {
		conn.debug("You will be player!")
		if room.addPlayer(conn) {
			done = true
		}
	} else if room.addObserver(conn) {
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
	room.setDone()

	room.wGroup.Wait()

	room.Status = StatusFinished
	go room.historyFree()
	go room.messagesFree()
	go room.Players.Free()
	go room.Observers.Free()
	go room.Field.Free()

	close(room.chanFinish)
	close(room.chanStatus)
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
	room.LeaveAll()
	room.lobby.CloseRoom(room)
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

	players := room.Players.Connections.RGet()
	for _, conn := range players {
		go room.Leave(conn, ActionDisconnect)
	}
	observers := room.Observers.RGet()
	for _, conn := range observers {
		go room.Leave(conn, ActionDisconnect)
	}
}

// Leave handle user going back to lobby
func (room *Room) Leave(conn *Connection, action int) (done bool) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	//conn.setDisconnected()
	if room.lobby.Metrics() {
		metrics.Players.WithLabelValues(room.ID, conn.User.Name).Dec()
	}
	pa := *room.addAction(conn.ID(), action)
	room.sendAction(pa, room.AllExceptThat(conn))
	fmt.Println("Left room")

	done = room.RemoveFromGame(conn, action == ActionDisconnect)

	return
}

// applyAction applies the effects of opening a cell
func (room *Room) applyAction(conn *Connection, cell *Cell) {
	index := conn.Index()

	fmt.Println("points:", room.Settings.Width, room.Settings.Height, 100*float64(cell.Value+1)/float64(room.Settings.Width*room.Settings.Height))
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
	if room.Status != StatusRunning {
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
	fmt.Println("len cell", len(cells))
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
	if room.Status == StatusFlagPlacing {
		room.SetFlag(conn, cell)
	} else if room.Status == StatusRunning {
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
	return room.Status == StatusFlagPlacing || room.Status == StatusRunning
}

// ActionHandle processes the Action came from the user
func (room *Room) ActionHandle(conn *Connection, action int) (done bool) {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()
	fmt.Println("action", action)

	switch action {
	case ActionGiveUp:
		if room.IsActive() {
			go room.GiveUp(conn)
			return true
		}
	case ActionRestart:
		if room.Status == StatusRunning || room.Status == StatusFlagPlacing {
			return false
		}
		conn.lobby.greet(conn)
		if room.Status == StatusFinished {
			room.Restart()
			room.lobby.addRoom(room)
		}
		if room.Status == StatusPeopleFinding {
			room.addPlayer(conn)
		}
		return true
	case ActionBackToLobby:
		room.lobby.LeaveRoom(conn, room, ActionBackToLobby)
		return true
	}
	return false
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

	fmt.Println("room handle conn", room.isAlive(conn), conn.Index(), conn.Disconnected())
	found := room.Search(conn)
	if found != nil {
		fmt.Println("room handle found", room.isAlive(found), found.Index(), found.Disconnected())
	}
	if rr.IsGet() {
		//go room.greet(conn)
	} else if rr.IsSend() {
		//done := false
		switch {
		case rr.Send.Messages != nil:
			Messages(conn, rr.Send.Messages, room.Messages())
		case rr.Send.Cell != nil:
			if room.isAlive(conn) {
				go room.CellHandle(conn, rr.Send.Cell)
			}
		case rr.Send.Action != nil:
			room.ActionHandle(conn, *rr.Send.Action)
		}
	} else if rr.Message != nil {
		if conn.Index() < 0 {
			rr.Message.Status = models.StatusObserver
		} else {
			rr.Message.Status = models.StatusPlayer
		}
		Message(room.lobby, conn, rr.Message, room.appendMessage,
			room.setMessage, room.removeMessage, room.findMessage,
			room.send, room.InGame, true, room.ID)
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

	room.Status = StatusFlagPlacing
	players := room.Players.Connections.RGet()
	for _, conn := range players {
		room.MakePlayer(conn, false)
	}
	observers := room.Observers.RGet()
	for _, conn := range observers {
		room.MakeObserver(conn, false)
	}
	room.Players.Init(room.Field)

	room.lobby.RoomStart(room)
	go room.runGame()

	room.Date = time.Now()
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
	room.Status = StatusRunning
	room.Date = time.Now()
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

	if room.Status == StatusFinished {
		fmt.Println("room.Status == StatusFinished!")
		return
	}
	// if !timer {
	// 	room.chanFinish <- struct{}{}
	// }
	fmt.Println(room.ID, "We finish room!", room.Status)

	room.Status = StatusFinished
	fmt.Println(room.ID, "We finish room?", room.Status)

	room.sendStatus(room.All)
	room.sendMessage("Battle finished!", room.All)
	room.sendGameOver(timer, room.All)
	room.Save()
	room.Players.Finish()

	playersConns := room.Players.Connections.RGet()
	for _, conn := range playersConns {
		room.Observers.Add(conn, false)
	}
	room.Players.RefreshConnections()
	room.lobby.roomFinish(room)
}

// initTimers launch game timers. Call it when flag placement starts
func (room *Room) initTimers() {
	room.prepare = time.NewTimer(time.Second *
		time.Duration(room.Settings.TimeToPrepare))
	room.play = time.NewTimer(time.Second *
		time.Duration(room.Settings.TimeToPlay))
	return
}

func (room *Room) launchGarbageCollector(timeoutPeopleFinding, timeoutPlayer, timeoutObserver, timeoutFinished float64) {
	//fmt.Println("launchGarbageCollector")

	if room.Status == StatusPeopleFinding {
		timeoutPlayer = timeoutPeopleFinding
		timeoutObserver = timeoutPeopleFinding
	}
	if room.Status == StatusFinished {
		timeoutPlayer = timeoutFinished
		timeoutObserver = timeoutFinished
	}
	i := 0
	for _, conn := range room.Players.Connections.RGet() {
		if conn == nil {
			continue
		}
		i++
		if conn.Disconnected() && time.Since(conn.time).Seconds() > timeoutPlayer {
			//fmt.Println(conn.User.Name, " - bad")
			room.Leave(conn, ActionTimeOver)
		} else {
			//fmt.Println(conn.User.Name, " - good", conn.Disconnected(), time.Since(conn.time).Seconds())
		}
	}
	for _, conn := range room.Observers.RGet() {
		if conn == nil {
			continue
		}
		i++
		if conn.Disconnected() && time.Since(conn.time).Seconds() > timeoutObserver {
			//fmt.Println(conn.User.Name, " - bad")
			room.Leave(conn, ActionTimeOver)
		} else {
			//fmt.Println(conn.User.Name, " - good", conn)
		}
	}
	if i == 0 {
		room.chanStatus <- StatusAborted
	}
}

/*
func (room *Room) runHistory() {
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

	room.StartFlagPlacing()

	ticker := time.NewTicker(time.Millisecond * 10)
	defer func() {
		room.Status = StatusHistory
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			if (actionsSize + cellsSize) == (actionsI + cellsI) {
				return
			}
			for actionsI < actionsSize{

			}

		}
	}
}
*/

func (room *Room) runRoom() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	// все в конфиг
	ticker := time.NewTicker(time.Second * 1)
	var timeoutPeopleFinding, timeoutPlayer, timeoutObserver, timeoutFinished float64
	timeoutPeopleFinding = 2
	timeoutPlayer = 60
	timeoutObserver = 5
	timeoutFinished = 20

	for {
		select {
		case <-ticker.C:
			go room.launchGarbageCollector(timeoutPeopleFinding, timeoutPlayer, timeoutObserver, timeoutFinished)
		case newStatus := <-room.chanStatus:
			if newStatus == room.Status || newStatus > StatusFinished {
				continue
			}
			switch newStatus {
			// case StatusPeopleFinding:
			// 	//
			case StatusFlagPlacing:
				room.StartFlagPlacing()
			case StatusRunning:
				room.prepare.Stop()
				room.StartGame()
			case StatusFinished:
				ok := room.play.Stop()
				ticker.Stop()
				room.FinishGame(!ok)
				return
			case StatusAborted:
				ticker.Stop()
				return
			}
		}
	}
}

// run - room goroutine
func (room *Room) runGame() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	ticker := time.NewTicker(time.Second * 4)

	room.initTimers()
	defer func() {
		ticker.Stop()
		room.prepare.Stop()
		room.play.Stop()
		fmt.Println("Room: Game is over!")
	}()

	room.Date = time.Now()

	for {
		select {
		case <-room.chanFinish:
			return
		case <-room.prepare.C:
			room.chanStatus <- StatusRunning
		case <-room.play.C:
			room.chanStatus <- StatusFinished
			//fmt.Println("finish!")
			//room.FinishGame(true)
			return
		case clock := <-ticker.C:
			//fmt.Println("clock!", room.ID)
			go room.sendMessage(clock.String()+" passed", room.All)
		}
	}
}
