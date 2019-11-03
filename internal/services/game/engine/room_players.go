package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

/*
addConnection add player to room and to metrics, notify other players, that new one connected
 and provide the connection(client) with the necessary json
*/
func (room *Room) addConnection(conn *Connection, isPlayer bool, needRecover bool) bool {
	if room.done() || conn.done() {
		return false
	}

	room.wGroup.Add(1)
	defer room.wGroup.Done()

	conn.wGroup.Add(1)
	defer conn.wGroup.Done()

	utils.Debug(false, "Room("+room.ID()+") wanna connect you mr ", conn.ID())

	// primary: add player to room
	if !room.Push(conn, isPlayer, needRecover) {
		return false
	}

	// secondary: notify other players, that new connected
	room.NotifyNewConnection(conn, isPlayer, needRecover)
	// primary: provide the connection(client) with the necessary json
	if !needRecover {
		room.wGroup.Add(1)
		room.lobby.sendRoomUpdate(room, All, room.wGroup)
	}
	room.send.StatusToOne(conn)
	room.send.greet(conn, isPlayer)

	utils.Debug(false, "user in room")

	return true
}

func (room *Room) NotifyNewConnection(conn *Connection, isPlayer bool, needRecover bool) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer room.wGroup.Done()

	var pa *PlayerAction
	if needRecover {
		pa = room.addAction(conn.ID(), ActionReconnect)
	} else if isPlayer {
		if !room.Players.EnoughPlace() {
			return
		}
		pa = room.addAction(conn.ID(), ActionConnectAsPlayer)
		// maybe delete it?
		//go room.sendPlayerEnter(*conn, room.AllExceptThat(conn))
	} else {
		if !room.Observers.EnoughPlace() {
			return
		}
		pa = room.addAction(conn.ID(), ActionConnectAsObserver)
		// maybe delete it?
		go room.send.ObserverEnter(conn, room.AllExceptThat(conn))
	}
	go room.send.Action(*pa, room.AllExceptThat(conn))
}

// Push add the connection to the room.
// isPlayer - if true, the connection will add as player, otherwise as observer
// needRecover - if true, then the connection has already added to the room and
// 	it must be restored
// Returns true if added, otherwise false
// If the game has already started, then the connection from waiter slice goes
// to player slice. Otherwise if the game is looking for people then
// the connection remains the waiter, but gets waiting room - this one
func (room *Room) Push(conn *Connection, isPlayer bool, needRecover bool) bool {
	if room.done() || conn.done() {
		return false
	}

	room.wGroup.Add(1)
	defer room.wGroup.Done()

	conn.wGroup.Add(1)
	defer conn.wGroup.Done()

	if isPlayer {
		if !needRecover && !room.Players.EnoughPlace() {
			return false
		}
		room.Players.Add(conn, room.Field.CreateRandomFlag(conn.ID()), needRecover)
		if !needRecover && !room.Players.EnoughPlace() {
			room.recruitingOver()
		}
	} else {
		if !needRecover && !room.Observers.EnoughPlace() {
			return false
		}
		room.Observers.Add(conn)
	}

	if room.Status() != StatusRecruitment {
		room.lobby.waiterToPlayer(conn, room)
	} else {
		conn.setWaitingRoom(room)
	}

	return true
}

// Search search the connection in players slice and observers slice of room
// return connection and flag isPlayer
func (room *Room) Search(find *Connection) (*Connection, bool) {
	i, found := room.Players.SearchConnection(find)
	if i >= 0 {
		return found, true
	}
	i, found = room.Observers.SearchByID(find.ID())
	if i >= 0 {
		return found, false
	}
	return nil, true
}
