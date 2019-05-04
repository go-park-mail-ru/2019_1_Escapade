package game

import "fmt"

// RecoverPlayer call it in lobby.join if player disconnected
func (room *Room) RecoverPlayer(newConn *Connection) {

	// add connection as player
	room.MakePlayer(newConn)
	pa := *room.addAction(newConn.ID(), ActionReconnect)
	room.sendAction(pa, room.AllExceptThat(newConn))
	room.greet(newConn)

	return
}

// RecoverObserver recover connection as observer
func (room *Room) RecoverObserver(oldConn *Connection, newConn *Connection) {

	room.MakeObserver(newConn)
	pa := *room.addAction(newConn.ID(), ActionReconnect)
	room.sendAction(pa, room.AllExceptThat(newConn))
	room.greet(newConn)

	return
}

// observe try to connect user as observer
func (room *Room) addObserver(conn *Connection) bool {
	// if we havent a place
	if !room.Observers.enoughPlace() {
		conn.debug("Room cant execute request ")
		return false
	}
	conn.debug("addObserver")
	room.MakeObserver(conn)

	room.addAction(conn.ID(), ActionConnectAsObserver)
	room.sendObserverEnter(*conn, room.AllExceptThat(conn))
	room.lobby.sendRoomUpdate(*room, All)
	room.greet(conn)

	return true
}

// EnterPlayer handle player try to enter room
func (room *Room) addPlayer(conn *Connection) bool {
	// if room have already started
	// if room.Status != StatusPeopleFinding {
	// 	return false
	// }

	conn.debug("Room(" + room.ID + ") wanna connect you")

	// if room hasnt got places
	if !room.Players.enoughPlace() {
		conn.debug("Room(" + room.ID + ") hasnt any place")
		return false
	}

	room.MakePlayer(conn)

	room.addAction(conn.ID(), ActionConnectAsPlayer)
	room.sendPlayerEnter(*conn, room.AllExceptThat(conn))
	room.lobby.sendRoomUpdate(*room, All)
	room.greet(conn)

	if !room.Players.enoughPlace() {
		room.startFlagPlacing()
	}

	return true
}

// MakePlayer mark connection as connected as Player
// add to players slice and set flag inRoom true
func (room *Room) MakePlayer(conn *Connection) {
	if room.Status != StatusPeopleFinding {
		room.lobby.waiterToPlayer(conn, room)
		conn.both = false
	} else {
		conn.both = true
	}
	room.Players.Add(conn, false)
	conn.PushToRoom(room)
}

// MakeObserver mark connection as connected as Observer
// add to observers slice and set flag inRoom true
func (room *Room) MakeObserver(conn *Connection) {
	if room.Status != StatusPeopleFinding {
		room.lobby.waiterToPlayer(conn, room)
		conn.both = false
	} else {
		conn.both = true
	}
	room.Observers.Add(conn, false)
	room.sendObservers(room.All)
	conn.PushToRoom(room)
}

func (room *Room) removeFromGame(conn *Connection, disconnected bool) {
	fmt.Println("removeDuringGame")
	i := room.Players.SearchIndexPlayer(conn)
	if i >= 0 {
		if room.Status == StatusRunning && !disconnected {
			fmt.Println("give up", i)
			room.GiveUp(conn)
		}

		room.Players.Remove(conn)
		room.sendPlayerExit(*conn, room.All)
	} else {
		room.sendObserverExit(*conn, room.All)
		room.Observers.Remove(conn)
		room.sendObservers(room.All)
	}
	fmt.Println("room.Players len", len(room.Players.Players))
	if room.Players.Empty() {
		fmt.Println("room.Players.Empty")
		room.Close()
	}
	fmt.Println("there")
}