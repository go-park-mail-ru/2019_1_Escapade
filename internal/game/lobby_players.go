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
	go lobby.sendWaiterEnter(*newConn, AllExceptThat(newConn))
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

	lobby.Waiting.Remove(conn)
	lobby.sendWaiterExit(*conn, All)
	conn.PushToRoom(room)
	lobby.Playing.Add(conn, false)
	lobby.sendPlayerEnter(*conn, All)
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

	lobby.Playing.Remove(conn)
	lobby.sendPlayerExit(*conn, All)
	conn.PushToLobby()
	lobby.addWaiter(conn)
}

// recoverInRoom return true if can find Connection in any room
// otherwise false
func (lobby *Lobby) recoverInRoom(newConn *Connection, disconnect bool) {

	_, room := lobby.allRoomsSearchPlayer(newConn, disconnect)
	if room == nil {
		return
	}
	conn, _ := room.Search(newConn)
	if conn == nil {
		return
	}
	if !room.done() {
		room.chanConnection <- ConnectionAction{
			conn:   newConn,
			action: ActionConnect,
		}
	}
}

// r
// new!!!!
// restore
// call it before enter connection
//
func (lobby *Lobby) restore(conn *Connection) bool {

	var found = lobby.Playing.Restore(conn)

	if found {
		room := conn.Room()
		if !room.done() {
			room.chanConnection <- ConnectionAction{
				conn:   conn,
				action: ActionReconnect, //bug::: не обрабатывается в руме, поправить
			}
		}
		return found
	}
	found = lobby.Waiting.Restore(conn)

	return found
}
