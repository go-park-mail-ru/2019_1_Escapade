package game

import (
	"encoding/json"
	"fmt"
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

type RoomsIterator struct {
	current int
	rooms   *Rooms
}

func NewRoomsIterator(rooms *Rooms) *RoomsIterator {
	return &RoomsIterator{rooms: rooms, current: -1}
}

// Remove -> FastRemove
func (rooms *Rooms) Remove(room *Room) bool {
	if room == nil {
		//panic("123123123")
		return false
	}
	i, room := rooms.Search(room.ID())
	fmt.Println("remove conn", i)
	if i < 0 {
		return false
	}
	rooms.remove(i)
	return true
}

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
