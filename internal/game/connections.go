package game

import (
	"encoding/json"
	"sync"
)

// Connections - slice of connections with capacity
type Connections struct {
	capacityM *sync.RWMutex
	_capacity int

	getM *sync.RWMutex
	_get []*Connection
}

// NewConnections create instance of Connections
func NewConnections(capacity int) (conns *Connections) {
	conns = &Connections{
		capacityM: &sync.RWMutex{},
		_capacity: capacity,

		getM: &sync.RWMutex{},
		_get: make([]*Connection, 0, capacity),
	}
	return conns
}

// Free free memory
func (conns *Connections) Free() {
	if conns == nil {
		return
	}
	conns.getM.Lock()
	for _, conn := range conns._get {
		conn.Free()
	}
	conns._get = nil
	conns.getM.Unlock()

	conns.capacityM.Lock()
	conns._capacity = 0
	conns.capacityM.Unlock()
}

// Remove delete element and decrement size if element
// exists in map
func (conns *Connections) Remove(conn *Connection, onlyIfDisconnected bool) bool {
	_, i := conns.SearchByID(conn.ID())
	if i < 0 {
		return false
	}
	if onlyIfDisconnected && !conns.RGet()[i].Disconnected() {
		return false
	}
	conns.remove(i)
	return true
	//sendError(conn, "Remove", "You disconnected ")
}

// Add try add element if its possible. Return bool result
// if element not exists it will be create, otherwise it will change its value
func (conns *Connections) Add(conn *Connection, kill bool) (i int) {
	_, i = conns.SearchByID(conn.ID())
	if i >= 0 {
		oldConn := conns.RGet()[i]
		if kill && !oldConn.Disconnected() {
			oldConn.Kill("Another connection found", true)
		}
		conn._room = oldConn.Room()
		conn._Index = oldConn.Index()

		conns.set(i, conn)
		i = oldConn.Index()
	} else if conns.EnoughPlace() {
		i = len(conns.RGet())
		conns.append(conn)
	} else {
		return -1
	}
	return i
}

type ConnectionsJSON struct {
	Capacity int           `json:"capacity"`
	Get      []*Connection `json:"get"`
}

func (conns *Connections) JSON() ConnectionsJSON {
	return ConnectionsJSON{
		Capacity: conns.Capacity(),
		Get:      conns.RGet(),
	}
}

func (conns *Connections) MarshalJSON() ([]byte, error) {
	return json.Marshal(conns.JSON())
}

func (conns *Connections) UnmarshalJSON(b []byte) error {
	temp := &ConnectionsJSON{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	conns.Set(temp.Get)
	conns.SetCapacity(temp.Capacity)

	return nil
}
