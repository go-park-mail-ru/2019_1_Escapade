package game

import (
	"escapade/internal/models"
	//re "escapade/internal/return_errors"

	"sync"
)

func NewRoom(rs *models.RoomSettings, name string, lobby *Lobby) *Room {

	room := &Room{
		Name:   name,
		Status: StatusPeopleFinding,

		players:   NewConnections(rs.Players),
		observers: NewConnections(10),

		history: make([]*PlayerAction, 0),
		flags:   make(map[*Connection]*models.Cell),

		lobby:       lobby,
		field:       models.NewField(rs),
		chanLeave:   make(chan *Connection),
		chanRequest: make(chan *RoomRequest),
	}
	return room
}

// observe try to connect user as observer
/* instruction to call
 first response will be as GameInfo(json)
if success, then PlayerAction will be returned
otherwise GameInfo

then if success be ready to receive Field and People models
*/
func (room *Room) enterObserver(conn *Connection) bool {
	// if we have a place
	if room.observers.enoughPlace() {
		room.observers.Add(conn, false)
		room.addAction(conn, ActionConnectAsObserver)
		room.sendTAIRPeople()
		room.sendTOCAll(conn)
		return true
	}
	return false
}

// addPlayer add Connection as player
func (room *Room) EnterPlayer(conn *Connection) bool {
	// if room have already started
	if room.Status != StatusPeopleFinding {
		return false
	}

	// if room hasnt got places
	if !room.players.enoughPlace() {
		return false
	}

	cell := room.field.RandomCell()
	cell.PlayerID = conn.GetPlayerID()
	room.players.Add(conn, false)

	room.addAction(conn, ActionConnectAsPlayer)
	room.sendTAIRPeople()
	room.sendTOCAll(conn)

	if !room.players.enoughPlace() {
		room.startFlagPlacing()
	}

	return true
}

// RecoverPlayer call it in lobby.join if player disconnected
func (room *Room) RecoverPlayer(old *Connection, new *Connection) (played bool) {

	room.players.Add(new, room.players.Get[old])
	room.players.Remove(old)

	sendError(old, "RecoverPlayer", "another enter in account ")
	room.addAction(new, ActionReconnect)
	room.sendTAIRPeople()
	room.sendTOCAll(new)

	return
}

// RecoverObserver call it in lobby.join if observer disconnected
func (room *Room) RecoverObserver(old *Connection, new *Connection) (played bool) {

	room.observers.Add(new, room.observers.Get[old])
	room.observers.Remove(old)

	sendError(old, "RecoverObserver", "another enter in account ")
	room.addAction(new, ActionReconnect)
	room.sendTAIRPeople()
	room.sendTOCAll(new)

	return
}

// alreadyPlaying
// If user disconnected, it will recover it
// or if somebody use second tab it will delete old
// and activate new
func (room *Room) alreadyPlaying(conn *Connection) (played bool) {
	for oldConn := range room.players.Get {
		if oldConn.GetPlayerID() == conn.GetPlayerID() {
			room.RecoverPlayer(oldConn, conn)
			break
		}
	}
	return
}

// room closes
// func (room *Room) close() {
// 	room.sendAllGameStatus(StatusClosed)
// 	room.Players = nil
// 	room.Observers = nil
// 	//delete(allRooms, room.ID)
// 	//delete(freeRooms, room.ID)
// }

// flagFound is called, when somebody find cell flag
func (room *Room) flagFound(found *models.Cell) {
	id := found.Value - models.CellIncrement
	for conn := range room.players.Get {
		if conn.GetPlayerID() == id {
			room.kill(conn)
		}
	}
}

// kill make user die, decrement size and check for finish battle
func (room *Room) kill(conn *Connection) {
	if !room.players.Get[conn] {
		room.players.Get[conn] = true
		room.players.Size--
		if room.players.Size <= 1 {
			room.lobby.roomFinish(room)
		}
		room.addAction(conn, ActionGiveUp)
		room.sendTAIRHistory()
	}
}

func (room *Room) GiveUp(conn *Connection) {
	room.kill(conn)
	room.addAction(conn, ActionGiveUp)
	room.sendTAIRHistory()
}

func (room *Room) removeBeforeLaunch(conn *Connection) {
	room.players.Remove(conn)
	// if room.PlayersSize == 0 {
	// 	room.close()
	// }
}

// removeFinishedGame
func (room *Room) removeAfterLaunch(conn *Connection) {

	if _, ok := room.observers.Get[conn]; ok {
		room.observers.Remove(conn)
	}
	return
}

func (room *Room) setFlags() {
	for conn, cell := range room.flags {
		room.field.SetFlag(cell.X, cell.Y, conn.GetPlayerID())
	}
}

func (room *Room) fillField() {
	room.setFlags()
	room.field.SetMines()
}

func (room *Room) sendToAllInRoom(info interface{}) {
	waitJobs := &sync.WaitGroup{}
	for conn := range room.players.Get {
		waitJobs.Add(1)
		conn.sendGroupInformation(info, waitJobs)
	}

	for conn := range room.observers.Get {
		waitJobs.Add(1)
		conn.sendGroupInformation(info, waitJobs)
	}
	waitJobs.Wait()
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendTAIRPeople() {
	get := &RoomGet{
		players:   true,
		observers: true,
		history:   true,
	}
	send := room.makeGetModel(get)
	room.sendToAllInRoom(send)
}

// sendTAIRField send field to all in room
func (room *Room) sendTAIRField() {
	get := &RoomGet{
		field: true,
	}
	send := room.makeGetModel(get)
	room.sendToAllInRoom(send)
}

// sendTAIRHistory send actions history to all in room
func (room *Room) sendTAIRHistory() {
	get := &RoomGet{
		history: true,
	}
	send := room.makeGetModel(get)
	room.sendToAllInRoom(send)
}

// sendTAIRPeople send only name and status to all in room
func (room *Room) sendTAIRStatus() {
	get := &RoomGet{}
	send := room.makeGetModel(get)
	room.sendToAllInRoom(send)
}

// sendTAIRAll send everything to one connection
func (room *Room) sendTOCAll(conn *Connection) {
	get := &RoomGet{
		players:   true,
		observers: true,
		field:     true,
		history:   true,
	}
	if room.Status == StatusPeopleFinding {
		get.field = false
	}
	send := room.makeGetModel(get)
	conn.SendInformation(send)
}

func (room *Room) requestGet(rr *RoomRequest) {
	send := room.makeGetModel(rr.Get)
	rr.Connection.SendInformation(send)
}

func (room *Room) makeGetModel(get *RoomGet) *Room {
	sendRoom := &Room{
		Name:   room.Name,
		Status: room.Status,
	}

	if get.players {
		sendRoom.players = room.players
	}
	if get.observers {
		sendRoom.observers = room.observers
	}
	if get.field {
		sendRoom.field = room.field
	}
	if get.history {
		sendRoom.history = room.history
	}
	return sendRoom
}
