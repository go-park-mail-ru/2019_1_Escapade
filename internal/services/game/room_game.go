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
		Players:   NewConnections(rs.Players),
		Observers: NewConnections(rs.Observers),

		History: make([]*PlayerAction, 0),
		flags:   make(map[*Connection]*Cell),

		lobby: lobby,
		Field: NewField(rs),
		//chanLeave: make(chan *Connection),
		//chanRequest: make(chan *RoomRequest),
	}
	return room
}

// flagFound is called, when somebody find cell flag
func (room *Room) flagFound(found *Cell) {
	thatID := found.Value - CellIncrement
	for _, conn := range room.Players.Get {
		if thatID == conn.GetPlayerID() {
			room.kill(conn, ActionFlagLost)
		}
	}
}

// kill make user die, decrement size and check for finish battle
func (room *Room) kill(conn *Connection, action int) {
	// cause all in pointers
	if !conn.Player.Finished {
		conn.Player.Finished = true
		room.Players.Capacity--
		if room.Players.Capacity <= 1 {
			room.lobby.roomFinish(room)
		}
		room.addAction(conn, action)
		room.sendHistory(room.all())
	}
}

// use it when somebody exit
func (room *Room) GiveUp(conn *Connection) {
	room.kill(conn, ActionGiveUp)
}

// Close clear all resources. Call it when no
//  observers and players inside
func (room *Room) Close() {
	room.Players.Clear()
	room.Observers.Clear()
	room.History = nil
	room.flags = nil
	room.Field.Clear()
}

func (room *Room) TryClose() {
	if room.Players.Empty() && room.Observers.Empty() {
		room.Close()
		room.lobby.CloseRoom(room)
	}
}

func (room *Room) setFlags() {
	for conn, cell := range room.flags {
		room.Field.SetFlag(cell.X, cell.Y, conn.GetPlayerID())
	}
}

func (room *Room) fillField() {
	fmt.Println("fillField", room.Field.Height, room.Field.Width, len(room.Field.Matrix))

	room.setFlags()
	room.Field.SetMines()

}
