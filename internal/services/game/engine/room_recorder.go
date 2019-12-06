package engine

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	action_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/action"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
)

// ActionRecorderI control access to actions history
// Proxy Pattern
type ActionRecorderI interface {
	Restart(conn *Connection)
	Disconnect(conn *Connection)
	FlagСonflict(conn *Connection)
	FlagSet(conn *Connection)
	Leave(conn *Connection, action int32, isPlayer bool)

	ModelActions() []models.Action
	history() []*action_.PlayerAction

	AddConnection(conn *Connection, isPlayer bool, needRecover bool)
	setHistory(history []*action_.PlayerAction)

	Kill(conn *Connection, action int32, isDeathmatch bool)

	configure(info []models.Action)
}

// RoomRecorder notify actions to room history and to users
// implements ActionRecorderProxyI
type RoomRecorder struct {
	s  synced.SyncI
	i  RoomInformationI
	p  PeopleI
	f  FieldProxyI
	se RSendI
	mo RModelsI

	historyM *sync.RWMutex
	_history []*action_.PlayerAction
}

// init struct's values
func (room *RoomRecorder) init() {
	room.historyM = &sync.RWMutex{}
	room.setHistory(make([]*action_.PlayerAction, 0))
}

// build components
func (room *RoomRecorder) build(builder RBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildInformation(&room.i)
	builder.BuildPeople(&room.p)
	builder.BuildField(&room.f)
	builder.BuildSender(&room.se)
	builder.BuildModelsAdapter(&room.mo)
}

// subscribe to room events
func (room *RoomRecorder) subscribe(builder RBuilderI) {
	var (
		events EventsI
		client RClientI
	)
	builder.BuildEvents(&events)
	builder.BuildConnectionEvents(&client)

	room.eventsSubscribe(events)
	room.peopleSubscribe(room.p)
	room.connectionSubscribe(client)
}

// Init configure dependencies with other components of the room
func (room *RoomRecorder) Init(builder RBuilderI) {
	room.init()
	room.build(builder)
	room.subscribe(builder)
}

func (room *RoomRecorder) finish() {
	room.historyFree()
}

// Leave update metainformation about user leaving room
func (room *RoomRecorder) Leave(conn *Connection, action int32, isPlayer bool) {
	if isPlayer {
		room.se.PlayerExit(conn)
	} else {
		room.se.ObserverExit(conn)
	}
	room.notifyAll(conn, action)
}

func (room *RoomRecorder) Disconnect(conn *Connection) {
	room.notifyAll(conn, action_.Disconnect)
}

func (room *RoomRecorder) FlagSet(conn *Connection) {
	room.notifyAll(conn, action_.FlagSet)
}

func (room *RoomRecorder) Restart(conn *Connection) {
	room.notifyAll(conn, action_.Restart)
}

func (room *RoomRecorder) flag(conn *Connection) {
	cell := room.p.Flag(conn.Index())
	room.f.SaveCell(&cell.Cell)
	room.se.NewCells(cell.Cell)
}

func (room *RoomRecorder) Kill(conn *Connection,
	action int32, isDeathmatch bool) {
	room.notifyAll(conn, action)
	room.s.Do(func() {
		if isDeathmatch {
			room.flag(conn)
		}
	})
}

func (room *RoomRecorder) ModelActions() []models.Action {
	history := room.history()
	actions := make([]models.Action, 0)
	room.s.Do(func() {
		for _, action := range history {
			actions = append(actions, action.ToModel())
		}
	})
	return actions
}

func (room *RoomRecorder) Reconnect(conn *Connection) {
	room.notifyAll(conn, action_.Reconnect)
}

func (room *RoomRecorder) AddPlayer(conn *Connection) {
	room.notifyAll(conn, action_.ConnectAsPlayer)
	room.se.PlayerEnter(conn)
}

func (room *RoomRecorder) AddObserver(conn *Connection) {
	room.notifyAll(conn, action_.ConnectAsObserver)
	room.se.ObserverEnter(conn)
}

func (room *RoomRecorder) AddConnection(conn *Connection, isPlayer bool, needRecover bool) {
	if needRecover {
		room.Reconnect(conn)
	} else if isPlayer {
		room.AddPlayer(conn)
	} else {
		room.AddObserver(conn)
	}
	room.se.StatusToOne(conn)
	room.se.Room(conn)
}

func (room *RoomRecorder) FlagСonflict(conn *Connection) {
	room.notifyAll(conn, action_.FlagСonflict)
}

func (room *RoomRecorder) notifyAll(conn *Connection, action int32) {
	room.s.Do(func() {
		pa := action_.NewPlayerAction(conn.ID(), action)
		room.appendAction(pa)
		room.se.Action(*pa, room.se.AllExceptThat(conn))
	})
}

func (room *RoomRecorder) configure(info []models.Action) {
	room.setHistory(make([]*action_.PlayerAction, 0))
	for _, actionDB := range info {
		var action = &action_.PlayerAction{}
		room.appendAction(action.FromModel(actionDB))
	}
}

/////////////////////////////// mutex

// history get the history of the actions that occurred in the room
func (room *RoomRecorder) history() []*action_.PlayerAction {
	room.historyM.RLock()
	v := room._history
	room.historyM.RUnlock()
	return v
}

// setHistory set the history of the actions that occurred in the room
func (room *RoomRecorder) setHistory(history []*action_.PlayerAction) {
	room.historyM.Lock()
	room._history = history
	room.historyM.Unlock()
}

// appendAction append action to the history of actions
func (room *RoomRecorder) appendAction(action *action_.PlayerAction) {
	room.historyM.Lock()
	defer room.historyM.Unlock()
	room._history = append(room._history, action)
}

// historyFree clear the history of actions
func (room *RoomRecorder) historyFree() {
	room.historyM.Lock()
	room._history = nil
	room.historyM.Unlock()
}

///////////////////////// subscripe

// actionBackToLobby is called when connection want to go to lobby
func (room *RoomRecorder) actionBackToLobby(msg synced.Msg) {
	if conn, ok := room.connectionCheck(msg); ok {
		room.Leave(conn, msg.Action, conn.IsPlayer())
	}
}

// actionRestart is called when connection want to restart
func (room *RoomRecorder) actionRestart(msg synced.Msg) {
	if conn, ok := room.connectionCheck(msg); ok {
		room.Restart(conn)
	}
}

// actionDisconnect is called when connection disconnected
func (room *RoomRecorder) actionDisconnect(msg synced.Msg) {
	if conn, ok := room.connectionCheck(msg); ok {
		room.Disconnect(conn)
	}
}

// peoplCheck check that people publisher send correct message
func (room *RoomRecorder) connectionCheck(msg synced.Msg) (*Connection, bool) {
	action, ok := msg.Extra.(ConnectionMsg)
	if !ok {
		return nil, ok
	}
	return action.connection, ok
}

// connectionSubscribe subscibe to events associated with room's connection requests
func (room *RoomRecorder) connectionSubscribe(c RClientI) {
	observer := synced.NewObserver(
		synced.NewPair(action_.BackToLobby, room.actionBackToLobby),
		synced.NewPair(action_.Restart, room.actionRestart),
		synced.NewPair(action_.Disconnect, room.actionDisconnect))
	c.Observe(observer.AddPublisherCode(room_.UpdateChat))
}

// peoplCheck check that people publisher send correct message
func (room *RoomRecorder) peoplCheck(msg synced.Msg) (*Connection, bool, bool) {
	action, ok := msg.Extra.(ConnectionMsg)
	if !ok {
		return nil, ok, ok
	}
	needRecover, ok := action.content.(bool)
	if !ok {
		return nil, ok, ok
	}
	return action.connection, needRecover, ok
}

// peoplePlayerEnter is called when connection join room as player
func (room *RoomRecorder) peoplePlayerEnter(msg synced.Msg) {
	if conn, recover, ok := room.peoplCheck(msg); ok {
		room.AddConnection(conn, true, recover)
	}
}

// peopleObserverEnter is called when connection join room as observer
func (room *RoomRecorder) peopleObserverEnter(msg synced.Msg) {
	if conn, recover, ok := room.peoplCheck(msg); ok {
		room.AddConnection(conn, false, recover)
	}
}

// peopleSubscribe subscibe to events associated with room's members
func (room *RoomRecorder) peopleSubscribe(p PeopleI) {
	observer := synced.NewObserver(
		synced.NewPair(room_.ObserverEnter, room.peopleObserverEnter),
		synced.NewPair(room_.PlayerEnter, room.peoplePlayerEnter))
	p.Observe(observer.AddPublisherCode(room_.UpdatePeople))
}

// eventsSubscribe subscibe to events associated with room's status
func (room *RoomRecorder) eventsSubscribe(e EventsI) {
	observer := synced.NewObserver(
		synced.NewPairNoArgs(room_.StatusFinished, room.finish))
	e.Observe(observer.AddPublisherCode(room_.UpdateStatus))
}

// 265 -> 291
