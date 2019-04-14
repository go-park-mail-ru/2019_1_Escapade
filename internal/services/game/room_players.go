package game

// RecoverPlayer call it in lobby.join if player disconnected
func (room *Room) RecoverPlayer(old *Connection, new *Connection) (played bool) {

	// if old in game, disconnect him
	if !old.disconnected {
		old.Kill([]byte("Another connection found"))
	}

	// copy information about player to new connection
	new.Player = old.Player

	// add connection as player
	room.MakePlayer(new)

	room.addAction(new, ActionReconnect)
	room.sendHistory(room.allExceptThat(new))
	room.sendRoom(new)

	return
}

func (room *Room) RecoverObserver(old *Connection, new *Connection) (played bool) {

	if !old.disconnected {
		old.Kill([]byte("Another connection found"))
	}

	room.MakeObserver(new)

	room.addAction(new, ActionReconnect)
	room.sendHistory(room.allExceptThat(new))
	room.sendRoom(new)

	return
}

// observe try to connect user as observer
func (room *Room) addObserver(conn *Connection) bool {
	// if we havent a place
	if !room.Observers.enoughPlace() {
		Answer(conn, []byte("Error. No place in room."))
		return false
	}
	room.MakeObserver(conn)

	room.addAction(conn, ActionConnectAsObserver)

	room.sendObservers(room.allExceptThat(conn))

	room.AnswerOK(conn)
	return true
}

// EnterPlayer handle player try to enter room
func (room *Room) addPlayer(conn *Connection) bool {
	// if room have already started
	// if room.Status != StatusPeopleFinding {
	// 	return false
	// }

	conn.debug("Room(" + room.Name + ") wanna connect you")

	// if room hasnt got places
	if !room.Players.enoughPlace() {
		conn.debug("Room(" + room.Name + ") hasnt any place")
		return false
	}

	cell := room.Field.RandomCell()
	cell.PlayerID = conn.GetPlayerID()

	room.MakePlayer(conn)

	room.addAction(conn, ActionConnectAsPlayer)
	Answer(conn, []byte("OK"))
	room.sendPlayers(room.all())

	if !room.Players.enoughPlace() {
		room.startFlagPlacing()
	}

	return true
}

// MakePlayer mark connection as connected as Player
// add to players slice and set flag inRoom true
func (room *Room) MakePlayer(conn *Connection) {
	conn.PushToRoom(room)
	conn.Player.SetAsPlayer()
	room.Players.Add(conn)
}

// MakePlayer mark connection as connected as Player
// add to players slice and set flag inRoom true
func (room *Room) MakeObserver(conn *Connection) {
	conn.PushToRoom(room)
	conn.Player.SetAsObserver()
	room.Observers.Add(conn)
}

func (room *Room) removeBeforeLaunch(conn *Connection) {
	room.Players.Remove(conn)
	conn.debug("you went back to lobby")
	if room.TryClose() {
		conn.debug("We closed room :С")
	}
}

func (room *Room) removeDuringGame(conn *Connection) {
	i := room.Players.Search(conn)
	if i >= 0 {
		room.GiveUp(conn)
		room.sendHistory(room.all())
		room.sendPlayers(room.all())
		room.TryClose()
		return
	}

	room.Observers.Remove(conn)
	room.sendObservers(room.all())
	room.TryClose()
}

// removeFinishedGame
// func (room *Room) removeAfterLaunch(conn *Connection) {
// 	i := room.Players.Search(conn)
// 	if i >= 0 {
// 		room.TryClose()
// 		return
// 	}

// 	room.Observers.Remove(conn)
// 	room.sendObservers(room.all())
// 	room.TryClose()
// 	return
// }
