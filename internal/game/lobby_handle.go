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

// Run the room in goroutine
func (lobby *Lobby) Run() {
	defer func() {
		utils.CatchPanic("lobby_handle.go Run()")
		lobby.Free()
	}()

	fmt.Println("Lobby run")
	for {
		select {
		case connection := <-lobby.chanJoin:
			go lobby.Join(connection, false)

		case message := <-lobby.chanBroadcast:
			go lobby.Analize(message)
		}
	}
}

// Join handle user join to lobby
func (lobby *Lobby) Join(newConn *Connection, disconnected bool) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	fmt.Println("join")
	defer func() {
		utils.CatchPanic("lobby_handle.go Join()")
		lobby.wGroup.Done()
	}()

	lobby.addWaiter(newConn)

	if !lobby.canCloseRooms && lobby.recoverInRoom(newConn, disconnected) {
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

	fmt.Println("LEAVE!")

	var disconnected bool

	if conn.Both() || conn.InRoom() {
		fmt.Println("delete from room")
		fmt.Println("both -  #", conn.ID())
		disconnected = lobby.LeaveRoom(conn, conn.Room(), ActionDisconnect)
	}

	if !conn.InRoom() {
		fmt.Println("delete from lobby")
		disconnected = lobby.Waiting.Remove(conn, true) //lobby.waitingRemove(conn)
		if disconnected {
			lobby.sendWaiterExit(*conn, AllExceptThat(conn))
		}
	}

	if disconnected {
		fmt.Println("disconnected -  #", conn.ID())
	}
	return
}

// LeaveRoom handle leave room
func (lobby *Lobby) LeaveRoom(conn *Connection, room *Room, action int) (done bool) {

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
		lobby.PlayerToWaiter(conn)
		done = room.Leave(conn, action) // exit to lobby
	} else {
		//go lobby.playingRemove(conn)
		found := room.Search(conn)
		if found != nil {
			found.setDisconnected()
		}
		go lobby.sendPlayerExit(*conn, AllExceptThat(conn))
	}
	if done && len(room.Players.Connections.RGet()) > 0 {
		lobby.sendRoomUpdate(*room, AllExceptThat(conn))
	}
	return done
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

	fmt.Println("see block!")
	conn.actionSem <- struct{}{}
	defer func() { <-conn.actionSem }()

	fmt.Println("EnterRoom", rs)
	if conn.InRoom() {
		fmt.Println("in room", rs)
		fmt.Println("EnterRoom ID compare", conn.Room().ID, rs.ID, rs)
		if conn.Room().ID == rs.ID {
			return
		}
		lobby.LeaveRoom(conn, conn.Room(), ActionBackToLobby)
		conn.debug("change room")
	}

	fmt.Println("not in room", rs.ID, rs.ID == "create", conn.ID())
	conn.debug("enter room" + rs.ID)
	if rs.ID == "create" {
		fmt.Println("you wanna crete room", rs)
		conn.debug("see you wanna create room?")
		lobby.CreateAndAddToRoom(rs, conn)
		return
	}

	if _, room := lobby.allRoomsSearch(rs.ID); room != nil {
		conn.debug("lobby found required room")
		room.Enter(conn)
	} else {
		conn.debug("lobby search room for you")
		lobby.PickUpRoom(conn, rs)
	}
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
		if room.addPlayer(conn) {
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

	fmt.Println("analyze", req.Connection.Both(), req.Connection.InRoom(), req.Connection.Disconnected())
	if req.Connection.Both() || !req.Connection.InRoom() {
		fmt.Println("lobby work")
		var send *LobbyRequest
		if err := json.Unmarshal(req.Message, &send); err != nil {
			req.Connection.SendInformation(err)
		} else {
			go lobby.HandleRequest(req.Connection, send)
		}
	}
	if req.Connection.Both() || req.Connection.InRoom() {
		fmt.Println("room work")
		if req.Connection.Room() == nil {
			fmt.Println("bot room")
			return
		}
		var send *RoomRequest
		if err := json.Unmarshal(req.Message, &send); err != nil {
			req.Connection.SendInformation(err)
		} else {
			fmt.Println("room is here")
			room := req.Connection.Room()
			if room == nil {
				fmt.Println("no room")
				return
			}
			room.HandleRequest(req.Connection, send)
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
			lobby.EnterRoom(conn, lr.Send.RoomSettings)
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
	inv.Message.Time = time.Now()
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
