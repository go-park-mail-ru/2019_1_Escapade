package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// sendToAllInRoom send info to those in room, whose predicate
// returns true
func (room *Room) send(info interface{}, predicate SendPredicate) {
	players := room.playersConnections()
	observers := room.observers()
	SendToConnections(info, predicate, players, observers)
}

func (room *Room) sendMessage(text string, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go sendMessage()")
	}()

	room.send("Room("+room.ID+"):"+text, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayerPoints(player Player, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go sendPlayerPoints()")
	}()

	response := models.Response{
		Type:  "RoomPlayerPoints",
		Value: player,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendGameOver(predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go sendGameOver()")
	}()

	cells := make([]Cell, 0)
	room.Field.OpenEverything(&cells)
	response := models.Response{
		Type: "RoomGameOver",
		Value: struct {
			Players []Player `json:"players"`
			Cells   []Cell   `json:"cells"`
		}{
			Players: room.players(),
			Cells:   cells,
		},
	}
	room.send(response, predicate)
}

func (room *Room) sendNewCells(cells []Cell, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go RoomNewCells()")
	}()

	response := models.Response{
		Type:  "RoomNewCells",
		Value: cells,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayerEnter(conn Connection, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go RoomPlayerEnter()")
	}()

	response := models.Response{
		Type:  "RoomPlayerEnter",
		Value: conn,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayerExit(conn Connection, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go RoomPlayerExit()")
	}()

	response := models.Response{
		Type:  "RoomPlayerExit",
		Value: conn,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendObserverEnter(conn Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomObserverEnter",
		Value: conn,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendObserverExit(conn Connection, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go RoomObserverExit()")
	}()

	response := models.Response{
		Type:  "RoomObserverExit",
		Value: conn,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendStatus(predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go RoomStatus()")
	}()

	response := models.Response{
		Type: "RoomStatus",
		Value: struct {
			ID     string `json:"id"`
			Status int    `json:"status"`
		}{
			ID:     room.ID,
			Status: room.Status,
		},
	}
	room.send(response, predicate)
}

// sendTAIRHistory send actions history to all in room
func (room *Room) sendAction(pa PlayerAction, predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go sendAction()")
	}()

	response := models.Response{
		Type:  "RoomAction",
		Value: pa,
	}
	room.send(response, predicate)
}

// sendTAIRHistory send actions history to all in room
func (room *Room) sendError(err error, conn Connection) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go sendError()")
	}()

	response := models.Response{
		Type:  "RoomError",
		Value: err.Error(),
	}
	conn.SendInformation(response)
}

// sendTAIRField send field to all in room
func (room *Room) sendField(predicate SendPredicate) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go sendField()")
	}()

	response := models.Response{
		Type:  "RoomField",
		Value: room.Field,
	}
	room.send(response, predicate)
}

// sendTAIRAll send everything to one connection
func (room *Room) greet(conn *Connection) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("room_send.go greet()")
	}()

	var flag *Cell
	if conn.Index() >= 0 {
		flag = room.setCell(conn)
	}
	response := models.Response{
		Type: "Room",
		Value: struct {
			Room *Room `json:"room"`
			Flag *Cell `json:"flag,omitempty"`
		}{
			Room: room,
			Flag: flag,
		},
	}
	conn.SendInformation(response)
}
