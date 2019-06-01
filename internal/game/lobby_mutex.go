package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

//setDone set done = true. It will finish all operaions on Lobby
func (lobby *Lobby) setDone() {
	lobby.doneM.Lock()
	lobby._done = true
	lobby.doneM.Unlock()
}

// done return '_done' field
func (lobby *Lobby) done() bool {
	lobby.doneM.RLock()
	v := lobby._done
	lobby.doneM.RUnlock()
	return v
}

// allRoomsFree free slice of all rooms
func (lobby *Lobby) allRoomsFree() {
	lobby.allRoomsM.Lock()
	defer lobby.allRoomsM.Unlock()
	lobby._allRooms.Free()
}

// freeRoomsFree free slice of free rooms
func (lobby *Lobby) freeRoomsFree() {
	lobby.freeRoomsM.Lock()
	defer lobby.freeRoomsM.Unlock()
	lobby._freeRooms.Free()
}

// allRoomsSearch search room by its ID
func (lobby *Lobby) allRoomsSearch(roomID string) (int, *Room) {
	lobby.allRoomsM.RLock()
	defer lobby.allRoomsM.RUnlock()
	index, room := lobby._allRooms.SearchRoom(roomID)
	return index, room
}

// TODO Зачем нам массив Playing, если у нас есть такие чудесные функции?
// allRoomsSearchPlayer search player in all rooms
func (lobby *Lobby) allRoomsSearchPlayer(conn *Connection, disconnect bool) (int, *Room) {
	lobby.allRoomsM.RLock()
	defer lobby.allRoomsM.RUnlock()
	index, room := lobby._allRooms.SearchPlayer(conn, disconnect)
	return index, room
}

// allRoomsSearchObserver search observer in all rooms
func (lobby *Lobby) allRoomsSearchObserver(conn *Connection) *Connection {
	lobby.allRoomsM.RLock()
	defer lobby.allRoomsM.RUnlock()
	old := lobby._allRooms.SearchObserver(conn)
	return old
}

// freeRoomsEmpty return flag is free rooms slice empty
func (lobby *Lobby) freeRoomsEmpty() bool {
	lobby.freeRoomsM.RLock()
	defer lobby.freeRoomsM.RUnlock()
	v := lobby._freeRooms.Empty()
	return v
}

// freeRooms return '_freeRooms' field
func (lobby *Lobby) freeRooms() []*Room {
	lobby.freeRoomsM.RLock()
	defer lobby.freeRoomsM.RUnlock()
	v := lobby._freeRooms.Get
	return v
}

// appendMessage append message to messages slice
func (lobby *Lobby) appendMessage(message *models.Message) {
	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	lobby._messages = append(lobby._messages, message)
}

// removeMessage remove message from messages slice
func (lobby *Lobby) removeMessage(i int) {
	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	if i < 0 {
		return
	}
	size := len(lobby._messages)

	lobby._messages[i], lobby._messages[size-1] = lobby._messages[size-1], lobby._messages[i]
	lobby._messages[size-1] = nil
	lobby._messages = lobby._messages[:size-1]
	return
}

// setMessage update message from messages slice with index i
func (lobby *Lobby) setMessage(i int, message *models.Message) {
	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	if i < 0 {
		return
	}
	lobby._messages[i] = message
	lobby._messages[i].Edited = true
	return
}

// findMessage search message by message ID
func (lobby *Lobby) findMessage(search *models.Message) int {
	lobby.messagesM.Lock()
	messages := lobby._messages
	lobby.messagesM.Unlock()

	for i, message := range messages {
		if message.ID == search.ID {
			return i
		}
	}
	return -1
}

// Messages return slice of messages
func (lobby *Lobby) Messages() []*models.Message {

	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	return lobby._messages
}

// freeRoomsRemove remove room from free rooms slice
func (lobby *Lobby) freeRoomsRemove(room *Room) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_mutex.go freeRoomsRemove()")
		lobby.wGroup.Done()
	}()

	lobby.freeRoomsM.Lock()
	defer lobby.freeRoomsM.Unlock()
	lobby._freeRooms.Remove(room)
}

// allRoomsRemove remove room from all rooms slice
func (lobby *Lobby) allRoomsRemove(room *Room) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_mutex.go allRoomsRemove()")
		lobby.wGroup.Done()
	}()

	lobby.allRoomsM.Lock()
	defer lobby.allRoomsM.Unlock()
	lobby._allRooms.Remove(room)
}

// freeRoomsAdd add new room to free rooms slice
func (lobby *Lobby) freeRoomsAdd(room *Room) bool {
	lobby.freeRoomsM.Lock()
	defer lobby.freeRoomsM.Unlock()
	v := lobby._freeRooms.Add(room)
	return v
}

// allRoomsAdd add new room to all rooms slice
func (lobby *Lobby) allRoomsAdd(room *Room) bool {
	lobby.allRoomsM.Lock()
	defer lobby.allRoomsM.Unlock()
	v := lobby._allRooms.Add(room)
	return v
}
