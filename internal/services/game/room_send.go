package game

import (
	"encoding/json"
	"sync"
)

// SendPredicate - returns true if the parcel send to that conn
type SendPredicate func(conn *Connection) bool

// sendToAllInRoom send info to those in room, whose predicate
// returns true
func (room *Room) sendToGroup(info interface{}, predicate SendPredicate) {
	waitJobs := &sync.WaitGroup{}
	bytes, _ := json.Marshal(info)
	for _, conn := range room.Players.Get {
		if predicate(conn) {
			waitJobs.Add(1)
			conn.sendGroupInformation(bytes, waitJobs)
		}
	}

	for _, conn := range room.Observers.Get {
		if predicate(conn) {
			waitJobs.Add(1)
			conn.sendGroupInformation(bytes, waitJobs)
		}
	}
	waitJobs.Wait()
}

// allExceptThat is predicat to sendToAllInRoom
// it will send everybody except selected one and disconnected
func (room *Room) allExceptThat(me *Connection) func(conn *Connection) bool {
	return func(conn *Connection) bool {
		return conn != me && conn.disconnected == false && conn.room == room
	}
}

// all is predicat to sendToAllInRoom
// it will send everybody except disconnected
func (room *Room) all() func(conn *Connection) bool {
	return func(conn *Connection) bool {
		return conn.disconnected == false && conn.room == room
	}
}

// onlyThat is predicat to sendToAllInRoom
// it will send to that, even it is disconnected
func (room *Room) onlyThat(me *Connection) func(conn *Connection) bool {
	return func(conn *Connection) bool {
		return conn == me
	}
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendPlayers(predicate SendPredicate) {
	get := &RoomGet{
		Players: true,
	}
	send := room.copyLast(get)
	room.sendToGroup(send, predicate)
}

func (room *Room) sendMessage(text string, predicate SendPredicate) {
	room.sendToGroup("Room("+room.Name+"):"+text, predicate)
}

func (room *Room) sendObservers(predicate SendPredicate) {
	get := &RoomGet{
		Observers: true,
	}
	send := room.copyLast(get)
	room.sendToGroup(send, predicate)
}

// sendTAIRField send field to all in room
func (room *Room) sendField(predicate SendPredicate) {
	get := &RoomGet{
		Field: true,
	}
	send := room.copyLast(get)
	room.sendToGroup(send, predicate)
}

// sendTAIRHistory send actions history to all in room
func (room *Room) sendHistory(predicate SendPredicate) {
	get := &RoomGet{
		History: true,
	}
	send := room.copyLast(get)
	room.sendToGroup(send, predicate)
}

/*
// sendTAIRPeople send only name and status to all in room
func (room *Room) sendTAIRStatus() {
	get := &RoomGet{}
	send := room.makeGetModel(get)
	room.sendToAllInRoom(send)
}
*/
// sendTAIRAll send everything to one connection
func (room *Room) sendRoom(conn *Connection) {
	get := &RoomGet{
		Players:   true,
		Observers: true,
		Field:     true,
		History:   true,
	}
	if room.Status == StatusPeopleFinding {
		get.Field = false
	}
	send := room.copy(get)
	bytes, _ := json.Marshal(send)
	conn.SendInformation(bytes)
}

func (room *Room) AnswerOK(conn *Connection) {
	Answer(conn, []byte("OK"))
	room.sendRoom(conn)
}

func Answer(conn *Connection, message []byte) {
	conn.SendInformation(message)
}

// copy returns full slices of selected fields
func (room *Room) copy(get *RoomGet) *Room {
	sendRoom := &Room{
		Name:   room.Name,
		Status: room.Status,
	}

	if get.Players {
		sendRoom.Players = room.Players
	}
	if get.Observers {
		sendRoom.Observers = room.Observers
	}
	if get.Field {
		sendRoom.Field = room.Field
	}
	if get.History {
		sendRoom.History = room.History
	}
	return sendRoom
}

// copyLast returns last element of slices of selected fields
func (room *Room) copyLast(get *RoomGet) *Room {
	sendRoom := &Room{
		Name:   room.Name,
		Status: room.Status,
	}

	if get.Players {
		sendRoom.Players = room.Players
	}
	if get.Observers {
		sendRoom.Observers = room.Observers
	}
	if get.Field {
		sendRoom.Field = room.Field
	}
	if get.History {
		sendRoom.History = room.History
	}
	return sendRoom
}
