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

	for _, conn := range lobby.Waiting.RGet() {
		if conn == nil {
			continue
		}
		if time.Since(conn.time).Seconds() > timeout {
			fmt.Println(conn.User.Name, " - bad")
			lobby.Leave(conn, "")
		} else {
			fmt.Println(conn.User.Name, " - good", conn.Disconnected(), time.Since(conn.time).Seconds())
		}
	}
}

// Run the room in goroutine
func (lobby *Lobby) Run() {

	ticker := time.NewTicker(time.Second * 10)
	//var timeout float64
	//timeout = 10

	defer func() {
		utils.CatchPanic("lobby_handle.go Run()")
		ticker.Stop()
		lobby.Free()
	}()

	for {
		select {
		//case <-ticker.C:
		//	go lobby.launchGarbageCollector(timeout)
		case connection := <-lobby.chanJoin:
			go lobby.Join(connection, false)

		case message := <-lobby.chanBroadcast:
			lobby.Analize(message)
		}
	}
}

// Join handle user join to lobby
func (lobby *Lobby) Join(newConn *Connection, disconnected bool) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go Join()")
		lobby.wGroup.Done()
	}()

	found, _ := lobby.Waiting.SearchByID(newConn.ID())
	if found != nil {
		sendAccountTaken(*found)
	}
	//sendAccountTaken(*oldConn)

	lobby.addWaiter(newConn)

	if lobby.canCloseRooms {

		lobby.recoverInRoom(newConn, disconnected)
		//go lobby.sendPlayerEnter(*newConn, AllExceptThat(newConn))
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

	var disconnected bool

	if conn.Both() || conn.InRoom() {
		if !conn.Room().done() {
			conn.Room().chanConnection <- ConnectionAction{
				conn:   conn,
				action: ActionDisconnect,
			}
		}
		//disconnected = lobby.LeaveRoom(conn, conn.Room(), ActionDisconnect)
	}

	if !conn.InRoom() {
		timeout := time.Duration(time.Second) * 20
		//fmt.Println("compate", time.Since(conn.time).Seconds(), timeout.Seconds())
		if time.Since(conn.time).Seconds() > timeout.Seconds() {
			disconnected = lobby.Waiting.FastRemove(conn) //lobby.waitingRemove(conn)
			if disconnected {
				lobby.sendWaiterExit(*conn, All)
			}
		}
	}

	if disconnected {
		fmt.Println("disconnected -  #", conn.ID())
	}
	return
}

// LeaveRoom handle leave room
func (lobby *Lobby) LeaveRoom(conn *Connection, action int) {

	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go LeaveRoom()")
		lobby.wGroup.Done()
	}()

	if conn.Room() == nil {
		return
	}

	fmt.Println("check", action, ActionDisconnect)
	if action != ActionDisconnect {
		lobby.PlayerToWaiter(conn)
	} else if len(conn.Room().Players.Connections.RGet()) > 0 {
		lobby.sendRoomUpdate(*conn.Room(), AllExceptThat(conn))
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

	if conn.InRoom() {
		if conn.Room().ID == rs.ID {
			return
		}
		conn.Room().processActionBackToLobby(conn)
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
		if room.addPlayer(conn, false) {
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

	//fmt.Println("analyze", req.Connection.Both(), req.Connection.InRoom(), req.Connection.Disconnected())
	if req.Connection.Both() || !req.Connection.InRoom() {
		fmt.Println("lobby work")
		var send *LobbyRequest
		if err := json.Unmarshal(req.Message, &send); err != nil {
			req.Connection.SendInformation(err)
		} else {
			lobby.HandleRequest(req.Connection, send)
		}
	}
	fmt.Println("req.Connection.InRoom()", req.Connection.InRoom())
	if req.Connection.Both() || req.Connection.InRoom() {
		var rsend *RoomRequest
		if err := json.Unmarshal(req.Message, &rsend); err != nil {
			fmt.Println("big json error")
		} else {
			fmt.Println("room is here")
			room := req.Connection.Room()
			if room == nil {
				fmt.Println("no room")
				return
			}
			room.HandleRequest(req.Connection, rsend)
		}
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
		waiting := lobby.Waiting.RGet()
		var find *Connection
		for _, conn := range waiting {
			if conn.User.Name == inv.To {
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
