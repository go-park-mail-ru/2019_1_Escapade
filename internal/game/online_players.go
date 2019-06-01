package game

// Rooms - slice of rooms with capacity
type Rooms struct {
	Capacity int     `json:"capacity"`
	Get      []*Room `json:"get"`
}

// Refresh reset the players, their connections and flags(generate new)
func (onlinePlayers *OnlinePlayers) Refresh(field Field) {
	size := onlinePlayers.Capacity()

	players := make([]Player, size)
	onlinePlayers.SetPlayers(players)

	flags := field.RandomFlags(players)
	onlinePlayers.SetFlags(flags)

	onlinePlayers.flagsLeft = size
	onlinePlayers.RefreshConnections()
}

// RefreshConnections reset connections
func (onlinePlayers *OnlinePlayers) RefreshConnections() {
	size := onlinePlayers.Capacity()
	onlinePlayers.Connections = *NewConnections(size)
}

// Init create players and flags
func (onlinePlayers *OnlinePlayers) Init(field *Field) {

	for i, conn := range onlinePlayers.Connections._get {
		if i > onlinePlayers.Capacity() {
			room := conn.Room()
			if room == nil {
				continue
			}
			room.LeavePlayer(conn)
			continue
		}
		onlinePlayers.SetPlayer(i, *NewPlayer(conn.User.ID))
		conn.SetIndex(i)
	}
	//onlinePlayers.Flags = field.RandomFlags(onlinePlayers.Players)

	return
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

// SearchIndexPlayer search connection index in the slice of Players
func (onlinePlayers *OnlinePlayers) SearchIndexPlayer(conn *Connection) (i int) {
	return sliceIndex(onlinePlayers.Capacity(), func(i int) bool {
		return onlinePlayers.Player(i).ID == conn.ID()
	})
}

// SearchConnection search connection index in the slice of connections
func (onlinePlayers *OnlinePlayers) SearchConnection(conn *Connection) (*Connection, int) {
	return onlinePlayers.Connections.SearchByID(conn.ID())
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

// Empty check rooms length is 0
func (rooms *Rooms) Empty() bool {
	return len(rooms.Get) == 0
}

// Empty check no connections connected
func (onlinePlayers *OnlinePlayers) Empty() bool {
	return onlinePlayers.Connections.Empty()
}

// Add try add element if its possible. Return bool result
// if element not exists it will be create, otherwise it will change its value
func (onlinePlayers *OnlinePlayers) Add(conn *Connection, cell Cell, kill bool, recover bool) bool {
	// if conn == nil {
	// 	panic(1)
	// }
	i := onlinePlayers.Connections.Add(conn, kill)
	if i < 0 {
		return false
	}
	onlinePlayers.SetPlayerID(i, conn.ID())
	conn.SetIndex(i)

	if !recover {

		if !onlinePlayers._flags[i].Set {
			onlinePlayers._flags[i] = Flag{
				Cell: cell,
				Set:  false,
			}
		}
	}

	//onlinePlayers.Connections.Get[i].SetIndex(i)
	return true
}

// EnoughPlace check that you can add more elements
func (onlinePlayers *OnlinePlayers) EnoughPlace() bool {
	return onlinePlayers.Connections.EnoughPlace()
}

// enoughPlace check that you can add more elements
func (rooms *Rooms) enoughPlace() bool {
	return len(rooms.Get) < rooms.Capacity
}

// SearchRoom find room with selected name and return it if success
// otherwise nil
func (rooms *Rooms) SearchRoom(id string) (i int, room *Room) {
	for i, room = range rooms.Get {
		if room.ID == id {
			return
		}
	}
	i, room = -1, nil
	return
}

// SearchPlayer find connection in rooms players and return it if success
// otherwise nil
func (rooms *Rooms) SearchPlayer(new *Connection, mustNotFinished bool) (int, *Room) {
	for _, room := range rooms.Get {
		i := room.Players.SearchIndexPlayer(new) //playersSearchIndexPlayer(new)

		// cant found
		if i < 0 {
			continue
		}
		if mustNotFinished && room.Players.Player(i).Finished {
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
		old, i := room.Observers.SearchByID(new.ID())
		if i < 0 {
			continue
		}

		return old
	}
	return nil
}

// Add try add element if its possible. Return bool result
// if element not exists it will be create, otherwise it will change its value
func (rooms *Rooms) Add(room *Room) bool {
	if i, _ := rooms.SearchRoom(room.ID); i >= 0 {
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
	i, _ := rooms.SearchRoom(room.ID)
	if i < 0 {
		return
	}
	rooms.Get[i], rooms.Get[size-1] = rooms.Get[size-1], rooms.Get[i]
	rooms.Get[size-1] = nil
	rooms.Get = rooms.Get[:size-1]
}
