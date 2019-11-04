package engine

import (
	"sync"
)

type OnlinePlayersMutex struct {
	capacityM *sync.RWMutex
	_capacity int32

	playersM *sync.RWMutex
	_players []Player

	flagsM     *sync.RWMutex
	_flags     []Flag
	_flagsLeft int32
}

// Capacity return '_capacity' field
func (onlinePlayers *OnlinePlayersMutex) Capacity() int32 {

	onlinePlayers.capacityM.RLock()
	defer onlinePlayers.capacityM.RUnlock()
	return onlinePlayers._capacity
}

// SetCapacity set capacity
func (onlinePlayers *OnlinePlayersMutex) SetCapacity(capacity int32) {

	onlinePlayers.capacityM.Lock()
	defer onlinePlayers.capacityM.Unlock()
	onlinePlayers._capacity = capacity
}

// SetPlayer set player in slice of players
func (onlinePlayers *OnlinePlayersMutex) SetPlayer(i int, player Player) {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	onlinePlayers._players[i] = player
	return
}

// SetPlayers set player slice
func (onlinePlayers *OnlinePlayersMutex) SetPlayers(players []Player) {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	onlinePlayers._players = players
	return
}

// IncreasePlayerPoints increase points of element of player slice where index i
func (onlinePlayers *OnlinePlayersMutex) IncreasePlayerPoints(index int, points float64) {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()

	if index >= len(onlinePlayers._players) || index < 0 {
		return
	}

	onlinePlayers._players[index].Points += points
	if onlinePlayers._players[index].Points < 0 {
		onlinePlayers._players[index].Points = 0
	}
}

// SetFlags set flag slice
func (onlinePlayers *OnlinePlayersMutex) SetFlags(flags []Flag) {

	onlinePlayers.flagsM.Lock()
	defer onlinePlayers.flagsM.Unlock()
	onlinePlayers._flags = flags
	return
}

// SetPlayerID sets the id of an player slice element with an index i
func (onlinePlayers *OnlinePlayersMutex) SetPlayerID(i int, id int32) {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	onlinePlayers._players[i].ID = id
	return
}

// SetFlag set flag which connection is conn
func (onlinePlayers *OnlinePlayersMutex) SetFlag(conn *Connection, cell Cell, prepareOver func()) {

	onlinePlayers.flagsM.Lock()
	defer onlinePlayers.flagsM.Unlock()
	index := conn.Index()
	if index < 0 {
		return
	}

	onlinePlayers._flags[index].Cell.X = cell.X
	onlinePlayers._flags[index].Cell.Y = cell.Y
	onlinePlayers._flags[index].Cell.PlayerID = conn.ID()
	onlinePlayers._flags[index].Cell.Value = FlagID(conn.ID())

	if !onlinePlayers._flags[conn.Index()].Set {
		onlinePlayers._flags[conn.Index()].Set = true
		onlinePlayers._flagsLeft--
		if onlinePlayers._flagsLeft == 0 {
			prepareOver()
		}
	}
}

// SetFlagByIndex set flag by index
func (onlinePlayers *OnlinePlayersMutex) SetFlagByIndex(i int, flag Flag) {

	onlinePlayers.flagsM.Lock()
	defer onlinePlayers.flagsM.Unlock()
	onlinePlayers._flags[i] = flag
}

// Flags return '_flags' field
func (onlinePlayers *OnlinePlayersMutex) Flags() []Flag {

	onlinePlayers.flagsM.Lock()
	defer onlinePlayers.flagsM.Unlock()

	return onlinePlayers._flags
}

// Finish set flag finish true to all players
func (onlinePlayers *OnlinePlayersMutex) Finish(wg *sync.WaitGroup) {

	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	// all players 'Finished' set true
	onlinePlayers.playersM.Lock()
	for _, player := range onlinePlayers._players {
		player.Finished = true
	}
	onlinePlayers.playersM.Unlock()

}

// Player return element from slice of players with index i
func (onlinePlayers *OnlinePlayersMutex) Player(i int) Player {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	return onlinePlayers._players[i]
}

// PlayerFinish set to player with index i in slice of players flags Finished and Died true
func (onlinePlayers *OnlinePlayersMutex) PlayerFinish(i int) {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	onlinePlayers._players[i].Finished = true
	onlinePlayers._players[i].Died = true
}

// Flag return element from slice of flags with index i
func (onlinePlayers *OnlinePlayersMutex) Flag(i int) Flag {

	onlinePlayers.flagsM.Lock()
	defer onlinePlayers.flagsM.Unlock()
	return onlinePlayers._flags[i]
}

// CopyPlayers return copy of slice of players
func (onlinePlayers *OnlinePlayersMutex) CopyPlayers() []Player {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	copySlice := append([]Player(nil), onlinePlayers._players...)
	return copySlice
}

// RPlayers return slice of players only for read
func (onlinePlayers *OnlinePlayersMutex) RPlayers() []Player {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	return onlinePlayers._players
}

// Free free memory
func (onlinePlayers *OnlinePlayersMutex) Free() {
	if onlinePlayers == nil {
		return
	}

	onlinePlayers.playersM.Lock()
	onlinePlayers._players = nil
	onlinePlayers.playersM.Unlock()

	onlinePlayers.flagsM.Lock()
	onlinePlayers._flags = nil
	onlinePlayers.flagsM.Unlock()
}
