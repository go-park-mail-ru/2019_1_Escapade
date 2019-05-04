package game

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

func (lobby *Lobby) addWaiter(newConn *Connection) {
	fmt.Println("addWaiter called")
	go lobby.waitingAdd(newConn)
	if !newConn.Both() {
		go lobby.greet(newConn)
	}
}

func (lobby *Lobby) addPlayer(newConn *Connection) {
	fmt.Println("addPlayer called")
	go lobby.playingAdd(newConn)
}

func (lobby *Lobby) waiterToPlayer(newConn *Connection) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go waiterToPlayer()")
		lobby.wGroup.Done()
	}()

	fmt.Println("waiterToPlayer called")
	lobby.waitingRemove(newConn)
	lobby.addPlayer(newConn)
}

func (lobby *Lobby) PlayerToWaiter(conn *Connection) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go PlayerToWaiter()")
		lobby.wGroup.Done()
	}()

	fmt.Println("PlayerToWaiter called")
	lobby.playingRemove(conn)
	lobby.addWaiter(conn)
	conn.PushToLobby()
}

func (lobby *Lobby) recoverInRoom(newConn *Connection) bool {
	// find such player
	i, room := lobby.allRoomsSearchPlayer(newConn)

	if i > 0 {
		fmt.Println("we found you in game!")
		room.RecoverPlayer(newConn)
		return true
	}

	// find such observer
	old := lobby.allRoomsSearchObserver(newConn)
	if old != nil {
		room = old.Room()
		if room == nil {
			return false
		}
		room.RecoverObserver(old, newConn)
		return true
	}
	return false
}
