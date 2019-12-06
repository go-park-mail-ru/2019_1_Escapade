package engine

import (
	"math/rand"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
	ctypes "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/database"
	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
)

// MessagesI control access to messages
// Proxy Pattern
type MessagesI interface {
	synced.PublisherI

	Fix(message *models.Message, conn *Connection)
	Proto(message *models.Message) (*chat.Message, error)

	Write(message *models.Message, send *chat.Message) error
	Update(message *models.Message, send *chat.Message) error
	Delete(message *models.Message, send *chat.Message) error

	Send(message *models.Message)
	HandleError(message *models.Message, send *chat.Message)

	setMessages(messages []*models.Message)
	Messages() []*models.Message

	ChatID() int32
	setChatID(id int32)
}

// RoomMessages implements MessagesProxyI
type RoomMessages struct {
	synced.PublisherBase

	s  synced.SyncI
	i  RoomInformationI
	l  LobbyProxyI
	se RSendI

	info *RoomMessageInfo
}

// init struct's values
func (room *RoomMessages) init(info *RoomMessageInfo) {
	room.info = info
	room.info.init()
}

// build components
func (room *RoomMessages) build(builder RBuilderI) {
	builder.BuildLobby(&room.l)
	builder.BuildSync(&room.s)
	builder.BuildSender(&room.se)
}

// subscribe to room events
func (room *RoomMessages) subscribe(builder RBuilderI) {
	var events EventsI
	builder.BuildEvents(&events)
	events.SubscribeRunnable(room)
}

// Init configure dependencies with other components of the room
func (room *RoomMessages) Init(builder RBuilderI,
	service clients.Chat, chatID int32, location *time.Location) {

	info := NewRoomMessageInfo(service, chatID, location)
	room.init(info)
	room.build(builder)
	room.subscribe(builder)
}

func (room *RoomMessages) ChatID() int32 {
	return room.info.dbChatID
}

func (room *RoomMessages) setChatID(id int32) {
	room.info.dbChatID = id
}

// func (room *RoomMessages) Proto(message *models.Message) (*chat.Message, error) {
// 	return ctypes.MessageToProto(message, room.dbChatID)
// }

func (room *RoomMessages) Fix(message *models.Message, conn *Connection) {
	room.s.Do(func() {
		if conn.Index() < 0 {
			message.Status = models.StatusObserver
		} else {
			message.Status = models.StatusPlayer
		}
		message.User = conn.User
		message.Time = time.Now().In(room.info.location)

		if message.Action == models.Write {
			rand.Seed(time.Now().UnixNano())
			message.ID = rand.Int31n(10000000) // в конфиг?
		}
	})
}

func (room *RoomMessages) HandleError(message *models.Message, send *chat.Message) {
	action := message.Action
	room.l.SaveMessages(&MessageWithAction{
		message, send, action, func() (int32, error) {
			if room.info.dbChatID != 0 {
				return room.info.dbChatID, nil
			}
			id, err := GetChatID(room.l.ChatService(), ctypes.RoomType, room.i.RoomID())
			if err != nil {
				room.info.dbChatID = id
			}
			return id, err
		}})
}

func (room *RoomMessages) setMessages(messages []*models.Message) {
	room.info.setMessages(messages)
}

func (room *RoomMessages) Messages() []*models.Message {
	return room.info.Messages()
}

func (room *RoomMessages) Proto(message *models.Message) (*chat.Message, error) {
	return ctypes.MessageToProto(message, room.info.dbChatID)
}

func (room *RoomMessages) Send(message *models.Message) {
	room.se.Message(*message)
}

func (room *RoomMessages) Write(message *models.Message, send *chat.Message) error {
	var err error
	room.s.Do(func() {
		room.info.appendMessage(message)
		var msgID *chat.MessageID
		msgID, err = room.info.service.AppendMessage(send)
		if err != nil {
			return
		}
		if msgID != nil {
			message.ID = msgID.Value
		}
		room.notify(room_.Add)
	})
	return err
}

func (room *RoomMessages) Update(message *models.Message, send *chat.Message) error {
	var err error
	room.s.Do(func() {
		found := room.info.findMessage(message.ID)
		if found <= 0 {
			err = re.InvalidMessageID()
			return
		}
		room.info.setMessage(found, message)

		_, err = room.info.service.UpdateMessage(send)
	})
	return err
}

func (room *RoomMessages) Delete(message *models.Message, send *chat.Message) error {
	var err error
	room.s.Do(func() {
		if message.ID <= 0 {
			err = re.InvalidMessageID()
			return
		}
		found := room.info.findMessage(message.ID)
		room.info.removeMessage(found)

		_, err = room.info.service.DeleteMessage(send)
		if err != nil {
			return
		}

		room.notify(room_.Delete)
	})
	return err
}

// messagesFree free message slice
func (room *RoomMessages) messagesFree() {
	room.info.messagesM.Lock()
	room.info._messages = nil
	room.info.messagesM.Unlock()
}

// notify all subscribers
func (room *RoomMessages) notify(action int32) {
	room.Notify(synced.Msg{
		Publisher: room_.UpdateChat,
		Action:    action,
	})
}

// start all goroutines
func (room *RoomMessages) start() {
	room.StartPublish()
}

// finish all goroutines and free memory
func (room *RoomMessages) stop() {
	room.messagesFree()
	room.StopPublish()
}

// 314
