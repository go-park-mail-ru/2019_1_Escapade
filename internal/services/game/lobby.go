package game

import (
	"escapade/internal/models"
	"sync"
)

//re "escapade/internal/return_errors"
//"math/rand"

type LobbyRequest struct {
	Connection *Connection `json:"connection"`
	Send       *LobbySend  `json:"send"`
	Get        *LobbyGet   `json:"get"`
}

func (lr *LobbyRequest) IsGet() bool {
	return lr.Get != nil
}

type LobbySend struct {
	RoomSettings *models.RoomSettings
}

type LobbyGet struct {
	allRooms  bool `json:"allRooms"`
	freeRooms bool `json:"freeRooms"`
	waiting   bool `json:"waiting"`
	playing   bool `json:"playing"`
}

// user will get Lobby
type Lobby struct {
	allRooms  *Rooms `json:"allRooms"`
	freeRooms *Rooms `json:"freeRooms"`

	// room cause they can observe game
	waiting map[*Connection]*Room `json:"waiting"`
	playing map[*Connection]*Room `json:"playing"`

	// connection joined lobby
	ChanJoin chan *Connection `json:"-"`
	// connection left lobby
	chanLeave   chan *Connection   `json:"-"`
	chanRequest chan *LobbyRequest `json:"-"`

	semJoin    chan bool `json:"-"`
	semRequest chan bool `json:"-"`
	//chanRoom  chan *Room       // room change status
}

func NewLobby(roomsCapacity int) *Lobby {

	// Вынести в конфиг
	maxJoin := 1
	maxRequest := 1
	lobby := &Lobby{
		allRooms:  NewRooms(roomsCapacity),
		freeRooms: NewRooms(roomsCapacity),

		waiting: make(map[*Connection]*Room),
		playing: make(map[*Connection]*Room),

		ChanJoin:    make(chan *Connection),
		chanLeave:   make(chan *Connection),
		chanRequest: make(chan *LobbyRequest),

		semJoin:    make(chan bool, maxJoin),
		semRequest: make(chan bool, maxRequest),
	}
	return lobby
}

func (lobby *Lobby) createRoom(rs *models.RoomSettings) *Room {

	name := RandString(16)
	room := NewRoom(rs, name, lobby)
	if !lobby.allRooms.Add(room, name) {
		return nil
	}
	lobby.freeRooms.Add(room, name)
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

	// maybe user disconnected and we need return him
	for _, room := range lobby.allRooms.Rooms {
		// work only when game launched, because
		// otherwise player delete from room
		for foundConn := range room.players.Get {
			if foundConn.GetPlayerID() == conn.GetPlayerID() {
				conn.Status = connectionPlayer
				room.RecoverPlayer(foundConn, conn)
				return
			}
		}
		// if the second account entered as observer
		for foundConn := range room.observers.Get {
			if foundConn.GetPlayerID() == conn.GetPlayerID() {
				conn.Status = connectionPlayer
				room.RecoverObserver(foundConn, conn)
				return
			}
		}
	}
	// player is new
	conn.Status = connectionLobby
	lobby.sendRooms(conn)
	lobby.waiting[conn] = nil
	go lobby.sendTAILPeople()
}

// Join handle user join to lobby
func (lobby *Lobby) Leave(conn *Connection) {

	lobby.removeWaiter(conn)
	lobby.sendTAILPeople()
	return
}

// ----- handle room status
func (lobby *Lobby) roomStart(room *Room) {
	lobby.freeRooms.Remove(room)
	go lobby.sendTAILRooms()
}

func (lobby *Lobby) roomFinish(room *Room) {
	room.Status = StatusFinished
	for conn := range room.players.Get {
		room.players.Get[conn] = true
		lobby.playerToWaiter(conn)
	}
	lobby.allRooms.Remove(room)
	go lobby.sendTAILRooms()
}

// -----

// ----- handle connection status
func (lobby *Lobby) addWaiter(conn *Connection) {
	conn.Status = connectionLobby
	lobby.playing[conn] = nil
}

func (lobby *Lobby) setWaiterRoom(conn *Connection, room *Room) {
	conn.Status = connectionRoomEnter
	conn.room = room
	lobby.waiting[conn] = room
}

func (lobby *Lobby) addPlayer(conn *Connection, room *Room) {
	conn.Status = connectionRoomEnter
	conn.room = room
	lobby.playing[conn] = room
}

func (lobby *Lobby) removeWaiter(conn *Connection) {
	delete(lobby.waiting, conn)
}

func (lobby *Lobby) removePlayer(conn *Connection) {
	delete(lobby.playing, conn)
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

func (lobby *Lobby) EnterFreeRoom(conn *Connection, rs *models.RoomSettings) (done bool) {
	// if there is no room
	if lobby.freeRooms.Empty() {
		// if room capacity ended return nil
		room := lobby.createRoom(rs)
		return room != nil
	}

	// lets find room for him
	for _, room := range lobby.freeRooms.Rooms {
		//if room.SameAs()
		if room.EnterPlayer(conn) {
			done = true
			break
		}
	}
	return done
}

// EnterBusyRoom try connect as observer
func (lobby *Lobby) EnterBusyRoom(conn *Connection) bool {

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
	if room, ok := lobby.allRooms.Rooms[rs.Name]; ok {
		conn.room = room
		done = lobby.EnterBusyRoom(conn)
	} else {
		done = lobby.EnterFreeRoom(conn, rs)
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
	conn.SendInformation(lobby.allRooms)
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
	for conn := range lobby.waiting {
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
		allRooms:  true,
		freeRooms: true,
	}
	send := lobby.makeGetModel(get)
	lobby.sendToAllInLobby(send)
}

func (lobby *Lobby) sendTAILPeople() {
	get := &LobbyGet{
		waiting: true,
		playing: true,
	}
	send := lobby.makeGetModel(get)
	lobby.sendToAllInLobby(send)
}

func (lobby *Lobby) makeGetModel(get *LobbyGet) *Lobby {
	sendLobby := &Lobby{}
	if get.allRooms {
		sendLobby.allRooms = lobby.allRooms
	}
	if get.freeRooms {
		sendLobby.freeRooms = lobby.freeRooms
	}
	if get.waiting {
		sendLobby.waiting = lobby.waiting
	}
	if get.playing {
		sendLobby.playing = lobby.playing
	}
	return sendLobby
}

func (lobby *Lobby) requestGet(lr *LobbyRequest) {
	sendLobby := lobby.makeGetModel(lr.Get)
	lr.Connection.SendInformation(sendLobby)
}
