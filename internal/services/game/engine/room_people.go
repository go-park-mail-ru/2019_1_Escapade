package engine

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	action_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/action"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
)

type PeopleI interface {
	synced.PublisherI

	Connections() []*Connections

	Enter(conn *Connection, isPlayer bool, needRecover bool) bool

	Remove(conn *Connection)
	Search(find *Connection) (*Connection, bool)

	isAlive(conn *Connection) bool

	ForEach(action func(c *Connection, isPlayer bool))

	players() *OnlinePlayers
	observers() *Connections

	Flag(index int) Flag

	getPlayer(conn *Connection) Player
	PlayersSlice() []Player

	IsWinner(index int) bool
	Winners() []int

	configure(info models.GameInformation)

	setFlag(conn *Connection, cell Cell)
	OpenCell(conn *Connection, cell *Cell)

	flagExists(cell Cell, this *Connection) (bool, *Connection)
}

type RoomPeople struct {
	synced.PublisherBase

	s  synced.SyncI
	se RSendI
	c  RClientI
	e  EventsI
	l  LobbyProxyI
	f  FieldProxyI

	pointsPerCellK float64

	winnersM *sync.Mutex
	_winners []int

	Players *OnlinePlayers

	isDeathmatch bool

	Observers *Connections
	killedM   *sync.RWMutex
	_killed   int32 //amount of killed users
}

func (room *RoomPeople) Enter(conn *Connection, isPlayer bool, needRecover bool) bool {
	var result bool
	room.s.DoWithOther(conn, func() {
		if !room.add(conn, isPlayer, needRecover) {
			return
		}
		if isPlayer {
			room.notify(action_.ConnectAsPlayer, conn)
		} else {
			room.notify(action_.ConnectAsObserver, conn)
		}
		result = true
	})
	return result
}

// init struct's values
func (room *RoomPeople) init(rs *models.RoomSettings) {
	sq := rs.Width * rs.Height
	room.pointsPerCellK = 1000 / float64(sq)

	room.winnersM = &sync.Mutex{}
	room._winners = nil
	room.killedM = &sync.RWMutex{}
	room.setKilled(0)
	room.Players = newOnlinePlayers(rs.Players)
	room.Observers = NewConnections(rs.Observers)
	room.isDeathmatch = rs.Deathmatch

	room.PublisherBase = *synced.NewPublisher()
}

// build components
func (room *RoomPeople) build(builder RBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildSender(&room.se)
	builder.BuildConnectionEvents(&room.c)
	builder.BuildEvents(&room.e)
	builder.BuildLobby(&room.l)
	builder.BuildField(&room.f)
}

func (room *RoomPeople) subscribe() {
	room.eventsSubscribe()
	room.connectionSubscribe(room.c)
}

// Init configure dependencies with other components of the room
func (room *RoomPeople) Init(builder RBuilderI, rs *models.RoomSettings) {
	room.init(rs)
	room.build(builder)
}

func (room *RoomPeople) start() {
	room.StartPublish()
}

func (room *RoomPeople) finish() {
	go room.Players.Free()
	go room.Observers.Free()
	room.PublisherBase.StopPublish()
}

func (room *RoomPeople) startGame() {
	room.s.Do(func() {
		room.ForEach(func(c *Connection, isPlayer bool) {
			room.se.Room(c)
		})

		room.f.Field().Fill(room.Players.m.Flags())
		room.Players.Init()
	})
}

func (room *RoomPeople) finishGame() {
	room.s.Do(func() {
		room.Players.m.Finish()
	})
}

func (room *RoomPeople) configure(info models.GameInformation) {
	room.s.Do(func() {
		room.setKilled(info.Game.Settings.Players)
		room.Players = newOnlinePlayers(info.Game.Settings.Players)
		for i, gamer := range info.Gamers {
			room.Players.m.SetPlayer(i, Player{
				ID:       gamer.ID,
				Points:   gamer.Score,
				Died:     gamer.Explosion,
				Finished: true,
			})
		}
	})
}

/* use ApplyActionToAll with func() instead of LeaveAll
func (room *RoomPeople) LeaveAll() {
	playersIterator := NewConnectionsIterator(room.Players.Connections)
	for playersIterator.Next() {
		player := playersIterator.Value()
		go room.r.connEvents.Leave(player)
	}

	observersIterator := NewConnectionsIterator(room.Observers)
	for observersIterator.Next() {
		observer := observersIterator.Value()
		go room.r.connEvents.Leave(observer)
	}
}*/

func (room *RoomPeople) Flag(index int) Flag {
	return room.Players.m.Flag(index)
}

func (room *RoomPeople) setFlag(conn *Connection, cell Cell) {
	room.Players.m.SetFlag(conn, cell, func() {
		room.e.UpdateStatus(room_.StatusRunning)
	})
}

func (room *RoomPeople) PlayersSlice() []Player {
	return room.Players.m.RPlayers()
}

func (room *RoomPeople) players() *OnlinePlayers {
	return room.Players
}

func (room *RoomPeople) observers() *Connections {
	return room.Observers
}

func (room *RoomPeople) ForEach(action func(c *Connection, isPlayer bool)) {
	room.s.Do(func() {
		playersIterator := NewConnectionsIterator(room.Players.Connections)
		for playersIterator.Next() {
			player := playersIterator.Value()
			action(player, true)
		}

		observersIterator := NewConnectionsIterator(room.Observers)
		for playersIterator.Next() {
			observer := observersIterator.Value()
			action(observer, false)
		}
	})
}

func (room *RoomPeople) Connections() []*Connections {
	players := room.Players.Connections
	observers := room.Observers
	people := make([]*Connections, 0)
	return append(people, players, observers)
}

// Empty check room has no people
func (room *RoomPeople) Empty() bool {
	var result = true
	room.s.Do(func() {
		result = room.Players.Connections.len()+room.Observers.len() == 0
	})
	return result
}

func (room *RoomPeople) IncreasePoints(index int, points float64) {
	room.Players.m.IncreasePlayerPoints(index, points)
}

func (room *RoomPeople) OpenCell(conn *Connection, cell *Cell) {
	room.s.DoWithOther(conn, func() {
		switch {
		case cell.Value < CellMine:
			room.openSafeCell(conn, cell)
		case cell.Value == CellMine:
			room.openMine(conn)
		case cell.Value > CellIncrement:
			room.openFlag(conn, cell)
		}
	})
}

// openFlag is called, when somebody find cell flag
func (room *RoomPeople) openSafeCell(conn *Connection, cell *Cell) {
	room.s.DoWithOther(conn, func() {
		points := float64(cell.Value) * room.pointsPerCellK
		room.IncreasePoints(conn.Index(), points)
	})
}

// openFlag is called, when somebody find cell flag
func (room *RoomPeople) openMine(conn *Connection) {
	room.s.DoWithOther(conn, func() {
		room.IncreasePoints(conn.Index(), float64(-1000))
		room.kill(conn, action_.Explode)
	})
}

// openFlag is called, when somebody find cell flag
func (room *RoomPeople) openFlag(founder *Connection, found *Cell) {
	room.s.Do(func() {
		var which int32
		room.Players.ForEachFlag(room.findFlagOwner(found, &which))
		if which == founder.ID() {
			return
		}

		room.Players.m.IncreasePlayerPoints(founder.Index(), 300)
		index, killConn := room.Players.Connections.SearchByID(which)
		if index >= 0 {
			room.kill(killConn, action_.FlagLost)
		}
	})
}

func (room *RoomPeople) kill(conn *Connection, action int32) {
	if !room.isAlive(conn) || !room.e.IsActive() {
		return
	}
	room.setFinished(conn)
	room.notify(action, ConnectionMsg{
		connection: conn,
		content:    room.isDeathmatch,
	})
}

func (room *RoomPeople) findFlagOwner(found *Cell, which *int32) func(int, Flag) {
	return func(index int, flag Flag) {
		if flag.Cell.X == found.X && flag.Cell.Y == found.Y {
			*which = flag.Cell.PlayerID
		}
	}
}

// Search the connection in players slice and observers slice of room
// return connection and flag isPlayer
func (room *RoomPeople) Search(find *Connection) (*Connection, bool) {
	i, found := room.Players.SearchConnection(find)
	if i >= 0 {
		return found, true
	}
	i, found = room.Observers.SearchByID(find.ID())
	if i >= 0 {
		return found, false
	}
	return nil, true
}

// add connection to room as player if 'isPlayer' is true, else as observer
//   if 'needRecover' is true, then the user will be found among the members
//   of the room and "restored to rights". Return true if the add or restore
//   was successful. Otherwise false
func (room *RoomPeople) add(conn *Connection, isPlayer bool, needRecover bool) bool {
	var result bool
	room.s.DoWithOther(conn, func() {
		result = room.push(conn, isPlayer, needRecover)
		if !result {
			return
		}
		var status int32
		if isPlayer {
			status = room_.PlayerEnter
		} else {
			status = room_.ObserverEnter
		}
		room.notify(status, ConnectionMsg{
			connection: conn,
			content:    needRecover,
		})
	})
	return result
}

// Push add the connection to the room.
// isPlayer - if true, the connection will add as player, otherwise as observer
// needRecover - if true, then the connection has already added to the room and
// 	it must be restored
// Returns true if added, otherwise false
// If the game has already started, then the connection from waiter slice goes
// to player slice. Otherwise if the game is looking for people then
// the connection remains the waiter, but gets waiting room - this one
func (room *RoomPeople) push(conn *Connection, isPlayer bool, needRecover bool) bool {
	var result bool
	room.s.DoWithOther(conn, func() {
		if isPlayer {
			if !needRecover && !room.Players.EnoughPlace() {
				result = false
				return
			}
			room.Players.Add(conn, room.f.Field().RandomFlag(conn.ID()), needRecover)
			if !needRecover && !room.Players.EnoughPlace() {
				room.e.UpdateStatus(room_.StatusFlagPlacing)

			}
		} else {
			if !needRecover && !room.Observers.EnoughPlace() {
				result = false
				return
			}
			room.Observers.Add(conn)
		}

		if room.e.Status() != room_.StatusRecruitment {
			room.l.WaiterToPlayer(conn)
		} else {
			room.l.setWaitingRoom(conn)
		}

		result = true
	})
	return result
}

func (room *RoomPeople) notify(code int32, extra interface{}) {
	room.Notify(synced.Msg{
		Publisher: room_.UpdatePeople,
		Action:    code,
		Extra:     extra,
	})
}

// isAlive check if connection is player and he is not died
func (room *RoomPeople) isAlive(conn *Connection) bool {
	index := conn.Index()
	return index >= 0 && !room.Players.m.Player(index).Finished
}

func (room *RoomPeople) getPlayer(conn *Connection) Player {
	return room.Players.m.Player(conn.Index())
}

func (room *RoomPeople) flagExists(cell Cell, this *Connection) (bool, *Connection) {
	var (
		player int
		found  bool
		flags  = room.Players.m.Flags()
	)
	for index, flag := range flags {
		if (flag.Cell.X == cell.X) && (flag.Cell.Y == cell.Y) {
			if this == nil || index != this.Index() {
				found = true
				player = index
			}
			break
		}
	}
	if !found {
		return false, nil
	}
	conn := room.Players.Connections.SearchByIndex(player)
	return found, conn
}

func (room *RoomPeople) Remove(conn *Connection) {
	room.s.Do(func() {
		if room.e.Status() == room_.StatusRecruitment {
			room.Players.Connections.Remove(conn)
		}
		room.Observers.Remove(conn)
		if room.Empty() {
			room.notify(room_.AllExit, nil)
		}
	})
}

// IsWinner is player wuth id playerID is winner
func (room *RoomPeople) IsWinner(index int) bool {
	var found bool
	room.s.Do(func() {
		winners := room.Winners()
		for _, i := range winners {
			utils.Debug(false, "compare:", i, index)
			if i == index {
				found = true
				return
			}
		}
	})
	return found
}

// Winners determine who won the game
func (room *RoomPeople) Winners() []int {
	var winners []int
	room.s.Do(func() {
		room.winnersM.Lock()
		if room._winners == nil {
			room._winners = make([]int, 0)
			room.Players.ForEach(room.addWinner(&room._winners))
		}
		winners = room._winners
		room.winnersM.Unlock()
	})
	return winners
}

func (room *RoomPeople) addWinner(winners *[]int) func(int, Player) {
	max := 0.
	return func(index int, player Player) {
		if !player.Died {
			if player.Points > max {
				max = player.Points
				*winners = []int{index}
			} else if player.Points == max {
				*winners = append(*winners, index)
			}
		}
	}
}

////////////// mutex

// setFinished set player finished
func (room *RoomPeople) setFinished(conn *Connection) {
	room.s.Do(func() {
		index := conn.Index()
		if index < 0 {
			return
		}
		room.Players.m.PlayerFinish(index)

		room.killedM.Lock()
		room._killed++
		room.killedM.Unlock()
	})
}

// done return '_killed' field
func (room *RoomPeople) killed() int32 {
	room.killedM.RLock()
	v := room._killed
	room.killedM.RUnlock()
	return v
}

// incrementKilled increment amount of killed
func (room *RoomPeople) incrementKilled() {
	var (
		allKilled bool
		capacity  = room.Players.m.Capacity()
	)
	room.killedM.Lock()
	room._killed++
	allKilled = room._killed >= capacity
	room.killedM.Unlock()
	if allKilled {
		room.notify(room_.AllDied, nil)
	}
}

// setKilled set new value of killed
func (room *RoomPeople) setKilled(killed int32) {
	room.killedM.Lock()
	room._killed = killed
	room.killedM.Unlock()
}

//////////////////////////// callbacks

// eventsSubscribe subscibe to events associated with room's status
func (room *RoomPeople) eventsSubscribe() {
	observer := synced.NewObserver(
		synced.NewPairNoArgs(room_.StatusRecruitment, room.start),
		synced.NewPairNoArgs(room_.StatusFlagPlacing, room.startGame),
		synced.NewPairNoArgs(room_.StatusFinished, room.finishGame),
		synced.NewPairNoArgs(room_.StatusAborted, room.finish))
	room.e.Observe(observer.AddPublisherCode(room_.UpdateStatus))
}

// connectionBackToLobby is called when user went to lobby
func (room *RoomPeople) connectionBackToLobby(msg synced.Msg) {
	action, ok := msg.Extra.(ConnectionMsg)
	if !ok {
		return
	}
	room.Remove(action.connection)
	room.kill(action.connection, msg.Action)
}

// connectionReconnect is called when user tryes to reconnect
func (room *RoomPeople) connectionReconnect(msg synced.Msg) {
	action, ok := msg.Extra.(ConnectionMsg)
	if !ok {
		return
	}
	isPlayer, ok := action.content.(bool)
	if !ok {
		return
	}
	room.add(action.connection, isPlayer, true)
}

// connectionConnectAsPlayer is called when user wants to connect as player
func (room *RoomPeople) connectionConnectAsPlayer(msg synced.Msg) {
	action, ok := msg.Extra.(ConnectionMsg)
	if !ok {
		return
	}
	room.add(action.connection, true, false)
}

// connectionGiveUp is called when user wants to connect as observer
func (room *RoomPeople) connectionConnectAsObserver(msg synced.Msg) {
	action, ok := msg.Extra.(ConnectionMsg)
	if !ok {
		return
	}
	room.add(action.connection, false, false)
}

// connectionGiveUp is called when user has surrendered
func (room *RoomPeople) connectionGiveUp(msg synced.Msg) {
	action, ok := msg.Extra.(ConnectionMsg)
	if !ok {
		return
	}
	room.kill(action.connection, msg.Action)
}

// connectionDisconnect is called when user disconnected
func (room *RoomPeople) connectionDisconnect(msg synced.Msg) {
	action, ok := msg.Extra.(ConnectionMsg)
	if !ok {
		return
	}
	found, _ := room.Search(action.connection)
	if found == nil {
		return
	}
	found.setDisconnected()
}

// connectionSubscribe subscibe to events associated with connection's events
func (room *RoomPeople) connectionSubscribe(c RClientI) {
	observer := synced.NewObserver(
		synced.NewPair(action_.BackToLobby, room.connectionBackToLobby),
		synced.NewPair(action_.Reconnect, room.connectionReconnect),
		synced.NewPair(action_.ConnectAsObserver, room.connectionConnectAsObserver),
		synced.NewPair(action_.ConnectAsPlayer, room.connectionConnectAsPlayer),
		synced.NewPair(action_.GiveUp, room.connectionGiveUp),
		synced.NewPair(action_.Disconnect, room.connectionDisconnect))
	c.Observe(observer.AddPublisherCode(room_.UpdateConnection))
}

// 445 -> 615
