package game

import "fmt"

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
	go room.MakePlayer(newConn)
	pa := *room.addAction(newConn.ID(), ActionReconnect)
	go room.sendAction(pa, room.AllExceptThat(newConn))
	go room.greet(newConn)

	return
}

// RecoverObserver recover connection as observer
func (room *Room) RecoverObserver(oldConn *Connection, newConn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	go room.MakeObserver(newConn)
	pa := *room.addAction(newConn.ID(), ActionReconnect)
	go room.sendAction(pa, room.AllExceptThat(newConn))
	go room.greet(newConn)

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

	// if we havent a place
	if !room.observersEnoughPlace() {
		conn.debug("Room cant execute request ")
		return false
	}
	conn.debug("addObserver")
	go room.MakeObserver(conn)

	go room.addAction(conn.ID(), ActionConnectAsObserver)
	go room.sendObserverEnter(*conn, room.AllExceptThat(conn))
	room.lobby.sendRoomUpdate(*room, All)

	return true
}

// EnterPlayer handle player try to enter room
func (room *Room) addPlayer(conn *Connection) bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	// if room have already started
	// if room.Status != StatusPeopleFinding {
	// 	return false
	// }

	conn.debug("Room(" + room.ID + ") wanna connect you")

	// if room hasnt got places
	if !room.playersEnoughPlace() {
		conn.debug("Room(" + room.ID + ") hasnt any place")
		return false
	}

	room.MakePlayer(conn)

	go room.addAction(conn.ID(), ActionConnectAsPlayer)
	go room.sendPlayerEnter(*conn, room.AllExceptThat(conn))
	go room.lobby.sendRoomUpdate(*room, All)

	fmt.Println("len", room._Players.Connections)
	if !room.playersEnoughPlace() {
		room.StartFlagPlacing()
	}

	return true
}

// MakePlayer mark connection as connected as Player
// add to players slice and set flag inRoom true
func (room *Room) MakePlayer(conn *Connection) {
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
	room.playersAdd(conn, false)
	go room.greet(conn)
	conn.PushToRoom(room)
}

// MakeObserver mark connection as connected as Observer
// add to observers slice and set flag inRoom true
func (room *Room) MakeObserver(conn *Connection) {
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
	room.observersAdd(conn, false)
	go room.greet(conn)
	conn.PushToRoom(room)
}

func (room *Room) RemoveFromGame(conn *Connection, disconnected bool) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	fmt.Println("removeDuringGame")
	fmt.Println("removeDuringGame before len", len(room._Players.Connections))

	i := room.playersSearchIndexPlayer(conn)
	if i >= 0 {
		if room.Status == StatusRunning && !disconnected {
			fmt.Println("give up", i)
			room.GiveUp(conn)
		}

		room.playersRemove(conn)
		go room.sendPlayerExit(*conn, room.All)
	} else {
		go room.sendObserverExit(*conn, room.All)
		room.observersRemove(conn)
	}

	fmt.Println("removeDuringGame after len", len(room._Players.Connections))
	fmt.Println("removeDuringGame system says", room.playersEmpty())
	if room.playersEmpty() {
		fmt.Println("room.Players.Empty")
		room.Close()
	}
	fmt.Println("there")
}
