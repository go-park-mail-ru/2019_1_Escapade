package game

import (
	"crypto/rand"
)

// Rooms - map with size and capacity
type Rooms struct {
	Capacity int
	Size     int
	Rooms    map[string]*Room
}

// NewRooms return new instance of Rooms
func NewRooms(capacity int) *Rooms {
	rooms := &Rooms{
		Capacity: capacity,
		Rooms:    make(map[string]*Room),
	}
	return rooms
}

// Empty check are rooms size 0
func (rooms *Rooms) Empty() bool {
	return rooms.Size == 0
}

// enoughPlace check that you can add more elements
func (rooms *Rooms) enoughPlace() bool {
	return rooms.Size < rooms.Capacity
}

// Add try add element if its possible. Return bool result
func (rooms *Rooms) Add(room *Room, name string) bool {
	if rooms.enoughPlace() {
		room.Name = name
		rooms.Rooms[room.Name] = room
		rooms.Size++
		return true
	}
	return false
}

// Remove delete element and decrement size if element
// exists in map
func (rooms *Rooms) Remove(room *Room) {
	if _, ok := rooms.Rooms[room.Name]; ok {
		delete(rooms.Rooms, room.Name)
		rooms.Size--
	}
}

// RandString create random string with n length
func RandString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
