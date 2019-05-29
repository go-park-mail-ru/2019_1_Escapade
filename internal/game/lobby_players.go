package game

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// addWaiter add connection to waiters slice and send to the connection LobbyJSON
func (lobby *Lobby) addWaiter(newConn *Connection) {
	if lobby.metrics {
		metrics.WaitingPlayers.Add(1)
	}
	fmt.Println("addWaiter called")
	lobby.Waiting.Add(newConn, false)
	if !newConn.Both() {
		lobby.greet(newConn)
	}
}

// Anonymous return anonymous id
func (lobby *Lobby) Anonymous() int {
	var id int
	lobby.anonymousM.Lock()
	id = lobby._anonymous
	lobby._anonymous--
	lobby.anonymousM.Unlock()
	return id
}

// addPlayer add connection to players slice
func (lobby *Lobby) addPlayer(newConn *Connection) {
	fmt.Println("addPlayer called")
	lobby.Playing.Add(newConn, false)
}

// waiterToPlayer turns the waiting into a player
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
	lobby.Waiting.Remove(newConn, true)
	//lobby.waitingRemove(newConn)
	lobby.addPlayer(newConn)
}

// PlayerToWaiter turns the player into a waiting
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
	lobby.Playing.Remove(conn, true)
	//lobby.playingRemove(conn)
	lobby.addWaiter(conn)
	conn.PushToLobby()
}

// recoverInRoom return true if can find Connection in any room
// otherwise false
func (lobby *Lobby) recoverInRoom(newConn *Connection, disconnect bool) bool {
	// find such player
	fmt.Println("somebody wants you to recover")

	i, room := lobby.allRoomsSearchPlayer(newConn, disconnect)

	if i >= 0 {
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
