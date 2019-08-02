package game

import (
	"context"
	"sync"
	"time"

	сhat "github.com/go-park-mail-ru/2019_1_Escapade/chat/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// SetImage function to set image
type SetImage func(users ...*models.UserPublicInfo)

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

	dbM *sync.RWMutex
	_db *database.DataBase

	configM *sync.RWMutex
	_config *config.GameConfig

	dbChatIDM *sync.RWMutex
	_dbChatID int32

	locationM *sync.RWMutex
	_location *time.Location

	SetImage SetImage
}

// NewLobby create new instance of Lobby
func NewLobby(config *config.GameConfig, db *database.DataBase,
	SetImage SetImage) *Lobby {

	context, cancel := context.WithCancel(context.Background())
	lobby := &Lobby{
		wGroup: &sync.WaitGroup{},

		doneM: &sync.RWMutex{},
		_done: false,

		allRooms:  NewRooms(config.RoomsCapacity),
		freeRooms: NewRooms(config.RoomsCapacity),

		Waiting: NewConnections(config.ConnectionCapacity),
		Playing: NewConnections(config.ConnectionCapacity),

		messagesM: &sync.Mutex{},
		_messages: make([]*models.Message, 0),

		anonymousM: &sync.Mutex{},
		_anonymous: -1,

		context: context,
		cancel:  cancel,

		dbM:       &sync.RWMutex{},
		configM:   &sync.RWMutex{},
		dbChatIDM: &sync.RWMutex{},
		locationM: &sync.RWMutex{},

		chanJoin:      make(chan *Connection),
		chanBroadcast: make(chan *Request),
		chanBreak:     make(chan interface{}),
	}
	lobby.SetConfiguration(config, db, SetImage)
	return lobby
}

func (lobby *Lobby) SetConfiguration(config *config.GameConfig, db *database.DataBase,
	SetImage SetImage) {

	var (
		messages []*models.Message
		err      error
		chatID   int32
		location *time.Location
	)
	location, err = time.LoadLocation(config.Location)
	if err != nil {
		utils.Debug(true, "cant set location!")
	}
	if db != nil {
		chatID, messages, err = GetChatIDAndMessages(location, сhat.ChatType_LOBBY, 0)
		if err != nil {
			utils.Debug(true, "cant load messages:", err.Error())
		}
		for _, message := range messages {
			SetImage(message.User)
		}
	} else {
		messages = make([]*models.Message, 0)
	}
	lobby.setConfig(config)
	lobby.setMessages(messages)
	lobby.setDB(db)
	lobby.setDBChatID(chatID)
	lobby.setLocation(location)
	lobby.SetImage = SetImage

	return
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
		utils.Debug(false, "stop called!")
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

	utils.Debug(false, "All resources clear!")

	go lobby.allRooms.Free()
	go lobby.freeRooms.Free()
	go lobby.Waiting.Free()
	go lobby.Playing.Free()

	lobby.cancel()

	close(lobby.chanJoin)
	close(lobby.chanBroadcast)
	lobby.setConfig(nil)
	lobby.setMessages(nil)
	lobby.setDB(nil)
	lobby.setLocation(nil)
	lobby = nil
}
