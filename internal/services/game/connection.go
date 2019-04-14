package game

import (
	"fmt"
	"log"
	"sync"
	"time"

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

	send chan []byte

	//chanRead   chan *Connection  `json:"-"`
	//chanWrite chan *RoomRequest `json:"-"`
}

func (conn *Connection) PushToRoom(room *Room) {
	conn.room = room
}

func (conn *Connection) PushToLobby() {
	conn.room = nil
}

func (conn *Connection) IsPlayerAlive() bool {
	return conn.Player.IsAlive()
}

// Kill send last image signals about killing, close websocket
// and close chanells
func (conn *Connection) Kill(message []byte) {
	conn.disconnected = true
	conn.SendInformation(message)
	conn.ws.Close()
	fmt.Println("killed with message:" + string(message))
	/*
		need some time before close. Maybe set timer?
	*/
	close(conn.send)
}

// NewConnection creates a new connection
func NewConnection(ws *websocket.Conn, player *Player, lobby *Lobby) *Connection {
	return &Connection{
		ws,
		player,
		lobby,
		nil,
		false,
		make(chan []byte),
	}
}

// NewConnection check is player in room
func (conn *Connection) InRoom() bool {

	return conn.room != nil
}

// GetPlayerID get player id
func (conn *Connection) GetPlayerID() int {
	return conn.Player.ID
}

// все в конфиг
const (
	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (conn *Connection) ReadConn() {
	defer func() {
		conn.lobby.chanLeave <- conn
		conn.ws.Close()
	}()
	conn.ws.SetReadLimit(maxMessageSize)
	conn.ws.SetReadDeadline(time.Now().Add(pongWait))
	conn.ws.SetPongHandler(
		func(string) error {
			conn.ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
	for {
		_, message, err := conn.ws.ReadMessage()
		conn.debug("read from conn")
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			conn.debug("error found")
			break
		}
		conn.lobby.chanBroadcast <- &Request{
			Connection: conn,
			Message:    message,
		}
	}
}

// write writes a message with the given message type and payload.
func (conn *Connection) write(mt int, payload []byte) error {
	conn.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.ws.WriteMessage(mt, payload)
}

// dont put conn.debug here
func (conn *Connection) WriteConn() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.ws.Close()
	}()
	for {
		select {
		// send here json!
		case message, ok := <-conn.send:
			if !ok {
				conn.write(websocket.CloseMessage, []byte{})
				return
			}

			conn.ws.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := conn.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			var newline = []byte{'\n'}
			n := len(conn.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-conn.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := conn.write(websocket.PingMessage, []byte{}); err != nil {
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
	Answer(conn, []byte(message))
}
