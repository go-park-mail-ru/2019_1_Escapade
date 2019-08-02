package game

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

//setDone set done = true. It will finish all operaions on Room
func (room *Room) setDone() {
	room.doneM.Lock()
	room._done = true
	room.doneM.Unlock()
}

// done return room readiness flag to free up resources
func (room *Room) done() bool {
	room.doneM.RLock()
	v := room._done
	room.doneM.RUnlock()
	return v
}

// Status return room's current status
func (room *Room) Status() int {
	room.statusM.RLock()
	v := room._status
	room.statusM.RUnlock()
	return v
}

// Name return the name of room given by its creator
func (room *Room) Name() string {
	room.nameM.RLock()
	v := room._name
	room.nameM.RUnlock()
	return v
}

// ID return room's unique identificator
func (room *Room) ID() string {
	room.idM.RLock()
	v := room._id
	room.idM.RUnlock()
	return v
}

// Date return date, when room was created
func (room *Room) Date() time.Time {
	room.dateM.RLock()
	v := room._date
	room.dateM.RUnlock()
	return v
}

func (room *Room) recruitmentTime() time.Duration {
	room.recruitmentTimeM.RLock()
	v := room._recruitmentTime
	room.recruitmentTimeM.RUnlock()
	return v
}

func (room *Room) setRecruitmentTime() {
	room.dateM.RLock()
	v := room._date
	room.dateM.RUnlock()

	t := time.Now().In(room.lobby.location())

	fmt.Println("compare these:", t, v)

	room.recruitmentTimeM.Lock()
	room._recruitmentTime = t.Sub(v)
	room.recruitmentTimeM.Unlock()

	room.dateM.Lock()
	room._date = t
	room.dateM.Unlock()
}

func (room *Room) playingTime() time.Duration {
	room.playingTimeM.RLock()
	v := room._playingTime
	room.playingTimeM.RUnlock()
	return v
}

func (room *Room) setPlayingTime() {
	room.dateM.RLock()
	v := room._date
	room.dateM.RUnlock()

	t := time.Now().In(room.lobby.location())

	room.playingTimeM.Lock()
	room._playingTime = t.Sub(v)
	room.playingTimeM.Unlock()

	room.dateM.Lock()
	room._date = t
	room.dateM.Unlock()
}

// Next return next room to whick players from this room will be
// sent in case of pressing the restart button
func (room *Room) Next() *Room {
	room.nextM.RLock()
	v := room._next
	room.nextM.RUnlock()
	return v
}

func (room *Room) setStatus(status int) {
	room.statusM.Lock()
	room._status = status
	room.statusM.Unlock()
}

func (room *Room) setName(name string) {
	room.nameM.Lock()
	room._name = name
	room.nameM.Unlock()
}

func (room *Room) setID(id string) {
	room.idM.Lock()
	room._id = id
	room.idM.Unlock()
}

func (room *Room) setDate(date time.Time) {
	room.dateM.Lock()
	room._date = date
	room.dateM.Unlock()
}

func (room *Room) setNext(next *Room) {
	room.nextM.Lock()
	room._next = next
	room.nextM.Unlock()
}

// history return '_history' field
func (room *Room) history() []*PlayerAction {
	room.historyM.RLock()
	v := room._history
	room.historyM.RUnlock()
	return v
}

func (room *Room) setHistory(history []*PlayerAction) {
	room.historyM.Lock()
	room._history = history
	room.historyM.Unlock()
}

// appendAction append action to action slice(history)
func (room *Room) appendAction(action *PlayerAction) {
	room.historyM.Lock()
	defer room.historyM.Unlock()
	room._history = append(room._history, action)
}

// historyFree free action slice
func (room *Room) historyFree() {
	room.historyM.Lock()
	room._history = nil
	room.historyM.Unlock()
}

// done return '_killed' field
func (room *Room) killed() int32 {
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
func (room *Room) setKilled(killed int32) {
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

/////////////////////// messages

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
func (room *Room) setMessages(messages []*models.Message) {
	room.messagesM.Lock()
	room._messages = messages
	room.messagesM.Unlock()
}

// Messages return slice of messages
func (room *Room) Messages() []*models.Message {
	room.messagesM.Lock()
	defer room.messagesM.Unlock()
	return room._messages
}

// messagesFree free message slice
func (room *Room) messagesFree() {
	room.messagesM.Lock()
	room._messages = nil
	room.messagesM.Unlock()
}
