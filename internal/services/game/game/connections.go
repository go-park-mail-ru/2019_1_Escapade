package game

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// Connections - slice of connections with capacity
type Connections struct {
	capacityM *sync.RWMutex
	_capacity int32

	getM *sync.RWMutex
	_get []*Connection
}

// NewConnections create instance of Connections
func NewConnections(capacity int32) (conns *Connections) {
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
	//conns.getM.Lock()
	// for _, conn := range conns._get {
	// 	conn.Free()
	// }
	//conns._get = nil
	//conns.getM.Unlock()

	conns.capacityM.Lock()
	conns._capacity = 0
	conns.capacityM.Unlock()
}

// ConnectionsIterator is the iterator for the slice of connections
type ConnectionsIterator struct {
	current int
	conns   *Connections
}

// NewConnectionsIterator create new instance of ConnectionsIterator
func NewConnectionsIterator(conns *Connections) *ConnectionsIterator {
	return &ConnectionsIterator{conns: conns, current: -1}
}

// Remove remove connecrion from the slice of connections
func (conns *Connections) Remove(conn *Connection) bool {
	if conn == nil {
		utils.Debug(true, "no conn")
		return false
	}
	i, conn := conns.SearchByID(conn.ID())
	utils.Debug(false, "remove conn", i)
	if i < 0 {
		return false
	}
	conns.remove(i)
	return true
}

// Restore connection
func (conns *Connections) Restore(conn *Connection) bool {
	i, found := conns.SearchByID(conn.ID())
	if i != -1 {
		conn.Restore(found)
		sendAccountTaken(found)
		conns.set(i, conn)
	}
	return i != -1
}

// ConnectionsJSON is a wrapper for sending Connections by JSON
//easyjson:json
type ConnectionsJSON struct {
	Capacity int32         `json:"capacity"`
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
	return conns.JSON().MarshalJSON()
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (conns *Connections) UnmarshalJSON(b []byte) error {
	temp := &ConnectionsJSON{}

	if err := temp.UnmarshalJSON(b); err != nil {
		return err
	}

	conns.Set(temp.Get)
	conns.SetCapacity(temp.Capacity)

	return nil
}
