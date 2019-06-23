package game

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

//setDone set done = true. It will finish all operaions on Lobby
func (lobby *Lobby) setDone() {
	lobby.doneM.Lock()
	lobby._done = true
	lobby.doneM.Unlock()
}

// done return '_done' field
func (lobby *Lobby) done() bool {
	lobby.doneM.RLock()
	v := lobby._done
	lobby.doneM.RUnlock()
	return v
}

// appendMessage append message to messages slice
func (lobby *Lobby) appendMessage(message *models.Message) {
	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	lobby._messages = append(lobby._messages, message)
}

// removeMessage remove message from messages slice
func (lobby *Lobby) removeMessage(i int) {
	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	if i < 0 {
		return
	}
	size := len(lobby._messages)

	lobby._messages[i], lobby._messages[size-1] = lobby._messages[size-1], lobby._messages[i]
	lobby._messages[size-1] = nil
	lobby._messages = lobby._messages[:size-1]
	return
}

// setMessage update message from messages slice with index i
func (lobby *Lobby) setMessage(i int, message *models.Message) {
	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	if i < 0 {
		return
	}
	lobby._messages[i] = message
	lobby._messages[i].Edited = true
	return
}

// findMessage search message by message ID
func (lobby *Lobby) findMessage(search *models.Message) int {
	lobby.messagesM.Lock()
	messages := lobby._messages
	lobby.messagesM.Unlock()

	for i, message := range messages {
		if message.ID == search.ID {
			return i
		}
	}
	return -1
}

// Messages return slice of messages
func (lobby *Lobby) Messages() []*models.Message {

	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	return lobby._messages
}
