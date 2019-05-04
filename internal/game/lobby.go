package game

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// Request connect Connection and his message
type Request struct {
	Connection *Connection
	Message    []byte
}

// Lobby there are all rooms and users placed
type Lobby struct {
	wGroup *sync.WaitGroup

	doneM *sync.RWMutex
	_done bool

	allRoomsM *sync.RWMutex
	_AllRooms *Rooms `json:"allRooms"`

	freeRoomsM *sync.RWMutex
	_FreeRooms *Rooms `json:"freeRooms"`

	waitingM *sync.RWMutex
	_Waiting *Connections `json:"waiting"`

	playingM *sync.RWMutex
	_Playing *Connections `json:"playing"`

	messagesM *sync.Mutex
	_Messages []*models.Message `json:"messages"`

	context context.Context
	cancel  context.CancelFunc

	// connection joined lobby
	ChanJoin chan *Connection `json:"-"`
	// connection left lobby
	chanLeave chan *Connection

	chanBroadcast chan *Request

	chanBreak chan interface{}

	db            *database.DataBase
	canCloseRooms bool
}

// NewLobby create new instance of Lobby
func NewLobby(connectionsCapacity, roomsCapacity,
	maxJoin, maxRequest int, db *database.DataBase,
	canCloseRooms bool) *Lobby {

	messages, err := db.LoadMessages(false, "")
	if err != nil {
		fmt.Println("cant load messages:", err.Error())
	}
	context, cancel := context.WithCancel(context.Background())
	lobby := &Lobby{
		wGroup: &sync.WaitGroup{},

		doneM: &sync.RWMutex{},
		_done: false,

		allRoomsM: &sync.RWMutex{},
		_AllRooms: NewRooms(roomsCapacity),

		freeRoomsM: &sync.RWMutex{},
		_FreeRooms: NewRooms(roomsCapacity),

		waitingM: &sync.RWMutex{},
		_Waiting: NewConnections(connectionsCapacity),

		playingM: &sync.RWMutex{},
		_Playing: NewConnections(connectionsCapacity),

		messagesM: &sync.Mutex{},
		_Messages: messages,

		context: context,
		cancel:  cancel,

		ChanJoin:      make(chan *Connection),
		chanLeave:     make(chan *Connection),
		chanBroadcast: make(chan *Request),
		chanBreak:     make(chan interface{}),

		db:            db,
		canCloseRooms: canCloseRooms,
	}
	return lobby
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

		go lobby.Run(nil)
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

	if lobby.done() {
		return
	}
	lobby.setDone()

	go lobby.sendLobbyMessage("server closed", All)

	lobby.wGroup.Wait()

	fmt.Println("All resources clear!")

	go lobby.allRoomsFree()
	go lobby.freeRoomsFree()
	go lobby.waitingFree()
	go lobby.playingFree()

	lobby.cancel()

	close(lobby.ChanJoin)
	close(lobby.chanLeave)
	close(lobby.chanBroadcast)
	lobby.db = nil
	lobby = nil
}
