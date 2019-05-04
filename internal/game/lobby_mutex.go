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
	lobby._AllRooms.Free()
	lobby.allRoomsM.Unlock()
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) freeRoomsFree() {
	lobby.freeRoomsM.Lock()
	lobby._FreeRooms.Free()
	lobby.freeRoomsM.Unlock()
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) waitingFree() {
	lobby.waitingM.Lock()
	lobby._Waiting.Free()
	lobby.waitingM.Unlock()
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) playingFree() {
	lobby.playingM.Lock()
	lobby._Playing.Free()
	lobby.playingM.Unlock()
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) waiting() []*Connection {
	lobby.waitingM.RLock()
	v := lobby._Waiting.Get
	lobby.waitingM.RUnlock()

	return v
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) playing() []*Connection {
	lobby.playingM.RLock()
	v := lobby._Playing.Get
	lobby.playingM.RUnlock()

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
	lobby._Playing.Remove(conn)
	lobby.playingM.Unlock()
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
	lobby._Playing.Add(conn, false)
	lobby.playingM.Unlock()
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
	lobby._Waiting.Remove(conn)
	lobby.waitingM.Unlock()
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
	lobby._Waiting.Add(conn, false)
	lobby.waitingM.Unlock()
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) allRoomsSearch(roomID string) (int, *Room) {
	lobby.allRoomsM.RLock()
	index, room := lobby._AllRooms.SearchRoom(roomID)
	lobby.allRoomsM.RUnlock()
	return index, room
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) allRoomsSearchPlayer(conn *Connection) (int, *Room) {
	lobby.allRoomsM.RLock()
	index, room := lobby._AllRooms.SearchPlayer(conn)
	lobby.allRoomsM.RUnlock()
	return index, room
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) allRoomsSearchObserver(conn *Connection) *Connection {
	lobby.allRoomsM.RLock()
	old := lobby._AllRooms.SearchObserver(conn)
	lobby.allRoomsM.RUnlock()
	return old
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) freeRoomsEmpty() bool {
	lobby.freeRoomsM.RLock()
	v := lobby._FreeRooms.Empty()
	lobby.freeRoomsM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) freeRooms() []*Room {
	lobby.freeRoomsM.RLock()
	v := lobby._FreeRooms.Get
	lobby.freeRoomsM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) setToMessages(message *models.Message) {
	lobby.messagesM.Lock()
	lobby._Messages = append(lobby._Messages, message)
	lobby.messagesM.Unlock()
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
	lobby._FreeRooms.Remove(room)
	lobby.freeRoomsM.Unlock()
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
	lobby._AllRooms.Remove(room)
	lobby.allRoomsM.Unlock()
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) freeRoomsAdd(room *Room) bool {
	lobby.freeRoomsM.Lock()
	v := lobby._FreeRooms.Add(room)
	lobby.freeRoomsM.Unlock()
	return v
}

// getMatrixValue get a value from matrix
func (lobby *Lobby) allRoomsAdd(room *Room) bool {
	lobby.allRoomsM.Lock()
	v := lobby._AllRooms.Add(room)
	lobby.allRoomsM.Unlock()
	return v
}
