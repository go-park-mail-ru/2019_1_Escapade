package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"fmt"
	"sync"
	"time"

	"context"

	"github.com/gorilla/websocket"
)

// Connection status
const (
	connectionLobby     = iota // can ask for rooms list
	connectionRoomEnter        // can ask for field and people
	connectionPlayer           // can send cell and get update field
	connectionObserver         // get update field
)

// Connection is a websocket of a player, that belongs to room
type Connection struct {
	User *models.UserPublicInfo `json:"user,omitempty"`

	ws           *websocket.Conn
	lobby        *Lobby
	room         *Room
	disconnected bool `json:"disconnected,omitempty"`
	both         bool

	index int

	cancel context.CancelFunc
	send   chan []byte
}

// PushToRoom set field 'room' to real room
func (conn *Connection) PushToRoom(room *Room) {
	conn.room = room
}

// PushToLobby set field 'room' to nil
func (conn *Connection) PushToLobby() {
	conn.room = nil
	conn.both = false
}

// IsConnected check player isnt disconnected
func (conn *Connection) IsConnected() bool {
	return conn.disconnected == false
}

// dirty make connection dirty. it make connection ID
// -1 and when connection try to leave lobby, lobby will not
// delete this connections from list, cause it will not find
// anybody with such id
func (conn *Connection) dirty() {
	conn.User.ID = -1
}

// Kill call context.CancFunc, that finish goroutines of
// writer and reader and free connection memory
func (conn *Connection) Kill(message string, makeDirty bool) {
	conn.SendInformation([]byte(message))
	if makeDirty {
		conn.dirty()
	}
	conn.disconnected = true
	conn.cancel()
}

// Free free memory, if flag disconnect true then connection and player will not become nil
func (conn *Connection) Free() {
	if conn == nil {
		return
	}
	// dont delete. conn = nil make pointer nil, but other pointers
	// arent nil. If conn.disconnected = true it is mean that all
	// resources are cleared, but pointer alive, so we only make pointer = nil
	if conn.lobby == nil {
		conn = nil
		return
	}
	conn.ws.Close()
	close(conn.send)
	// dont delete. conn = nil make pointer nil, but other pointers
	// arent nil and we make 'conn.disconnected = true' for them
	conn.disconnected = true
	conn.lobby = nil
	conn.room = nil
	conn = nil

	fmt.Println("conn free memory")
}

// NewConnection creates a new connection
func NewConnection(ws *websocket.Conn, user *models.UserPublicInfo, lobby *Lobby) *Connection {
	return &Connection{
		ws:           ws,
		index:        -1,
		User:         user,
		lobby:        lobby,
		room:         nil,
		disconnected: false,
		send:         make(chan []byte),
	}
}

// InRoom check is player in room
func (conn *Connection) InRoom() bool {
	return conn.room != nil
}

// Launch run the writer and reader goroutines and wait them to free memory
func (conn *Connection) Launch(ws config.WebSocketSettings) {

	if lobby == nil {
		fmt.Println("lobby nil!")
		return
	}

	all := &sync.WaitGroup{}
	var connContext context.Context
	connContext, conn.cancel = context.WithCancel(lobby.Context)

	conn.lobby.ChanJoin <- conn
	all.Add(1)
	go conn.WriteConn(connContext, ws, all)
	all.Add(1)
	go conn.ReadConn(connContext, ws, all)

	all.Wait()
	fmt.Println("conn finished")
	conn.lobby.Leave(conn, "finished")
	conn.Free()
}

// ReadConn connection goroutine to read messages from websockets
func (conn *Connection) ReadConn(parent context.Context, wsc config.WebSocketSettings, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		utils.CatchPanic("connection.go ReadConn()")
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
				conn.Kill("Client websocket died", false)
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
	defer func() {
		wg.Done()
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
			if !ok {
				conn.write(websocket.CloseMessage, []byte{}, wsc)
				return
			}

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
func (conn *Connection) SendInformation(bytes []byte) {
	if !conn.disconnected {
		conn.send <- bytes
	}
}

// sendGroupInformation send info with WaitGroup
func (conn *Connection) sendGroupInformation(bytes []byte, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		utils.CatchPanic("connection.go sendGroupInformation()")
	}()
	conn.SendInformation(bytes)
}

// ID return players id
func (conn *Connection) ID() int {
	if conn.User == nil {
		return -1
	}
	return conn.User.ID
}

// debug print devug information to console and websocket
func (conn *Connection) debug(message string) {
	fmt.Println("Connection #", conn.ID(), "-", message)
	conn.SendInformation([]byte(message))
}
