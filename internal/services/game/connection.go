package game

import (
	"fmt"
	"sync"

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

	//chanRead   chan *Connection  `json:"-"`
	//chanWrite chan *RoomRequest `json:"-"`
}

// NewConnection creates a new connection and run it
func NewConnection(ws *websocket.Conn, player *Player, lobby *Lobby) *Connection {
	conn := &Connection{ws, player, lobby, nil, connectionLobby}
	go conn.run()
	return conn
}

// GetPlayerID get player id
func (conn *Connection) GetPlayerID() int {
	return conn.Player.ID
}

func (conn *Connection) lobbyWork() bool {
	var request = &LobbyRequest{}
	err := conn.ws.ReadJSON(request)
	if err != nil {
		fmt.Println("Error reading json.", err)
		return false
	}
	request.Connection = conn
	conn.lobby.chanRequest <- request
	return true
}

func (conn *Connection) roomWork() bool {
	var request = &RoomRequest{}
	err := conn.ws.ReadJSON(request)
	if err != nil {
		fmt.Println("Error reading json.", err)
		return false
	}
	conn.room.chanRequest <- request
	return true
}

func (conn *Connection) run() {
	for {
		if conn.Status == connectionLobby {
			if !conn.lobbyWork() {
				break
			}
		} else {
			if !conn.roomWork() {
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

// SendInformation send info
func (conn *Connection) SendInformation(info interface{}) {
	if err := conn.ws.WriteJSON(info); err != nil {
		fmt.Println(err)
	}
}

func (conn *Connection) sendGroupInformation(info interface{}, wg *sync.WaitGroup) {

	defer wg.Done()
	conn.SendInformation(info)
}

func (conn *Connection) debug(place string, channel string, function string, message string) {
	fmt.Println(conn.GetPlayerID(), place+"- channel:", channel, ". Handle by ", function, message)
}
