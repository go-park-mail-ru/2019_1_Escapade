package game

import (
	"encoding/json"
	"sync"
)

// Rooms - slice of rooms with capacity
type Rooms struct {
	capacityM *sync.RWMutex
	_capacity int

	getM *sync.RWMutex
	_get []*Room
}

// NewRooms create instance of Rooms
func NewRooms(capacity int) *Rooms {
	return &Rooms{
		capacityM: &sync.RWMutex{},
		_capacity: capacity,

		getM: &sync.RWMutex{},
		_get: make([]*Room, 0, capacity),
	}
}

// RoomsIterator is the iterator for the slice of rooms
type RoomsIterator struct {
	current int
	rooms   *Rooms
}

// NewRoomsIterator create new instance of RoomsIterator
func NewRoomsIterator(rooms *Rooms) *RoomsIterator {
	return &RoomsIterator{rooms: rooms, current: -1}
}

// RoomsJSON is a wrapper for sending Rooms by JSON
type RoomsJSON struct {
	Capacity int     `json:"capacity"`
	Get      []*Room `json:"get"`
}

// JSON convert Connections to ConnectionsJSON
func (rooms *Rooms) JSON() RoomsJSON {
	return RoomsJSON{
		Capacity: rooms.Capacity(),
		Get:      rooms.RGet(),
	}
}

// MarshalJSON - overriding the standard method json.Marshal
func (rooms *Rooms) MarshalJSON() ([]byte, error) {
	return json.Marshal(rooms.JSON())
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (rooms *Rooms) UnmarshalJSON(b []byte) error {
	temp := &RoomsJSON{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	rooms.Set(temp.Get)
	rooms.SetCapacity(temp.Capacity)

	return nil
}
