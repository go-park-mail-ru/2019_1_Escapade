package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// RoomStart - room remove from free
func (lobby *Lobby) RoomStart(room *Room, roomID string) {
	lobby.s.DoWithOther(room, func() {
		lobby.removeFromFreeRooms(roomID)
		lobby.sendRoomUpdate(room, All)
	})
}

// roomFinish - room remove from all
func (lobby *Lobby) roomFinish(roomID string) {
	lobby.s.Do(func() {
		lobby.removeFromAllRooms(roomID)
		lobby.sendRoomDelete(roomID, All)
	})
}

// CloseRoom free room resources
func (lobby *Lobby) CloseRoom(roomID string) {
	lobby.s.Do(func() {
		lobby.removeFromFreeRooms(roomID)
		lobby.removeFromAllRooms(roomID)
		lobby.sendRoomDelete(roomID, All)
	})
}

// CreateAndAddToRoom create room and add player to it
func (lobby *Lobby) CreateAndAddToRoom(rs *models.RoomSettings, conn *Connection) (*Room, error) {
	var (
		room *Room
		err  error
	)
	lobby.s.DoWithOther(conn, func() {
		room, err = lobby.createRoom(rs)
		if err == nil {
			utils.Debug(false, "We create your own room, cool!", conn.ID())
			room.people.add(conn, true, false)
		} else {
			utils.Debug(true, "cant create. Why?", conn.ID(), err.Error())
		}
	})
	return room, err
}

// createRoom create room, add to all and free rooms
// and run it
func (lobby *Lobby) createRoom(rs *models.RoomSettings) (*Room, error) {
	var (
		room *Room
		err  error
	)
	lobby.s.Do(func() {
		id := utils.RandomString(lobby.rconfig().IDLength)
		game := &models.Game{Settings: rs}
		room, err = NewRoom(lobby.rconfig(), lobby, game, id)
		if err != nil {
			return
		}
		if err = lobby.addRoom(room); err != nil {
			return
		}
	})
	return room, err
}

// LoadRooms load rooms from database
func (lobby *Lobby) LoadRooms(URLs []string) error {
	var err = re.ErrorLobbyDone()
	lobby.s.Do(func() {
		err = nil
		for _, URL := range URLs {
			room, err := lobby.Load(URL)
			if err != nil {
				return
			}
			if err = lobby.addRoom(room); err != nil {
				return
			}
		}
	})
	return err
}

// addRoom add room to slice of all and free lobby rooms
func (lobby *Lobby) addRoom(room *Room) error {
	var err = re.ErrorLobbyDone()
	lobby.s.Do(func() {
		if err = lobby.addToAllRooms(room); err != nil {
			return
		}
		if err = lobby.addToFreeRooms(room); err != nil {
			return
		}
		lobby.sendRoomCreate(room, All) // inform all about new room
	})
	return err
}
