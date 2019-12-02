package engine

import (
	"time"

	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
)

// Disconnected return   '_disconnected' field
func (conn *Connection) Disconnected() bool {
	var v bool
	conn.s.Do(func() {
		conn.disconnectedM.RLock()
		v = conn._disconnected
		conn.disconnectedM.RUnlock()
	})
	return v
}

// setDisconnected set _disconnected true
func (conn *Connection) setDisconnected() {
	conn.s.Do(func() {
		conn.disconnectedM.Lock()
		conn._disconnected = true
		conn.disconnectedM.Unlock()
		//conn.time = time.Now()
	})
}

// SetConnected set_disconnected false and update last connection time
func (conn *Connection) SetConnected() {
	conn.s.Do(func() {
		conn.disconnectedM.Lock()
		conn._disconnected = false
		conn.disconnectedM.Unlock()

		conn.timeM.Lock()
		conn._time = time.Now()
		conn.timeM.Unlock()
	})
}

// Time return the last time the connection sent a 'pong' message
func (conn *Connection) Time() time.Time {
	var v time.Time
	conn.s.Do(func() {
		conn.timeM.RLock()
		v = conn._time
		conn.timeM.RUnlock()
	})
	return v
}

// PlayingRoom return the pointer to the room in which the
// connection is playing
func (conn *Connection) PlayingRoom() *Room {
	var v *Room
	conn.s.Do(func() {
		conn.playingRoomM.RLock()
		v = conn._playingRoom
		conn.playingRoomM.RUnlock()
	})
	return v
}

// WaitingRoom return the pointer to the room in which the
// connection is waiting other players to connect
func (conn *Connection) WaitingRoom() *Room {
	var v *Room
	conn.s.Do(func() {
		conn.waitingRoomM.RLock()
		v = conn._waitingRoom
		conn.waitingRoomM.RUnlock()
	})
	return v
}

// Index return   '_index' field
func (conn *Connection) Index() int {
	var v int
	conn.s.Do(func() {
		conn.indexM.RLock()
		v = conn._index
		conn.indexM.RUnlock()
	})
	return v
}

// SetIndex set '_index' - index in slice of players
func (conn *Connection) SetIndex(value int) {
	conn.s.Do(func() {
		conn.indexM.Lock()
		conn._index = value
		conn.indexM.Unlock()
	})
}

// setRoom set a pointer to the room in which the connection is located
func (conn *Connection) setPlayingRoom(room *Room) {
	conn.s.Do(func() {
		conn.playingRoomM.Lock()
		conn._playingRoom = room
		conn.playingRoomM.Unlock()
	})
}

// setBoth sets the flag whether the connection belongs to both the lobby and the room
func (conn *Connection) setWaitingRoom(room *Room) {
	conn.s.Do(func() {
		conn.waitingRoomM.Lock()
		conn._waitingRoom = room
		conn.waitingRoomM.Unlock()
	})
}

// websocket

func (conn *Connection) wsInit(wsc config.WebSocket) {
	conn.s.Do(func() {
		conn.wsM.Lock()
		conn._ws.SetReadLimit(wsc.MaxMessageSize)
		conn._ws.SetReadDeadline(time.Now().Add(wsc.PongWait.Duration))
		conn._ws.SetPongHandler(
			func(string) error {
				conn._ws.SetReadDeadline(time.Now().Add(wsc.PongWait.Duration))
				conn.SetConnected()
				return nil
			})
		conn.wsM.Unlock()
	})
}

func (conn *Connection) wsReadMessage() (messageType int, p []byte, err error) {
	//conn.wsM.Lock()
	//defer conn.wsM.Unlock()
	messageType, p, err = conn._ws.ReadMessage()
	return
}

func (conn *Connection) wsWriteMessage(mt int, payload []byte, wsc config.WebSocket) error {
	conn.wsM.Lock()
	defer conn.wsM.Unlock()
	conn._ws.SetWriteDeadline(time.Now().Add(wsc.WriteWait.Duration))
	return conn._ws.WriteMessage(mt, payload)
}

func (conn *Connection) wsClose() error {
	conn.wsM.Lock()
	defer conn.wsM.Unlock()
	return conn._ws.Close()
}

func (conn *Connection) setWs(ws *websocket.Conn) {
	conn.wsM.Lock()
	defer conn.wsM.Unlock()
	conn._ws = ws
}

func (conn *Connection) wsWriteInWriter(message []byte, wsc config.WebSocket) error {
	conn.wsM.Lock()
	conn._ws.SetWriteDeadline(time.Now().Add(wsc.WriteWait.Duration))
	w, err := conn._ws.NextWriter(websocket.TextMessage)
	conn.wsM.Unlock()
	if err != nil {
		return err
	}
	w.Write(message)

	return w.Close()
}
