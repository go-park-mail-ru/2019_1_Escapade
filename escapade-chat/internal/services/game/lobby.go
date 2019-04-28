package game

import (
	"context"
	"encoding/json"
	"escapade/internal/config"
	"escapade/internal/models"
	"fmt"
	"time"
)

// Request connect Connection and his message
type Request struct {
	Connection *Connection
	Message    []byte
}

// Lobby there are all rooms and users placed
type Lobby struct {
	Type string `json:"type,omitempty"`

	AllRooms  *Rooms `json:"allRooms,omitempty"`
	FreeRooms *Rooms `json:"freeRooms,omitempty"`

	Waiting *Connections `json:"waiting,omitempty"`
	Playing *Connections `json:"playing,omitempty"`

	Context context.Context `json:"-"`

	// connection joined lobby
	ChanJoin chan *Connection `json:"-"`
	// connection left lobby
	chanLeave chan *Connection
	//chanRequest   chan *LobbyRequest
	chanBroadcast chan *Request

	chanBreak chan interface{}

	semJoin    chan bool
	semRequest chan bool
}

// lobby singleton
var (
	lobby *Lobby
)

// Launch launchs lobby goroutine
func Launch(gc *config.GameConfig) {

	if lobby == nil {
		lobby = newLobby(gc.RoomsCapacity,
			gc.LobbyJoin, gc.LobbyRequest)

		go lobby.Run()
	}
}

// GetLobby create lobby if it is nil and get it
func GetLobby() *Lobby {
	return lobby
}

// Stop lobby goroutine
func (lobby *Lobby) Stop() {
	if lobby != nil {
		fmt.Println("Stop called!")
		lobby.chanBreak <- nil
	}
}

// Free delete all rooms and conenctions. Inform all players
// about closing
func (lobby *Lobby) Free() {
	if lobby == nil {
		return
	}
	fmt.Println("All resources clear!")
	SendToConnections("server closed", All, lobby.Waiting.Get, lobby.Playing.Get)

	lobby.AllRooms.Free()
	lobby.FreeRooms.Free()
	lobby.Waiting.Free()
	lobby.Playing.Free()
	close(lobby.ChanJoin)
	close(lobby.chanLeave)
	close(lobby.chanBroadcast)
	lobby = nil
}

// newLobby create new instance of Lobby
func newLobby(roomsCapacity, maxJoin, maxRequest int) *Lobby {

	connectionsCapacity := 500
	lobby := &Lobby{
		Type: "Lobby",

		AllRooms:  NewRooms(roomsCapacity),
		FreeRooms: NewRooms(roomsCapacity),

		Waiting: NewConnections(connectionsCapacity),
		Playing: NewConnections(connectionsCapacity),

		ChanJoin:      make(chan *Connection),
		chanLeave:     make(chan *Connection),
		chanBroadcast: make(chan *Request),
		chanBreak:     make(chan interface{}),

		semJoin:    make(chan bool, maxJoin),
		semRequest: make(chan bool, maxRequest),
	}
	return lobby
}

func NewMessage(req *Request) (message *models.Message) {
	jsonType := "Message"
	message = &models.Message{}
	err := json.Unmarshal(req.Message, &message)
	if err != nil {
		jsonType = "Error"
		message.Message = " cant send text"
	}
	message = &models.Message{
		Type:    jsonType,
		User:    req.Connection.User,
		Message: message.Message,
		Time:    time.Now(),
	}
	// send to db
	return message
}
