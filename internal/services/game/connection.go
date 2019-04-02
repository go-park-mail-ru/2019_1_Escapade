package game

import (
	"escapade/internal/models"
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
	player *models.Player
	lobby  *Lobby
	room   *Room
	Status int
}

// NewConnection creates a new connection and run it
func NewConnection(ws *websocket.Conn, player *models.Player, lobby *Lobby) *Connection {
	conn := &Connection{ws, player, lobby, nil, connectionLobby}
	go conn.run()
	return conn
}

func (conn *Connection) GetPlayerID() int {
	return conn.player.ID
}

func (conn *Connection) GiveUp() {
	conn.player.LastAction = models.ActionGiveUp
	conn.room.chanLeave <- conn
}

func (conn *Connection) run() {
	for {
		var cell *models.Cell
		err := conn.ws.ReadJSON(cell)
		if err != nil {
			fmt.Println("Error reading json.", err)
			return
		}
		conn.room.chanRequest <- NewRequest(conn, cell)
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

func (conn *Connection) SendInformation(info interface{}) {
	if err := conn.ws.WriteJSON(info); err != nil {
		fmt.Println(err)
	}
}

func (conn *Connection) sendGroupInformation(info interface{}, wg *sync.WaitGroup) {

	defer wg.Done()
	conn.SendInformation(info)
}
