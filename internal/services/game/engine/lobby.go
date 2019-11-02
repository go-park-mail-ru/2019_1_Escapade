package engine

import (
	"context"
	"sync"
	"time"

	config "github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	database "github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	models "github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat"
	utils "github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

/*
SetImage - a function that accesses the photo service. It passes to the service
 the name of the photo to be specified in the PhotoURL field of the
 UserPublicInfo structure. In response, it gets a key to access the image it is
 looking for and places key in the FileKey field.
	The entrance accepts an unlimited number of users. The above action applies
	to each of them. If for some user it will not be possible to get an
	image(an error will be returned), the function will continue to work,
	skipping this user. Does not return errors even if they have occurred
*/
type SetImage func(users ...*models.UserPublicInfo)

/*
Request - Structure for messaging with websocket connections

	Connection - message sender/receiver.

	Message - the information to be sent.
*/
type Request struct {
	Connection *Connection
	Message    []byte
}

/*
MessageWithAction  - structure required to save actions related to messages in
 case they cannot be added to the database

 	origin - a pointer to an item in the lobby(or room) message array.
		If the user has created a message, but it could not get into the database,
		it is assigned a random ID(instead of getting the ID from the database).
		With it you can find this message to perform operations on it(update,
		delete). After the message is delivered to the database, its identifier
		is updated in accordance with the specified in the database(set new origin.id)

	message - proto-struct, that is sent to chat service

	action - defines the operation to be performed on the message(create,
		edit, delete). All available types are in package 'models'

	getChatID - function to get the chat ID(if the corresponding message field
		is not initialized)
*/
type MessageWithAction struct {
	origin    *models.Message
	message   *chat.Message
	action    int32
	getChatID func() (int32, error)
}

/*
Lobby - the structure that controls the slices of players and rooms.

	wGroup - any goroutine that is associated with the Lobby at the beginning
	 of execution incremented counter of wGroup and before ending decrements it.
	 A function to clear memory (which will be invoked in the event of
	 destruction of the Lobby) will wait until wGroup will not become equal to 0.
	 This ensures that all associated with this structure goroutines will have
	 time to finish their job

	_done - determines whether the memory cleanup function of the object has
	 been called. Each goroutine before started check this field. If done is
	 true, the goroutine terminates. Otherwise, increments the wGroup counter
	 described above

	doneM - provides exclusive access to the field _done. Available actions
	 you can see in lobby_mutex.go

	allRooms - slice of all active rooms(rooms whose gameplay is not over)
	freeRooms - slice of free rooms(rooms that search for players)

	Waiting - an slice of connections located in the lobby(IMPORTANT: the
	 connection can be both in the lobby and in the rooms. This is possible
	 if the connection is in a room that hosts a set of players. This
	 connection will also be in this slice)

	Playing - an slice of connections located only in room(connections with any
	 status are meant: both players, and observers)

	_messages - an slice of messages from the lobby chat

	messagesM - provides exclusive access to the field _messages. Available actions
	 you can see in lobby_mutex.go

	_notSavedMessages - an slice of actions related to messages that have not
	 yet been stored in the database. This field is protected by a
	 mutex(notSavedMessagesM)

	__notSavedGames - an slice of games that have not yet been stored in the
	 database. This field is protected by a mutex(notSavedGamesM)

	_anonymous - the ID of the next anonymous connection. To distinguish
	 anonymous connections from the non-anonumous connections for anonymous
	 connections are specified by negative numbers(at the time, as identicator
	 compounds registered users are defined as the ID of the corresponding user
	 from the database). This field is protected by a mutex(anonymousM), so to
	 get this field you need to call the anonymous getter function, which will
	 return the value of this field and decrement the value in this field. This
	 ensures the uniqueness of the IDs of anonymous connections(as long as the
	 object of the Lobby will not be rebuilt). The uniqueness of the identifiers
	 is required for the correct operation of the games, where players are
	 identified with the help of user identifier.

	context, cancel - TODO delete or normally implement

	chanJoin - the channel in which the connections joining the lobby are
	 transmitted.

	chanJoin - the channel for transmitting information from connections (Not
	 for writing to connections. To write to a connection, use the channel of
	 the connection to which you want to send information)

	chanBreak - the channel transmitting the signal of completion of the work

	_db - instance of database struct. This field is protected by a
	 mutex(dbM)

	_config - configuration of game. This field is protected by a
	 mutex(configM)

	_dbChatID - the chat ID(obtained from the database), that is used by the
	 Lobby. It determines in which chat, operations on messages(recording,
	 editing, deleting) will be saved. It is set during initialization. If
	 install failed(because of the inability to connect to the chat service),
	 will be installed later in the goroutine sending unsent actions in the chat
	 service. This field is protected by a mutex(dbChatIDM)

	_location - time.Location. It is necessary to set the time zone when
	 specifying dates(in rooms, messages, actions, etc.). It is set during
	 initialization. This field is protected by a mutex(locationM)

	SetImage - instance of SetImage, described above. It is necessary to obtain
	 user images for chat
*/
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

	notSavedMessagesM *sync.Mutex
	_notSavedMessages []*MessageWithAction

	notSavedGamesM *sync.Mutex
	_notSavedGames []*models.GameInformation

	anonymousM *sync.Mutex
	_anonymous int32

	context context.Context
	cancel  context.CancelFunc

	chanJoin      chan *Connection
	chanBroadcast chan *Request
	chanBreak     chan interface{}

	dbM *sync.RWMutex
	_db *database.DataBase

	configM *sync.RWMutex
	_config *config.Game

	dbChatIDM *sync.RWMutex
	_dbChatID int32

	locationM *sync.RWMutex
	_location *time.Location

	SetImage SetImage
}

// NewLobby create new instance of Lobby
func NewLobby(config *config.Game, db *database.DataBase,
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

		notSavedMessagesM: &sync.Mutex{},
		_notSavedMessages: make([]*MessageWithAction, 0),

		notSavedGamesM: &sync.Mutex{},
		_notSavedGames: make([]*models.GameInformation, 0),

		chanJoin:      make(chan *Connection),
		chanBroadcast: make(chan *Request),
		chanBreak:     make(chan interface{}),
	}
	lobby.SetConfiguration(config, db, SetImage)
	return lobby
}

// SetConfiguration set lobby configuration
func (lobby *Lobby) SetConfiguration(config *config.Game, db *database.DataBase,
	setImage SetImage) {

	var (
		err      error
		chatID   int32
		location *time.Location
	)
	location, err = time.LoadLocation(config.Location)
	if err != nil {
		utils.Debug(true, "cant set location!")
	}
	lobby.setMessages(make([]*models.Message, 0))
	lobby.setConfig(config)
	lobby.setDB(db)
	lobby.setDBChatID(chatID)
	lobby.setLocation(location)
	lobby.SetImage = setImage

	return
}

// lobby singleton
var (
	LOBBY *Lobby
)

// Launch launchs lobby goroutine
func Launch(gc *config.Game, db *database.DataBase, si SetImage) {

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

// Free clean the memory allocated to the structure of the lobby
func (lobby *Lobby) Free() {

	if lobby.checkAndSetCleared() {
		return
	}

	groupWaitTimeout := 80 * time.Second // TODO в конфиг
	utils.WaitWithTimeout(lobby.wGroup, groupWaitTimeout)

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
	lobby.db().Db.Close()
	lobby.setDB(nil)
	lobby.setLocation(nil)
	lobby = nil
}
