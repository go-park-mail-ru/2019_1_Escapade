package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"fmt"
	"time"
)

// Enter handle user joining as player or observer
func (room *Room) Enter(conn *Connection) bool {
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
	if room == nil || room.History == nil {
		return
	}
	room.wGroup.Wait()
	room.wGroup = nil
	fmt.Println("room free")
	room.Status = StatusFinished
	room.History = nil
	close(room.chanFinish)
	room.Players.Free()
	room.Observers.Free()
	for _, action := range room.History {
		action.Free()
	}
	room.Field.Clear()
	room = nil
}

// Close drives away players out of the room, free resources
// and inform lobby, that rooms closes
func (room *Room) Close() bool {
	fmt.Println("Can close?", room.lobby.canCloseRooms)
	if !room.lobby.canCloseRooms {
		return false
	}
	fmt.Println("We closed room :С")
	room.LeaveAll()
	room.lobby.CloseRoom(room)
	room.Free()
	return true
}

// LeaveAll make every room connection to leave
func (room *Room) LeaveAll() {
	for _, conn := range room.Players.Connections {
		room.Leave(conn, ActionDisconnect)
	}
	for _, conn := range room.Observers.Get {
		room.Leave(conn, ActionDisconnect)
	}
}

// Leave handle user going back to lobby
func (room *Room) Leave(conn *Connection, action int) {

	room.removeFromGame(conn, action == ActionDisconnect)

	go func() {
		pa := *room.addAction(conn.ID(), action)
		room.sendAction(pa, room.AllExceptThat(conn))
	}()
}

// openCell open cell
func (room *Room) openCell(conn *Connection, cell *Cell) (roomFinished bool) {
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

	// set who try open cell(for history)
	cell.PlayerID = conn.ID()
	cells := room.Field.OpenCell(cell)
	fmt.Println("len cell", len(cells))
	if len(cells) == 1 {
		newCell := cells[0]
		fmt.Println("newCell value", newCell.Value)
		if newCell.Value < CellMine {
			room.Players.Players[conn.Index].Points += 1 + newCell.Value
		} else if newCell.Value == CellMine {
			roomFinished = room.kill(conn, ActionExplode)
			room.Players.Players[conn.Index].Points -= 100
		} else if newCell.Value >= CellIncrement {
			room.flagFound(*conn, &newCell)
		} else if newCell.Value == CellOpened {
			return
		}
	} else {
		for _, foundCell := range cells {
			if foundCell.Value < CellMine {
				room.Players.Players[conn.Index].Points += 1 + foundCell.Value
			}
		}
	}

	go room.sendPlayerPoints(room.Players.Players[conn.Index], room.All)
	go room.sendNewCells(cells, room.All)

	if !roomFinished && room.Field.IsCleared() {
		roomFinished = true
		//room.finishGame(true)
	}
	return
}

func (room *Room) cellHandle(conn *Connection, cell *Cell) (done bool) {
	fmt.Println("cellHandle")
	if room.Status == StatusFlagPlacing {
		room.setFlag(conn, cell)
	} else if room.Status == StatusRunning {
		done = room.openCell(conn, cell)
	}
	return
}

// IsActive check if game is started and results not known
func (room *Room) IsActive() bool {
	return room.Status == StatusFlagPlacing || room.Status == StatusRunning
}

func (room *Room) actionHandle(conn *Connection, action int) (done bool) {
	if room.IsActive() {
		if action == ActionGiveUp {
			conn.debug("we see you wanna give up?")
			room.GiveUp(conn)
			return true
		}
	}
	if action == ActionBackToLobby {
		conn.debug("we see you wanna back to lobby?")
		room.lobby.LeaveRoom(conn, room, ActionBackToLobby)
		return true
	}

	return false
}

// handleRequest
func (room *Room) handleRequest(conn *Connection, rr *RoomRequest) {
	if (room == nil || room.Status == StatusFinished) {
		return
	}
	room.wGroup.Add(1)
	defer room.wGroup.Done()
	conn.debug("room handle conn")
	if rr.IsGet() {
		room.requestGet(conn, rr)
	} else if rr.IsSend() {
		//done := false
		if rr.Send.Cell != nil {
			if room.isAlive(conn) {
				room.cellHandle(conn, rr.Send.Cell)
			}
		} else if rr.Send.Action != nil {

			room.actionHandle(conn, *rr.Send.Action)
		}
		//if done {
			//room.finishGame(true)
		//}
	} else if rr.Message != nil {
		Message(lobby, conn, rr.Message, &room.Messages,
			room.send, room.InGame, true, room.ID)
	}
}

func (room *Room) startFlagPlacing() {
	room.Status = StatusFlagPlacing
	for _, conn := range room.Players.Connections {
		room.MakePlayer(conn)
	}
	for _, conn := range room.Observers.Get {
		room.MakeObserver(conn)
	}
	room.Players.Init(room.Field)

	go room.lobby.roomStart(room)
	go room.run()

	room.sendStatus(room.All)
	room.sendField(room.All)
	room.sendMessage("Battle will be start soon! Set your flag!", room.All)
}

func (room *Room) startGame() {
	room.Status = StatusRunning
	room.fillField()
	room.sendStatus(room.All)
	room.sendMessage("Battle began! Destroy your enemy!", room.All)
}

func (room *Room) finishGame(needStop bool) {
	if room.Status == StatusFinished {
		return
	}
	if needStop {
		room.chanFinish <- struct{}{}
	}
	fmt.Println(room.ID, "We finish room!", room.Status)

	room.Status = StatusFinished
	fmt.Println(room.ID, "We finish room?", room.Status)

	room.sendStatus(room.All)
	room.sendMessage("Battle finished!", room.All)
	room.sendGameOver(room.All)
	room.Save()
	go room.lobby.roomFinish(room)
	for _, player := range room.Players.Players {
		player.Finished = true
	}
}

// initTimers launch game timers. Call it when flag placement starts
func (room *Room) initTimers() (prepare, play *time.Timer) {
	prepare = time.NewTimer(time.Second *
		time.Duration(room.settings.TimeToPrepare))
	play = time.NewTimer(time.Second *
		time.Duration(room.settings.TimeToPlay))
	return
}

func (room *Room) run() {
	defer utils.CatchPanic("room_handle.go run()")
	ticker := time.NewTicker(time.Second * 4)

	timerToPrepare, timerToPlay := room.initTimers()
	defer func() {
		ticker.Stop()
		timerToPrepare.Stop()
		timerToPlay.Stop()
		fmt.Println("Room: Game is over!")
	}()

	for {
		select {
		case <-room.chanFinish:
			return
		case <-timerToPrepare.C:
			room.startGame()
		case <-timerToPlay.C:
			room.finishGame(false)
			return
		case clock := <-ticker.C:
			fmt.Println("clock!", room.ID)
			room.sendMessage(clock.String()+" passed", room.All)
		}
	}
}

func (room *Room) requestGet(conn *Connection, rr *RoomRequest) {
	room.greet(conn)
}
