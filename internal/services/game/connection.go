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
	ws     *websocket.Conn
	Player *Player `json:"player"`
	lobby  *Lobby
	room   *Room
	Status int `json:"status"`

	send chan []byte

	//chanRead   chan *Connection  `json:"-"`
	//chanWrite chan *RoomRequest `json:"-"`
}

// NewConnection creates a new connection
func NewConnection(ws *websocket.Conn, player *Player, lobby *Lobby) *Connection {
	conn := &Connection{ws, player, lobby, nil, connectionLobby, make(chan []byte)}
	//go conn.run()
	return conn
}

// GetPlayerID get player id
func (conn *Connection) GetPlayerID() int {
	return conn.Player.ID
}

/*
func (conn *Connection) lobbyWork() bool {
	var request = &LobbyRequest{}
	err := conn.ws.ReadJSON(request)

	if err != nil {
		fmt.Println("Error reading json.", err)
		return false
	}
	conn.debug("lobbyWork", "lobbyWork", "lobbyWork", "lobbyWork")
	request.Connection = conn
	conn.lobby.chanRequest <- request
	conn.debug("lobbyWork done", "lobbyWork done", "lobbyWork done", "lobbyWork done")
	return true
}

// roomWork reed from websocket on
func (conn *Connection) roomWork() bool {
	var request = &RoomRequest{}

	err := conn.ws.ReadJSON(request)
	conn.debug("roomWork ReadJSON", "roomWork", "roomWork", "roomWork")
	if err != nil {
		fmt.Println("Error reading json.", err)
		return false

	}
	conn.debug("roomWork", "roomWork", "roomWork", "roomWork")
	request.Connection = conn
	conn.room.chanRequest <- request
	conn.debug("roomWork done", "roomWork done", "roomWork done", "roomWork done")
	return true
}
*/
// все в конфиг
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

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
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
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

/*
// run launch connection
func (conn *Connection) run() {
	for {
		//if conn.Status == connectionLobby {
		if !conn.lobbyWork() {
			break

		} else {
			if !(conn.roomWork()) {
				break
			}
		}
	}
	switch conn.Status {
	case connectionLobby:
		conn.lobby.chanLeave <- conn
	case connectionPlayer:
		conn.room.chanLeave <- conn
	case connectionObserver:
		conn.lobby.chanLeave <- conn
		conn.room.chanLeave <- conn
	}
	conn.ws.Close()
	return
}
*/
func (conn *Connection) WriteConn() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.ws.Close()
		conn.lobby.chanLeave <- conn
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
	conn.send <- bytes
}

func (conn *Connection) sendGroupInformation(bytes []byte, wg *sync.WaitGroup) {

	defer wg.Done()
	conn.SendInformation(bytes)
}

func (conn *Connection) debug(message string) {
	fmt.Println("Connection #", conn.GetPlayerID(), "-", message)
}
