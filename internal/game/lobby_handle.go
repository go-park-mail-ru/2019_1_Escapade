package game

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"encoding/json"
	"fmt"
)

// JoinConn is the wrapper in order to put the connection in the channel chanJoin
func (lobby *Lobby) JoinConn(conn *Connection, d int) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go Join()")
		lobby.wGroup.Done()
	}()
	fmt.Println("call me")
	lobby.chanJoin <- conn
}

/*
А зачем нам рум гарбаж коллекток, если можно ходить по слайсу playing?
*/
func (lobby *Lobby) launchGarbageCollector(timeout float64) {
	fmt.Println("lobby launchGarbageCollector")
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go Join()")
		lobby.wGroup.Done()
	}()

	it := NewConnectionsIterator(lobby.Waiting)
	for it.Next() {
		waiter := it.Value()
		if waiter == nil {
			panic("why nill")
		}
		t := waiter.Time()
		if time.Since(t).Seconds() > timeout {
			fmt.Println(waiter.User.Name, " - bad")
			lobby.Leave(waiter, "")
		} else {
			fmt.Println(waiter.User.Name, " - good", waiter.Disconnected(), time.Since(t).Seconds())
		}
		// if waiter.isClosed() {
		// 	lobby.Leave(waiter, "")
		// }
	}
	/*
		for _, conn := range lobby.Waiting.RGet() {
			if conn == nil {
				continue
			}
			if conn.isClosed() {
				lobby.Leave(conn, "")
			}
			// if time.Since(conn.time).Seconds() > timeout {
			// 	fmt.Println(conn.User.Name, " - bad")
			// 	lobby.Leave(conn, "")
			// } else {
			// 	fmt.Println(conn.User.Name, " - good", conn.Disconnected(), time.Since(conn.time).Seconds())
			// }
		}
	*/
}

/*
Run accepts connections and messages from them
Goroutine. When it is finished, the lobby will be cleared
*/
func (lobby *Lobby) Run() {
	defer func() {
		utils.CatchPanic("lobby_handle.go Run()")
		lobby.Free()
	}()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	var timeout float64
	timeout = 10

	for {
		select {
		case <-ticker.C:
			go lobby.launchGarbageCollector(timeout)
		case connection := <-lobby.chanJoin:
			go lobby.Join(connection)
		case message := <-lobby.chanBroadcast:
			lobby.Analize(message)
		case <-lobby.chanBreak:
			return
		}
	}
}

// Join handle user join to lobby
func (lobby *Lobby) Join(conn *Connection) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go Join()")
		lobby.wGroup.Done()
	}()

	// try restore user
	if lobby.restore(conn) {
		fmt.Println("lobby.restore", conn.ID(), conn.PlayingRoom(), conn.WaitingRoom(), conn.Index())
		lobby.greet(conn)
		return
	}

	lobby.addWaiter(conn)

	conn.debug("new waiter")
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

	var disconnected bool

	// check
	waiter, _ := lobby.Waiting.SearchByID(conn.ID())
	if waiter != nil {
		if waiter.UUID != conn.UUID {
			return
		} /*else {
			fmt.Println("waiter closed ws:", waiter.isClosed())
		}*/
		// err := waiter.ws.Close()
		// if err != nil {
		// 	fmt.Println("cant leave:", err.Error())
		// 	return
		// }
	} else {
		player, _ := lobby.Playing.SearchByID(conn.ID())
		if player != nil {
			if player.UUID != conn.UUID {
				return
			}
			/*
				if !player.IsConnected() {
					return
				} else {
					fmt.Println("player closed ws:", player.isClosed())
				}*/
			// err := player.ws.Close()
			// if err != nil {
			// 	fmt.Println("cant leave:", err.Error())
			// 	return
			// }
		}
	}
	//

	if conn.PlayingRoom() != nil {
		if !conn.PlayingRoom().done() {
			conn.PlayingRoom().chanConnection <- ConnectionAction{
				conn:   conn,
				action: ActionDisconnect,
			}
		}
		// dont delete from lobby, because player not in lobby
		return
	} else if conn.WaitingRoom() != nil {
		if !conn.WaitingRoom().done() {
			conn.WaitingRoom().chanConnection <- ConnectionAction{
				conn:   conn,
				action: ActionDisconnect,
			}
		}
		// continue, because player in lobby
	}

	fmt.Println("see disc -  #", conn.ID())
	disconnected = lobby.Waiting.Remove(conn)
	if disconnected {
		go lobby.sendWaiterExit(*conn, All)
	}

	if disconnected {
		fmt.Println("disconnected -  #", conn.ID())
	}
	return
}

// LeaveRoom handle leave room
func (lobby *Lobby) LeaveRoom(conn *Connection, action int, room *Room) {

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
		if room.Status != StatusPeopleFinding {
			lobby.PlayerToWaiter(conn)
		}
	} else if room.Players.Connections.len() > 0 {
		lobby.sendRoomUpdate(*room, AllExceptThat(conn))
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

	conn.actionSem <- struct{}{}
	defer func() { <-conn.actionSem }()

	if conn.WaitingRoom() != nil {
		if conn.WaitingRoom().ID == rs.ID {
			return
		}
		conn.WaitingRoom().processActionBackToLobby(conn)
	}

	if rs.ID == "create" {
		fmt.Println("you wanna crete room", rs)
		conn.debug("see you wanna create room?")
		lobby.CreateAndAddToRoom(rs, conn)
		return
	}

	if _, room := lobby.allRoomsSearch(rs.ID); room != nil {
		conn.debug("lobby found required room")
		room.Enter(conn)
	}

	// panic there sometimes... rarely
	// } else {
	// 	conn.debug("lobby search room for you")
	// 	lobby.PickUpRoom(conn, rs)
	// }
}

// PickUpRoom find room for player
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
		if room.addConnection(conn, true, false) { //room.addPlayer(conn, false) {
			return
		}
	}
	// oh we cant find room, so lets create one
	lobby.CreateAndAddToRoom(rs, conn)
	return
}

// Analize handle where the connection sends the request
func (lobby *Lobby) Analize(req *Request) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go Analize()")
		lobby.wGroup.Done()
	}()

	if req.Connection.PlayingRoom() != nil {
		var rsend *RoomRequest
		if err := json.Unmarshal(req.Message, &rsend); err != nil {
			fmt.Println("big json error")
		} else {
			req.Connection.PlayingRoom().HandleRequest(req.Connection, rsend)
		}
		return
		// not in lobby
		return
	} else if req.Connection.WaitingRoom() != nil {
		var rsend *RoomRequest
		if err := json.Unmarshal(req.Message, &rsend); err != nil {
			fmt.Println("big json error")
		} else {
			req.Connection.WaitingRoom().HandleRequest(req.Connection, rsend)
		}
		// in lobby
	}

	var send *LobbyRequest
	if err := json.Unmarshal(req.Message, &send); err != nil {
		req.Connection.SendInformation(err)
	} else {
		lobby.HandleRequest(req.Connection, send)
	}
}

// HandleRequest handle any request sent to lobby
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
		//go lobby.greet(conn)
	} else if lr.IsSend() {
		switch {
		case lr.Send.Messages != nil:
			Messages(conn, lr.Send.Messages, lobby.Messages())
		case lr.Send.RoomSettings != nil:
			go lobby.EnterRoom(conn, lr.Send.RoomSettings)
		}
	} else if lr.Message != nil {
		Message(lobby, conn, lr.Message,
			lobby.appendMessage, lobby.setMessage,
			lobby.removeMessage, lobby.findMessage,
			lobby.send, All, false, "")
	}
}

// Invite invite people to your room
func (lobby *Lobby) Invite(conn *Connection, inv *Invitation) {

	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go EnterRoom()")
		lobby.wGroup.Done()
	}()

	inv.From = conn.User
	inv.Message.User = conn.User
	loc, _ := time.LoadLocation("Europe/Moscow")
	inv.Message.Time = time.Now().In(loc)
	if inv.All {
		lobby.sendInvitation(inv, All)
		lobby.sendInvitationCallback(conn, nil)
	} else {

		var find *Connection
		it := NewConnectionsIterator(lobby.Waiting)
		for it.Next() {
			waiter := it.Value()
			if waiter.User.Name == inv.To {
				find = conn
				break
			}
		}
		if find != nil {
			lobby.sendInvitation(inv, Me(find))
			lobby.sendInvitationCallback(conn, nil)
		} else {
			lobby.sendInvitationCallback(conn, re.ErrorUserNotFound())
		}
	}

}
