package game

import "escapade/internal/models"

//re "escapade/internal/return_errors"
//"math/rand"

type Lobby struct {
	allRoomsCapacity int

	allRoomsSize  int
	allRooms      map[int]*Room
	freeRoomsSize int
	freeRooms     map[int]*Room

	waiting map[*Connection]*Room // room cause they can observe game
	playing map[*Connection]*Room

	ChanJoin    chan *Connection // connection joined lobby
	chanLeave   chan *Connection // connection left lobby
	chanRequest chan *Request
	//chanRoom  chan *Room       // room change status
}

func NewLobby() *Lobby {

	lobby := &Lobby{
		allRoomsCapacity: 500, // вынести в конфиг
		allRoomsSize:     0,
		allRooms:         make(map[int]*Room),

		freeRoomsSize: 500, // вынести в конфиг
		freeRooms:     make(map[int]*Room),

		waiting: make(map[*Connection]*Room),
		playing: make(map[*Connection]*Room),

		ChanJoin:    make(chan *Connection),
		chanLeave:   make(chan *Connection),
		chanRequest: make(chan *Request),
	}
	return lobby
}

func (lobby *Lobby) createRoomID() (id int) {

	id = 0
	for _, ok := lobby.allRooms[id]; ok; {
		id++
	}
	return
}

func (lobby *Lobby) createRoom(rs *models.RoomSettings) *Room {
	if lobby.allRoomsSize == lobby.allRoomsCapacity {
		return nil
	}

	id := lobby.createRoomID()
	room := NewRoom(rs, id, lobby)
	lobby.addFreeRoom(room)
	lobby.addRoom(room)
	go room.run()
	return room
}

// Join handle user join to lobby
func (lobby *Lobby) Join(conn *Connection) {

	// maybe user disconnected and we need return him
	for _, room := range lobby.allRooms {
		// work only when game launched, because
		// otherwise player delete from room
		for foundConn := range room.Players {
			if foundConn.GetPlayerID() == conn.GetPlayerID() {
				conn.Status = connectionPlayer
				room.RecoverPlayer(foundConn, conn)
				return
			}
		}
		// if the second account entered as observer
		for foundConn := range room.Observers {
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

	return
}

// Join handle user join to lobby
func (lobby *Lobby) Leave(conn *Connection) {

	lobby.removeWaiter(conn)
	return
}

// ----- handle allRooms and freeRooms
func (lobby *Lobby) addFreeRoom(room *Room) {
	lobby.freeRooms[room.ID] = room
	lobby.freeRoomsSize++
}

func (lobby *Lobby) addRoom(room *Room) {
	lobby.allRooms[room.ID] = room
	lobby.allRoomsSize++
}

func (lobby *Lobby) removeFreeRoom(room *Room) {
	delete(lobby.freeRooms, room.ID)
}

func (lobby *Lobby) removeRoom(room *Room) {
	delete(lobby.allRooms, room.ID)
}

// -----

// ----- handle room status
func (lobby *Lobby) roomStart(room *Room) {
	lobby.removeFreeRoom(room)
}

func (lobby *Lobby) roomFinish(room *Room) {
	room.Status = StatusFinished
	for conn, playing := range room.Players {
		playing.Finished = true
		lobby.playerToWaiter(conn)
	}
	lobby.removeRoom(room)
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

func (lobby *Lobby) EnterFreeRoom(conn *Connection) (done bool) {
	// if there is no room
	if lobby.freeRoomsSize == 0 {
		// if room capacity ended return nil
		room := lobby.createRoom(models.NewUsualRoom())
		return room != nil
	}

	// lets find room for him
	for _, room := range lobby.freeRooms {
		if room.EnterPlayer(conn) {
			done = true
			break
		}
	}
	return done
}

// EnterBusyRoom try connect as observer
func (lobby *Lobby) EnterBusyRoom(conn *Connection) bool {
	return conn.room.enterObserver(conn)
}

// EnterRoom handle user join to room
func (lobby *Lobby) EnterRoom(request *Request) {

	done := false
	if room, ok := lobby.allRooms[request.Data.RoomSettings.ID]; ok {
		request.Connection.room = room
		done = lobby.EnterBusyRoom(request.Connection)
	} else {
		done = lobby.EnterFreeRoom(request.Connection)
	}

	if done {
		lobby.waiterToPlayer(request.Connection)
	} else {
		sendNotAllowed(request.Connection)
	}
	return
}

// sendRooms send rooms info for user
func (lobby *Lobby) sendRooms(conn *Connection) {

	roomsSlice := make([]Room, lobby.allRoomsSize)

	i := 0
	for _, room := range lobby.allRooms {
		roomsSlice[i] = *room
		i++
	}

	rooms := &Rooms{
		Size:  lobby.allRoomsSize,
		Rooms: roomsSlice,
	}
	conn.SendInformation(rooms)
}

// Run the room in goroutine
func (lobby *Lobby) Run() {

	for {
		select {
		case connection := <-lobby.ChanJoin:
			lobby.Join(connection)

		case request := <-lobby.chanRequest:
			lobby.EnterRoom(request)

		case connection := <-lobby.chanLeave:
			lobby.Leave(connection)
		}
	}
}

func sendNotAllowed(conn *Connection) {
	gameInfo := models.GameInfo{
		Send:   models.SendGameStatus,
		Status: StatusAborted,
	}
	conn.SendInformation(gameInfo)
}

func sendDisconnected(conn *Connection) {
	conn.player.LastAction = models.ActionDisconnect
	gameInfo := models.GameInfo{
		Send:         models.SendPlayerAction,
		PlayerAction: *conn.player,
	}
	conn.SendInformation(gameInfo)
}
