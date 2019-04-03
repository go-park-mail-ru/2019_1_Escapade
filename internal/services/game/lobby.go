package game

import (
	"escapade/internal/models"
	"sync"
)

// LobbyRequest - client send it by websocket to
// send/get information from Lobby
type LobbyRequest struct {
	Connection *Connection `json:"connection"`
	Send       *LobbySend  `json:"send"`
	Get        *LobbyGet   `json:"get"`
}

// IsGet checks wanna client get info
func (lr *LobbyRequest) IsGet() bool {
	return lr.Get != nil
}

// LobbySend - Information, that client can send to lobby
type LobbySend struct {
	RoomSettings *models.RoomSettings
}

// LobbyGet - Information, that client can get from lobby
type LobbyGet struct {
	AllRooms  bool `json:"allRooms"`
	FreeRooms bool `json:"freeRooms"`
	Waiting   bool `json:"waiting"`
	Playing   bool `json:"playing"`
}

// Lobby there are all rooms and users placed
type Lobby struct {
	AllRooms  *Rooms `json:"allRooms"`
	FreeRooms *Rooms `json:"freeRooms"`

	// room cause they can observe game
	Waiting map[int]*Connection `json:"waiting"`
	Playing map[int]*Connection `json:"playing"`

	// connection joined lobby
	ChanJoin chan *Connection `json:"-"`
	// connection left lobby
	chanLeave   chan *Connection
	chanRequest chan *LobbyRequest

	semJoin    chan bool
	semRequest chan bool
	//chanRoom  chan *Room       // room change status
}

// NewLobby create new instance of Lobby
func NewLobby(roomsCapacity, maxJoin, maxRequest int) *Lobby {

	lobby := &Lobby{
		AllRooms:  NewRooms(roomsCapacity),
		FreeRooms: NewRooms(roomsCapacity),

		Waiting: make(map[int]*Connection),
		Playing: make(map[int]*Connection),

		ChanJoin:    make(chan *Connection),
		chanLeave:   make(chan *Connection),
		chanRequest: make(chan *LobbyRequest),

		semJoin:    make(chan bool, maxJoin),
		semRequest: make(chan bool, maxRequest),
	}
	return lobby
}

// createRoom create room, add to all and free rooms
// and run it
func (lobby *Lobby) createRoom(rs *models.RoomSettings) *Room {

	name := RandString(16)
	room := NewRoom(rs, name, lobby)
	if !lobby.AllRooms.Add(room, name) {
		return nil
	}

	lobby.FreeRooms.Add(room, name)
	go lobby.sendTAILRooms() // inform all about new room
	go room.run()
	return room
}

// Join handle user join to lobby
func (lobby *Lobby) Join(conn *Connection) {
	conn.debug("lobby", "ChanJoin", "Join", "waiting for semJoin")
	lobby.semJoin <- true
	conn.debug("lobby", "ChanJoin", "Join", "taken semJoin")
	defer func() {
		conn.debug("lobby", "ChanJoin", "Join", "free semJoin")
		<-lobby.semJoin
	}()

	thatID := conn.GetPlayerID()
	// maybe user disconnected and we need return him
	for _, room := range lobby.AllRooms.Rooms {
		// work only when game launched, because
		// otherwise player delete from room
		for id, foundConn := range room.Players.Get {
			if id == thatID {
				conn.Status = connectionPlayer
				room.RecoverPlayer(foundConn, conn)
				return
			}
		}
		// if the second account entered as observer
		for id, foundConn := range room.Observers.Get {
			if id == thatID {
				conn.Status = connectionPlayer
				room.RecoverObserver(foundConn, conn)
				return
			}
		}
	}
	// player is new
	conn.Status = connectionLobby
	lobby.sendRooms(conn)
	lobby.addWaiter(conn)
	go lobby.sendTAILPeople()
}

// Leave handle user leave lobby
func (lobby *Lobby) Leave(conn *Connection) {

	lobby.removeWaiter(conn)
	lobby.sendTAILPeople()
	return
}

// ----- handle room status
// roomStart - room remove from free
func (lobby *Lobby) roomStart(room *Room) {
	lobby.FreeRooms.Remove(room)

	go lobby.sendTAILRooms()
}

// roomFinish - room remove from all
func (lobby *Lobby) roomFinish(room *Room) {
	room.Status = StatusFinished
	for _, conn := range room.Players.Get {
		conn.Player.Finished = true
		lobby.playerToWaiter(conn)
	}
	lobby.AllRooms.Remove(room)
	go lobby.sendTAILRooms()
}

// -----

// ----- handle connection status
func (lobby *Lobby) addWaiter(conn *Connection) {
	conn.Status = connectionLobby
	lobby.Playing[conn.GetPlayerID()] = conn
}

func (lobby *Lobby) setWaiterRoom(conn *Connection, room *Room) {
	conn.Status = connectionRoomEnter
	conn.room = room
	lobby.Waiting[conn.GetPlayerID()] = conn
}

func (lobby *Lobby) addPlayer(conn *Connection, room *Room) {
	conn.Status = connectionRoomEnter
	conn.room = room
	lobby.Playing[conn.GetPlayerID()] = conn
}

func (lobby *Lobby) removeWaiter(conn *Connection) {
	delete(lobby.Waiting, conn.GetPlayerID())
}

func (lobby *Lobby) removePlayer(conn *Connection) {
	delete(lobby.Playing, conn.GetPlayerID())
}

func (lobby *Lobby) waiterToPlayer(conn *Connection) {
	lobby.removeWaiter(conn)
	lobby.addPlayer(conn, conn.room)
}

func (lobby *Lobby) playerToWaiter(conn *Connection) {
	lobby.removePlayer(conn)
	lobby.addWaiter(conn) // it send finished player to lobby
}

// -----

func (lobby *Lobby) enterFreeRoom(conn *Connection, rs *models.RoomSettings) (done bool) {
	// if there is no room
	if lobby.FreeRooms.Empty() {
		// if room capacity ended return nil
		conn.debug("enterFreeRoom before", "enterFreeRoom before", "enterFreeRoom before", "enterFreeRoom before")
		room := lobby.createRoom(rs)
		if room != nil {
			room.Players.Add(conn)
		}
		return room != nil
	}

	// lets find room for him
	for _, room := range lobby.FreeRooms.Rooms {
		//if room.SameAs()
		if room.EnterPlayer(conn) {
			conn.debug("enterFreeRoom", "enterFreeRoom", "enterFreeRoom", "enterFreeRoom")
			done = true
			break
		}
	}
	return done
}

func (lobby *Lobby) enterBusyRoom(conn *Connection) bool {

	return conn.room.Join(conn)
}

func (lobby *Lobby) isInvalid(lr *LobbyRequest) bool {
	return lr == nil || lr.Connection == nil || (lr.Get == nil && lr.Send == nil)
}

// handleRequest
func (lobby *Lobby) handleRequest(lr *LobbyRequest) {
	if lobby.isInvalid(lr) {
		return
	}

	lr.Connection.debug("lobby", "ChanRequest", "handleRequest", "waiting semRequest")
	lobby.semRequest <- true
	lr.Connection.debug("lobby", "ChanRequest", "handleRequest", "taken semRequest")
	defer func() {
		lr.Connection.debug("lobby", "ChanRequest", "handleRequest", "free semRequest")
		<-lobby.semRequest
	}()

	if lr.IsGet() {
		lobby.requestGet(lr)
	} else {
		lobby.EnterRoom(lr.Connection, lr.Send.RoomSettings)
	}
}

// EnterRoom handle user join to room
func (lobby *Lobby) EnterRoom(conn *Connection, rs *models.RoomSettings) {

	done := false
	if room, ok := lobby.AllRooms.Rooms[rs.Name]; ok {
		conn.room = room
		done = lobby.enterBusyRoom(conn)
	} else {
		done = lobby.enterFreeRoom(conn, rs)
	}

	if done {
		lobby.waiterToPlayer(conn)
		go lobby.sendTAILPeople()
	} else {
		sendError(conn, "EnterRoom", "cant enter room")
	}
}

// sendRooms send rooms info for user
func (lobby *Lobby) sendRooms(conn *Connection) {
	conn.SendInformation(lobby.AllRooms)
}

// Run the room in goroutine
func (lobby *Lobby) Run() {

	for {
		select {
		case connection := <-lobby.ChanJoin:
			go lobby.Join(connection)

		case request := <-lobby.chanRequest:
			go lobby.handleRequest(request)

		case connection := <-lobby.chanLeave:
			go lobby.Leave(connection)
		}
	}
}

func sendError(conn *Connection, place, message string) {
	res := &models.Result{
		Place:   place,
		Success: false,
		Message: message,
	}
	conn.SendInformation(res)
}

func (lobby *Lobby) sendToAllInLobby(info interface{}) {
	waitJobs := &sync.WaitGroup{}
	for _, conn := range lobby.Waiting {
		conn.debug("sendToAllInLobby", "sendToAllInLobby", "sendToAllInLobby", "sendToAllInLobby")
		if conn.Status == connectionLobby {
			waitJobs.Add(1)
			conn.sendGroupInformation(info, waitJobs)
		}
	}
	waitJobs.Wait()
}

// send to all in lobby
func (lobby *Lobby) sendTAILRooms() {
	get := &LobbyGet{
		AllRooms:  true,
		FreeRooms: true,
	}
	send := lobby.makeGetModel(get)
	lobby.sendToAllInLobby(send)
}

func (lobby *Lobby) sendTAILPeople() {
	get := &LobbyGet{
		Waiting: true,
		Playing: true,
	}
	send := lobby.makeGetModel(get)
	lobby.sendToAllInLobby(send)
}

func (lobby *Lobby) makeGetModel(get *LobbyGet) *Lobby {
	sendLobby := &Lobby{}
	if get.AllRooms {
		sendLobby.AllRooms = lobby.AllRooms
	}
	if get.FreeRooms {
		sendLobby.FreeRooms = lobby.FreeRooms
	}
	if get.Waiting {
		sendLobby.Waiting = lobby.Waiting
	}
	if get.Playing {
		sendLobby.Playing = lobby.Playing
	}
	return sendLobby
}

func (lobby *Lobby) requestGet(lr *LobbyRequest) {
	sendLobby := lobby.makeGetModel(lr.Get)
	lr.Connection.SendInformation(sendLobby)
}
