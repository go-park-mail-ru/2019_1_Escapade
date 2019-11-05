package engine

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
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
	lobby.chanJoin <- conn

	go lobby.mUserWelcome(conn.IsAnonymous())
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
		utils.Debug(false, "lobby.restore", conn.ID(), conn.PlayingRoom(), conn.WaitingRoom(), conn.Index())
		if conn.PlayingRoom() == nil {
			lobby.greet(conn)
		}
		return
	}

	lobby.addWaiter(conn)

	utils.Debug(false, "new waiter")
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

	// check
	_, waiter := lobby.Waiting.SearchByID(conn.ID())
	if waiter != nil {
		if waiter.UUID != conn.UUID {
			return
		}
	} else {
		_, player := lobby.Playing.SearchByID(conn.ID())
		if player != nil {
			if player.UUID != conn.UUID {
				return
			}
		}
	}
	if conn.PlayingRoom() != nil {
		conn.PlayingRoom().connEvents.Disconnect(conn)
		// dont delete from lobby, because player not in lobby
		return
	} else if conn.WaitingRoom() != nil {
		conn.PlayingRoom().connEvents.Disconnect(conn)
		// continue, because player in lobby
	}

	utils.Debug(false, "see disc -  #", conn.ID())
	disconnected := lobby.Waiting.Remove(conn)
	if disconnected {
		go lobby.sendWaiterExit(conn, All)
		go lobby.mUserBye(conn.IsAnonymous())
		utils.Debug(false, "disconnected -  #", conn.ID())
	}

	return
}

// LeaveRoom handle leave room
func (lobby *Lobby) LeaveRoom(conn *Connection, action int) {
	defer utils.CatchPanic("lobby_handle.go LeaveRoom()")

	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	if action != ActionDisconnect {
		if conn.PlayingRoom() != nil {
			lobby.PlayerToWaiter(conn)
		} else {
			conn.PushToLobby()
		}
	}
	lobby.greet(conn)

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

	if !rs.AnonymousCheck(conn.ID() < 0) {
		utils.Debug(false, "anonymous check fault")
		return
	}

	conn.actionSem <- struct{}{}
	defer func() { <-conn.actionSem }()

	if conn.WaitingRoom() != nil {
		if conn.WaitingRoom().info.ID() == rs.ID {
			return
		}
		conn.WaitingRoom().connEvents.Leave(conn)
	}

	if rs.ID == "create" {
		utils.Debug(false, "see you wanna create room?", rs)
		lobby.PickUpRoom(conn, rs)
		//lobby.CreateAndAddToRoom(rs, conn)
		return
	}

	if room := lobby.allRooms.Search(rs.ID); room != nil {
		utils.Debug(false, "lobby found required room")
		room.connEvents.Enter(conn)
	} else {
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

	freeRoomsIterator := NewRoomsIterator(lobby.freeRooms)
	for freeRoomsIterator.Next() {
		freeRoom := freeRoomsIterator.Value()
		if freeRoom.info.Settings().Similar(rs) && freeRoom.people.add(conn, true, false) {
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
		var rsend RoomRequest
		if err := rsend.UnmarshalJSON(req.Message); err != nil {
			utils.Debug(true, "json error")
		} else {
			req.Connection.PlayingRoom().api.Handle(req.Connection, &rsend)
		}
		// not in lobby
		return
	} else if req.Connection.WaitingRoom() != nil {
		var rsend RoomRequest
		if err := rsend.UnmarshalJSON(req.Message); err != nil {
			utils.Debug(true, "json error")
		} else {
			req.Connection.WaitingRoom().api.Handle(req.Connection, &rsend)
		}
		// in lobby
	}

	var send LobbyRequest
	if err := send.UnmarshalJSON(req.Message); err != nil {
		req.Connection.SendInformation(&models.Result{
			Success: false,
			Message: err.Error(),
		})
	} else {
		lobby.HandleRequest(req.Connection, &send)
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
			lobby.send, All, nil, lobby.dbChatID())
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
	inv.Message.Time = time.Now().In(lobby.location())
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
