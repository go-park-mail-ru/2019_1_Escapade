package game

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
)

// Request connect Connection and his message
type Request struct {
	Connection *Connection
	Message    []byte
}

// Lobby there are all rooms and users placed
type Lobby struct {
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

	db            *database.DataBase
	canCloseRooms bool

	semJoin    chan bool
	semRequest chan bool
}

// lobby singleton
var (
	lobby *Lobby
)

// Launch launchs lobby goroutine
func Launch(gc *config.GameConfig, db *database.DataBase) {

	if lobby == nil {
		lobby = NewLobby(gc.ConnectionCapacity, gc.RoomsCapacity,
			gc.LobbyJoin, gc.LobbyRequest, db, gc.CanClose)

		go lobby.Run(false)
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
	lobby.db = nil
	lobby = nil
}

// NewLobby create new instance of Lobby
func NewLobby(connectionsCapacity, roomsCapacity,
	maxJoin, maxRequest int, db *database.DataBase,
	canCloseRooms bool) *Lobby {

	//connectionsCapacity := 500
	lobby := &Lobby{

		AllRooms:  NewRooms(roomsCapacity),
		FreeRooms: NewRooms(roomsCapacity),

		Waiting: NewConnections(connectionsCapacity),
		Playing: NewConnections(connectionsCapacity),

		ChanJoin:      make(chan *Connection),
		chanLeave:     make(chan *Connection),
		chanBroadcast: make(chan *Request),
		chanBreak:     make(chan interface{}),

		db:            db,
		canCloseRooms: canCloseRooms,

		semJoin:    make(chan bool, maxJoin),
		semRequest: make(chan bool, maxRequest),
	}
	return lobby
}
