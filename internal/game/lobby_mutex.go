package game

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
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

func (lobby *Lobby) setDB(newDB *database.DataBase) {
	lobby.dbM.Lock()
	lobby._db = newDB
	lobby.dbM.Unlock()
}

func (lobby *Lobby) db() *database.DataBase {
	lobby.dbM.RLock()
	v := lobby._db
	lobby.dbM.RUnlock()
	return v
}

func (lobby *Lobby) setConfig(c *config.GameConfig) {
	lobby.configM.Lock()
	lobby._config = c
	lobby.configM.Unlock()
}

func (lobby *Lobby) config() *config.GameConfig {
	lobby.configM.RLock()
	v := lobby._config
	lobby.configM.RUnlock()
	return v
}

func (lobby *Lobby) setDBChatID(id int32) {
	lobby.dbChatIDM.Lock()
	lobby._dbChatID = id
	lobby.dbChatIDM.Unlock()
}

func (lobby *Lobby) dbChatID() int32 {
	lobby.dbChatIDM.RLock()
	v := lobby._dbChatID
	lobby.dbChatIDM.RUnlock()
	return v
}

func (lobby *Lobby) setLocation(newLocation *time.Location) {
	lobby.locationM.Lock()
	lobby._location = newLocation
	lobby.locationM.Unlock()
}

func (lobby *Lobby) location() *time.Location {
	lobby.locationM.RLock()
	v := lobby._location
	lobby.locationM.RUnlock()
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

func (lobby *Lobby) setMessages(messages []*models.Message) {
	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	lobby._messages = messages
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

// Anonymous return anonymous id
func (lobby *Lobby) Anonymous() int {
	var id int
	lobby.anonymousM.Lock()
	id = lobby._anonymous
	lobby._anonymous--
	lobby.anonymousM.Unlock()
	return id
}
