package game

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// SetImage function to set image
type SetImage func(users ...*models.UserPublicInfo) (err error)

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

	allRooms  *Rooms
	freeRooms *Rooms

	Waiting *Connections
	Playing *Connections

	messagesM *sync.Mutex
	_messages []*models.Message

	anonymousM *sync.Mutex
	_anonymous int

	context context.Context
	cancel  context.CancelFunc

	// connection joined lobby
	chanJoin chan *Connection
	// connection send some JSON
	chanBroadcast chan *Request

	chanBreak chan interface{}

	db *database.DataBase

	config *config.GameConfig

	SetImage SetImage
}

// NewLobby create new instance of Lobby
func NewLobby(config *config.GameConfig, db *database.DataBase,
	SetImage SetImage) *Lobby {

	var (
		messages []*models.Message
		err      error
	)
	if db != nil {
		if messages, err = db.LoadMessages(false, ""); err != nil {
			fmt.Println("cant load messages:", err.Error())
		}
		for _, message := range messages {
			SetImage(message.User)
		}
	} else {
		messages = make([]*models.Message, 0)
	}
	context, cancel := context.WithCancel(context.Background())
	lobby := &Lobby{
		wGroup: &sync.WaitGroup{},

		doneM: &sync.RWMutex{},
		_done: false,

		config: config,

		allRooms:  NewRooms(config.RoomsCapacity),
		freeRooms: NewRooms(config.RoomsCapacity),

		Waiting: NewConnections(config.ConnectionCapacity),
		Playing: NewConnections(config.ConnectionCapacity),

		messagesM: &sync.Mutex{},
		_messages: messages,

		anonymousM: &sync.Mutex{},
		_anonymous: -1,

		context: context,
		cancel:  cancel,

		chanJoin:      make(chan *Connection),
		chanBroadcast: make(chan *Request),
		chanBreak:     make(chan interface{}),

		db:       db,
		SetImage: SetImage,
	}
	return lobby
}

// lobby singleton
var (
	LOBBY *Lobby
)

// Launch launchs lobby goroutine
func Launch(gc *config.GameConfig, db *database.DataBase, si SetImage) {

	if LOBBY == nil {
		LOBBY = NewLobby(gc, db, si)
		go LOBBY.Run()
		//LOBBY.stress(10)
	}
}

// GetLobby create lobby if it is nil and get it
func GetLobby() *Lobby {
	return LOBBY
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

	go lobby.allRooms.Free()
	go lobby.freeRooms.Free()
	go lobby.Waiting.Free()
	go lobby.Playing.Free()

	lobby.cancel()

	close(lobby.chanJoin)
	close(lobby.chanBroadcast)
	lobby.db = nil
	lobby = nil
}
