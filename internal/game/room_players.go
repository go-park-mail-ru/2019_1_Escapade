package game

// RecoverPlayer call it in lobby.join if player disconnected
/*
func (room *Room) RecoverPlayer(newConn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	// add connection as player
	// room.MakePlayer(newConn, true)
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
*/

// observe try to connect user as observer
/*
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
}*/

// EnterPlayer handle player try to enter room
func (room *Room) addConnection(conn *Connection, isPlayer bool, needRecover bool) bool {
	//fmt.Println("addPlayer", recover)
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	// later return back!!!!!
	// if room.lobby.Metrics() {
	// 	metrics.Players.WithLabelValues(room.ID, conn.User.Name).Inc()
	// }

	conn.debug("Room(" + room.ID + ") wanna connect you")

	// if room hasnt got places
	if !room.Push(conn, isPlayer, needRecover) {
		return false
	}

	var pa *PlayerAction
	if needRecover {
		pa = room.addAction(conn.ID(), ActionReconnect)
	} else if isPlayer {
		if !room.Players.EnoughPlace() {
			return false
		}
		pa = room.addAction(conn.ID(), ActionConnectAsPlayer)
		// maybe delete it?
		go room.sendPlayerEnter(*conn, room.AllExceptThat(conn))
	} else {
		if !room.Observers.EnoughPlace() {
			return false
		}
		pa = room.addAction(conn.ID(), ActionConnectAsObserver)
		// maybe delete it?
		go room.sendObserverEnter(*conn, room.AllExceptThat(conn))
	}
	go room.sendAction(*pa, room.AllExceptThat(conn))

	if !needRecover {
		room.lobby.sendRoomUpdate(*room, All)

		if !room.Players.EnoughPlace() {
			room.chanStatus <- StatusFlagPlacing
		}
	} else {
		room.sendStatusOne(*conn)
	}

	return true
}

func (room *Room) Push(conn *Connection, isPlayer bool, needRecover bool) bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	if isPlayer {
		if !room.Players.EnoughPlace() {
			return false
		}
		room.Players.Add(conn, room.Field.CreateRandomFlag(conn.ID()), false, needRecover)
		if !room.Players.EnoughPlace() {
			room.StartFlagPlacing()
		}
	} else {
		if !room.Observers.EnoughPlace() {
			return false
		}
		room.Observers.Add(conn)
	}

	room.greet(conn, isPlayer)
	if room.Status != StatusPeopleFinding {
		room.lobby.waiterToPlayer(conn, room)
	} else {
		conn.setWaitingRoom(room)
		//conn.setBoth(true)
	}

	return true

}

/*
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
		room.lobby.waiterToPlayer(conn, room)
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
		room.lobby.waiterToPlayer(conn, room)
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
*/

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
