package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"fmt"
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

	go lobby.freeRoomsRemove(room)
	go lobby.sendRoomUpdate(*room, All)
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

	go lobby.allRoomsRemove(room)
	go lobby.sendRoomDelete(*room, All)
}

// CloseRoom free room resources
func (lobby *Lobby) CloseRoom(room *Room) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		lobby.wGroup.Done()
		utils.CatchPanic("lobby_room.go roomFinish()")
	}()

	// if not in freeRooms nothing bad will happen
	// there is check inside, it will just return without errors
	lobby.freeRoomsRemove(room)
	lobby.allRoomsRemove(room)
	fmt.Println("sendRoomDelete")
	go lobby.sendRoomDelete(*room, All)
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
		conn.debug("We create your own room, cool!")
		room.addPlayer(conn, false)
	} else {
		conn.debug("cant create. Why?" + err.Error())
		panic(":(")
	}
	return
}

// createRoom create room, add to all and free rooms
// and run it
func (lobby *Lobby) createRoom(rs *models.RoomSettings) (room *Room, err error) {

	id := utils.RandomString(16) // вынести в кофиг
	if room, err = NewRoom(rs, id, lobby); err != nil {
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
	if lobby.metrics {
		metrics.Rooms.Add(1)
		metrics.FreeRooms.Add(1)
	}

	if !lobby.allRoomsAdd(room) {
		err = re.ErrorLobbyCantCreateRoom()
		fmt.Println("cant add to all")
		return err
	}

	if !lobby.freeRoomsAdd(room) {
		err = re.ErrorLobbyCantCreateRoom()
		fmt.Println("cant add to free")
		return err
	}

	lobby.sendRoomCreate(*room, All) // inform all about new room
	return
}
