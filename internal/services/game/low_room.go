package game

import (
	"escapade/internal/models"
	//re "escapade/internal/return_errors"

	"sync"
)

func NewRoom(rs *models.RoomSettings, id int, lobby *Lobby) *Room {

	room := &Room{
		ID:              id,
		Status:          StatusPeopleFinding,
		PlayersCapacity: rs.Players,
		PlayersSize:     0,
		Players:         make(map[*Connection]*Playing),

		ObserversCapacity: 10, // вынести в конфиг
		ObserversSize:     0,
		Observers:         make(map[*Connection]*models.Player),

		lobby: lobby,
		Field: models.NewField(rs),
		//chanUpdateAll: make(chan *struct{}),
		//chanJoin:      make(chan *Connection),
		chanLeave:   make(chan *Connection),
		chanRequest: make(chan *Request),
	}
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
	conn.SendInformation(people)
}

// sendPeople send field info for user
func (room *Room) sendField(conn *Connection) {

	conn.SendInformation(*room.Field)
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

// observe try to connect user as observer
/* instruction to call
 first response will be as GameInfo(json)
if success, then PlayerAction will be returned
otherwise GameInfo

then if success be ready to receive Field and People models
*/
func (room *Room) enterObserver(conn *Connection) bool {
	// if we have a place
	if room.ObserversSize < room.ObserversCapacity {
		room.addObserver(conn)
		room.sendAllPlayerAction(conn, models.ActionConnectAsObserver)
		room.sendPrepareInfo(conn)
		return true
	}
	return false
}

// addPlayer add Connection as player
func (room *Room) EnterPlayer(conn *Connection) bool {
	// if room have already started
	// if room.Status != models.StatusPeopleFinding {
	// 	return false
	// }

	// if room hasnt got places
	if room.PlayersSize == room.PlayersCapacity {
		return false
	}

	cell := room.Field.RandomCell()
	cell.PlayerID = conn.GetPlayerID()
	playing := NewPlaying(conn.player, cell)
	room.addPlayer(conn, playing)
	room.sendAllPlayerAction(conn, models.ActionConnectAsPlayer)

	if room.PlayersSize == room.PlayersCapacity {
		room.startFlagPlacing()
	}

	return true
}

// RecoverPlayer call it in lobby.join if player disconnected
func (room *Room) RecoverPlayer(old *Connection, new *Connection) (played bool) {

	value := room.Players[old]
	room.addPlayer(new, value)
	room.removePlayer(old)

	room.sendAllPlayerAction(new, models.ActionReconnect)
	room.sendPrepareInfo(new)

	return
}

// RecoverPlayer call it in lobby.join if player disconnected
func (room *Room) RecoverObserver(old *Connection, new *Connection) (played bool) {

	room.addObserver(new)
	room.removeObserver(old)

	room.sendAllPlayerAction(new, models.ActionReconnect)
	room.sendPrepareInfo(new)

	return
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
	room.sendAllGameStatus(StatusClosed)
	room.Players = nil
	room.Observers = nil
	//delete(allRooms, room.ID)
	//delete(freeRooms, room.ID)
}

// removePlayer
func (room *Room) removePlayer(conn *Connection) {
	sendDisconnected(conn)
	delete(room.Players, conn)
	room.PlayersSize--
}

// removeObserver
func (room *Room) removeObserver(conn *Connection) {
	sendDisconnected(conn)
	delete(room.Observers, conn)
	room.ObserversSize--
}

// addPlayer
func (room *Room) addPlayer(conn *Connection, value *Playing) {
	room.Players[conn] = value
	room.PlayersSize++
}

// addObserver
func (room *Room) addObserver(conn *Connection) {
	room.Observers[conn] = conn.player
	room.ObserversSize++
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
	room.sendAllGameStatus(StatusAborted)
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
		conn.sendGroupInformation(info, waitJobs)
	}

	for conn := range room.Observers {
		waitJobs.Add(1)
		conn.sendGroupInformation(info, waitJobs)
	}
	waitJobs.Wait()
}

func (room *Room) setFlags() {
	for _, playing := range room.Players {
		room.Field.SetFlag(playing.Flag.X, playing.Flag.Y, playing.Flag.PlayerID)
	}
}

func (room *Room) fillField() {
	room.setFlags()
	room.Field.SetMines()
}
