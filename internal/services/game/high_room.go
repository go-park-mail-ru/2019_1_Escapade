package game

import (
	"escapade/internal/models"
	"fmt"
)

// Game status
const (
	StatusPeopleFinding = iota
	StatusAborted       // in case of error
	StatusFlagPlacing
	StatusRunning
	StatusFinished
	StatusClosed
)

type Room struct {
	Name   string `json:"name"`
	Status int    `json:"status"`

	Players   *Connections `json:"players"`
	Observers *Connections `json:"observers"`

	History []*PlayerAction `json:"history"`

	flags map[*Connection]*models.Cell

	lobby *Lobby
	Field *models.Field `json:"field"`

	chanLeave chan *Connection
	//chanRequest chan *RoomRequest
}

func (room *Room) addAction(conn *Connection, action int) {
	pa := NewPlayerAction(conn.Player, action)
	room.History = append(room.History, pa)
}

// SameAs compare  one room with another
func (room *Room) SameAs(another *Room) bool {
	return room.Field.SameAs(another.Field)
}

// Join handle user joining as player or observer
func (room *Room) Join(conn *Connection) bool {

	// if game not finish, lets check is that conn already in game
	if room.Status != StatusFinished {
		if room.alreadyPlaying(conn) {
			return true
		}
	}

	// reset old points
	conn.Player.Reset()

	// if room is searching new players
	if room.Status == StatusPeopleFinding {
		if room.EnterPlayer(conn) {
			return true
		}
	}

	// if you cant play, try observe
	if room.enterObserver(conn) {
		return true
	}

	return false
}

func (room *Room) Leave(conn *Connection) {

	// cant delete players, cause they always need
	// if game began
	if room.Status == StatusPeopleFinding {
		room.removeBeforeLaunch(conn)
	} else {
		room.removeAfterLaunch(conn)
	}
	room.addAction(conn, ActionDisconnect)
	room.sendTAIRPeople()
	return
}

func (room *Room) setFlag(conn *Connection, cell *models.Cell) bool {
	// if user try set flag after game launch
	if room.Status != StatusFlagPlacing {
		return false
	}

	if !room.Field.IsInside(cell) {
		return false
	}
	room.flags[conn] = cell
	return true
}

// nanfle openCell
func (room *Room) openCell(conn *Connection, cell *models.Cell) bool {
	// if user try set open cell before game launch
	if room.Status != StatusRunning {
		return false
	}

	// if wrong cell
	if !room.Field.IsInside(cell) {
		return false
	}

	// if user died
	if conn.Player.Finished == true {
		return false
	}

	// set who try open cell(for history)
	cell.PlayerID = conn.GetPlayerID()
	room.Field.OpenCell(cell)

	room.sendTAIRField()

	if room.Field.IsCleared() {
		room.lobby.roomFinish(room)
	}
	return true
}

func (room *Room) cellHandle(conn *Connection, cell *models.Cell) (done bool) {
	fmt.Println("cellHandle")
	if cell.Value == models.CellFlag {
		done = room.setFlag(conn, cell)
	} else {
		done = room.openCell(conn, cell)
	}
	return
}

func (room *Room) actionHandle(conn *Connection, action int) (done bool) {
	if action == ActionGiveUp {
		room.GiveUp(conn)
		return true
	}
	return false
}

// handleRequest
func (room *Room) handleRequest(conn *Connection, rr *RoomRequest) {

	if rr.IsGet() {
		room.requestGet(conn, rr)
	} else if rr.IsSend() {
		done := false
		if rr.Send.Cell != nil {
			done = room.cellHandle(conn, rr.Send.Cell)
		} else if rr.Send.Action != nil {
			done = room.actionHandle(conn, *rr.Send.Action)
		}
		if !done {
			sendError(conn, "room request", "Cant execute request ")
		}
	}
}

func (room *Room) startFlagPlacing() {
	room.Status = StatusRunning //StatusFlagPlacing
	fmt.Println("startFlagPlacing 1 ")
	room.lobby.roomStart(room)
	fmt.Println("startFlagPlacing 2 ")
	room.fillField()
	fmt.Println("startFlagPlacing 3 ")
	room.sendTAIRField()
}

func (room *Room) startGame() {
	room.Status = StatusRunning
	room.fillField()
}

// Run the room in goroutine
func (room *Room) run() {
	//timer := time.NewTimer()
	for {
		select {
		case connection := <-room.chanLeave:
			room.Leave(connection)
		}
	}
}
