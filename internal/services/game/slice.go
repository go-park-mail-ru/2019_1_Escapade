package game

import "fmt"

// OnlinePlayers online players
type OnlinePlayers struct {
	Capacity    int           `json:"capacity"`
	Players     []Player      `json:"players"`
	Flags       []Cell        `json:"flags"`
	Connections []*Connection `json:"-"`
}

// Connections - slice of connections with capacity
type Connections struct {
	Capacity int           `json:"capacity"`
	Get      []*Connection `json:"get"`
}

// Rooms - slice of rooms with capacity
type Rooms struct {
	Capacity int     `json:"capacity"`
	Get      []*Room `json:"get"`
}

// NewConnections create instance of Connections
func newOnlinePlayers(size int) *OnlinePlayers {
	return &OnlinePlayers{
		Capacity:    size,
		Players:     make([]Player, size),
		Flags:       make([]Cell, size),
		Connections: make([]*Connection, 0, size),
	}
}

// Init create players and flags
func (onlinePlayers *OnlinePlayers) Init(field *Field) {

	for i, conn := range onlinePlayers.Connections {
		if i > onlinePlayers.Capacity {
			conn.room.Leave(conn)
			continue
		}
		onlinePlayers.Players[i] = *NewPlayer(conn.User.ID)
		conn.index = i
	}
	onlinePlayers.Flags = field.RandomFlags(onlinePlayers.Players)

	return
}

// NewConnections create instance of Connections
func NewConnections(capacity int) *Connections {
	return &Connections{capacity,
		make([]*Connection, 0, capacity)}
}

// NewRooms create instance of Rooms
func NewRooms(capacity int) *Rooms {
	return &Rooms{capacity,
		make([]*Room, 0, capacity)}
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

// Free free memory
func (onlinePlayers *OnlinePlayers) Free() {

	if onlinePlayers == nil {
		return
	}
	for i := 0; i < len(onlinePlayers.Connections); i++ {
		onlinePlayers.Connections[i] = nil
	}

	onlinePlayers.Players = nil
	onlinePlayers.Flags = nil
	onlinePlayers.Connections = nil
	onlinePlayers = nil
}

// Search find connection in slice and return its index if success
// otherwise -1
func (onlinePlayers *OnlinePlayers) Search(conn *Connection) (i int) {
	return sliceIndex(onlinePlayers.Capacity, func(i int) bool {
		return onlinePlayers.Players[i].ID == conn.ID()
	})
}

// Free free memory
func (rooms *Rooms) Free() {
	if rooms == nil {
		return
	}
	for _, room := range rooms.Get {
		room.Free()
	}
	rooms.Get = nil
	rooms.Capacity = 0
	rooms = nil
}

// Free free memory
func (conns *Connections) Free() {
	if conns == nil {
		return
	}
	for _, conn := range conns.Get {
		conn.Free()
	}
	conns.Get = nil
	conns.Capacity = 0
}

// Empty check rooms length is 0
func (rooms *Rooms) Empty() bool {
	return len(rooms.Get) == 0
}

// Empty check no connections connected
func (onlinePlayers *OnlinePlayers) Empty() bool {
	return len(onlinePlayers.Connections) == 0
}

// Add try add element if its possible. Return bool result
// if element not exists it will be create, otherwise it will change its value
func (onlinePlayers *OnlinePlayers) Add(conn *Connection) bool {
	var i int
	if i = onlinePlayers.Search(conn); i >= 0 {
		onlinePlayers.Connections[i] = conn
	} else if onlinePlayers.enoughPlace() {
		i = len(onlinePlayers.Connections)
		onlinePlayers.Connections = append(onlinePlayers.Connections, conn)
	} else {
		return false
	}
	onlinePlayers.Players[i].ID = onlinePlayers.Connections[i].ID()
	onlinePlayers.Connections[i].index = i
	return false
}

// Remove delete element and decrement size if element
// exists in map
func (onlinePlayers *OnlinePlayers) Remove(conn *Connection) {
	size := len(onlinePlayers.Connections)
	i := onlinePlayers.Search(conn)
	if i < 0 {
		fmt.Println("cant found", i, size)
		return
	}
	fmt.Println("i and size,", i, size)
	onlinePlayers.Connections[i], onlinePlayers.Connections[size-1] = onlinePlayers.Connections[size-1], onlinePlayers.Connections[i]
	onlinePlayers.Connections[size-1] = nil
	onlinePlayers.Connections = onlinePlayers.Connections[:size-1]
	//sendError(conn, "Remove", "You disconnected ")
}

// enoughPlace check that you can add more elements
func (onlinePlayers *OnlinePlayers) enoughPlace() bool {
	return len(onlinePlayers.Connections) < onlinePlayers.Capacity
}

// enoughPlace check that you can add more elements
func (rooms *Rooms) enoughPlace() bool {
	return len(rooms.Get) < rooms.Capacity
}

// SearchRoom find room with selected name and return it if success
// otherwise nil
func (rooms *Rooms) SearchRoom(name string) (i int, room *Room) {
	for i, room = range rooms.Get {
		if room.Name == name {
			return
		}
	}
	i, room = -1, nil
	return
}

// SearchPlayer find connection in rooms players and return it if success
// otherwise nil
func (rooms *Rooms) SearchPlayer(new *Connection) (int, *Room) {
	for _, room := range rooms.Get {
		i := room.Players.Search(new)
		// cant found
		if i < 0 {
			continue
		}
		player := room.Players.Players[i]
		if player.Finished {
			continue
		}

		return i, room
	}
	return -1, nil
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
	if i, _ := rooms.SearchRoom(room.Name); i >= 0 {
		rooms.Get[i] = room
	} else if rooms.enoughPlace() {
		rooms.Get = append(rooms.Get, room)
	} else {
		return false
	}
	return true
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

// Empty check rooms capacity is 0
// it will happen, when finish is over, cause
// when somebody explodes, the capacity decrements
func (conns *Connections) Empty() bool {
	return len(conns.Get) == 0
}

// Search find connection in slice and return its index if success
// otherwise -1
func (conns *Connections) Search(conn *Connection) (i int) {
	return sliceIndex(len(conns.Get), func(i int) bool { return conns.Get[i].ID() == conn.ID() })
}

// enoughPlace check that you can add more elements
func (conns *Connections) enoughPlace() bool {
	return len(conns.Get) < conns.Capacity
}

// Add try add element if its possible. Return bool result
// if element not exists it will be create, otherwise it will change its value
func (conns *Connections) Add(conn *Connection) bool {
	if i := conns.Search(conn); i >= 0 {
		conns.Get[i] = conn
	} else if conns.enoughPlace() {
		conns.Get = append(conns.Get, conn)
	} else {
		return false
	}
	return true
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
