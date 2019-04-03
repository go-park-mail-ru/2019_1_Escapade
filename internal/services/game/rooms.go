package game

//re "escapade/internal/return_errors"
//"math/rand"

import (
	"crypto/rand"
)

type Rooms struct {
	Capacity int
	Size     int
	Rooms    map[string]*Room
}

func NewRooms(capacity int) *Rooms {
	rooms := &Rooms{
		Capacity: capacity,
		Rooms:    make(map[string]*Room),
	}
	return rooms
}

func (rooms *Rooms) Empty() bool {
	return rooms.Size == 0
}

func (rooms *Rooms) enoughPlace() bool {
	return rooms.Size < rooms.Capacity
}

func (rooms *Rooms) Add(room *Room, name string) bool {
	if rooms.enoughPlace() {
		room.Name = name
		rooms.Rooms[room.Name] = room
		rooms.Size++
		return true
	}
	return false
}

func (rooms *Rooms) Remove(room *Room) {
	delete(rooms.Rooms, room.Name)
	rooms.Size--
}

func RandString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
