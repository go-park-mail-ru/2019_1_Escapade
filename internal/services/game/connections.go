package game

//re "escapade/internal/return_errors"
//"math/rand"

type Connections struct {
	Capacity int
	Size     int
	Get      map[*Connection]bool
}

func NewConnections(capacity int) *Connections {
	conns := &Connections{
		Capacity: capacity,
		Get:      make(map[*Connection]bool),
	}
	return conns
}

func (conns *Connections) Empty() bool {
	return conns.Size == 0
}

func (conns *Connections) enoughPlace() bool {
	return conns.Size < conns.Capacity
}

func (conns *Connections) Add(conn *Connection, b bool) bool {
	if conns.enoughPlace() {
		conns.Get[conn] = b
		conns.Size++
		return true
	}
	return false
}

func (conns *Connections) Remove(conn *Connection) {
	if _, ok := conns.Get[conn]; ok {
		sendError(conn, "Remove", "You disconnected ")

		delete(conns.Get, conn)
		conns.Size--
	}

}
