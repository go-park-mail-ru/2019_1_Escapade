package engine

import "sync"

// Flag contaion Cell and flag, that it was set by User
//easyjson:json
type Flag struct {
	Cell Cell `json:"cell"`
	Set  bool `json:"set"`
}

// OnlinePlayers online players
type OnlinePlayers struct {
	m *OnlinePlayersMutex

	Connections *Connections
}

// NewConnections create instance of Connections
func newOnlinePlayers(size int32, field Field) *OnlinePlayers {
	players := make([]Player, size)
	flags := make([]Flag, size)
	var m = &OnlinePlayersMutex{
		capacityM: &sync.RWMutex{},
		_capacity: size,

		playersM: &sync.RWMutex{},
		_players: players,

		flagsM:     &sync.RWMutex{},
		_flags:     flags,
		_flagsLeft: size,
	}
	return &OnlinePlayers{
		m:           m,
		Connections: NewConnections(size),
	}
}

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
		onlinePlayers.m.SetPlayer(i, *NewPlayer(conn.User.ID))
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
	return sliceIndex(onlinePlayers.m.Capacity(), func(i int) bool {
		return onlinePlayers.m.Player(i).ID == conn.ID()
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
func (onlinePlayers *OnlinePlayers) Add(conn *Connection, cell Cell, recover bool) bool {
	i := onlinePlayers.Connections.Add(conn)
	if i < 0 {
		return false
	}
	onlinePlayers.m.SetPlayerID(i, conn.ID())
	conn.SetIndex(i)

	if !recover {
		if !onlinePlayers.m.Flag(i).Set {
			onlinePlayers.m.SetFlagByIndex(i, Flag{
				Cell: cell,
				Set:  false,
			})
		}
	}

	//onlinePlayers.Connections.Get[i].SetIndex(i)
	return true
}

// EnoughPlace check that you can add more elements
func (onlinePlayers *OnlinePlayers) EnoughPlace() bool {
	return onlinePlayers.Connections.EnoughPlace()
}

// Free free memory
func (onlinePlayers *OnlinePlayers) Free() {

	if onlinePlayers == nil {
		return
	}
	onlinePlayers.Connections.Free()
	onlinePlayers.m.Free()
}
