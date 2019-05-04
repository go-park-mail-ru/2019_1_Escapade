package game

import (
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

// setMatrixValue set a value to matrix
func (conn *Connection) setDone() {
	conn.doneM.Lock()
	conn._done = true
	conn.doneM.Unlock()
}

// getMatrixValue get a value from matrix
func (conn *Connection) done() bool {
	conn.doneM.RLock()
	v := conn._done
	conn.doneM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (conn *Connection) disconnected() bool {
	conn.disconnectedM.RLock()
	v := conn._Disconnected
	conn.disconnectedM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (conn *Connection) setDisconnected() {
	conn.disconnectedM.Lock()
	conn._Disconnected = true
	conn.disconnectedM.Unlock()
}

// getMatrixValue get a value from matrix
func (conn *Connection) Room() *Room {
	if conn.done() {
		return nil
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.roomM.RLock()
	v := conn._room
	conn.roomM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (conn *Connection) RoomID() string {
	if conn.done() {
		return re.ErrorConnectionDone().Error()
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.roomM.RLock()
	v := conn._room.ID
	conn.roomM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (conn *Connection) Both() bool {
	if conn.done() {
		return false
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.bothM.RLock()
	v := conn._both
	conn.bothM.RUnlock()
	return v
}

// setMatrixValue set a value to matrix
func (conn *Connection) setRoom(room *Room) {
	conn.roomM.Lock()
	conn._room = room
	conn.roomM.Unlock()
}

// setMatrixValue set a value to matrix
func (conn *Connection) setBoth(both bool) {
	conn.bothM.Lock()
	conn._both = both
	conn.bothM.Unlock()
}
