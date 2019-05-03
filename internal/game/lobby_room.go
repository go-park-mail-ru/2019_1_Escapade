package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"fmt"
)

// ----- handle room status
// roomStart - room remove from free
func (lobby *Lobby) roomStart(room *Room) {
	defer utils.CatchPanic("lobby_room.go roomStart()")
	lobby.FreeRooms.Remove(room)
	lobby.sendRoomUpdate(*room, All)
}

// roomFinish - room remove from all
func (lobby *Lobby) roomFinish(room *Room) {
	defer utils.CatchPanic("lobby_room.go finishGame()")
	lobby.AllRooms.Remove(room)
	lobby.sendRoomUpdate(*room, All)
}

// CloseRoom free room resources
func (lobby *Lobby) CloseRoom(room *Room) {
	// if not in freeRooms nothing bad will happen
	// there is check inside, it will just return without errors
	lobby.FreeRooms.Remove(room)
	lobby.AllRooms.Remove(room)
	fmt.Println("sendRoomDelete")
	lobby.sendRoomDelete(*room, All)
}

// createAndAddToRoom create room and add player to it
func (lobby *Lobby) createAndAddToRoom(rs *models.RoomSettings, conn *Connection) (room *Room, err error) {
	if room, err = lobby.createRoom(rs); err == nil {
		conn.debug("We create your own room, cool!")
		room.addPlayer(conn)
	} else {
		conn.debug("cant create. Why?")
		room.sendError(err, *conn)
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
	if !lobby.AllRooms.Add(room) {
		err = re.ErrorLobbyCantCreateRoom()
		fmt.Println("cant create room")
		return
	}

	lobby.FreeRooms.Add(room)
	lobby.sendRoomCreate(*room, All) // inform all about new room
	return
}

// LoadRooms load rooms from database
func (lobby *Lobby) LoadRooms(URLs []string) error {

	for _, URL := range URLs {
		room, err := lobby.Load(URL)
		if err != nil {
			return err
		}

		if !lobby.AllRooms.Add(room) {
			fmt.Println("cant create room")
			return re.ErrorLobbyCantCreateRoom()
		}
		lobby.AllRooms.Add(room)
		lobby.sendRoomCreate(*room, All)
	}
	return nil
}
