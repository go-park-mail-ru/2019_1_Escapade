package engine

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

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

	lobby.removeWaiter(conn)
	conn.PushToRoom(room)
	lobby.addPlayer(conn)
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

	lobby.removePlayer(conn)
	conn.PushToLobby()
	lobby.addWaiter(conn)
}

// restore
// call it before enter connection
func (lobby *Lobby) restore(conn *Connection) bool {

	var found = lobby.Playing.Restore(conn)
	var room *Room

	if found {
		room = conn.PlayingRoom()
	} else {
		found = lobby.Waiting.Restore(conn)
		if found {
			room = conn.WaitingRoom()
		}
	}

	if room != nil {
		room.chanConnection <- &ConnectionAction{
			conn:   conn,
			action: ActionReconnect,
		}
	}
	return found
}
