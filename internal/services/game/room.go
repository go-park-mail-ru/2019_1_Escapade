package game

import (
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"math/rand"
)

type Request struct {
	Connection *Connection
	Cell *models.Cell
}

type Response struct {
	Connection *Connection
	Information *models.GameInfo
}

type Room struct {
	ID int
	Size int
	Count int
	Started bool
	People map[*Connection]*Playing
	Field *Field
	UpdateAll chan *struct{}
	join chan *Connection
	leave chan *Connection
	request chan *Request
	response chan *Response
}

// we use map, not array, cause in future will add name of rooms
var allRooms = make(map[int]*Room)
var freeRooms = make(map[int]*Room)
var roomsCount int

// вынести в конфиг
var roomIDMax = 10000

func NewRoom(rs *models.RoomSettings) *Room {
	id := rand.Intn(roomIDMax)
	
	// find id, that doesnt exist
	for elem, ok := allRooms[id]; ok; {
	} 
	
	room := &Room{
		ID:        id,
		Size: rs.Players,
		Count: 0,
		People: make(map[*Connection]*Playing),
		Field: NewField(rs),
		UpdateAll:   make(chan *struct{}),
		join:        make(chan *Connection),
		leave:       make(chan *Connection),
		request:        make(chan *Request),
		response:       make(chan *Response),
	}

	allRooms[id] = room
	freeRooms[id] = room

	// run room
	go room.run()

	roomsCount += 1

	return room
}

func (room *Room) roomJoin(conn *Connection) {
	// if room is full return
	if room.Started {
		roomJoining := models.RoomJoiningFault(re.ErrorRoomIsFull().Error())
		conn.sendInformation(roomJoining)
		return
	}

	// init player info
	cell := room.Field.randomCell()
	flag := models.NewFlag(cell, conn.player.ID)
	room.People[conn] = NewPlaying(conn.player,flag)

	// update room info
	room.Count++
	if (room.Count == room.Size) {
		room.Started = true	
	}

	// send user, that he connected
	roomJoining := models.RoomJoiningSuccess()
	conn.sendInformation(roomJoining)
}

func (room *Room) messageStatus() (){
	// send to everybody, that new user joined
	status := models.NewStatus(room.Count, room.Started)
	room.updateAllPlayers(status)
	return
}

func (room *Room) roomLeave(conn *Connection) {
	// if player in game
	if room.Started {
		room.People[conn] = 
		roomJoining := models.RoomJoiningFault(re.ErrorRoomIsFull().Error())
		conn.sendInformation(roomJoining)
		return
	}

	// init player info
	cell := room.Field.randomCell()
	flag := models.NewFlag(cell, conn.player.ID)
	room.People[conn] = NewPlaying(conn.player,flag)

	// update room info
	room.Count++
	if (room.Count == room.Size) {
		room.Started = true	
	}

	// send user, that he connected
	roomJoining := models.RoomJoiningSuccess()
	conn.sendInformation(roomJoining)
}

// Run the room in goroutine
func (room *Room) run() {
	for {
		select {
		case connection := <-room.join:
			room.roomJoin(connection)
			room.messageStatus()

		case c := <-r.leave:
			c.GiveUp()
			r.updateAllPlayers()
			delete(r.playerConns, c)
			if len(r.playerConns) == 0 {
				goto Exit
			}
		case <-r.updateAll:
			r.updateAllPlayers()
		}
	}

Exit:

	// delete room
	delete(allRooms, r.name)
	delete(freeRooms, r.name)
	roomsCount -= 1
	log.Print("Room closed:", r.name)
}

func (r *Room) updateAllPlayers(info interface{}) {
	waitJobs := &sync.WaitGroup{}
	for c := range r.People {
		wait_jobs.Add(1)
		c.sendGroupInformation(info)
	}
	waitJobs.Wait()
}