package game

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Connection is a websocket of a player, that belongs to room
type Connection struct {
	ws     *websocket.Conn
	player *Player
	room   *Room
}

// NewConnection creates a new connection and run it
func NewConnection(ws *websocket.Conn, player *Player, room *Room) *Connection {
	conn := &Connection{ws, player, room}
	go conn.run()
	return conn
}

func (conn *Connection) run() (returnError error) {
	for {
		messageType, command, err := conn.ws.ReadMessage()
		if err != nil {
			returnError = err
			break
		}

		//var cell Cell

		// cell.Value = 10 - setFlag
		// otherwise       - openCell
		// _ = json.NewDecoder(command).Decode(&cell)

		// if (cell.Value == "10") {

		// }

		// // execute a command
		// con.player.Command(string(command)
		// // update all conn
		// con.room.updateAll <- true
	}
	if returnError != nil {
		return
	}
	conn.room.leave <- conn
	conn.ws.Close()
	return
}

func (conn *Connection) sendInformation(info interface{}) {

	err := conn.ws.ReadJSON(&info)
	if err != nil {
		fmt.Println("Error reading json.", err)
	}

	fmt.Printf("Got message: %#v\n", info)

	if err = conn.ws.WriteJSON(info); err != nil {
		fmt.Println(err)
	}
}

func (conn *Connection) sendGroupInformation(info interface{}, wg *sync.WaitGroup) {

	defer wg.Done()
	conn.sendInformation(info)
}
