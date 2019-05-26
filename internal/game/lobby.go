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
	_AllRooms *Rooms

	freeRoomsM *sync.RWMutex
	_FreeRooms *Rooms

	waitingM *sync.RWMutex
	_Waiting *Connections

	playingM *sync.RWMutex
	_Playing *Connections

	messagesM *sync.Mutex
	_Messages []*models.Message

	context context.Context
	cancel  context.CancelFunc

	// connection joined lobby
	chanJoin chan *Connection
	// connection left lobby
	chanLeave chan *Connection

	chanBroadcast chan *Request

	chanBreak chan interface{}

	db            *database.DataBase
	canCloseRooms bool
	metrics       bool
}

// NewLobby create new instance of Lobby
func NewLobby(connectionsCapacity, roomsCapacity int,
	db *database.DataBase, canCloseRooms bool, metrics bool) *Lobby {

	var (
		messages []*models.Message
		err      error
	)
	if db != nil {
		if messages, err = db.LoadMessages(false, ""); err != nil {
			fmt.Println("cant load messages:", err.Error())
		}
	} else {
		messages = make([]*models.Message, 0)
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

		chanJoin:      make(chan *Connection),
		chanLeave:     make(chan *Connection),
		chanBroadcast: make(chan *Request),
		chanBreak:     make(chan interface{}),

		db:            db,
		canCloseRooms: canCloseRooms,
		metrics:       metrics,
	}
	return lobby
}

// lobby singleton
var (
	lobby *Lobby
)

// Launch launchs lobby goroutine
func Launch(gc *config.GameConfig, db *database.DataBase, metrics bool) {

	if lobby == nil {
		lobby = NewLobby(gc.ConnectionCapacity, gc.RoomsCapacity,
			db, gc.CanClose, metrics)

		go lobby.Run()
	}
}

// GetLobby create lobby if it is nil and get it
func GetLobby() *Lobby {
	return lobby
}

func (lobby *Lobby) Metrics() bool {
	return lobby.metrics
}

// Stop lobby goroutine
func (lobby *Lobby) Stop() {
	if lobby != nil {
		fmt.Println("Stop called!")
		lobby.chanBreak <- nil
	}
}

// Free delete all rooms and connections. Inform all players
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

	close(lobby.chanJoin)
	close(lobby.chanLeave)
	close(lobby.chanBroadcast)
	lobby.db = nil
	lobby = nil
}
