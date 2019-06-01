package game

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
)

// RecoverPlayer call it in lobby.join if player disconnected
func (room *Room) RecoverPlayer(newConn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	// add connection as player
	room.MakePlayer(newConn, true)
	pa := *room.addAction(newConn.ID(), ActionReconnect)
	room.addPlayer(newConn, true)
	room.sendAction(pa, room.AllExceptThat(newConn))
	//room.greet(newConn, true)

	return
}

// RecoverObserver recover connection as observer
func (room *Room) RecoverObserver(newConn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	go room.MakeObserver(newConn, true)
	pa := *room.addAction(newConn.ID(), ActionReconnect)
	go room.sendAction(pa, room.AllExceptThat(newConn))
	//go room.greet(newConn, false)

	return
}

// observe try to connect user as observer
func (room *Room) addObserver(conn *Connection) bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	if room.lobby.Metrics() {
		metrics.Players.WithLabelValues(room.ID, conn.User.Name).Inc()
	}

	// if we havent a place
	if !room.Observers.EnoughPlace() {
		conn.debug("Room cant execute request ")
		return false
	}
	conn.debug("addObserver")
	room.MakeObserver(conn, true)

	go room.addAction(conn.ID(), ActionConnectAsObserver)
	go room.sendObserverEnter(*conn, room.AllExceptThat(conn))
	room.lobby.sendRoomUpdate(*room, All)

	return true
}

// EnterPlayer handle player try to enter room
func (room *Room) addPlayer(conn *Connection, recover bool) bool {
	fmt.Println("addPlayer", recover)
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	if room.lobby.Metrics() {
		metrics.Players.WithLabelValues(room.ID, conn.User.Name).Inc()
	}

	conn.debug("Room(" + room.ID + ") wanna connect you")

	// if room hasnt got places
	if !recover && !room.Players.EnoughPlace() {
		conn.debug("Room(" + room.ID + ") hasnt any place")
		return false
	}

	room.MakePlayer(conn, recover)

	go room.addAction(conn.ID(), ActionConnectAsPlayer)
	go room.sendPlayerEnter(*conn, room.AllExceptThat(conn))

	if !recover {
		room.lobby.sendRoomUpdate(*room, All)
		room.lobby.sendRoomToOne(*room, *conn)

		if !room.Players.EnoughPlace() {
			room.chanStatus <- StatusFlagPlacing
		}
	}

	return true
}

// MakePlayer mark connection as connected as Player
// add to players slice and set flag inRoom true
func (room *Room) MakePlayer(conn *Connection, recover bool) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	if room.Status != StatusPeopleFinding {
		room.lobby.waiterToPlayer(conn)
		conn.setBoth(false)
	} else {
		conn.setBoth(true)
	}
	room.Players.Add(conn, room.Field.CreateRandomFlag(conn.ID()), false, recover)
	fmt.Println("MakePlayer", recover)
	room.greet(conn, true)
	if recover {
		room.sendStatusOne(*conn)
	}
	conn.PushToRoom(room)
}

// MakeObserver mark connection as connected as Observer
// add to observers slice and set flag inRoom true
func (room *Room) MakeObserver(conn *Connection, recover bool) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	if room.Status != StatusPeopleFinding {
		room.lobby.waiterToPlayer(conn)
		conn.setBoth(false)
	} else {
		conn.setBoth(true)
	}
	room.Observers.Add(conn, false)
	room.greet(conn, false)
	if recover {
		room.sendStatus(Me(conn))
	}
	conn.PushToRoom(room)
}

// Search search connection in players and observers of room
// return connection and flag isPlayer
func (room *Room) Search(find *Connection) (*Connection, bool) {
	found, i := room.Players.SearchConnection(find)
	if i >= 0 {
		return found, true
	}
	found, i = room.Observers.SearchByID(find.ID())
	if i >= 0 {
		return found, false
	}
	return nil, true
}
