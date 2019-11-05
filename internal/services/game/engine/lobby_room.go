package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// RoomStart - room remove from free
func (lobby *Lobby) RoomStart(room *Room, roomID string) {
	if lobby.done() {
		return
	}

	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	go lobby.removeFromFreeRooms(roomID)
	go lobby.sendRoomUpdate(room, All)
}

// roomFinish - room remove from all
func (lobby *Lobby) roomFinish(roomID string) {
	defer utils.CatchPanic("lobby_room.go roomFinish()")

	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	go lobby.removeFromAllRooms(roomID)

	go lobby.sendRoomDelete(roomID, All)
}

// CloseRoom free room resources
func (lobby *Lobby) CloseRoom(roomID string) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	go lobby.removeFromFreeRooms(roomID)
	go lobby.removeFromAllRooms(roomID)

	go lobby.sendRoomDelete(roomID, All)
}

// CreateAndAddToRoom create room and add player to it
func (lobby *Lobby) CreateAndAddToRoom(rs *models.RoomSettings, conn *Connection) (*Room, error) {
	defer utils.CatchPanic("lobby_room.go CreateAndAddToRoom()")
	if lobby.done() || conn.done() {
		return nil, re.ErrorLobbyDone()
	}

	lobby.wGroup.Add(1)
	defer lobby.wGroup.Done()

	conn.wGroup.Add(1)
	defer conn.wGroup.Done()

	room, err := lobby.createRoom(rs)
	if err == nil {
		utils.Debug(false, "We create your own room, cool!", conn.ID())
		room.people.add(conn, true, false)
	} else {
		utils.Debug(true, "cant create. Why?", conn.ID(), err.Error())
	}
	return room, err
}

// createRoom create room, add to all and free rooms
// and run it
func (lobby *Lobby) createRoom(rs *models.RoomSettings) (room *Room, err error) {

	id := utils.RandomString(16) // вынести в кофиг
	room, err = NewRoom(lobby.config().Field, lobby,
		&models.Game{Settings: rs}, id)
	if err != nil {
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

// addRoom add room to slice of all and free lobby rooms
func (lobby *Lobby) addRoom(room *Room) (err error) {

	if err = lobby.addToAllRooms(room); err != nil {
		return
	}

	if err = lobby.addToFreeRooms(room); err != nil {
		return
	}

	lobby.sendRoomCreate(room, All) // inform all about new room
	return
}
