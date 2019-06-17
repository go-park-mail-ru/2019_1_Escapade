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

	lobby.Waiting.Add(newConn /*, false*/)
	lobby.greet(newConn)
	go lobby.sendWaiterEnter(newConn, AllExceptThat(newConn))
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

// waiterToPlayer turns the waiting into a player
func (lobby *Lobby) waiterToPlayer(conn *Connection, room *Room) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go waiterToPlayer()")
		lobby.wGroup.Done()
	}()

	fmt.Println("waiterToPlayer called for ", conn.ID())

	lobby.Waiting.Remove(conn)
	lobby.sendWaiterExit(conn, All)
	conn.PushToRoom(room)
	lobby.Playing.Add(conn)
	lobby.sendPlayerEnter(conn, All)
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

	fmt.Println("PlayerToWaiter called for ", conn.ID())

	lobby.Playing.Remove(conn)
	lobby.sendPlayerExit(conn, All)
	conn.PushToLobby()
	lobby.addWaiter(conn)
}

// r
// new!!!!
// restore
// call it before enter connection
//
func (lobby *Lobby) restore(conn *Connection) bool {

	var found = lobby.Playing.Restore(conn)
	var room *Room

	fmt.Println("restore try")

	if found {
		fmt.Println("found in game")
		room = conn.PlayingRoom()
	} else {
		found = lobby.Waiting.Restore(conn)
		if found {
			fmt.Println("found in waiting room")
			room = conn.WaitingRoom()
		}
	}

	if room != nil {
		fmt.Println("send ActionReconnect")
		room.chanConnection <- ConnectionAction{
			conn:   conn,
			action: ActionReconnect,
		}
	}
	return found
}
