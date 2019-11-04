package engine

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

type RoomMessages struct {
	dbChatID int32

	messagesM *sync.Mutex
	_messages []*models.Message
}

func (room *RoomMessages) Init(chatID int32) {
	room.messagesM = &sync.Mutex{}
	room.dbChatID = chatID
	room.setMessages(make([]*models.Message, 0))
}

func (room *RoomMessages) Free() {
	room.messagesFree()
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
func (room *RoomMessages) findMessage(search *models.Message) int {

	room.messagesM.Lock()
	messages := room._messages
	room.messagesM.Unlock()

	for i, message := range messages {
		if message.ID == search.ID {
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
