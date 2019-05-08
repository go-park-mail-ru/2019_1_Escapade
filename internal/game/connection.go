package game

import (
	"encoding/json"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"fmt"
	"sync"
	"time"

	"context"

	"github.com/gorilla/websocket"
)

// NewConnection creates a new connection
func NewConnection(ws *websocket.Conn, user *models.UserPublicInfo, lobby *Lobby) *Connection {
	if ws == nil || user == nil || lobby == nil {
		return nil
	}

	context, cancel := context.WithCancel(lobby.context)

	return &Connection{
		wGroup: &sync.WaitGroup{},

		doneM: &sync.RWMutex{},
		_done: false,

		roomM: &sync.RWMutex{},
		_room: nil,

		disconnectedM: &sync.RWMutex{},
		_Disconnected: false,

		bothM: &sync.RWMutex{},
		_both: false,

		indexM: &sync.RWMutex{},
		_Index: -1,

		User: user,

		ws:    ws,
		lobby: lobby,

		context: context,
		cancel:  cancel,

		send:      make(chan []byte),
		actionSem: make(chan struct{}, 1),
	}
}

// PushToRoom set field 'room' to real room
func (conn *Connection) PushToRoom(room *Room) {
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.setRoom(room)
}

// PushToLobby set field 'room' to nil
func (conn *Connection) PushToLobby() {
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	conn.setRoom(nil)
	conn.setBoth(false)
}

// IsConnected check player isnt disconnected
func (conn *Connection) IsConnected() bool {
	if conn.done() {
		return false
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()
	return conn.Disconnected() == false
}

// dirty make connection dirty. it make connection ID
// -1 and when connection try to leave lobby, lobby will not
// delete this connections from list, cause it will not find
// anybody with such id
func (conn *Connection) Dirty() {
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()
	conn.User.ID = -1
}

// Kill call context.CancFunc, that finish goroutines of
// writer and reader and free connection memory
func (conn *Connection) Kill(message string, makeDirty bool) {
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	fmt.Println("SendInformation")
	conn.SendInformation(message)
	if makeDirty {
		conn.Dirty()
	}
	fmt.Println("setDisconnected")
	conn.setDisconnected()
	fmt.Println("cancel")
	conn.cancel()
	fmt.Println("done")
}

// Free free memory, if flag disconnect true then connection and player will not become nil
func (conn *Connection) Free() {

	if conn.done() {
		return
	}
	conn.setDone()

	conn.wGroup.Wait()

	// dont delete. conn = nil make pointer nil, but other pointers
	// arent nil. If conn.disconnected = true it is mean that all
	// resources are cleared, but pointer alive, so we only make pointer = nil
	if conn.lobby == nil {
		return
	}

	conn.setDisconnected()

	conn.ws.Close()
	close(conn.send)
	close(conn.actionSem)
	// dont delete. conn = nil make pointer nil, but other pointers
	// arent nil and we make 'conn.disconnected = true' for them

	conn.lobby = nil
	conn.setRoom(nil)

	fmt.Println("conn free memory")
}

// InRoom check is player in room
func (conn *Connection) InRoom() bool {
	if conn.done() {
		return false
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()
	return conn.Room() != nil
}

// Launch run the writer and reader goroutines and wait them to free memory
func (conn *Connection) Launch(ws config.WebSocketSettings) {

	// dont place there conn.wGroup.Add(1)
	if conn.lobby == nil || conn.lobby.context == nil {
		fmt.Println("lobby nil or hasnt context!")
		return
	}

	all := &sync.WaitGroup{}

	fmt.Println("JoinConn!")
	conn.lobby.JoinConn(conn)
	all.Add(1)
	go conn.WriteConn(conn.context, ws, all)
	all.Add(1)
	go conn.ReadConn(conn.context, ws, all)

	fmt.Println("Wait!")
	all.Wait()
	fmt.Println("conn finished")
	conn.lobby.Leave(conn, "finished")
	conn.Free()
}

// ReadConn connection goroutine to read messages from websockets
func (conn *Connection) ReadConn(parent context.Context, wsc config.WebSocketSettings, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
		utils.CatchPanic("connection.go WriteConn()")
	}()

	conn.ws.SetReadLimit(wsc.MaxMessageSize)
	conn.ws.SetReadDeadline(time.Now().Add(wsc.PongWait))
	conn.ws.SetPongHandler(
		func(string) error {
			conn.ws.SetReadDeadline(time.Now().Add(wsc.PongWait))
			return nil
		})
	for {
		select {
		case <-parent.Done():
			fmt.Println("ReadConn done catched")
			return
		default:
			_, message, err := conn.ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					fmt.Println("IsUnexpectedCloseError:" + err.Error())
				} else {
					fmt.Println("expected error:" + err.Error())
				}
				//conn.Kill("Client websocket died", false)
				return
			}
			conn.debug("read from conn")
			conn.lobby.chanBroadcast <- &Request{
				Connection: conn,
				Message:    message,
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (conn *Connection) write(mt int, payload []byte, wsc config.WebSocketSettings) error {
	conn.ws.SetWriteDeadline(time.Now().Add(wsc.WriteWait))
	return conn.ws.WriteMessage(mt, payload)
}

// WriteConn connection goroutine to write messages to websockets
// dont put conn.debug here
func (conn *Connection) WriteConn(parent context.Context, wsc config.WebSocketSettings, wg *sync.WaitGroup) {
	defer wg.Done()

	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
		utils.CatchPanic("connection.go WriteConn()")
	}()

	ticker := time.NewTicker(wsc.PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-parent.Done():
			fmt.Println("WriteConn done catched")
			return
		case message, ok := <-conn.send:

			fmt.Println("saw!")
			//fmt.Println("server wrote:", string(message))
			if !ok {
				fmt.Println("errrrrr!")
				conn.write(websocket.CloseMessage, []byte{}, wsc)
				return
			}
			fmt.Println("send something")

			conn.ws.SetWriteDeadline(time.Now().Add(wsc.WriteWait))
			w, err := conn.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := conn.write(websocket.PingMessage, []byte{}, wsc); err != nil {
				return
			}
		}
	}
}

// SendInformation send info
func (conn *Connection) SendInformation(value interface{}) {
	if conn.done() {
		return
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()

	if !conn.Disconnected() {
		var (
			bytes []byte
			err   error
		)

		bytes, err = json.Marshal(value)

		if err != nil {
			fmt.Println("cant send information", err.Error())
		} else {
			fmt.Println("server wrote to", conn.ID(), ":", string(bytes))
			conn.send <- bytes
			fmt.Println("move!")
		}
	}
}

// sendGroupInformation send info with WaitGroup
func (conn *Connection) sendGroupInformation(value interface{}, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		utils.CatchPanic("connection.go sendGroupInformation()")
	}()
	conn.SendInformation(value)
}

// ID return players id
func (conn *Connection) ID() int {
	if conn.done() {
		return conn.User.ID
	}
	conn.wGroup.Add(1)
	defer func() {
		conn.wGroup.Done()
	}()
	if conn.User == nil {
		return -1
	}
	return conn.User.ID
}

// debug print devug information to console and websocket
func (conn *Connection) debug(message string) {
	fmt.Println("Connection #", conn.ID(), "-", message)
	conn.SendInformation(message)
}
