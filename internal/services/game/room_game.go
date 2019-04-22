package game

import (
	"escapade/internal/models"

	"fmt"
)

// NewRoom return new instance of room
func NewRoom(rs *models.RoomSettings, name string, lobby *Lobby) *Room {
	fmt.Println("NewRoom rs = ", *rs)
	room := &Room{
		Name:      name,
		Status:    StatusPeopleFinding,
		Players:   newOnlinePlayers(rs.Players),
		Observers: NewConnections(rs.Observers),

		History: make([]*PlayerAction, 0),

		lobby:  lobby,
		Field:  NewField(rs),
		killed: 0,
	}
	return room
}

// flagFound is called, when somebody find cell flag
func (room *Room) flagFound(found *Cell) {
	thatID := found.Value - CellIncrement
	for i, player := range room.Players.Players {
		if thatID == player.ID {
			room.kill(room.Players.Connections[i], ActionFlagLost)
		}
	}
}

// kill make user die, decrement size and check for finish battle
func (room *Room) kill(conn *Connection, action int) {
	// cause all in pointers
	if !conn.Player.Finished {
		conn.Player.Finished = true
		room.killed++
		if room.Players.Capacity <= room.killed+1 {
			// остановить таймеры в run!!!
			room.lobby.roomFinish(room)
		}
		room.addAction(conn, action)
		room.sendHistory(room.all())
		conn.debug("give up. Check history")
	}
}

// use it when somebody exit
func (room *Room) GiveUp(conn *Connection) {
	room.kill(conn, ActionGiveUp)
}

// Close clear all resources. Call it when no
//  observers and players inside
func (room *Room) Free() {
	room.Players.Free()
	room.Observers.Free()
	for _, action := range room.History {
		action.Free()
	}
	room.History = nil
	room.Players.Free()
	room.Field.Clear()
}

func (room *Room) Close() bool {
	//.... leave all....
	room.lobby.CloseRoom(room)
	room.Free()
	return false
}

func (room *Room) setFlags() {
	for _, cell := range room.Players.Flags {
		room.Field.SetFlag(cell.X, cell.Y, cell.PlayerID)
	}
}

func (room *Room) fillField() {
	fmt.Println("fillField", room.Field.Height, room.Field.Width, len(room.Field.Matrix))

	room.setFlags()
	room.Field.SetMines()

}
