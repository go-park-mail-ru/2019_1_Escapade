package game

import (
	"escapade/internal/models"
	//re "escapade/internal/return_errors"
	//"math/rand"
)

type Request struct {
	Connection *Connection
	Cell       *models.Cell
}

func NewRequest(conn *Connection, cell *models.Cell) *Request {
	request := &Request{
		Connection: conn,
		Cell:       cell,
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

func (room *Room) SetFlag(req *Request) {
	// if user try set flag after game launch
	if room.Status != StatusFlagPlacing {
		sendNotAllowed(req.Connection)
		return
	}

	if !room.Field.IsInside(req.Cell) {
		sendNotAllowed(req.Connection)
		return
	}

	room.Players[req.Connection].Flag = req.Cell
	req.Connection.SendInformation(models.ActionFlagSet)
}

func (room *Room) GetRequest(req *Request) {
	if req.Cell.Value == models.CellFlag {
		room.SetFlag(req)
	} else {
		room.OpenCell(req)
	}
}

func (room *Room) OpenCell(req *Request) {
	// if user try set open cell before game launch
	if room.Status != StatusRunning {
		sendNotAllowed(req.Connection)
		return
	}

	if !room.Field.IsInside(req.Cell) {
		sendNotAllowed(req.Connection)
		return
	}

	// set who try open cell(for history)
	req.Cell.PlayerID = req.Connection.GetPlayerID()
	room.Field.OpenCell(req.Cell)

	req.Connection.SendInformation(models.ActionFlagSet)
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
	for {
		select {

		case connection := <-room.chanLeave:
			room.Leave(connection)
		case request := <-room.chanRequest:
			room.GetRequest(request)
		}
	}
}
