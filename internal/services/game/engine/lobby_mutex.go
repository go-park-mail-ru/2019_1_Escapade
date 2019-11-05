package engine

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// checkAndSetCleared checks if the cleanup function was called. This check is
// based on 'done'. If it is true, then the function has already been called.
// If not, set done to True and return false.
// IMPORTANT: this function must only be called in the cleanup function
func (lobby *Lobby) checkAndSetCleared() bool {
	lobby.doneM.Lock()
	defer lobby.doneM.Unlock()
	if lobby._done {
		return true
	}
	lobby._done = true
	return false
}

// done return '_done' field
func (lobby *Lobby) done() bool {
	lobby.doneM.RLock()
	v := lobby._done
	lobby.doneM.RUnlock()
	return v
}

// AddNotSavedMessage add message to slice of unsaved messages
func (lobby *Lobby) AddNotSavedMessage(mwa *MessageWithAction) {
	lobby.notSavedMessagesM.Lock()
	lobby._notSavedMessages = append(lobby._notSavedMessages, mwa)
	lobby.notSavedMessagesM.Unlock()
}

// NotSavedMessagesGetAndClear returns an array of unsaved messages,
// zeroing the corresponding lobby field
func (lobby *Lobby) NotSavedMessagesGetAndClear() []*MessageWithAction {
	lobby.notSavedMessagesM.Lock()
	slice := lobby._notSavedMessages
	lobby._notSavedMessages = make([]*MessageWithAction, 0)
	lobby.notSavedMessagesM.Unlock()
	return slice
}

// AddNotSavedGame add game to slice of unsaved games
func (lobby *Lobby) AddNotSavedGame(game *models.GameInformation) {
	lobby.notSavedGamesM.Lock()
	lobby._notSavedGames = append(lobby._notSavedGames, game)
	lobby.notSavedGamesM.Unlock()
}

// NotSavedGamesGetAndClear returns an array of unsaved games,
// zeroing the corresponding lobby field
func (lobby *Lobby) NotSavedGamesGetAndClear() []*models.GameInformation {
	lobby.notSavedGamesM.Lock()
	slice := lobby._notSavedGames
	lobby._notSavedGames = make([]*models.GameInformation, 0)
	lobby.notSavedGamesM.Unlock()
	return slice
}

// setDB set the database object that the lobby is working with
func (lobby *Lobby) setDB(newDB database.GameUseCaseI) {
	lobby.dbM.Lock()
	lobby._db = newDB
	lobby.dbM.Unlock()
}

// dbChatID get the database object that the lobby is working with
func (lobby *Lobby) db() database.GameUseCaseI {
	lobby.dbM.RLock()
	v := lobby._db
	lobby.dbM.RUnlock()
	return v
}

// dbChatID set the configuration of lobby
func (lobby *Lobby) setConfig(c *config.Game) {
	lobby.configM.Lock()
	lobby._config = c
	lobby.configM.Unlock()
}

// dbChatID get the congifuration of lobby
func (lobby *Lobby) config() *config.Game {
	lobby.configM.RLock()
	v := lobby._config
	lobby.configM.RUnlock()
	return v
}

// dbChatID set the ID of the chat associated with the lobby
func (lobby *Lobby) setDBChatID(id int32) {
	lobby.dbChatIDM.Lock()
	lobby._dbChatID = id
	lobby.dbChatIDM.Unlock()
}

// dbChatID get the ID of the chat associated with the lobby
func (lobby *Lobby) dbChatID() int32 {
	lobby.dbChatIDM.RLock()
	v := lobby._dbChatID
	lobby.dbChatIDM.RUnlock()
	return v
}

// setLocation set lobby location
func (lobby *Lobby) setLocation(newLocation *time.Location) {
	lobby.locationM.Lock()
	lobby._location = newLocation
	lobby.locationM.Unlock()
}

// location get lobby location
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
	if i < 0 {
		return
	}
	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	size := len(lobby._messages)

	lobby._messages[i], lobby._messages[size-1] = lobby._messages[size-1], lobby._messages[i]
	lobby._messages[size-1] = nil
	lobby._messages = lobby._messages[:size-1]
	return
}

// setMessage updates the text of an lobby messages slice
// element with an index i. Also it is marked as edited
func (lobby *Lobby) setMessage(i int, message *models.Message) {
	if i < 0 {
		return
	}
	lobby.messagesM.Lock()
	defer lobby.messagesM.Unlock()
	lobby._messages[i] = message
	lobby._messages[i].Edited = true
	return
}

// addMessages add specified messages to an existing slice of lobby
// messages. IMPORTANT: messages are added to the beginning of the array.
// Return the old message lobby slice
func (lobby *Lobby) insertMessages(messages []*models.Message) []*models.Message {
	lobby.messagesM.Lock()
	v := lobby._messages
	lobby._messages = append(messages, lobby._messages...)
	lobby.messagesM.Unlock()
	return v
}

// setMessages set specified messages slice as lobby's messages slice
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
func (lobby *Lobby) Anonymous() int32 {
	lobby.anonymousM.Lock()
	id := lobby._anonymous
	lobby._anonymous--
	lobby.anonymousM.Unlock()
	return id
}
