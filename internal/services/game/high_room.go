package game

import (
	"escapade/internal/models"
	//re "escapade/internal/return_errors"
	//"math/rand"
)

type Request struct {
	Connection *Connection
	Data       *models.ClientData
}

func NewRequest(conn *Connection, data *models.ClientData) *Request {
	request := &Request{
		Connection: conn,
		Data:       data,
	}
	return request
}

type Rooms struct {
	Size  int
	Rooms []Room
}

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
	ID     int
	Status int

	PlayersCapacity int
	PlayersSize     int
	Players         map[*Connection]*Playing

	ObserversCapacity int
	ObserversSize     int
	Observers         map[*Connection]*models.Player

	lobby *Lobby
	Field *models.Field
	//chanUpdateAll chan *struct{}

	//chanJoin    chan *Connection
	chanLeave   chan *Connection
	chanRequest chan *Request
}

// sendPrepareInfo send preparing info
func (room *Room) sendPrepareInfo(conn *Connection) {

	room.sendPeople(conn)
	room.sendField(conn)
}

/*
// join handle user joining as player or observer
func (room *Room) Join(conn *Connection) {

	// if game not finish, lets check is that conn already in game
	if room.Status != models.StatusFinished {
		if room.alreadyPlaying(conn) {
			return
		}
	}

	// reset old action and old points
	conn.player.Reset()

	// if room is searching new players
	if room.Status == models.StatusPeopleFinding {
		if room.enterPlayer(conn) {
			return
		}
	}

	// if you cant play, try observe
	if room.enterObserver(conn) {
		return
	}

	// room not ready to accept you
	sendNotAllowed(conn)
}
*/

func (room *Room) Leave(conn *Connection) {

	// cant delete players, cause they always need
	// if game began
	switch room.Status {
	case StatusPeopleFinding:
		room.removeBeforeLaunch(conn)
	case StatusRunning:
		room.removeDuringGame(conn)
	default:
		room.removeAfterFinish(conn)
	}
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

	room.Players[conn].Flag = cell
	// send for this user, that his flag posision accepted
	conn.SendInformation(models.ActionFlagSet)
	return true
}

func (room *Room) kill(conn *Connection) {
	if !room.Players[conn].Finished {
		room.Players[conn].Finished = true
		room.PlayersSize--
		if room.PlayersSize <= 1 {
			room.lobby.roomFinish(room)
		}
		room.sendAllPlayerAction(conn, models.ActionLose)
	}
}

func (room *Room) flagFound(found *models.Cell) {
	id := found.Value - models.CellIncrement
	for conn, _ := range room.Players {
		if conn.GetPlayerID() == id {
			room.kill(conn)
		}
	}
}

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
	if room.Players[conn].Finished == true {
		return false
	}

	// set who try open cell(for history)
	cell.PlayerID = conn.GetPlayerID()
	cells := room.Field.OpenCell(cell)
	if len(cells) == 1 {
		if cells[0].Value == models.CellMine {
			room.kill(conn)
		} else if cells[0].Value == models.CellFlag {
			room.flagFound(&cells[0])
		}

	}
	room.sendCells(cells)

	// send for this user, that his flag posision accepted
	conn.SendInformation(models.ActionFlagSet)
	return true
}

func (room *Room) cellHandle(conn *Connection, cell *models.Cell) (done bool) {
	if cell.Value == models.CellFlag {
		done = room.setFlag(conn, cell)
	} else {
		done = room.openCell(conn, cell)
	}
	return
}

func (room *Room) actionHandle(conn *Connection, action int) (done bool) {
	if action == models.ActionGiveUp {
		room.Players[conn].Finished = true
		return true
	}
	return false
}

func (room *Room) GetRequest(req *Request) {
	done := false
	switch req.Data.Send {
	//case models.SendRoomSettings:
	//
	case models.SendPlayerAction:
		done = room.actionHandle(req.Connection, req.Data.PlayerAction)
	case models.SendCells:
		done = room.cellHandle(req.Connection, req.Data.Cell)
	}
	if !done {
		sendNotAllowed(req.Connection)
	}
}

func (room *Room) startFlagPlacing() {
	room.Status = StatusFlagPlacing
	room.lobby.roomStart(room)
	room.fillField()
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
		case request := <-room.chanRequest:
			room.GetRequest(request)
		}
	}
}
