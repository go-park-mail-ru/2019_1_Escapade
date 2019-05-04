package game

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"encoding/json"
	"fmt"
)

// Run the room in goroutine
func (lobby *Lobby) Run(wg *sync.WaitGroup) {
	defer func() {
		utils.CatchPanic("lobby_handle.go Run()")
		lobby.Free()
	}()

	//var lobbyCancel context.CancelFunc
	//lobby.Context, lobbyCancel = context.WithCancel(context.Background())
	if wg != nil {
		wg.Done()
	}
	fmt.Println("Lobby run")
	for {
		select {
		case connection := <-lobby.ChanJoin:
			go lobby.Join(connection)

		case message := <-lobby.chanBroadcast:
			go lobby.Analize(message)

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
			return
		}
	}
}

// Join handle user join to lobby
func (lobby *Lobby) Join(newConn *Connection) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go Join()")
		lobby.wGroup.Done()
	}()

	lobby.addWaiter(newConn)

	if lobby.recoverInRoom(newConn) {
		go lobby.sendPlayerEnter(*newConn, AllExceptThat(newConn))
		return
	}

	go lobby.sendWaiterEnter(*newConn, AllExceptThat(newConn))

	newConn.debug("new waiter")
}

// Leave handle user leave lobby
func (lobby *Lobby) Leave(conn *Connection, message string) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go Leave()")
		lobby.wGroup.Done()
	}()

	fmt.Println("disconnected -  #", conn.ID())

	if !conn.InRoom() {
		go lobby.waitingRemove(conn)
		go lobby.sendWaiterExit(*conn, AllExceptThat(conn))
	}

	fmt.Println("here ", conn.both, conn.InRoom())
	if conn.both || conn.InRoom() {
		fmt.Println("both -  #", conn.ID())
		go lobby.LeaveRoom(conn, conn.room, ActionDisconnect)
	}
	return
}

// LeaveRoom handle leave room
func (lobby *Lobby) LeaveRoom(conn *Connection, room *Room, action int) {

	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go LeaveRoom()")
		lobby.wGroup.Done()
	}()

	fmt.Println("check", action, ActionDisconnect)
	if action != ActionDisconnect {
		go lobby.PlayerToWaiter(conn)
	} else {
		//go lobby.playingRemove(conn)
		go lobby.sendPlayerExit(*conn, AllExceptThat(conn))
	}

	fmt.Println("lobby.Playing.Get before", len(room.Players.Connections))
	room.Leave(conn, action) // exit to lobby
	fmt.Println("lobby.Playing.Get after", len(room.Players.Connections))
	if len(room.Players.Connections) > 0 {
		go lobby.sendRoomUpdate(*room, AllExceptThat(conn))
	}
}

// EnterRoom handle user join to room
func (lobby *Lobby) EnterRoom(conn *Connection, rs *models.RoomSettings) {

	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go EnterRoom()")
		lobby.wGroup.Done()
	}()

	fmt.Println("EnterRoom")
	if conn.InRoom() {
		fmt.Println("EnterRoom ID compare", conn.room.ID, rs.ID, rs)
		if conn.room.ID == rs.ID {
			return
		}
		go lobby.LeaveRoom(conn, conn.room, ActionBackToLobby)
		conn.debug("change room")
	}

	conn.debug("enter room" + rs.ID)
	if rs.ID == "create" {
		conn.debug("see you wanna create room?")
		go lobby.CreateAndAddToRoom(rs, conn)
		return
	}

	if _, room := lobby.allRoomsSearch(rs.ID); room != nil {
		conn.debug("lobby found required room")
		room.Enter(conn)
	} else {
		conn.debug("lobby search room for you")
		go lobby.PickUpRoom(conn, rs)
	}

}

// pickUpRoom find room for player
func (lobby *Lobby) PickUpRoom(conn *Connection, rs *models.RoomSettings) (room *Room) {

	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go PickUpRoom()")
		lobby.wGroup.Done()
	}()

	// lets find room for user
	FreeRooms := lobby.freeRooms()
	for _, room = range FreeRooms {
		//if room.SameAs()
		if room.addPlayer(conn) {
			return
		}
	}
	// oh we cant find room, so lets create one
	go lobby.CreateAndAddToRoom(rs, conn)
	return
}

// analize handle where the connection sends the request
func (lobby *Lobby) Analize(req *Request) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go Analize()")
		lobby.wGroup.Done()
	}()

	if req.Connection.both || !req.Connection.InRoom() {
		var send *LobbyRequest
		if err := json.Unmarshal(req.Message, &send); err != nil {
			req.Connection.SendInformation(err)
		} else {
			go lobby.HandleRequest(req.Connection, send)
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

// handleRequest handle any request sent to lobby
func (lobby *Lobby) HandleRequest(conn *Connection, lr *LobbyRequest) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go HandleRequest()")
		lobby.wGroup.Done()
	}()

	if lr.IsGet() {
		go lobby.greet(conn)
	} else if lr.IsSend() {
		if lr.Send.RoomSettings == nil {
			conn.debug("lobby cant execute request")
			return
		}
		lobby.EnterRoom(conn, lr.Send.RoomSettings)
	} else if lr.Message != nil {
		Message(lobby, conn, lr.Message, lobby.setToMessages,
			lobby.send, All, false, "")
	}
}
