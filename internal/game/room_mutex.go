package game

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

// setMatrixValue set a value to matrix
func (room *Room) setDone() {
	room.doneM.Lock()
	room._done = true
	room.doneM.Unlock()
}

// getMatrixValue get a value from matrix
func (room *Room) done() bool {
	room.doneM.RLock()
	v := room._done
	room.doneM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (room *Room) killed() int {
	room.killedM.RLock()
	v := room._killed
	room.killedM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (room *Room) incrementKilled() {
	room.killedM.Lock()
	room._killed++
	room.killedM.Unlock()
}

// getMatrixValue get a value from matrix
func (room *Room) setKilled(killed int) {
	room.killedM.Lock()
	room._killed = killed
	room.killedM.Unlock()
}

// getMatrixValue get a value from matrix
// func (room *Room) observers() (v []*Connection) {
// 	room.observersM.RLock()
// 	defer room.observersM.RUnlock()

// 	observers := room._Observers
// 	if observers == nil {
// 		return
// 	}
// 	v = room._Observers.Get
// 	return
// }

// getMatrixValue get a value from matrix
// func (room *Room) players() (v []Player) {
// 	room.playersM.RLock()
// 	defer room.playersM.RUnlock()

// 	players := room._Players
// 	if players == nil {
// 		return
// 	}
// 	v = room._Players.Players
// 	return
// }

// func (room *Room) SetFlagCoordinates(conn Connection, cell Cell) {
// 	if room.done() {
// 		return
// 	}
// 	room.wGroup.Add(1)
// 	defer func() {
// 		room.wGroup.Done()
// 	}()

// 	room.playersM.Lock()
// 	room._Players.Flags[conn.Index()].X = cell.X
// 	room._Players.Flags[conn.Index()].Y = cell.Y
// 	room.playersM.Unlock()
// }

// getMatrixValue get a value from matrix
// func (room *Room) RPlayersConnections() (v []*Connection) {
// 	room.playersM.RLock()
// 	defer room.playersM.RUnlock()

// 	if room._Players == nil {
// 		return
// 	}
// 	v = room._Players.Connections.RGet()
// 	return
// }

// getMatrixValue get a value from matrix
// func (room *Room) RPConnections() (v Connections) {
// 	room.playersM.RLock()
// 	defer room.playersM.RUnlock()

// 	if room._Players == nil {
// 		return
// 	}
// 	v = room._Players.Connections
// 	return
// }

// -> room.Players.Connections.Set(*NewConnections(room.Players.Capacity()))
// getMatrixValue get a value from matrix
// func (room *Room) zeroPlayers() {
// 	room.playersM.Lock()
// 	defer room.playersM.Unlock()

// 	room._Players.Connections = *NewConnections(room._Players.Capacity)

// }

// getMatrixValue get a value from matrix
// func (room *Room) playersFlags() (v []Cell) {
// 	room.playersM.RLock()
// 	defer room.playersM.RUnlock()

// 	if room._Players == nil {
// 		return
// 	}
// 	v = room._Players.Flags
// 	return
// }

// getMatrixValue get a value from matrix
// func (room *Room) playersCapacity() int {
// 	room.playersM.RLock()
// 	defer room.playersM.RUnlock()

// 	if room._Players == nil {
// 		return 0
// 	}
// 	v := room._Players.Capacity
// 	return v
// }

// getMatrixValue get a value from matrix
// func (room *Room) IncreasePlayerPoints(index, points int) {
// 	if room.done() {
// 		return
// 	}
// 	room.wGroup.Add(1)
// 	defer func() {
// 		room.wGroup.Done()
// 	}()

// 	room.playersM.Lock()
// 	defer room.playersM.Unlock()

// 	players := room._Players
// 	if players == nil {
// 		return
// 	}
// 	if index >= len(room._Players.Players) {
// 		return
// 	}
// 	room._Players.Players[index].Points += points
// }

// getMatrixValue get a value from matrix
// func (room *Room) playerFinished(index int) bool {

// 	room.playersM.RLock()
// 	defer room.playersM.RUnlock()

// 	players := room._Players
// 	if players == nil {
// 		return false
// 	}
// 	if index >= len(room._Players.Players) {
// 		return false
// 	}
// 	v := room._Players.Players[index].Finished
// 	return v
// }

// getMatrixValue get a value from matrix
// func (room *Room) player(index int) (v Player) {

// 	room.playersM.RLock()
// 	defer room.playersM.RUnlock()

// 	players := room._Players
// 	if players == nil {
// 		return
// 	}
// 	if index >= len(room._Players.Players) {
// 		return
// 	}
// 	v = room._Players.Players[index]
// 	return v
// }

// getMatrixValue get a value from matrix
// func (room *Room) playersInit() {

// 	room.playersM.Lock()
// 	defer room.playersM.Unlock()
// 	room._Players.Init(room.Field)
// }

// getMatrixValue get a value from matrix
// func (room *Room) playerFlag(index int) (v Cell) {

// 	room.playersM.RLock()
// 	defer room.playersM.RUnlock()

// 	players := room._Players
// 	if players == nil {
// 		return
// 	}
// 	if index >= len(room._Players.Players) {
// 		return
// 	}
// 	v = room._Players.Flags[index]
// 	return v
// }

// SetFinished increment amount of killed
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

// SetFinished increment amount of killed
// func (room *Room) setCell(conn *Connection) (flag *Cell) {

// 	index := conn.Index()
// 	if index < 0 {
// 		return
// 	}
// 	room.playersM.Lock()
// 	room._Players.Flags[index].PlayerID = conn.ID()
// 	room._Players.Flags[index].Value = conn.ID() + CellIncrement
// 	flag = &room._Players.Flags[index]
// 	room.playersM.Unlock()
// 	return
// }

// getMatrixValue get a value from matrix
func (room *Room) history() []*PlayerAction {
	room.historyM.RLock()
	v := room._History
	room.historyM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (room *Room) setToHistory(action *PlayerAction) {
	room.historyM.Lock()
	defer room.historyM.Unlock()
	room._History = append(room._History, action)
}

// getMatrixValue get a value from matrix
func (room *Room) setToMessages(message *models.Message) {
	room.messagesM.Lock()
	defer room.messagesM.Unlock()
	room._Messages = append(room._Messages, message)
}

func (room *Room) historyFree() {
	room.historyM.Lock()
	room._History = nil
	room.historyM.Unlock()
}

func (room *Room) messagesFree() {
	room.messagesM.Lock()
	room._Messages = nil
	room.messagesM.Unlock()
}

// func (room *Room) playersFree() {
// 	room.playersM.Lock()
// 	room._Players = nil
// 	room.playersM.Unlock()
// }

// func (room *Room) observersFree() {
// 	room.observersM.Lock()
// 	room._Observers = nil
// 	room.observersM.Unlock()
// }

// func (room *Room) observersEnoughPlace() bool {
// 	room.observersM.RLock()
// 	v := room._Observers.enoughPlace()
// 	room.observersM.RUnlock()
// 	return v
// }

// func (room *Room) playersEnoughPlace() bool {
// 	room.playersM.RLock()
// 	v := room._Players.enoughPlace()
// 	room.playersM.RUnlock()
// 	return v
// }

// getMatrixValue get a value from matrix
// func (room *Room) playersAdd(conn *Connection, kill bool) {

// 	room.playersM.Lock()
// 	defer room.playersM.Unlock()
// 	room._Players.Add(conn, kill)
// }

// getMatrixValue get a value from matrix
// func (room *Room) observersAdd(conn *Connection, kill bool) {

// 	room.observersM.Lock()
// 	defer room.observersM.Unlock()
// 	room._Observers.Add(conn, kill)
// }

// getMatrixValue get a value from matrix
// func (room *Room) playersRemove(conn *Connection, disconnect bool) bool {

// 	room.playersM.Lock()
// 	defer room.playersM.Unlock()
// 	return room._Players.Remove(conn, disconnect)
// }

// getMatrixValue get a value from matrix
// func (room *Room) observersRemove(conn *Connection, disconnect bool) bool {

// 	room.observersM.Lock()
// 	defer room.observersM.Unlock()
// 	return room._Observers.Remove(conn, disconnect)
// }

// getMatrixValue get a value from matrix
// func (room *Room) playersSearchIndexPlayer(conn *Connection) int {

// 	room.playersM.RLock()
// 	defer room.playersM.RUnlock()
// 	i := room._Players.SearchIndexPlayer(conn)
// 	return i
// }

// getMatrixValue get a value from matrix
// func (room *Room) playersEmpty() bool {

// 	room.playersM.RLock()
// 	defer room.playersM.RUnlock()
// 	v := room._Players.Empty()
// 	return v
// }

// getMatrixValue get a value from matrix
// func (room *Room) observersSearch(conn *Connection) int {

// 	room.playersM.RLock()
// 	defer room.playersM.RUnlock()
// 	i := room._Observers.Search(conn)
// 	return i
// }
