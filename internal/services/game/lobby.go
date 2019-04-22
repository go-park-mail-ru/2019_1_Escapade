package game

import (
	"context"
	"encoding/json"
	"escapade/internal/config"
	"escapade/internal/models"
	"escapade/internal/utils"
	"fmt"
)

// Request connect Connection and his message
type Request struct {
	Connection *Connection
	Message    []byte
}

// Lobby there are all rooms and users placed
type Lobby struct {
	AllRooms  *Rooms `json:"allRooms,omitempty"`
	FreeRooms *Rooms `json:"freeRooms,omitempty"`

	Waiting *Connections `json:"waiting,omitempty"`
	Playing *Connections `json:"playing,omitempty"`

	Context context.Context

	// connection joined lobby
	ChanJoin chan *Connection `json:"-"`
	// connection left lobby
	chanLeave chan *Connection
	//chanRequest   chan *LobbyRequest
	chanBroadcast chan *Request

	chanBreak chan interface{}

	semJoin    chan bool
	semRequest chan bool
}

// lobby singleton
var (
	lobby *Lobby
)

// Launch launchs lobby goroutine
func Launch(gc *config.GameConfig) {

	if lobby == nil {
		lobby = newLobby(gc.RoomsCapacity,
			gc.LobbyJoin, gc.LobbyRequest)
		go lobby.Run()
	}
}

// GetLobby create lobby if it is nil and get it
func GetLobby() *Lobby {
	return lobby
}

// Stop lobby goroutine
func (lobby *Lobby) Stop() {
	if lobby != nil {
		fmt.Println("Stop called!")
		lobby.chanBreak <- nil
	}
}

// Run the room in goroutine
func (lobby *Lobby) Run() {
	defer func() {
		fmt.Println("defer run!")
		lobby.Free()
		if r := recover(); r != nil {
			fmt.Println("Recovered Run", r)
		}
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

		case connection := <-lobby.chanLeave:
			go lobby.Leave(connection, "You disconnected!")
		case <-lobby.chanBreak:
			fmt.Println("Stop saw!")
			lobbyCancel()
			return
		}
	}
}

// Free delete all rooms and conenctions. Inform all players
// about closing
func (lobby *Lobby) Free() {
	if lobby == nil {
		return
	}
	fmt.Println("All resources clear!")
	SendToConnections("server closed", All(), lobby.Waiting.Get, lobby.Playing.Get)

	lobby.AllRooms.Free()
	lobby.FreeRooms.Free()
	lobby.Waiting.Free()
	lobby.Playing.Free()
	close(lobby.ChanJoin)
	close(lobby.chanLeave)
	close(lobby.chanBroadcast)
	lobby = nil
}

// newLobby create new instance of Lobby
func newLobby(roomsCapacity, maxJoin, maxRequest int) *Lobby {

	connectionsCapacity := 500
	lobby := &Lobby{
		AllRooms:  NewRooms(roomsCapacity),
		FreeRooms: NewRooms(roomsCapacity),

		Waiting: NewConnections(connectionsCapacity),
		Playing: NewConnections(connectionsCapacity),

		ChanJoin:      make(chan *Connection),
		chanLeave:     make(chan *Connection),
		chanBroadcast: make(chan *Request),
		chanBreak:     make(chan interface{}),

		semJoin:    make(chan bool, maxJoin),
		semRequest: make(chan bool, maxRequest),
	}
	return lobby
}

// CloseRoom free room resources
func (lobby *Lobby) CloseRoom(room *Room) {
	// if not in freeRooms nothing bad will happen
	// there is check inside, it will just return without errors
	lobby.FreeRooms.Remove(room)
	lobby.AllRooms.Remove(room)
	lobby.sendTAILRooms()
}

// createRoom create room, add to all and free rooms
// and run it
func (lobby *Lobby) createRoom(rs *models.RoomSettings) *Room {

	name := utils.RandomString(16) // вынести в кофиг
	room := NewRoom(rs, name, lobby)
	if !lobby.AllRooms.Add(room) {
		fmt.Println("cant create room")
		return nil
	}

	lobby.FreeRooms.Add(room)
	lobby.sendTAILRooms() // inform all about new room
	return room
}

func (lobby *Lobby) addWaiter(newConn *Connection) {
	lobby.Waiting.Add(newConn)
	lobby.greet(newConn)
}

func (lobby *Lobby) addPlayer(newConn *Connection, room *Room) {
	lobby.Playing.Add(newConn)
	room.greet(newConn)
}

func (lobby *Lobby) waiterToPlayer(newConn *Connection, room *Room) {
	lobby.Waiting.Remove(newConn)
	lobby.addPlayer(newConn, room)
}

func (lobby *Lobby) playerToWaiter(conn *Connection) {
	lobby.Playing.Remove(conn)
	lobby.addWaiter(conn)
	conn.PushToLobby()
}

func (lobby *Lobby) recoverInLobby(newConn *Connection) bool {
	who := lobby.Waiting.Search(newConn)

	if who >= 0 {
		deleteConn := lobby.Waiting.Get[who]
		deleteConn.Kill("Someone logged into your account")

		return true
	}
	return false
}

func (lobby *Lobby) recoverInRoom(newConn *Connection) bool {
	// find such player
	i, room := lobby.AllRooms.SearchPlayer(newConn)

	if i > 0 {
		room.RecoverPlayer(i, newConn)
		return true
	}

	// find such observer
	old := lobby.AllRooms.SearchObserver(newConn)
	if old != nil {
		old.room.RecoverObserver(old, newConn)
		return true
	}
	return false
}

// Join handle user join to lobby
func (lobby *Lobby) Join(newConn *Connection) {
	// lobby.semJoin <- true
	// defer func() {
	// 	<-lobby.semJoin
	// }()

	if lobby.recoverInLobby(newConn) {
		return
	}

	lobby.addWaiter(newConn)

	if lobby.recoverInRoom(newConn) {
		return
	}

	lobby.sendToWaiters(lobby.Waiting, AllExceptThat(newConn))

	newConn.debug("new waiter")
}

// Leave handle user leave lobby
func (lobby *Lobby) Leave(conn *Connection, message string) {

	fmt.Println("disconnected -  #", conn.ID())

	if conn.both || !conn.InRoom() {
		lobby.Waiting.Remove(conn)
		lobby.sendTAILPeople()
	}
	if conn.both || conn.InRoom() {
		lobby.Playing.Remove(conn)
		if !conn.room.IsActive() {
			conn.room.removeBeforeLaunch(conn)
		}
		conn.room.addAction(conn, ActionDisconnect)
		conn.room.sendHistory(conn.room.all())
	}
	conn.Kill(message)
	return
}

// SendMessage sends message to Connection from lobby
func (lobby *Lobby) SendMessage(conn *Connection, message string) {
	conn.SendInformation([]byte("Lobby message: " + message))
}

// ----- handle room status
// roomStart - room remove from free
func (lobby *Lobby) roomStart(room *Room) {
	lobby.FreeRooms.Remove(room)
	lobby.sendTAILRooms()
}

// roomFinish - room remove from all
func (lobby *Lobby) roomFinish(room *Room) {
	room.Status = StatusFinished
	lobby.AllRooms.Remove(room)
	lobby.sendTAILRooms()
}

// -----

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

// handleRequest
func (lobby *Lobby) handleRequest(conn *Connection, lr *LobbyRequest) {
	conn.debug("lobby handle conn")
	lobby.semRequest <- true
	defer func() {
		<-lobby.semRequest
	}()

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
		lobby.sendToWaiters(lobby, All())
	} else {
		conn.debug("lobby cant execute request")
	}
}

// sendRooms send rooms info for user
func (lobby *Lobby) sendRooms(conn *Connection) {
	bytes, _ := json.Marshal(lobby.AllRooms)
	conn.SendInformation(bytes)
}

func (lobby *Lobby) analize(req *Request) {
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

func (lobby *Lobby) requestGet(conn *Connection, lr *LobbyRequest) {
	sendLobby := lobby.makeGetModel(lr.Get)
	bytes, _ := json.Marshal(sendLobby)
	conn.debug("lobby execute get request")
	conn.SendInformation(bytes)
}
