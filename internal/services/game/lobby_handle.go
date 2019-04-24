package game

import (
	"context"
	"encoding/json"
	"escapade/internal/models"
	"escapade/internal/utils"
	"fmt"
)

// Run the room in goroutine
func (lobby *Lobby) Run() {
	defer func() {
		utils.CatchPanic("lobby_handle.go Run()")
		lobby.Free()
	}()

	var lobbyCancel context.CancelFunc
	lobby.Context, lobbyCancel = context.WithCancel(context.Background())
	fmt.Println("create context!")
	for {
		select {
		case connection := <-lobby.ChanJoin:
			go lobby.Join(connection)

		case message := <-lobby.chanBroadcast:
			go lobby.analize(message)

			// TODO delete chanleavem cause Leave call direcrly
		case connection := <-lobby.chanLeave:
			go lobby.Leave(connection, "You disconnected!")
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
		lobby.send(lobby.Playing, AllExceptThat(newConn))
		return
	}

	lobby.sendWaiting(AllExceptThat(newConn))

	newConn.debug("new waiter")
}

// Leave handle user leave lobby
func (lobby *Lobby) Leave(copy *Connection, message string) {
	defer utils.CatchPanic("lobby_handle.go Leave()")
	fmt.Println("disconnected -  #", copy.ID())
	conn := copy

	if conn.both || !conn.InRoom() {
		fmt.Println("lobby delete ", conn.ID())
		lobby.Waiting.Remove(conn)
		lobby.sendWaiting(AllExceptThat(conn))
	}
	if conn.both || conn.InRoom() {
		fmt.Println("room delete ", conn.ID())
		_, room := lobby.AllRooms.SearchPlayer(conn)
		fmt.Println("room id ", room.Name)
		lobby.Playing.Remove(conn)
		if !room.IsActive() {
			fmt.Println("removeBeforeLaunch")
			go room.removeBeforeLaunch(conn)
		} else {
			go room.removeDuringGame(conn)
		}
		go func() {
			room.addAction(conn, ActionDisconnect)
			room.sendHistory(conn.room.All)
		}()
	}
	return
}

// pickUpRoom find room for player
func (lobby *Lobby) pickUpRoom(conn *Connection, rs *models.RoomSettings) (done bool) {
	// if there is no room
	if lobby.FreeRooms.Empty() {
		// if room capacity ended return nil
		room := lobby.createRoom(rs)
		if room != nil {
			conn.debug("We create your own room, cool!")
			room.addPlayer(conn)
		} else {
			conn.debug("cant create. Why?")
		}
		return room != nil
	}
	conn.debug("We have some rooms!")

	// lets find room for him
	for _, room := range lobby.FreeRooms.Get {
		//if room.SameAs()
		if room.addPlayer(conn) {
			done = true
			break
		}
	}
	return done
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
	}
}

// EnterRoom handle user join to room
func (lobby *Lobby) EnterRoom(conn *Connection, rs *models.RoomSettings) {

	var done bool
	if conn.InRoom() {
		conn.debug("lobby cant execute request")
		return
	}
	if rs.Name == "create" {
		conn.debug("try create")
		room := lobby.createRoom(rs)
		done = room != nil
		if done {
			room.addPlayer(conn)
		}
	} else {
		if _, room := lobby.AllRooms.SearchRoom(rs.Name); room != nil {
			conn.debug("lobby found required room")
			done = room.Enter(conn)
		} else {
			conn.debug("lobby search room for you")
			done = lobby.pickUpRoom(conn, rs)
		}
	}

	if done {
		conn.debug("lobby done")
		lobby.send(lobby, All)
	} else {
		conn.debug("lobby cant execute request")
	}
}

// analize handle where the connection sends the request
func (lobby *Lobby) analize(req *Request) {
	defer utils.CatchPanic("lobby_handle.go analize()")

	if req.Connection.both || !req.Connection.InRoom() {
		var send *LobbyRequest
		if err := json.Unmarshal(req.Message, &send); err != nil {
			bytes, _ := json.Marshal(err)
			req.Connection.SendInformation(bytes)
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
			bytes, _ := json.Marshal(err)
			req.Connection.SendInformation(bytes)
		} else {
			req.Connection.room.handleRequest(req.Connection, send)
		}
	}
}

// requestGet handle get request to lobby
func (lobby *Lobby) requestGet(conn *Connection, lr *LobbyRequest) {
	sendLobby := lobby.makeGetModel(lr.Get)
	bytes, _ := json.Marshal(sendLobby)
	conn.debug("lobby execute get request")
	conn.SendInformation(bytes)
}
