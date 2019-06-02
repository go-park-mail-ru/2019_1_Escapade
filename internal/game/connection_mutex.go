package game

import (
	"fmt"
	"time"
)

// setDone set done = true. It will finish all operaions on Connection
func (conn *Connection) setDone() {
	conn.doneM.Lock()
	conn._done = true
	conn.doneM.Unlock()
}

// done return '_done' field
func (conn *Connection) done() bool {
	if conn == nil {
		fmt.Println("conn nil")
	}
	conn.doneM.RLock()
	v := conn._done
	conn.doneM.RUnlock()
	return v
}

// Disconnected return   '_disconnected' field
func (conn *Connection) Disconnected() bool {
	if conn.done() {
		return conn._disconnected
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.disconnectedM.RLock()
	v := conn._disconnected
	conn.disconnectedM.RUnlock()
	return v
}

// setDisconnected set _disconnected true
func (conn *Connection) setDisconnected() {
	conn.disconnectedM.Lock()
	conn._disconnected = true
	conn.disconnectedM.Unlock()
	conn.time = time.Now()
}

// SetConnected set _disconnected false
func (conn *Connection) SetConnected() {
	if conn._disconnected && conn.InRoom() {
		_, isPlayer := conn.Room().Search(conn)
		if isPlayer {
			pa := *conn.Room().addAction(conn.ID(), ActionConnectAsPlayer)
			conn.Room().sendAction(pa, conn.Room().All)
			//conn.Room().sendPlayerEnter(*conn, conn.Room().All)
		} else {
			pa := *conn.Room().addAction(conn.ID(), ActionConnectAsObserver)
			conn.Room().sendAction(pa, conn.Room().All)
			//conn.Room().sendObserverEnter(*conn, conn.Room().All)
		}
	}
	conn.disconnectedM.Lock()
	conn._disconnected = false
	conn.disconnectedM.Unlock()
	//fmt.Println("!!!!!!!!!!!!!!!!!!!1connected", time.Now())
	conn.time = time.Now()
}

// Room return   '_room' field
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
// func (conn *Connection) RoomID() string {
// 	if conn.done() {
// 		return re.ErrorConnectionDone().Error()
// 	}
// 	conn.wGroup.Add(1)
// 	defer func() {
// 		conn.wGroup.Done()
// 	}()

// 	conn.roomM.RLock()
// 	v := conn._room.ID
// 	conn.roomM.RUnlock()
// 	return v
// }

// Both return   '_both' field
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

// Index return   '_index' field
func (conn *Connection) Index() int {
	if conn.done() {
		return conn._index
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.indexM.RLock()
	v := conn._index
	conn.indexM.RUnlock()
	return v
}

// SetIndex set '_index' - index in slice of players
func (conn *Connection) SetIndex(value int) {
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.indexM.Lock()
	conn._index = value
	conn.indexM.Unlock()
}

// setRoom set a pointer to the room in which the connection is located
func (conn *Connection) setRoom(room *Room) {
	conn.roomM.Lock()
	conn._room = room
	conn.roomM.Unlock()
}

// setBoth sets the flag whether the connection belongs to both the lobby and the room
func (conn *Connection) setBoth(both bool) {
	conn.bothM.Lock()
	conn._both = both
	conn.bothM.Unlock()
}
