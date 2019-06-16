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

type ConnectionsIterator struct {
	current int
	conns   *Connections
}

func NewConnectionsIterator(conns *Connections) *ConnectionsIterator {
	return &ConnectionsIterator{conns: conns, current: -1}
}

// Remove -> FastRemove
func (conns *Connections) Remove(conn *Connection) bool {
	if conn == nil {
		//panic("123123123")
		return true
	}
	conn, i := conns.SearchByID(conn.ID())
	fmt.Println("remove conn", i)
	if i < 0 {
		return false
	}
	conns.remove(i)
	return true
}

// Restore connection
func (conns *Connections) Restore(conn *Connection) bool {
	found, i := conns.SearchByID(conn.ID())
	if i != -1 {
		conn.Restore(found)
		sendAccountTaken(*found)
		conns.set(i, conn)
	}
	return i != -1
}

// Add try add element if its possible. Return bool result
// if element not exists it will be create, otherwise it will change its value
func (conns *Connections) Add(conn *Connection /*, copy bool*/) int {
	_, i := conns.SearchByID(conn.ID())
	if i >= 0 {
		conns.set(i, conn)
	} else if conns.EnoughPlace() {
		i = conns.len()
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
		Get:      conns.NEWRGet(),
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
