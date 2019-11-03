package engine

import (
	"sync"

	handlers "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

type RoomSender struct {
	r *Room
}

// sendToAllInRoom send info to those in room, whose predicate returns true
// sendUnsafe is goroutine unsafe. Use sendAll for goroutine safe use
func (room *RoomSender) sendAllUnsafe(info handlers.JSONtype, predicate SendPredicate) {
	players := room.r.Players.Connections
	observers := room.r.Observers
	SendToConnections(info, predicate, players, observers)
}

func (room *RoomSender) sendAll(info handlers.JSONtype, predicate SendPredicate) {
	if room.r.done() {
		return
	}
	room.r.wGroup.Add(1)
	defer func() {
		room.r.wGroup.Done()
		utils.CatchPanic("room_send.go sendMessage()")
	}()
	room.sendAll(info, predicate)
}

func (room *RoomSender) Message(text string, predicate SendPredicate) {
	room.sendAll(&models.Result{
		Message: "Room(" + room.r.ID() + "):" + text}, predicate)
}

func (room *RoomSender) PlayerPoints(player Player, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomPlayerPoints",
		Value: player,
	}
	room.sendAll(&response, predicate)
}

func (room *RoomSender) GameOver(timer bool, predicate SendPredicate,
	cells []Cell, wg *sync.WaitGroup) {
	if room.r.done() {
		return
	}
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	response := room.r.models.responseRoomGameOver(timer, cells)
	room.sendAll(response, predicate)
}

func (room *RoomSender) NewCells(predicate SendPredicate, cells ...Cell) {
	response := models.Response{
		Type:  "RoomNewCells",
		Value: cells,
	}
	room.sendAll(&response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *RoomSender) PlayerEnter(conn *Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomPlayerEnter",
		Value: conn,
	}
	room.sendAll(&response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *RoomSender) PlayerExit(conn *Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomPlayerExit",
		Value: conn,
	}
	room.sendAll(&response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *RoomSender) ObserverEnter(conn *Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomObserverEnter",
		Value: conn,
	}
	room.sendAll(&response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *RoomSender) ObserverExit(conn *Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomObserverExit",
		Value: conn,
	}
	room.sendAll(&response, predicate)
}

func (room *RoomSender) StatusToAll(predicate SendPredicate, status int, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	if room.r.done() {
		return
	}

	response := room.r.models.responseRoomStatus(status)
	room.sendAll(response, predicate)
}

func (room *RoomSender) StatusToOne(conn *Connection) {
	if room.r.done() {
		return
	}
	status := room.r.Status()
	response := room.r.models.responseRoomStatus(status)
	conn.SendInformation(response)
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
	response := models.Response{
		Type:  "RoomError",
		Value: err.Error(),
	}
	conn.SendInformation(&response)
}

// sendTAIRField send field to all in room
func (room *RoomSender) Field(predicate SendPredicate) {
	if room.r.done() {
		return
	}

	response := models.Response{
		Type:  "RoomField",
		Value: room.Field,
	}
	room.sendAll(&response, predicate)
}

// sendTAIRAll send everything to one connection
func (room *RoomSender) greet(conn *Connection, isPlayer bool) {
	if room.r.done() {
		return
	}
	room.r.wGroup.Add(1)
	defer func() {
		room.r.wGroup.Done()
		utils.CatchPanic("room_send.go greet()")
	}()
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer conn.wGroup.Done()
	conn.SendInformation(room.r.models.responseRoom(conn, isPlayer))
}

// 302 -> 180
