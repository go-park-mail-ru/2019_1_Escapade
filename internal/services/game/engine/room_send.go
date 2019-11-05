package engine

import (
	"sync"

	handlers "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

type RoomSender struct {
	s SyncI
	e *RoomEvents
	p *RoomPeople
	c *RoomConnectionEvents
	i *RoomInformation
	m *RoomModelsConverter
}

func (room *RoomSender) Init(s SyncI, e *RoomEvents, p *RoomPeople,
	c *RoomConnectionEvents, i *RoomInformation, m *RoomModelsConverter) {
	room.s = s
	room.e = e
	room.p = p
	room.c = c
	room.i = i
	room.m = m
}

func (room *RoomSender) sendAll(info handlers.JSONtype, predicate SendPredicate) {
	room.s.do(func() {
		people := room.p.Connections()
		SendToConnections(info, predicate, people...)
	})
}

func (room *RoomSender) Text(text string, predicate SendPredicate) {
	room.sendAll(&models.Result{
		Message: "Room(" + room.i.ID() + "):" + text}, predicate)
}

func (room *RoomSender) Message(message models.Message) {
	response := models.Response{
		Type:  "GameMessage",
		Value: message,
	}
	room.sendAll(&response, room.All)
}

func (room *RoomSender) PlayerPoints(player Player) {
	response := models.Response{
		Type:  "RoomPlayerPoints",
		Value: player,
	}
	room.sendAll(&response, room.All)
}

func (room *RoomSender) GameOver(timer bool, predicate SendPredicate,
	cells []Cell, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	response := room.m.responseRoomGameOver(timer, cells)
	room.sendAll(response, predicate)
}

func (room *RoomSender) NewCells(cells ...Cell) {
	response := models.Response{
		Type:  "RoomNewCells",
		Value: cells,
	}
	room.sendAll(&response, room.All)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *RoomSender) PlayerEnter(conn *Connection) {
	response := models.Response{
		Type:  "RoomPlayerEnter",
		Value: conn,
	}
	room.sendAll(&response, room.AllExceptThat(conn))
}

// sendTAIRPeople send players, observers and history to all in room
func (room *RoomSender) PlayerExit(conn *Connection) {
	response := models.Response{
		Type:  "RoomPlayerExit",
		Value: conn,
	}
	room.sendAll(&response, room.AllExceptThat(conn))
}

// sendTAIRPeople send players, observers and history to all in room
func (room *RoomSender) ObserverEnter(conn *Connection) {
	response := models.Response{
		Type:  "RoomObserverEnter",
		Value: conn,
	}
	room.sendAll(&response, room.AllExceptThat(conn))
}

// sendTAIRPeople send players, observers and history to all in room
func (room *RoomSender) ObserverExit(conn *Connection) {
	response := models.Response{
		Type:  "RoomObserverExit",
		Value: conn,
	}
	room.sendAll(&response, room.AllExceptThat(conn))
}

func (room *RoomSender) StatusToAll(predicate SendPredicate, status int, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	room.s.do(func() {
		response := room.m.responseRoomStatus(status)
		room.sendAll(response, predicate)
	})
}

func (room *RoomSender) StatusToOne(conn *Connection) {
	room.s.doWithConn(conn, func() {
		status := room.e.Status()
		response := room.m.responseRoomStatus(status)
		conn.SendInformation(response)
	})
}

// Action send actions history to all in room
func (room *RoomSender) Action(pa PlayerAction, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomAction",
		Value: pa,
	}
	room.sendAll(&response, predicate)
}

func (room *RoomSender) Error(err error, conn *Connection) {
	room.s.doWithConn(conn, func() {
		response := models.Response{
			Type:  "RoomError",
			Value: err.Error(),
		}
		conn.SendInformation(&response)
	})
}

// FailFlagSet is called when room cant set flag
func (room *RoomSender) FailFlagSet(conn *Connection, value interface{},
	err error) {
	room.s.doWithConn(conn, func() {
		response := models.Response{
			Type:    "FailFlagSet",
			Message: err.Error(),
			Value:   value,
		}
		conn.SendInformation(&response)
	})
}

// RandomFlagSet is called when any player set his flag at the same as any other
func (room *RoomSender) RandomFlagSet(conn *Connection, value interface{}) {
	room.s.doWithConn(conn, func() {
		response := models.Response{
			Type:    "ChangeFlagSet",
			Message: "The cell you have selected is chosen by another person.",
			Value:   value,
		}
		conn.SendInformation(&response)
	})
}

// sendTAIRField send field to all in room
func (room *RoomSender) Field(predicate SendPredicate) {
	room.s.do(func() {
		response := models.Response{
			Type:  "RoomField",
			Value: room.Field,
		}
		room.sendAll(&response, predicate)
	})
}

// sendTAIRAll send everything to one connection
func (room *RoomSender) Room(conn *Connection) {
	room.s.doWithConn(conn, func() {
		isPlayer := room.c.isPlayer(conn)
		conn.SendInformation(room.m.responseRoom(conn, isPlayer))
	})
}

// 302 -> 180
