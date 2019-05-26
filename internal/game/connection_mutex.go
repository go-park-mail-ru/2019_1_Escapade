package game

import (
	"fmt"

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
	if conn == nil {
		fmt.Println("conn nil")
	}
	conn.doneM.RLock()
	v := conn._done
	conn.doneM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (conn *Connection) Disconnected() bool {
	if conn.done() {
		return false
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

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
		return conn._room
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
		return conn._both
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

// getMatrixValue get a value from matrix
func (conn *Connection) Index() int {
	if conn.done() {
		return conn._Index
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.indexM.RLock()
	v := conn._Index
	conn.indexM.RUnlock()
	return v
}

func (conn *Connection) SetIndex(value int) {
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.indexM.Lock()
	conn._Index = value
	conn.indexM.Unlock()
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
