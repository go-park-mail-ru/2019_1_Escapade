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

	if done {
		if room.Status != StatusPeopleFinding {
			room.lobby.waiterToPlayer(conn, room)
		}
	} else {
		conn.debug("No way!")
	}

	return done
}

// Free clear all resources. Call it when no
//  observers and players inside
func (room *Room) Free() {
	if room == nil || room.History == nil {
		return
	}
	room.Status = StatusFinished
	room.History = nil
	close(room.chanFinish)
	room.Players.Free()
	room.Observers.Free()
	for _, action := range room.History {
		action.Free()
	}
	room.Players.Free()
	room.Field.Clear()
	room = nil
}

// Close drives away players out of the room, free resources
// and inform lobby, that rooms closes
func (room *Room) Close() bool {
	room.LeaveAll()
	room.lobby.CloseRoom(room)
	room.Free()
	return false
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

	if !room.IsActive() || action == ActionDisconnect {
		room.removeBeforeLaunch(conn)
	} else {
		room.removeDuringGame(conn)
		conn.debug("Welcome back to lobby!")
	}
	go func() {
		room.addAction(conn.ID(), action)
		room.sendHistory(room.All)
	}()
}

// openCell open cell
func (room *Room) openCell(conn *Connection, cell *Cell) bool {
	// if user try set open cell before game launch
	if room.Status != StatusRunning {
		return false
	}

	// if wrong cell
	if !room.Field.IsInside(cell) {
		return false
	}

	// if user died
	if !room.isAlive(conn) {
		return false
	}

	// set who try open cell(for history)
	cell.PlayerID = conn.ID()
	cells := room.Field.OpenCell(cell)
	if len(cells) == 1 {
		newCell := cells[0]
		if newCell.Value == CellMine {
			room.kill(conn, ActionExplode)
			room.Players.Players[conn.Index].Points -= 100
		} else if newCell.Value > CellIncrement {
			room.flagFound(&newCell)
			room.Players.Players[conn.Index].Points += 100000
		}
	} else {
		for _, foundCell := range cells {
			if foundCell.Value < CellMine {
				room.Players.Players[conn.Index].Points += 1 + foundCell.Value
			}
		}
	}

	go room.sendPlayers(room.All)
	go room.sendField(room.All)

	if room.Field.IsCleared() {
		room.finishGame()
	}
	return true
}

func (room *Room) cellHandle(conn *Connection, cell *Cell) (done bool) {
	fmt.Println("cellHandle")
	if room.Status == StatusFlagPlacing {
		done = room.setFlag(conn, cell)
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

	conn.debug("room handle conn")
	if rr.IsGet() {
		room.requestGet(conn, rr)
	} else if rr.IsSend() {
		done := false
		if rr.Send.Cell != nil {
			if room.isAlive(conn) {
				done = room.cellHandle(conn, rr.Send.Cell)
			}
		} else if rr.Send.Action != nil {

			done = room.actionHandle(conn, *rr.Send.Action)
		}
		if !done {
			conn.debug("Room cant execute request")
		}
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

	room.sendField(room.All)
	room.sendMessage("Battle will be start soon! Set your flag!", room.All)
}

func (room *Room) startGame() {
	room.Status = StatusRunning
	room.fillField()
	room.sendMessage("Battle began! Destroy your enemy!", room.All)
}

func (room *Room) finishGame() {
	if room.Status == StatusFinished {
		return
	}
	fmt.Println(room.Name, "We finish room!", room.Status)
	room.chanFinish <- nil
	room.Status = StatusFinished
	fmt.Println(room.Name, "We finish room?", room.Status)
	go room.lobby.roomFinish(room)
	for _, player := range room.Players.Players {
		player.Finished = true
	}
	room.sendMessage("Battle finished!", room.All)
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
	ticker := time.NewTicker(time.Second * 20)

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
			room.finishGame()
			return
		case clock := <-ticker.C:
			fmt.Println("clock!", room.Name)
			room.sendMessage(clock.String()+" passed", room.All)
		}
	}
}

func (room *Room) requestGet(conn *Connection, rr *RoomRequest) {
	send := room.copy(rr.Get)
	conn.SendInformation(send)
}
