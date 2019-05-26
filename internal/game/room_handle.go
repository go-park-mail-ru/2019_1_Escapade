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
	go room.playersFree()
	go room.observersFree()
	go room.Field.Free()

	close(room.chanFinish)
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

	players := room.playersConnections()
	for _, conn := range players {
		go room.Leave(conn, ActionDisconnect)
	}
	observers := room.observers()
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

	if room.lobby.Metrics() {
		metrics.Players.WithLabelValues(room.ID, conn.User.Name).Dec()
	}
	pa := *room.addAction(conn.ID(), action)
	room.sendAction(pa, room.AllExceptThat(conn))
	fmt.Println("Left room")

	return room.RemoveFromGame(conn, action == ActionDisconnect)
}

// openCell open cell
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
	index := conn.Index()

	// set who try open cell(for history)
	cell.PlayerID = conn.ID()
	cells := room.Field.OpenCell(cell)
	fmt.Println("len cell", len(cells))
	if len(cells) == 1 {
		newCell := cells[0]
		fmt.Println("newCell value", newCell.Value)
		if newCell.Value < CellMine {
			room.IncreasePlayerPoints(index, 1+newCell.Value)
		} else if newCell.Value == CellMine {
			go room.IncreasePlayerPoints(index, -100) // в конфиг
			room.Kill(conn, ActionExplode)
		} else if newCell.Value >= CellIncrement {
			room.FlagFound(*conn, &newCell)
		} else if newCell.Value == CellOpened {
			return
		}
	} else {
		for _, foundCell := range cells {
			value := foundCell.Value
			if value < CellMine {
				go room.IncreasePlayerPoints(index, 1+value)
			}
		}
	}

	if len(cells) > 0 {
		go room.sendPlayerPoints(room.player(index), room.All)
		go room.sendNewCells(cells, room.All)
	}
	if room.Field.IsCleared() {
		go room.FinishGame(false)
	}
	return
}

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

func (room *Room) ActionHandle(conn *Connection, action int) (done bool) {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()
	fmt.Println("action", action)
	if room.IsActive() {
		if action == ActionGiveUp {
			conn.debug("we see you wanna give up?")
			go room.GiveUp(conn)
			return true
		}
	} else {
		if action == ActionRestart {
			room.addPlayer(conn)
			return true
		}
	}
	if action == ActionBackToLobby {
		conn.debug("we see you wanna back to lobby?")
		room.lobby.LeaveRoom(conn, room, ActionBackToLobby)
		conn.debug("we did it")
		return true
	}

	return false
}

// handleRequest
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

	conn.debug("room handle conn")
	if rr.IsGet() {
		go room.greet(conn)
	} else if rr.IsSend() {
		//done := false
		if rr.Send.Cell != nil {
			if room.isAlive(conn) {
				go room.CellHandle(conn, rr.Send.Cell)
			}
		} else if rr.Send.Action != nil {
			fmt.Println("action")
			room.ActionHandle(conn, *rr.Send.Action)
		}
		//if done {
		//room.finishGame(true)
		//}
	} else if rr.Message != nil {
		i := room.observersSearch(conn)
		if i > 0 {
			rr.Message.Status = models.StatusObserver
		} else {
			rr.Message.Status = models.StatusPlayer
		}
		Message(lobby, conn, rr.Message, room.setToMessages,
			room.send, room.InGame, true, room.ID)
	}
}

func (room *Room) StartFlagPlacing() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	//fmt.Println("StartFlagPlacing")

	room.Status = StatusFlagPlacing
	players := room.playersConnections()
	for _, conn := range players {
		room.MakePlayer(conn)
	}
	observers := room.observers()
	for _, conn := range observers {
		room.MakeObserver(conn)
	}
	room.playersInit()

	go room.lobby.RoomStart(room)
	go room.run()

	go room.sendStatus(room.All)
	go room.sendField(room.All)
	go room.sendMessage("Battle will be start soon! Set your flag!", room.All)
}

func (room *Room) StartGame() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	room.Status = StatusRunning
	go room.sendStatus(room.All)
	go room.sendMessage("Battle began! Destroy your enemy!", room.All)
	room.FillField()
}

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
	if !timer {
		room.chanFinish <- struct{}{}
	}
	fmt.Println(room.ID, "We finish room!", room.Status)

	room.Status = StatusFinished
	fmt.Println(room.ID, "We finish room?", room.Status)

	room.sendStatus(room.All)
	room.sendMessage("Battle finished!", room.All)
	room.sendGameOver(timer, room.All)
	players := room.players()
	for _, player := range players {
		player.Finished = true
	}

	playersConns := room.playersConnections()
	for i, conn := range playersConns {
		fmt.Println("found", i)
		if conn == nil {
			fmt.Println("fine nil")
			//continue
		} else {
			fmt.Println("id", conn.ID())
		}
		//room.playersRemove(conn)
		room.observersAdd(conn, false)
	}
	room.zeroPlayers()
	room.lobby.roomFinish(room)
	room.Save()
}

// initTimers launch game timers. Call it when flag placement starts
func (room *Room) initTimers() (prepare, play *time.Timer) {
	prepare = time.NewTimer(time.Second *
		time.Duration(room.Settings.TimeToPrepare))
	play = time.NewTimer(time.Second *
		time.Duration(room.Settings.TimeToPlay))
	return
}

func (room *Room) run() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("room_handle.go run()")
		room.wGroup.Done()
	}()

	ticker := time.NewTicker(time.Second * 4)

	timerToPrepare, timerToPlay := room.initTimers()
	defer func() {
		ticker.Stop()
		timerToPrepare.Stop()
		timerToPlay.Stop()
		fmt.Println("Room: Game is over!")
	}()

	room.Date = time.Now()

	for {
		select {
		case <-room.chanFinish:
			return
		case <-timerToPrepare.C:
			go room.StartGame()
		case <-timerToPlay.C:
			fmt.Println("finish!")
			room.FinishGame(true)
			return
		case clock := <-ticker.C:
			//fmt.Println("clock!", room.ID)
			go room.sendMessage(clock.String()+" passed", room.All)
		}
	}
}
