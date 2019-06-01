package game

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

//setDone set done = true. It will finish all operaions on Room
func (room *Room) setDone() {
	room.doneM.Lock()
	room._done = true
	room.doneM.Unlock()
}

// done return '_done' field
func (room *Room) done() bool {
	room.doneM.RLock()
	v := room._done
	room.doneM.RUnlock()
	return v
}

// done return '_killed' field
func (room *Room) killed() int {
	room.killedM.RLock()
	v := room._killed
	room.killedM.RUnlock()
	return v
}

// incrementKilled increment amount of killed
func (room *Room) incrementKilled() {
	room.killedM.Lock()
	room._killed++
	room.killedM.Unlock()
}

// setKilled set new value of killed
func (room *Room) setKilled(killed int) {
	room.killedM.Lock()
	room._killed = killed
	room.killedM.Unlock()
}

// SetFinished set player finished
func (room *Room) SetFinished(conn *Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	index := conn.Index()
	if index < 0 {
		return
	}
	room.Players.PlayerFinish(index)

	room.killedM.Lock()
	room._killed++
	room.killedM.Unlock()
}

// history return '_history' field
func (room *Room) history() []*PlayerAction {
	room.historyM.RLock()
	v := room._history
	room.historyM.RUnlock()
	return v
}

// appendAction append action to action slice(history)
func (room *Room) appendAction(action *PlayerAction) {
	room.historyM.Lock()
	defer room.historyM.Unlock()
	room._history = append(room._history, action)
}

// appendMessage append message to message slice
func (room *Room) appendMessage(message *models.Message) {
	room.messagesM.Lock()
	defer room.messagesM.Unlock()
	room._messages = append(room._messages, message)
}

// removeMessage remove message from messages slice
func (room *Room) removeMessage(i int) {
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
func (room *Room) setMessage(i int, message *models.Message) {

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
func (room *Room) findMessage(search *models.Message) int {

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
func (room *Room) Messages() []*models.Message {

	room.messagesM.Lock()
	defer room.messagesM.Unlock()
	return room._messages
}

// historyFree free action slice
func (room *Room) historyFree() {
	room.historyM.Lock()
	room._history = nil
	room.historyM.Unlock()
}

// messagesFree free message slice
func (room *Room) messagesFree() {
	room.messagesM.Lock()
	room._messages = nil
	room.messagesM.Unlock()
}
