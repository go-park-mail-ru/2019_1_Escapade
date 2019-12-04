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
	Connections() []*Connections

	Remove(conn *Connection)
	add(conn *Connection, isPlayer bool, needRecover bool) bool
	Search(find *Connection) (*Connection, bool)

	isAlive(conn *Connection) bool
	SetFinished(conn *Connection)

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

	EventsSub() synced.SubscriberI
	ConnectionSub() synced.SubscriberI
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

func (room *RoomPeople) Init(builder RBuilderI, rs *models.RoomSettings) {
	builder.BuildSync(&room.s)
	builder.BuildSender(&room.se)
	builder.BuildConnectionEvents(&room.c)
	builder.BuildEvents(&room.e)
	builder.BuildLobby(&room.l)
	builder.BuildField(&room.f)

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
	room.PublisherBase.Start(room.e.PeopleSub())
}

func (room *RoomPeople) free() {
	go room.Players.Free()
	go room.Observers.Free()
	room.PublisherBase.Stop()
}

func (room *RoomPeople) start() {
	room.s.Do(func() {
		room.ForEach(func(c *Connection, isPlayer bool) {
			room.se.Room(c)
		})

		room.f.Field().Fill(room.Players.m.Flags())
		room.Players.Init()
	})
}

func (room *RoomPeople) finish() {
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
	room.Players.m.SetFlag(conn, cell, room.e.PrepareOver)
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
	room.SetFinished(conn)
	room.notifyWithConn(conn, action, room.isDeathmatch)
}

func (room *RoomPeople) findFlagOwner(found *Cell, which *int32) func(int, Flag) {
	return func(index int, flag Flag) {
		if flag.Cell.X == found.X && flag.Cell.Y == found.Y {
			*which = flag.Cell.PlayerID
		}
	}
}

// search the connection in players slice and observers slice of room
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

func (room *RoomPeople) add(conn *Connection, isPlayer bool, needRecover bool) bool {
	var result bool
	room.s.DoWithOther(conn, func() {
		result = room.push(conn, isPlayer, needRecover)
		if !result {
			return
		}
		if isPlayer {
			room.notifyWithConn(conn, room_.PlayerEnter, needRecover)
		} else {
			room.notifyWithConn(conn, room_.ObserverEnter, needRecover)
		}
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
				room.e.RecruitingOver()
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

func (room *RoomPeople) notify(code int) {
	room.Notify(synced.Msg{
		Code:    room_.UpdatePeople,
		Content: code,
	})
}

func (room *RoomPeople) notifyWithConn(conn *Connection, code int32, content interface{}) {
	room.Notify(synced.Msg{
		Code: room_.UpdatePeople,
		Content: ConnectionMsg{
			connection: conn,
			code:       code,
			content:    content,
		},
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
			room.notify(room_.AllExit)
		}
	})
}

////////////// mutex

// SetFinished set player finished
func (room *RoomPeople) SetFinished(conn *Connection) {
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
		capacity  = room.players().m.Capacity()
	)
	room.killedM.Lock()
	room._killed++
	allKilled = room._killed >= capacity
	room.killedM.Unlock()
	if allKilled {
		room.notify(room_.AllDied)
	}
}

// setKilled set new value of killed
func (room *RoomPeople) setKilled(killed int32) {
	room.killedM.Lock()
	room._killed = killed
	room.killedM.Unlock()
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

//////////////////////////// callbacks

func (room *RoomPeople) EventsSub() synced.SubscriberI {
	return synced.NewSubscriber(room.eventsCallback)
}

func (room *RoomPeople) eventsCallback(msg synced.Msg) {
	if msg.Code != room_.UpdateStatus {
		return
	}
	code, ok := msg.Content.(int)
	if !ok {
		return
	}
	switch code {
	case room_.StatusFinished:
		room.finish()
	case room_.StatusAborted:
		room.free()
	case room_.StatusFlagPlacing:
		room.start()
	}
}

func (room *RoomPeople) ConnectionSub() synced.SubscriberI {
	return synced.NewSubscriber(room.connectionCallback)
}

func (room *RoomPeople) connectionCallback(msg synced.Msg) {
	if msg.Code != room_.UpdateConnection {
		return
	}
	action, ok := msg.Content.(ConnectionMsg)
	if !ok {
		return
	}
	switch action.code {
	case action_.BackToLobby:
		room.Remove(action.connection)
		room.kill(action.connection, action.code)
	case action_.Reconnect:
		isPlayer, ok := action.content.(bool)
		if !ok {
			return
		}
		room.add(action.connection, isPlayer, true)
	case action_.ConnectAsPlayer:
		room.add(action.connection, true, false)
	case action_.ConnectAsObserver:
		room.add(action.connection, false, false)
	case action_.GiveUp:
		room.kill(action.connection, action.code)
	}

}

// 445
