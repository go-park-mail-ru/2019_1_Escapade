package game

import (
	"sync"
)

// Rooms - slice of rooms with capacity
type Rooms struct {
	capacityM *sync.RWMutex
	_capacity int32

	getM *sync.RWMutex
	_get []*Room
}

// NewRooms create instance of Rooms
func NewRooms(capacity int32) *Rooms {
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
//easyjson:json
type RoomsJSON struct {
	Capacity int32   `json:"capacity"`
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
	return rooms.JSON().MarshalJSON()
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (rooms *Rooms) UnmarshalJSON(b []byte) error {
	temp := &RoomsJSON{}

	if err := temp.UnmarshalJSON(b); err != nil {
		return err
	}

	rooms.Set(temp.Get)
	rooms.SetCapacity(temp.Capacity)

	return nil
}
