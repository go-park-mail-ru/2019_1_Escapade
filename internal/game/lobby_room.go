package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// RoomStart - room remove from free
func (lobby *Lobby) RoomStart(room *Room) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby_room.go RoomStart()")
	}()

	go lobby.freeRooms.Remove(room.ID())
	go lobby.sendRoomUpdate(room, All)
}

// roomFinish - room remove from all
func (lobby *Lobby) roomFinish(room *Room) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby_room.go roomFinish()")
	}()

	go lobby.allRooms.Remove(room.ID())
	go lobby.sendRoomDelete(room, All)
}

// CloseRoom free room resources
func (lobby *Lobby) CloseRoom(room *Room) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	room.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		room.wGroup.Done()
		utils.CatchPanic("lobby_room.go roomFinish()")
	}()

	// if not in freeRooms nothing bad will happen
	// there is check inside, it will just return without errors
	lobby.freeRooms.Remove(room.ID())
	lobby.allRooms.Remove(room.ID())
	utils.Debug(false, "sendRoomDelete")
	lobby.sendRoomDelete(room, All)
}

// CreateAndAddToRoom create room and add player to it
func (lobby *Lobby) CreateAndAddToRoom(rs *models.RoomSettings, conn *Connection) (room *Room, err error) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_room.go CreateAndAddToRoom()")
		lobby.wGroup.Done()
	}()

	if room, err = lobby.createRoom(rs); err == nil {
		utils.Debug(false, "We create your own room, cool!", conn.ID())
		room.addConnection(conn, true, false)
	} else {
		utils.Debug(true, "cant create. Why?", conn.ID(), err.Error())
	}
	return
}

// createRoom create room, add to all and free rooms
// and run it
func (lobby *Lobby) createRoom(rs *models.RoomSettings) (room *Room, err error) {

	id := utils.RandomString(16) // вынести в кофиг
	if room, err = NewRoom(lobby.config.Field, lobby, rs, id); err != nil {
		return
	}
	if err = lobby.addRoom(room); err != nil {
		return
	}
	return
}

// LoadRooms load rooms from database
func (lobby *Lobby) LoadRooms(URLs []string) error {

	if lobby.done() {
		return re.ErrorLobbyDone()
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_room.go LoadRooms()")
		lobby.wGroup.Done()
	}()

	for _, URL := range URLs {
		room, err := lobby.Load(URL)
		if err != nil {
			return err
		}
		if err = lobby.addRoom(room); err != nil {
			return err
		}
	}
	return nil
}

func (lobby *Lobby) addRoom(room *Room) (err error) {
	if lobby.config.Metrics {
		metrics.Rooms.Add(1)
		metrics.FreeRooms.Add(1)
	}

	if !lobby.allRooms.Add(room) {
		err = re.ErrorLobbyCantCreateRoom()
		utils.Debug(false, "cant add to all rooms")
		return err
	}

	if !lobby.freeRooms.Add(room) {
		err = re.ErrorLobbyCantCreateRoom()
		utils.Debug(false, "cant add to free rooms")
		return err
	}

	lobby.sendRoomCreate(room, All) // inform all about new room
	return
}
