package game

import "fmt"

// RGet return connections slice only for Read!
func (onlinePlayers *OnlinePlayers) Capacity() int {

	onlinePlayers.capacityM.RLock()
	defer onlinePlayers.capacityM.RUnlock()
	return onlinePlayers._capacity
}

// RGet return connections slice only for Read!
func (onlinePlayers *OnlinePlayers) SetCapacity(capacity int) {

	onlinePlayers.capacityM.Lock()
	defer onlinePlayers.capacityM.Unlock()
	onlinePlayers._capacity = capacity
}

// RGet return connections slice only for Read!
func (onlinePlayers *OnlinePlayers) SetPlayer(i int, player Player) {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	onlinePlayers._players[i] = player
	return
}

func (onlinePlayers *OnlinePlayers) SetPlayers(players []Player) {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	onlinePlayers._players = players
	return
}

func (onlinePlayers *OnlinePlayers) IncreasePlayerPoints(index, points int) {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()

	if index >= len(onlinePlayers._players) || index < 0 {
		return
	}

	onlinePlayers._players[index].Points += points
}

func (onlinePlayers *OnlinePlayers) SetFlags(flags []Flag) {

	onlinePlayers.flagsM.Lock()
	defer onlinePlayers.flagsM.Unlock()
	onlinePlayers._flags = flags
	return
}

// RGet return connections slice only for Read!
func (onlinePlayers *OnlinePlayers) SetPlayerID(i int, id int) {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	onlinePlayers._players[i].ID = id
	return
}

func (onlinePlayers *OnlinePlayers) SetFlag(conn Connection, cell Cell) bool {

	onlinePlayers.flagsM.Lock()
	defer onlinePlayers.flagsM.Unlock()
	index := conn.Index()
	if index < 0 {
		return false
	}
	fmt.Println("somebody set flag")
	onlinePlayers._flags[index].Cell.X = cell.X
	onlinePlayers._flags[index].Cell.Y = cell.Y
	onlinePlayers._flags[index].Cell.PlayerID = conn.ID()
	onlinePlayers._flags[index].Cell.Value = conn.ID() + CellIncrement

	if !onlinePlayers._flags[conn.Index()].Set {
		onlinePlayers._flags[conn.Index()].Set = true
		onlinePlayers.flagsLeft--
	}
	return onlinePlayers.flagsLeft == 0
}

func (onlinePlayers *OnlinePlayers) Flags() []Flag {

	onlinePlayers.flagsM.Lock()
	defer onlinePlayers.flagsM.Unlock()

	return onlinePlayers._flags
}

func (onlinePlayers *OnlinePlayers) Finish() {

	// all players 'Finished' set true
	onlinePlayers.playersM.Lock()
	for _, player := range onlinePlayers._players {
		player.Finished = true
	}
	onlinePlayers.playersM.Unlock()

}

// RGet return connections slice only for Read!
func (onlinePlayers *OnlinePlayers) Player(i int) Player {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	return onlinePlayers._players[i]
}

func (onlinePlayers *OnlinePlayers) PlayerFinish(i int) {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	onlinePlayers._players[i].Finished = true
	onlinePlayers._players[i].Died = true
}

func (onlinePlayers *OnlinePlayers) Flag(i int) Flag {

	onlinePlayers.flagsM.Lock()
	defer onlinePlayers.flagsM.Unlock()
	return onlinePlayers._flags[i]
}

// RGet return connections slice only for Read!
func (onlinePlayers *OnlinePlayers) RPlayers() []Player {

	onlinePlayers.playersM.Lock()
	defer onlinePlayers.playersM.Unlock()
	return onlinePlayers._players
}

// Free free memory
func (onlinePlayers *OnlinePlayers) Free() {

	if onlinePlayers == nil {
		return
	}
	onlinePlayers.Connections.Free()

	onlinePlayers.playersM.Lock()
	onlinePlayers._players = nil
	onlinePlayers.playersM.Unlock()

	onlinePlayers.flagsM.Lock()
	onlinePlayers._flags = nil
	onlinePlayers.flagsM.Unlock()
}
