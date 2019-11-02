package engine

// Init create players and flags
func (onlinePlayers *OnlinePlayers) Init(field *Field) {

	for i, conn := range onlinePlayers.Connections._get {
		// if i > onlinePlayers.Capacity() {
		// 	room := conn.PlayingRoom()
		// 	if room == nil {
		// 		continue
		// 	}
		// 	room.Leave(conn, true)
		// 	continue
		// }
		onlinePlayers.SetPlayer(i, *NewPlayer(conn.User.ID))
		conn.SetIndex(i)
	}

	return
}

// search element in slice
func sliceIndex(limit32 int32, predicate func(i int) bool) int {
	limit := int(limit32)
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
func (onlinePlayers *OnlinePlayers) SearchConnection(conn *Connection) (int, *Connection) {
	return onlinePlayers.Connections.SearchByID(conn.ID())
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
	i := onlinePlayers.Connections.Add(conn)
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
