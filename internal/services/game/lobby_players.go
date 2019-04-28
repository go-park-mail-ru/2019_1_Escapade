package game

import (
	"fmt"
)

func (lobby *Lobby) addWaiter(newConn *Connection) {
	fmt.Println("addWaiter called")
	lobby.Waiting.Add(newConn, false)
	lobby.greet(newConn)
}

func (lobby *Lobby) addPlayer(newConn *Connection, room *Room) {
	fmt.Println("addPlayer called")
	lobby.Playing.Add(newConn, false)
	room.greet(newConn)
}

func (lobby *Lobby) waiterToPlayer(newConn *Connection, room *Room) {
	fmt.Println("waiterToPlayer called")
	lobby.Waiting.Remove(newConn)
	lobby.addPlayer(newConn, room)
}

func (lobby *Lobby) playerToWaiter(conn *Connection) {
	fmt.Println("playerToWaiter called")
	lobby.Playing.Remove(conn)
	lobby.addWaiter(conn)
	conn.PushToLobby()
}

func (lobby *Lobby) recoverInRoom(newConn *Connection) bool {
	// find such player
	i, room := lobby.AllRooms.SearchPlayer(newConn)

	if i > 0 {
		fmt.Println("we found you in game!")
		room.RecoverPlayer(newConn)
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
