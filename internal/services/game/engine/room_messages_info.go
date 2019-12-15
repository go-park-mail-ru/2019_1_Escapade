package engine

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
)

type RoomMessageInfo struct {
	messagesM *sync.Mutex
	_messages []*models.Message

	dbChatID int32
	location *time.Location
	service  clients.ChatI
}

// NewRoomMessageInfo create new instanse of RoomMessageInfo
func NewRoomMessageInfo(c clients.ChatI, id int32, l *time.Location) *RoomMessageInfo {
	return &RoomMessageInfo{
		dbChatID: id,
		location: l,
		service:  c,
	}
}

// init struct's values
func (info *RoomMessageInfo) init() {
	info.messagesM = &sync.Mutex{}
	info.setMessages(make([]*models.Message, 0))
}

// appendMessage append message to message slice
func (info *RoomMessageInfo) appendMessage(message *models.Message) {
	info.messagesM.Lock()
	defer info.messagesM.Unlock()
	info._messages = append(info._messages, message)
}

// removeMessage remove message from messages slice
func (info *RoomMessageInfo) removeMessage(i int) {
	info.messagesM.Lock()
	defer info.messagesM.Unlock()
	if i < 0 {
		return
	}
	size := len(info._messages)

	info._messages[i], info._messages[size-1] = info._messages[size-1], info._messages[i]
	info._messages[size-1] = nil
	info._messages = info._messages[:size-1]
	return
}

// setMessage update message from messages slice with index i
func (info *RoomMessageInfo) setMessage(i int, message *models.Message) {

	info.messagesM.Lock()
	defer info.messagesM.Unlock()
	if i < 0 {
		return
	}
	info._messages[i] = message
	info._messages[i].Edited = true
	return
}

// findMessage search message by message ID
func (info *RoomMessageInfo) findMessage(searchID int32) int {
	if searchID <= 0 {
		return -1
	}
	info.messagesM.Lock()
	messages := info._messages
	info.messagesM.Unlock()

	for i, message := range messages {
		if message.ID == searchID {
			return i
		}
	}
	return -1
}

// Messages return slice of messages
func (info *RoomMessageInfo) setMessages(messages []*models.Message) {
	info.messagesM.Lock()
	info._messages = messages
	info.messagesM.Unlock()
}

// Messages return slice of messages
func (info *RoomMessageInfo) Messages() []*models.Message {
	info.messagesM.Lock()
	defer info.messagesM.Unlock()
	return info._messages
}
