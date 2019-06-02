package game

import (
	"encoding/json"
	"fmt"
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

// Remove -> FastRemove
func (conns *Connections) FastRemove(conn *Connection) bool {
	if conn == nil {
		panic("123123123")
		return false
	}
	conn, i := conns.SearchByID(conn.ID())
	fmt.Println("remove conn", i)
	if i < 0 {
		return false
	}
	conns.remove(i)
	return true
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
		conn.setRoom(oldConn.Room())
		conn.SetIndex(oldConn.Index())

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

// ConnectionsJSON is a wrapper for sending Connections by JSON
type ConnectionsJSON struct {
	Capacity int           `json:"capacity"`
	Get      []*Connection `json:"get"`
}

// JSON convert Connections to ConnectionsJSON
func (conns *Connections) JSON() ConnectionsJSON {
	return ConnectionsJSON{
		Capacity: conns.Capacity(),
		Get:      conns.RGet(),
	}
}

// MarshalJSON - overriding the standard method json.Marshal
func (conns *Connections) MarshalJSON() ([]byte, error) {
	return json.Marshal(conns.JSON())
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (conns *Connections) UnmarshalJSON(b []byte) error {
	temp := &ConnectionsJSON{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	conns.Set(temp.Get)
	conns.SetCapacity(temp.Capacity)

	return nil
}
