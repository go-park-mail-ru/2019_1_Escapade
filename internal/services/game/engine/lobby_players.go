package engine

// waiterToPlayer turns the waiting into a player
func (lobby *Lobby) waiterToPlayer(conn *Connection, room *Room) {
	lobby.s.Do(func() {
		lobby.removeWaiter(conn)
		conn.PushToRoom(room)
		lobby.addPlayer(conn)
	})
}

// PlayerToWaiter turns the player into a waiting
func (lobby *Lobby) PlayerToWaiter(conn *Connection) {
	lobby.s.Do(func() {
		lobby.removePlayer(conn)
		conn.PushToLobby()
		lobby.addWaiter(conn)
	})
}

// restore
// call it before enter connection
func (lobby *Lobby) restore(conn *Connection) bool {
	var room *Room
	lobby.s.Do(func() {
		found := lobby.Playing.Restore(conn)

		if found {
			room = conn.PlayingRoom()
		} else {
			found = lobby.Waiting.Restore(conn)
			if found {
				room = conn.WaitingRoom()
			}
		}

		if room != nil {
			room.connEvents.Reconnect(conn)
		}
	})
	return room != nil
}
