package game

import (
	"fmt"
)

func (lobby *Lobby) addWaiter(newConn *Connection) {
	lobby.Waiting.Add(newConn)
	lobby.greet(newConn)
}

func (lobby *Lobby) addPlayer(newConn *Connection, room *Room) {
	lobby.Playing.Add(newConn)
	room.greet(newConn)
}

func (lobby *Lobby) waiterToPlayer(newConn *Connection, room *Room) {
	who := lobby.Waiting.Search(newConn)
	lobby.Waiting.Remove(who)
	lobby.addPlayer(newConn, room)
}

func (lobby *Lobby) playerToWaiter(conn *Connection) {
	who := lobby.Waiting.Search(conn)
	lobby.Playing.Remove(who)
	lobby.addWaiter(conn)
	conn.PushToLobby()
}

func (lobby *Lobby) recoverInLobby(newConn *Connection) bool {
	who := lobby.Waiting.Search(newConn)

	if who >= 0 {
		fmt.Println("we found you in lobby!")
		foundConn := lobby.Waiting.Get[who]
		lobby.Waiting.Remove(who)
		//foundConn.SendInformation([]byte("Another connection found"))
		foundConn.Kill("Another connection found", true)
		lobby.addWaiter(newConn)
		return true
	}
	return false
}

func (lobby *Lobby) recoverInRoom(newConn *Connection) bool {
	// find such player
	i, room := lobby.AllRooms.SearchPlayer(newConn)

	if i > 0 {
		fmt.Println("we found you in game!")
		room.RecoverPlayer(i, newConn)
		return true
	}

	// find such observer
	old := lobby.AllRooms.SearchObserver(newConn)
	if old != nil {
		old.room.RecoverObserver(old, newConn)
		return true
	}
	return false
}
