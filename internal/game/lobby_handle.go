package game

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"context"
	"encoding/json"
	"fmt"
)

// Run the room in goroutine
func (lobby *Lobby) Run(wg *sync.WaitGroup) {
	defer func() {
		utils.CatchPanic("lobby_handle.go Run()")
		lobby.Free()
	}()

	var lobbyCancel context.CancelFunc
	lobby.Context, lobbyCancel = context.WithCancel(context.Background())
	if wg != nil {
		wg.Done()
	}
	fmt.Println("Lobby run")
	for {
		select {
		case connection := <-lobby.ChanJoin:
			go lobby.Join(connection)

		case message := <-lobby.chanBroadcast:
			go lobby.analize(message)

			// TODO delete chanleavem cause Leave call direcrly
		case connection := <-lobby.chanLeave:
			lobby.Leave(connection, "You disconnected!")
			// if need_stop {
			// 	if len(lobby.Playing.Get)+len(lobby.Waiting.Get) == 0 {
			// 		fmt.Println("Nobody there!")
			// 		lobbyCancel()
			// 		return
			// 	}
			// }
		case <-lobby.chanBreak:
			fmt.Println("Stop saw!")
			lobbyCancel()
			return
		}
	}
}

// Join handle user join to lobby
func (lobby *Lobby) Join(newConn *Connection) {
	defer utils.CatchPanic("lobby_handle.go Join()")
	// lobby.semJoin <- true
	// defer func() {
	// 	<-lobby.semJoin
	// }()

	lobby.addWaiter(newConn)

	if lobby.recoverInRoom(newConn) {
		lobby.sendPlayerEnter(*newConn, AllExceptThat(newConn))
		return
	}

	lobby.sendWaiterEnter(*newConn, AllExceptThat(newConn))

	newConn.debug("new waiter")
}

// Leave handle user leave lobby
func (lobby *Lobby) Leave(copy *Connection, message string) {
	defer utils.CatchPanic("lobby_handle.go Leave()")
	fmt.Println("disconnected -  #", copy.ID())
	conn := copy

	if conn.InRoom() {
		lobby.Playing.Remove(conn)
		lobby.sendPlayerExit(*conn, AllExceptThat(conn))
	} else {
		lobby.Waiting.Remove(conn)
		lobby.sendWaiterExit(*conn, AllExceptThat(conn))
	}
	if conn.both || conn.InRoom() {
		lobby.LeaveRoom(conn, conn.room, ActionDisconnect)
	}
	return
}

// LeaveRoom handle leave room
func (lobby *Lobby) LeaveRoom(conn *Connection, room *Room, action int) {

	if action != ActionDisconnect {
		lobby.playerToWaiter(conn)
	} else {
		lobby.Playing.Remove(conn)
	}
	fmt.Println("lobby.Playing.Get before", len(room.Players.Connections))
	room.Leave(conn, action) // exit to lobby
	fmt.Println("lobby.Playing.Get agter", len(room.Players.Connections))
	if len(room.Players.Connections) > 0 {
		go func() {
			lobby.sendRoomUpdate(*room, AllExceptThat(conn))
		}()
	}
}

// pickUpRoom find room for player
func (lobby *Lobby) pickUpRoom(conn *Connection, rs *models.RoomSettings) (room *Room) {
	// if there is no room
	if lobby.FreeRooms.Empty() {

		lobby.createAndAddToRoom(rs, conn)
		return
	}
	conn.debug("We have some rooms!")

	// lets find room for him
	for _, room = range lobby.FreeRooms.Get {
		//if room.SameAs()
		if room.addPlayer(conn) {
			return
		}
	}
	return
}

// handleRequest handle any request sent to lobby
func (lobby *Lobby) handleRequest(conn *Connection, lr *LobbyRequest) {
	conn.debug("lobby handle conn")
	lobby.semRequest <- true
	defer func() {
		<-lobby.semRequest
	}()
	conn.debug("sem throw")

	if lr.IsGet() {
		lobby.requestGet(conn, lr)
	} else if lr.IsSend() {
		if lr.Send.RoomSettings == nil {
			conn.debug("lobby cant execute request")
			return
		}
		lobby.EnterRoom(conn, lr.Send.RoomSettings)
	} else if lr.Message != nil {
		Message(lobby, conn, lr.Message, &lobby.Messages,
			lobby.send, All, false, "")
	}
}

// EnterRoom handle user join to room
func (lobby *Lobby) EnterRoom(conn *Connection, rs *models.RoomSettings) {

	if conn.InRoom() {
		lobby.LeaveRoom(conn, conn.room, ActionBackToLobby)
		conn.debug("change room")
	}

	if rs.ID == "create" {
		lobby.createAndAddToRoom(rs, conn)
		return
	}

	if _, room := lobby.AllRooms.SearchRoom(rs.ID); room != nil {
		conn.debug("lobby found required room")
		room.Enter(conn)
	} else {
		conn.debug("lobby search room for you")
		lobby.pickUpRoom(conn, rs)
	}

}

// analize handle where the connection sends the request
func (lobby *Lobby) analize(req *Request) {
	defer utils.CatchPanic("lobby_handle.go analize()")

	if req.Connection.both || !req.Connection.InRoom() {
		var send *LobbyRequest
		if err := json.Unmarshal(req.Message, &send); err != nil {
			req.Connection.SendInformation(err)
		} else {
			lobby.handleRequest(req.Connection, send)
		}
	}
	if req.Connection.both || req.Connection.InRoom() {
		if req.Connection.room == nil {
			return
		}
		var send *RoomRequest
		if err := json.Unmarshal(req.Message, &send); err != nil {
			req.Connection.SendInformation(err)
		} else {
			req.Connection.room.handleRequest(req.Connection, send)
		}
	}
}

// requestGet handle get request to lobby
func (lobby *Lobby) requestGet(conn *Connection, lr *LobbyRequest) {
	lobby.greet(conn)
}
