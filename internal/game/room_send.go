package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// sendToAllInRoom send info to those in room, whose predicate
// returns true
func (room *Room) send(info interface{}, predicate SendPredicate) {
	SendToConnections(info, predicate, room.Players.Connections,
		room.Observers.Get)
}

func (room *Room) sendMessage(text string, predicate SendPredicate) {
	room.send("Room("+room.ID+"):"+text, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayerPoints(player Player, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomPlayerPoints",
		Value: player,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayers(predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomGameOver",
		Value: room.Players.Players,
	}
	room.send(response, predicate)
}

func (room *Room) sendNewCells(cells []Cell, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomNewCells",
		Value: cells,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayerEnter(conn Connection, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomPlayerEnter",
		Value: conn,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayerExit(conn Connection, predicate SendPredicate) {
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
	response := models.Response{
		Type:  "RoomObserverExit",
		Value: conn,
	}
	room.send(response, predicate)
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendStatus(predicate SendPredicate) {
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

func (room *Room) sendObservers(predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomConnectionsObservers",
		Value: room.Observers,
	}
	room.send(response, predicate)
}

// sendTAIRField send field to all in room
func (room *Room) sendField(predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomField",
		Value: room.Field,
	}
	room.send(response, predicate)
}

// sendTAIRHistory send actions history to all in room
func (room *Room) sendAction(pa PlayerAction, predicate SendPredicate) {
	response := models.Response{
		Type:  "RoomAction",
		Value: pa,
	}
	room.send(response, predicate)
}

// sendTAIRAll send everything to one connection
func (room *Room) greet(conn *Connection) {

	var flag *Cell
	if conn.Index >= 0 {
		flag = &room.Players.Flags[conn.Index]
	}
	response := models.Response{
		Type: "Room",
		Value: struct {
			Room *Room `json:"room"`
			Flag *Cell `json:"flag, omitempty"`
		}{
			Room: room,
			Flag: flag,
		},
	}
	conn.SendInformation(response)
}
