package game

import (
	"fmt"
	"sync"
	"time"

	"escapade/internal/config"

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
	ws           *websocket.Conn
	Player       *Player `json:"player"`
	lobby        *Lobby
	room         *Room
	disconnected bool

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
}

// IsConnected check player isnt disconnected
func (conn *Connection) IsConnected() bool {
	return conn.disconnected == false
}

// IsPlayerAlive call player's IsAlive
func (conn *Connection) IsPlayerAlive() bool {
	return conn.Player.IsAlive()
}

func (conn *Connection) Kill(message string) {
	conn.SendInformation([]byte(message))
	conn.cancel()
}

// Free free memory, if flag disconnect true then connection and player will not become nil
func (conn *Connection) Free(disconnect bool) {
	if conn == nil {
		return
	}
	//conn.ws.Close()
	close(conn.send)
	if disconnect {
		conn.disconnected = true
	} else {
		conn.Player = nil //player doenst need Free()
		conn.lobby = nil
		conn.room = nil
		conn = nil
	}
}

// NewConnection creates a new connection
func NewConnection(ws *websocket.Conn, player *Player, lobby *Lobby) *Connection {
	return &Connection{
		ws:           ws,
		Player:       player,
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

// GetPlayerID get player id
func (conn *Connection) GetPlayerID() int {
	return conn.Player.ID
}

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
	go conn.WriteConn(ws, all, connContext)
	all.Add(1)
	go conn.ReadConn(ws, all, connContext)

	all.Wait()
	fmt.Println("conn finished")
	conn.lobby.chanLeave <- conn
	conn.Free(true)
}

// ReadConn connection goroutine to read messages from websockets
func (conn *Connection) ReadConn(wsc config.WebSocketSettings, wg *sync.WaitGroup, parent context.Context) {
	defer func() {
		wg.Done()
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
			conn.ws.Close()
			return
		default:
			_, message, err := conn.ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					fmt.Println("IsUnexpectedCloseError:" + err.Error())
				} else {
					fmt.Println("expected error:" + err.Error())
				}
				conn.Kill("Client websocket died")
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
func (conn *Connection) WriteConn(wsc config.WebSocketSettings, wg *sync.WaitGroup, parent context.Context) {
	ticker := time.NewTicker(wsc.PingPeriod)
	defer func() {
		wg.Done()
		ticker.Stop()
	}()
	for {
		select {
		case <-parent.Done():
			fmt.Println("WriteConn done catched")
			conn.ws.Close()
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

func (conn *Connection) sendGroupInformation(bytes []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	conn.SendInformation(bytes)
}

func (conn *Connection) debug(message string) {
	fmt.Println("Connection #", conn.GetPlayerID(), "-", message)
	conn.SendInformation([]byte(message))
}
