package engine

import (
	"sync"

	handlers "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/synced"
)

// SendStrategyI handle controls the distribution of responses to clients
// Strategy Pattern
type SendStrategyI interface {
	Room(conn *Connection)

	ObserverExit(conn *Connection)
	PlayerExit(conn *Connection)
	PlayerEnter(conn *Connection)
	ObserverEnter(conn *Connection)
	StatusToOne(conn *Connection)

	GameOver(timer bool, predicate SendPredicate, cells []Cell, wg *sync.WaitGroup)

	NewCells(cells ...Cell)
	Text(text string, predicate SendPredicate)
	Field(predicate SendPredicate)

	FailFlagSet(conn *Connection, cell *Cell, err error)
	RandomFlagSet(conn *Connection, cell *Cell)

	PlayerPoints(player Player)

	Message(message models.Message)
	Action(pa PlayerAction, predicate SendPredicate)

	All(conn *Connection) bool
	AllExceptThat(me *Connection) func(*Connection) bool
	StatusToAll(predicate SendPredicate, status int, wg *sync.WaitGroup)
}

// RoomSender implements SendStrategyI
type RoomSender struct {
	s synced.SyncI
	e EventsI
	p PeopleI
	c ConnectionEventsStrategyI
	i RoomInformationI
	m ModelsAdapterI
	f FieldProxyI
}

// Init configure dependencies with other components of the room
func (room *RoomSender) Init(builder ComponentBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildEvents(&room.e)
	builder.BuildPeople(&room.p)
	builder.BuildConnectionEvents(&room.c)
	builder.BuildInformation(&room.i)
	builder.BuildModelsAdapter(&room.m)
	builder.BuildField(&room.f)
}

func (room *RoomSender) sendAll(info handlers.JSONtype, predicate SendPredicate) {
	room.s.Do(func() {
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
	room.s.Do(func() {
		response := room.m.responseRoomStatus(status)
		room.sendAll(response, predicate)
	})
}

func (room *RoomSender) StatusToOne(conn *Connection) {
	room.s.DoWithOther(conn, func() {
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
	room.s.DoWithOther(conn, func() {
		response := models.Response{
			Type:  "RoomError",
			Value: err.Error(),
		}
		conn.SendInformation(&response)
	})
}

// FailFlagSet is called when room cant set flag
func (room *RoomSender) FailFlagSet(conn *Connection, cell *Cell, err error) {
	room.s.DoWithOther(conn, func() {
		response := models.Response{
			Type:    "FailFlagSet",
			Message: err.Error(),
			Value:   cell,
		}
		conn.SendInformation(&response)
	})
}

// RandomFlagSet is called when any player set his flag at the same as any other
func (room *RoomSender) RandomFlagSet(conn *Connection, cell *Cell) {
	room.s.DoWithOther(conn, func() {
		response := models.Response{
			Type:    "ChangeFlagSet",
			Message: "The cell you have selected is chosen by another person.",
			Value:   cell,
		}
		conn.SendInformation(&response)
	})
}

// Field send field to all in room
func (room *RoomSender) Field(predicate SendPredicate) {
	room.s.Do(func() {
		response := models.Response{
			Type:  "RoomField",
			Value: room.f.Field().JSON(),
		}
		room.sendAll(&response, predicate)
	})
}

// sendTAIRAll send everything to one connection
func (room *RoomSender) Room(conn *Connection) {
	room.s.DoWithOther(conn, func() {
		isPlayer := room.c.isPlayer(conn)
		conn.SendInformation(room.m.responseRoom(conn, isPlayer))
	})
}

// 302 -> 180
