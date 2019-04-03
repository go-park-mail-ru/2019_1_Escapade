package game

//re "escapade/internal/return_errors"
//"math/rand"

// Connections - map of connections with size and capacity
type Connections struct {
	Capacity int                 `json:"capacity"`
	Size     int                 `json:"size"`
	Get      map[int]*Connection `json:"get"`
}

// NewConnections create new instance of connections
func NewConnections(capacity int) *Connections {
	conns := &Connections{
		Capacity: capacity,
		Get:      make(map[int]*Connection),
	}
	return conns
}

// Empty check is size 0
func (conns *Connections) Empty() bool {
	return conns.Size == 0
}

func (conns *Connections) enoughPlace() bool {
	return conns.Size < conns.Capacity
}

// Add try add element if its possible. Return bool result
func (conns *Connections) Add(conn *Connection) bool {
	if conns.enoughPlace() {
		conns.Get[conn.GetPlayerID()] = conn
		conns.Size++
		return true
	}
	return false
}

// Change change connection by player id
func (conns *Connections) Change(conn *Connection) bool {
	if _, ok := conns.Get[conn.GetPlayerID()]; ok {
		conns.Get[conn.GetPlayerID()] = conn
		return true
	}
	return false
}

// Remove delete element and decrement size if element
// exists in map
func (conns *Connections) Remove(conn *Connection) {
	if _, ok := conns.Get[conn.GetPlayerID()]; ok {
		sendError(conn, "Remove", "You disconnected ")
		delete(conns.Get, conn.GetPlayerID())
		conns.Size--
	}

}
