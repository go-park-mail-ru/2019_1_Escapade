package engine

import (
	"math/rand"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
	ctypes "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/database"
	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
)

// MessagesI control access to messages
// Proxy Pattern
type MessagesI interface {
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

	Free()
}

// RoomMessages implements MessagesProxyI
type RoomMessages struct {
	dbChatID int32

	messagesM *sync.Mutex
	_messages []*models.Message

	service clients.Chat

	s  synced.SyncI
	i  RoomInformationI
	l  LobbyProxyI
	se RSendI
}

// Init configure dependencies with other components of the room
func (room *RoomMessages) Init(builder RBuilderI,
	service clients.Chat, chatID int32) {
	builder.BuildInformation(&room.i)
	builder.BuildLobby(&room.l)
	builder.BuildSync(&room.s)
	builder.BuildSender(&room.se)

	room.service = service
	room.messagesM = &sync.Mutex{}
	room.dbChatID = chatID
	room.setMessages(make([]*models.Message, 0))
}

func (room *RoomMessages) Free() {
	room.messagesFree()
}

func (room *RoomMessages) ChatID() int32 {
	return room.dbChatID
}

func (room *RoomMessages) setChatID(id int32) {
	room.dbChatID = id
}

func (room *RoomMessages) Proto(message *models.Message) (*chat.Message, error) {
	return ctypes.MessageToProto(message, room.dbChatID)
}

func (room *RoomMessages) Fix(message *models.Message, conn *Connection) {
	room.s.Do(func() {
		if conn.Index() < 0 {
			message.Status = models.StatusObserver
		} else {
			message.Status = models.StatusPlayer
		}
		message.User = conn.User
		message.Time = room.l.Date()

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
			if room.dbChatID != 0 {
				return room.dbChatID, nil
			}
			id, err := GetChatID(room.l.ChatService(), ctypes.RoomType, room.i.RoomID())
			if err != nil {
				room.dbChatID = id
			}
			return id, err
		}})
}

func (room *RoomMessages) Send(message *models.Message) {
	room.se.Message(*message)
}

func (room *RoomMessages) Write(message *models.Message, send *chat.Message) error {
	var err error
	room.s.Do(func() {
		room.appendMessage(message)
		var msgID *chat.MessageID
		msgID, err = room.service.AppendMessage(send)
		if err != nil {
			return
		}
		if msgID != nil {
			message.ID = msgID.Value
		}
		if room.l.metricsEnabled() {
			metrics.RoomsMessages.Inc()
		}
	})
	return err
}

func (room *RoomMessages) Update(message *models.Message, send *chat.Message) error {
	var err error
	room.s.Do(func() {
		found := room.findMessage(message.ID)
		if found <= 0 {
			err = re.InvalidMessageID()
			return
		}
		room.setMessage(found, message)

		_, err = room.service.UpdateMessage(send)
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
		found := room.findMessage(message.ID)
		room.removeMessage(found)

		_, err = room.service.DeleteMessage(send)
		if err != nil {
			return
		}

		if room.l.metricsEnabled() {
			metrics.RoomsMessages.Dec()
		}
	})
	return err
}

// appendMessage append message to message slice
func (room *RoomMessages) appendMessage(message *models.Message) {
	room.messagesM.Lock()
	defer room.messagesM.Unlock()
	room._messages = append(room._messages, message)
}

// removeMessage remove message from messages slice
func (room *RoomMessages) removeMessage(i int) {
	room.messagesM.Lock()
	defer room.messagesM.Unlock()
	if i < 0 {
		return
	}
	size := len(room._messages)

	room._messages[i], room._messages[size-1] = room._messages[size-1], room._messages[i]
	room._messages[size-1] = nil
	room._messages = room._messages[:size-1]
	return
}

// setMessage update message from messages slice with index i
func (room *RoomMessages) setMessage(i int, message *models.Message) {

	room.messagesM.Lock()
	defer room.messagesM.Unlock()
	if i < 0 {
		return
	}
	room._messages[i] = message
	room._messages[i].Edited = true
	return
}

// findMessage search message by message ID
func (room *RoomMessages) findMessage(searchID int32) int {
	if searchID <= 0 {
		return -1
	}
	room.messagesM.Lock()
	messages := room._messages
	room.messagesM.Unlock()

	for i, message := range messages {
		if message.ID == searchID {
			return i
		}
	}
	return -1
}

// Messages return slice of messages
func (room *RoomMessages) setMessages(messages []*models.Message) {
	room.messagesM.Lock()
	room._messages = messages
	room.messagesM.Unlock()
}

// Messages return slice of messages
func (room *RoomMessages) Messages() []*models.Message {
	room.messagesM.Lock()
	defer room.messagesM.Unlock()
	return room._messages
}

// messagesFree free message slice
func (room *RoomMessages) messagesFree() {
	room.messagesM.Lock()
	room._messages = nil
	room.messagesM.Unlock()
}

// Message send message to connections
/*
func (room *RoomMessages) Handle(conn *Connection, message *models.Message) error {
	var err error
	room.s.do(func() {
		room.Fix(message, conn)

		var msg *chat.Message
		msg, err = chat.MessageToProto(message, room.dbChatID)
		if err != nil {
			return
		}

		// ignore models.StartWrite, models.FinishWrite
		switch message.Action {
		case models.Write:
			err = room.Write(message, msg)
		case models.Update:
			err = room.Update(message, msg)
		case models.Delete:
			err = room.Delete(message, msg)
		}
		if err != nil {
			room.HandleError(message, msg)
		} else {
			room.Send(message)
		}
	})
	return err
}*/
