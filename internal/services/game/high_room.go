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

type Room struct {
	ID     int
	Status int

	PlayersCapacity int
	PlayersSize     int
	Players         map[*Connection]*Playing

	ObserversCapacity int
	ObserversSize     int
	Observers         map[*Connection]*models.Player

	Field         *models.Field
	chanUpdateAll chan *struct{}

	chanJoin    chan *Connection
	chanLeave   chan *Connection
	chanRequest chan *Request
}

// we use map, not array, cause in future will add name of rooms
var allRooms = make(map[int]*Room)
var freeRooms = make(map[int]*Room)
var roomsCount int

// вынести в конфиг
var roomIDMax = 10000

// sendPrepareInfo send preparing info
func (room *Room) sendPrepareInfo(conn *Connection) {

	room.sendPeople(conn)
	room.sendField(conn)
}

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
		if room.addPlayer(conn) {
			return
		}
	}

	// if you cant play, try observe
	if room.addObserver(conn) {
		return
	}

	// room not ready to accept you
	room.sendRoomIsBlocked(conn)
}

func (room *Room) Leave(conn *Connection) {

	// cant delete players, cause they always need
	// if game began
	switch room.Status {
	case models.StatusPeopleFinding:
		room.removeBeforeLaunch(conn)
	case models.StatusRunning:
		room.removeDuringGame(conn)
	default:
		room.removeAfterFinish(conn)
	}
	return
}

func (room *Room) SetFlag(conn *Connection) {
	if room.Status != models.StatusFlagPlacing {

	}
}

// Run the room in goroutine
func (room *Room) run() {
	for {
		select {
		case connection := <-room.chanJoin:
			room.Join(connection)

		case connection := <-room.chanLeave:
			room.Leave(connection)
		}
	}
}
