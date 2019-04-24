package game

import (
	"fmt"
)

// flagFound is called, when somebody find cell flag
func (room *Room) flagFound(found *Cell) {
	thatID := found.Value - CellIncrement
	for i, player := range room.Players.Players {
		if thatID == player.ID {
			room.kill(room.Players.Connections[i], ActionFlagLost)
		}
	}
}

// isAlive check if connection is player and he is not died
func (room *Room) isAlive(conn *Connection) bool {
	return conn.index >= 0 && !room.Players.Players[conn.index].Finished
}

// setFinished increment amount of killed
func (room *Room) setFinished(conn *Connection) {
	room.Players.Players[conn.index].Finished = true
	room.killed++
}

// kill make user die and check for finish battle
func (room *Room) kill(conn *Connection, action int) {
	// cause all in pointers
	if room.isAlive(conn) {
		room.setFinished(conn)
		if room.Players.Capacity <= room.killed+1 {
			room.finishGame()
		}
		room.addAction(conn, action)
		room.sendHistory(room.All)
	}
}

// GiveUp kill connection, that call it
func (room *Room) GiveUp(conn *Connection) {
	room.kill(conn, ActionGiveUp)
}

// setFlag handle user wanna set flag
func (room *Room) setFlag(conn *Connection, cell *Cell) bool {
	// if user try set flag after game launch
	if room.Status != StatusFlagPlacing {
		return false
	}

	if !room.Field.IsInside(cell) {
		return false
	}

	if !room.isAlive(conn) {
		return false
	}

	room.Players.Flags[conn.index].X = cell.X
	room.Players.Flags[conn.index].Y = cell.Y
	return true
}

// setFlags set players flags to field
// call it if game has already begun
func (room *Room) setFlags() {
	for _, cell := range room.Players.Flags {
		room.Field.SetFlag(cell.X, cell.Y, cell.PlayerID)
	}
}

// fillField set flags and mines
func (room *Room) fillField() {
	fmt.Println("fillField", room.Field.Height, room.Field.Width, len(room.Field.Matrix))

	room.setFlags()
	room.Field.SetMines()

}
