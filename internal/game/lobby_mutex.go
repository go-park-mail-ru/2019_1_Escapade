package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// setMatrixValue set a value to matrix
func (lobby *Lobby) setDone() {
	lobby.doneM.Lock()
	lobby._done = true
	lobby.doneM.Unlock()
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) done() bool {
	lobby.doneM.RLock()
	v := lobby._done
	lobby.doneM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) allRoomsFree() {
	lobby.allRoomsM.Lock()
	defer lobby.allRoomsM.Unlock()
	lobby._AllRooms.Free()
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) freeRoomsFree() {
	lobby.freeRoomsM.Lock()
	defer lobby.freeRoomsM.Unlock()
	lobby._FreeRooms.Free()
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) waitingFree() {
	lobby.waitingM.Lock()
	defer lobby.waitingM.Unlock()
	lobby._Waiting.Free()
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) playingFree() {
	lobby.playingM.Lock()
	defer lobby.playingM.Unlock()
	lobby._Playing.Free()
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) waiting() []*Connection {
	lobby.waitingM.RLock()
	defer lobby.waitingM.RUnlock()
	v := lobby._Waiting.Get

	return v
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) playing() []*Connection {
	lobby.playingM.RLock()
	defer lobby.playingM.RUnlock()
	v := lobby._Playing.Get

	return v
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) playingRemove(conn *Connection) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_mutex.go playingRemove()")
		lobby.wGroup.Done()
	}()

	lobby.playingM.Lock()
	defer lobby.playingM.Unlock()
	lobby._Playing.Remove(conn)
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) playingAdd(conn *Connection) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_mutex.go playingRemove()")
		lobby.wGroup.Done()
	}()

	lobby.playingM.Lock()
	defer lobby.playingM.Unlock()
	lobby._Playing.Add(conn, false)
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) waitingRemove(conn *Connection) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_mutex.go waitingRemove()")
		lobby.wGroup.Done()
	}()

	lobby.waitingM.Lock()
	defer lobby.waitingM.Unlock()
	lobby._Waiting.Remove(conn)
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) waitingAdd(conn *Connection) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_mutex.go waitingAdd()")
		lobby.wGroup.Done()
	}()

	lobby.waitingM.Lock()
	defer lobby.waitingM.Unlock()
	lobby._Waiting.Add(conn, false)
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) allRoomsSearch(roomID string) (int, *Room) {
	lobby.allRoomsM.RLock()
	defer lobby.allRoomsM.RUnlock()
	index, room := lobby._AllRooms.SearchRoom(roomID)
	return index, room
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) allRoomsSearchPlayer(conn *Connection) (int, *Room) {
	lobby.allRoomsM.RLock()
	defer lobby.allRoomsM.RUnlock()
	index, room := lobby._AllRooms.SearchPlayer(conn)
	return index, room
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) allRoomsSearchObserver(conn *Connection) *Connection {
	lobby.allRoomsM.RLock()
	defer lobby.allRoomsM.RUnlock()
	old := lobby._AllRooms.SearchObserver(conn)
	return old
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) freeRoomsEmpty() bool {
	lobby.freeRoomsM.RLock()
	defer lobby.freeRoomsM.RUnlock()
	v := lobby._FreeRooms.Empty()
	return v
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) freeRooms() []*Room {
	lobby.freeRoomsM.RLock()
	defer lobby.freeRoomsM.RUnlock()
	v := lobby._FreeRooms.Get
	return v
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) setToMessages(message *models.Message) {
	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	lobby._Messages = append(lobby._Messages, message)
}

// getMatrixValue get a value from matrix
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
	lobby._FreeRooms.Remove(room)
}

// getMatrixValue get a value from matrix
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
	lobby._AllRooms.Remove(room)
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) freeRoomsAdd(room *Room) bool {
	lobby.freeRoomsM.Lock()
	defer lobby.freeRoomsM.Unlock()
	v := lobby._FreeRooms.Add(room)
	return v
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) allRoomsAdd(room *Room) bool {
	lobby.allRoomsM.Lock()
	defer lobby.allRoomsM.Unlock()
	v := lobby._AllRooms.Add(room)
	return v
}
