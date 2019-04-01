package game

import (
	"escapade/internal/models"
	//re "escapade/internal/return_errors"
	"math/rand"
	"sync"
)

func NewRoom(rs *models.RoomSettings) *Room {
	id := rand.Intn(roomIDMax)

	// find id, that doesnt exist
	for elem, ok := allRooms[id]; ok; {
	}

	room := &Room{
		ID:              id,
		Status:          models.StatusPeopleFinding,
		PlayersCapacity: rs.Players,
		PlayersSize:     0,
		Players:         make(map[*Connection]*Playing),

		ObserversCapacity: 10, // вынести в конфиг
		ObserversSize:     0,
		Observers:         make(map[*Connection]*models.Player),

		Field:         models.NewField(rs),
		chanUpdateAll: make(chan *struct{}),
		chanJoin:      make(chan *Connection),
		chanLeave:     make(chan *Connection),
		chanRequest:   make(chan *Request),
	}

	allRooms[id] = room
	freeRooms[id] = room

	// run room
	go room.run()

	roomsCount++

	return room
}

// sendPeople send people info for user
func (room *Room) sendPeople(conn *Connection) {
	players := make([]models.Player, len(room.Players))
	observers := make([]models.Player, len(room.Observers))

	i := 0
	for _, playing := range room.Players {
		players[i] = *playing.Player
		i++
	}

	i = 0
	for _, player := range room.Observers {
		observers[i] = *player
		i++
	}

	people := &models.People{
		PlayersCapacity: room.PlayersCapacity,
		PlayersSize:     room.PlayersSize,
		Players:         players,

		ObserversCapacity: room.ObserversCapacity, // вынести в конфиг
		ObserversSize:     room.ObserversSize,
		Observers:         observers,
	}
	conn.sendInformation(people)
}

// sendPeople send field info for user
func (room *Room) sendField(conn *Connection) {

	conn.sendInformation(*room.Field)
}

// sendAllPlayerAction send users action to everybody
func (room *Room) sendAllPlayerAction(conn *Connection, action int) {
	conn.player.LastAction = action

	gameInfo := models.GameInfo{
		Send:         models.SendPlayerAction,
		PlayerAction: *conn.player,
	}
	room.sendAllPlayers(gameInfo)
}

// sendAllGameStatus send gamestatus to everybody
func (room *Room) sendAllGameStatus(status int) {
	room.Status = status

	gameInfo := models.GameInfo{
		Send:   models.SendGameStatus,
		Status: status,
	}
	room.sendAllPlayers(gameInfo)
}

func (room *Room) sendNotAllowed(conn *Connection) {
	gameInfo := models.GameInfo{
		Send:   models.SendGameStatus,
		Status: models.StatusAborted,
	}
	conn.sendInformation(gameInfo)
}

func (room *Room) sendRoomIsBlocked(conn *Connection) {
	gameInfo := models.GameInfo{
		Send:   models.SendGameStatus,
		Status: models.StatusBlock,
	}
	conn.sendInformation(gameInfo)
}

// observe try to connect user as observer
/* instruction to call
 first response will be as GameInfo(json)
if success, then PlayerAction will be returned
otherwise GameInfo

then if success be ready to receive Field and People models
*/
func (room *Room) addObserver(conn *Connection) bool {
	// if we have a place
	if room.ObserversSize < room.ObserversCapacity {
		room.Observers[conn] = conn.player
		room.ObserversSize++
		room.sendAllPlayerAction(conn, models.ActionConnectAsObserver)
		room.sendPrepareInfo(conn)
		return true
	}
	return false
}

// addPlayer add Connection as player
func (room *Room) addPlayer(conn *Connection) bool {
	// if room have already started
	if room.Status != models.StatusPeopleFinding {
		return false
	}

	// if room hasnt got places
	if room.PlayersSize == room.PlayersCapacity {
		return false
	}

	cell := room.Field.RandomCell()
	flag := models.NewFlag(cell, conn.player.ID)
	room.Players[conn] = NewPlaying(conn.player, flag)
	room.PlayersSize++

	room.sendAllPlayerAction(conn, models.ActionConnectAsPlayer)

	return true
}

// alreadyPlaying
// If user disconnected, it will recover it
// or if somebody use second tab it will delete old
// and activate new
func (room *Room) alreadyPlaying(conn *Connection) (played bool) {
	for oldConn, value := range room.Players {
		if oldConn.player.ID == conn.player.ID {
			// update connection
			room.Players[conn] = value
			delete(room.Players, oldConn)
			played = true
			break
		}
	}
	if played {
		room.sendAllPlayerAction(conn, models.ActionReconnect)
		room.sendPrepareInfo(conn)
	}
	return
}

// room closes
func (room *Room) close() {
	room.sendAllGameStatus(models.StatusClosed)
	room.Players = nil
	room.Observers = nil
	delete(allRooms, room.ID)
	delete(freeRooms, room.ID)
}

// removePlayer
func (room *Room) removePlayer(conn *Connection) {
	delete(room.Players, conn)
	room.PlayersSize--
}

// removeObserver
func (room *Room) removeObserver(conn *Connection) {
	room.ObserversSize--
	delete(room.Observers, conn)
}

func (room *Room) removeBeforeLaunch(conn *Connection) {
	room.removePlayer(conn)
	if room.PlayersSize == 0 {
		room.close()
	} else {
		room.sendAllPlayerAction(conn, models.ActionDisconnect)
	}
}

// removeDuringGame
func (room *Room) removeDuringGame(conn *Connection) {

	// if it is player, maybe he will return
	if playing, ok := room.Players[conn]; ok {
		// if he decided to give up
		if conn.player.LastAction == models.ActionGiveUp {
			playing.Finished = true
			room.sendAllPlayerAction(conn, models.ActionGiveUp)
		} else {
			room.sendAllPlayerAction(conn, models.ActionDisconnect)
		}
		return
	}

	// if it is observer, let him go
	if _, ok := room.Observers[conn]; ok {
		room.removeObserver(conn)
		room.sendAllPlayerAction(conn, models.ActionDisconnect)
		return
	}

	// somebody broke, fix it
	room.sendAllGameStatus(models.StatusAborted)
	return
}

// removeFinishedGame
func (room *Room) removeAfterFinish(conn *Connection) {

	if _, ok := room.Observers[conn]; ok {
		room.removeObserver(conn)
	}
	room.sendAllPlayerAction(conn, models.ActionDisconnect)
	return
}

func (room *Room) sendAllPlayers(info interface{}) {
	waitJobs := &sync.WaitGroup{}
	for conn := range room.Players {
		waitJobs.Add(1)
		conn.sendGroupInformation(info)
	}

	for conn := range room.Observers {
		waitJobs.Add(1)
		conn.sendGroupInformation(info)
	}
	waitJobs.Wait()
}
