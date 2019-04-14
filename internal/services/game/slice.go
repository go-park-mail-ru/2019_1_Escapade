package game

import (
	"crypto/rand"
	"fmt"
) //

// Connections - slice of connections with capacity
type Connections struct {
	Capacity int           `json:"capacity"`
	Get      []*Connection `json:"get"`
}

// NewConnections create instance of Connections
func NewConnections(capacity int) *Connections {
	return &Connections{capacity, make([]*Connection, 0, capacity)}
}

// Rooms - slice of rooms with capacity
type Rooms struct {
	Capacity int     `json:"capacity"`
	Get      []*Room `json:"get"`
}

// NewRooms create instance of Rooms
func NewRooms(capacity int) *Rooms {
	return &Rooms{capacity, make([]*Room, 0, capacity)}
}

// search element in slice
func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
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

////////////////// Rooms //////////////////////

// Empty check rooms length is 0
func (rooms *Rooms) Empty() bool {
	return len(rooms.Get) == 0
}

// Clear set slice to nil
func (rooms *Rooms) Clear() {
	rooms.Get = nil
}

// enoughPlace check that you can add more elements
func (rooms *Rooms) enoughPlace() bool {
	return len(rooms.Get) < rooms.Capacity
}

// SearchRoom find room with selected name and return it if success
// otherwise nil
func (rooms *Rooms) SearchRoom(name string) (i int, room *Room) {
	fmt.Println("search length", len(rooms.Get))
	for i, room = range rooms.Get {
		fmt.Println("try compare these:", room.Name, name)
		if room.Name == name {
			return
		}
		fmt.Println("no")
	}
	i, room = -1, nil
	return
}

// SearchPlayer find connection in rooms players and return it if success
// otherwise nil
func (rooms *Rooms) SearchPlayer(new *Connection) (old *Connection) {
	for _, room := range rooms.Get {
		i := room.Players.Search(new)
		// cant found
		if i < 0 {
			continue
		}
		/*
			Players are never deleted from rooms, because the room always
			must know its players. So, if you exit from room to lobby
			and refresh page you can return to room. To solve the  problem
			the following is done:
			When you exit to lobby, the room marks your connection
			as 'without room'(Connection has field - pointer to room.
			When it is nil, the connection connect to lobby). And
			when you disconnected from page(for example because of
			problems on the interner), the room doesnt change your room pointer
		*/
		foundConn := room.Players.Get[i]
		if foundConn.room != room {
			continue
		}

		return foundConn
	}
	return nil
}

// SearchObserver find connection in rooms obserers and return it if success
// otherwise nil
func (rooms *Rooms) SearchObserver(new *Connection) (old *Connection) {
	for _, room := range rooms.Get {
		i := room.Observers.Search(new)
		if i < 0 {
			continue
		}

		return room.Observers.Get[i]
	}
	return nil
}

// Add try add element if its possible. Return bool result
// if element not exists it will be create, otherwise it will change its value
func (rooms *Rooms) Add(room *Room) bool {
	if rooms.enoughPlace() {
		if i, _ := rooms.SearchRoom(room.Name); i < 0 {
			rooms.Get = append(rooms.Get, room)
		} else {
			rooms.Get[i] = room
		}
		return true
	}
	return false
}

// Remove delete element and decrement size if element
// exists in map
func (rooms *Rooms) Remove(room *Room) {
	size := len(rooms.Get)
	i, _ := rooms.SearchRoom(room.Name)
	if i < 0 {
		return
	}
	rooms.Get[i], rooms.Get[size-1] = rooms.Get[size-1], rooms.Get[i]
	rooms.Get[size-1] = nil
	rooms.Get = rooms.Get[:size-1]
}

// CopyLast copy Last element of 'from' to that as slice
func (rooms *Rooms) CopyLast(from *Rooms) {
	rooms.Capacity = 1
	rooms.Get = from.Get[len(from.Get)-1:]
}

////////////////// Connections //////////////////////

// CopyLast copy Last element of 'from' to that as slice
func (conns *Connections) CopyLast(from *Connections) {
	conns.Capacity = 1
	conns.Get = from.Get[len(from.Get)-1:]
}

// Clear set slice to nil
func (conns *Connections) Clear() {
	conns.Get = nil
}

// Empty check rooms capacity is 0
// it will happen, when finish is over, cause
// when somebody explodes, the capacity decrements
func (conns *Connections) Empty() bool {
	return len(conns.Get) == 0
}

// Search find connection in slice and return its index if success
// otherwise -1
func (conns *Connections) Search(conn *Connection) (i int) {
	return sliceIndex(len(conns.Get), func(i int) bool { return conns.Get[i].GetPlayerID() == conn.GetPlayerID() })
}

// enoughPlace check that you can add more elements
func (conns *Connections) enoughPlace() bool {
	return len(conns.Get) < conns.Capacity
}

// Add try add element if its possible. Return bool result
// if element not exists it will be create, otherwise it will change its value
func (conns *Connections) Add(conn *Connection) bool {
	if conns.enoughPlace() {
		if i := conns.Search(conn); i < 0 {
			conns.Get = append(conns.Get, conn)
		} else {
			conns.Get[i] = conn
		}
		return true
	}
	return false
}

// Remove delete element and decrement size if element
// exists in map
func (conns *Connections) Remove(conn *Connection) {
	size := len(conns.Get)
	i := conns.Search(conn)
	if i < 0 {
		return
	}
	conns.Get[i], conns.Get[size-1] = conns.Get[size-1], conns.Get[i]
	conns.Get[size-1] = nil
	conns.Get = conns.Get[:size-1]
	//sendError(conn, "Remove", "You disconnected ")
}
