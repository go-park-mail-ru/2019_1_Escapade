package game

import (
	"encoding/json"
	"escapade/internal/models"

	"fmt"

	"sync"
)

// NewRoom return new instance of room
func NewRoom(rs *models.RoomSettings, name string, lobby *Lobby) *Room {
	fmt.Println("NewRoom rs = ", *rs)
	room := &Room{
		Name:      name,
		Status:    StatusPeopleFinding,
		Players:   NewConnections(rs.Players),
		Observers: NewConnections(10),

		History: make([]*PlayerAction, 0),
		flags:   make(map[*Connection]*models.Cell),

		lobby:     lobby,
		Field:     models.NewField(rs),
		chanLeave: make(chan *Connection),
		//chanRequest: make(chan *RoomRequest),
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
	if room.Observers.enoughPlace() {
		room.Observers.Add(conn)
		room.addAction(conn, ActionConnectAsObserver)
		room.sendTAIRPeople()
		room.sendTOCAll(conn)
		return true
	}
	return false
}

// EnterPlayer handle player try to enter room
func (room *Room) EnterPlayer(conn *Connection) bool {
	// if room have already started
	if room.Status != StatusPeopleFinding {
		return false
	}

	// if room hasnt got places
	if !room.Players.enoughPlace() {
		return false
	}

	cell := room.Field.RandomCell()
	cell.PlayerID = conn.GetPlayerID()
	conn.Player.Reset()
	room.Players.Add(conn)

	room.addAction(conn, ActionConnectAsPlayer)
	room.sendTAIRPeople()

	if !room.Players.enoughPlace() {
		room.startFlagPlacing()
	}

	return true
}

// RecoverPlayer call it in lobby.join if player disconnected
func (room *Room) RecoverPlayer(old *Connection, new *Connection) (played bool) {

	new.Player = old.Player
	room.Players.Change(new)

	sendError(old, "RecoverPlayer", "another enter in account ")
	room.addAction(new, ActionReconnect)
	room.sendTAIRPeople()
	room.sendTOCAll(new)

	return
}

// RecoverObserver call it in lobby.join if observer disconnected
func (room *Room) RecoverObserver(old *Connection, new *Connection) (played bool) {

	room.Observers.Change(new)

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
	thatID := conn.GetPlayerID()
	for id, oldConn := range room.Players.Get {
		if id == thatID {
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
	thatID := found.Value - models.CellIncrement
	for id, conn := range room.Players.Get {
		if thatID == id {
			room.kill(conn)
		}
	}
}

// kill make user die, decrement size and check for finish battle
func (room *Room) kill(conn *Connection) {
	// cause all in pointers
	if conn.Player.Finished {
		conn.Player.Finished = true
		room.Players.Size--
		if room.Players.Size <= 1 {
			room.lobby.roomFinish(room)
		}
		room.addAction(conn, ActionGiveUp)
		room.sendTAIRHistory()
	}
}

// GiveUp kill conn
func (room *Room) GiveUp(conn *Connection) {
	room.kill(conn)
	room.addAction(conn, ActionGiveUp)
	room.sendTAIRHistory()
}

func (room *Room) removeBeforeLaunch(conn *Connection) {
	room.Players.Remove(conn)
	// if room.PlayersSize == 0 {
	// 	room.close()
	// }
}

// removeFinishedGame
func (room *Room) removeAfterLaunch(conn *Connection) {

	if _, ok := room.Observers.Get[conn.GetPlayerID()]; ok {
		room.Observers.Remove(conn)
	}
	return
}

func (room *Room) setFlags() {
	for conn, cell := range room.flags {
		room.Field.SetFlag(cell.X, cell.Y, conn.GetPlayerID())
	}
}

func (room *Room) fillField() {
	fmt.Println("fillField", room.Field.Height, room.Field.Width, len(room.Field.Matrix))

	room.setFlags()
	room.Field.SetMines()

}

func (room *Room) sendToAllInRoom(info interface{}) {
	waitJobs := &sync.WaitGroup{}
	bytes, _ := json.Marshal(info)
	for _, conn := range room.Players.Get {
		waitJobs.Add(1)
		conn.sendGroupInformation(bytes, waitJobs)
	}

	for _, conn := range room.Observers.Get {
		waitJobs.Add(1)
		conn.sendGroupInformation(bytes, waitJobs)
	}
	waitJobs.Wait()
}

// sendTAIRPeople send players, observers and history to all in room
func (room *Room) sendTAIRPeople() {
	get := &RoomGet{
		Players:   true,
		Observers: true,
		History:   true,
	}
	send := room.makeGetModel(get)
	room.sendToAllInRoom(send)
}

// sendTAIRField send field to all in room
func (room *Room) sendTAIRField() {
	get := &RoomGet{
		Field: true,
	}
	send := room.makeGetModel(get)

	room.sendToAllInRoom(send)
}

// sendTAIRHistory send actions history to all in room
func (room *Room) sendTAIRHistory() {
	get := &RoomGet{
		History: true,
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
		Players:   true,
		Observers: true,
		Field:     true,
		History:   true,
	}
	if room.Status == StatusPeopleFinding {
		get.Field = false
	}
	send := room.makeGetModel(get)
	bytes, _ := json.Marshal(send)
	conn.SendInformation(bytes)
}

func (room *Room) requestGet(conn *Connection, rr *RoomRequest) {
	send := room.makeGetModel(rr.Get)
	fmt.Println("here you go?", rr.Get)
	bytes, _ := json.Marshal(send)
	conn.SendInformation(bytes)
}

func (room *Room) makeGetModel(get *RoomGet) *Room {
	sendRoom := &Room{
		Name:   room.Name,
		Status: room.Status,
	}

	if get.Players {
		sendRoom.Players = room.Players
	}
	if get.Observers {
		sendRoom.Observers = room.Observers
	}
	if get.Field {
		sendRoom.Field = room.Field
	}
	if get.History {
		sendRoom.History = room.History
	}
	return sendRoom
}
