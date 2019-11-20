package engine

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
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
	history() []*PlayerAction

	AddConnection(conn *Connection, isPlayer bool, needRecover bool)
	setHistory(history []*PlayerAction)

	Kill(conn *Connection, action int32, isDeathmatch bool)

	configure(info []models.Action)

	Free()
}

// RoomRecorder notify actions to room history and to users
// implements ActionRecorderProxyI
type RoomRecorder struct {
	s  synced.SyncI
	i  RoomInformationI
	l  LobbyProxyI
	p  PeopleI
	f  FieldProxyI
	se RSendI
	mo RModelsI

	historyM *sync.RWMutex
	_history []*PlayerAction
}

func (room *RoomRecorder) Init(builder RBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildInformation(&room.i)
	builder.BuildLobby(&room.l)
	builder.BuildPeople(&room.p)
	builder.BuildField(&room.f)
	builder.BuildSender(&room.se)
	builder.BuildModelsAdapter(&room.mo)

	room.historyM = &sync.RWMutex{}
	room.setHistory(make([]*PlayerAction, 0))
}

func (room *RoomRecorder) Free() {
	room.historyFree()
}

// LeaveMeta update metainformation about user leaving room
func (room *RoomRecorder) Leave(conn *Connection, action int32, isPlayer bool) {
	if isPlayer {
		room.se.PlayerExit(conn)
	} else {
		room.se.ObserverExit(conn)
	}
	room.notifyAll(conn, action)
}

func (room *RoomRecorder) Disconnect(conn *Connection) {
	room.notifyAll(conn, ActionDisconnect)
}

func (room *RoomRecorder) FlagSet(conn *Connection) {
	room.notifyAll(conn, ActionFlagSet)
}

func (room *RoomRecorder) Restart(conn *Connection) {
	room.notifyAll(conn, ActionRestart)
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
		for _, actionHistory := range history {
			action := room.mo.toModelPlayerAction(actionHistory)
			actions = append(actions, action)
		}
	})
	return actions
}

func (room *RoomRecorder) Reconnect(conn *Connection) {
	room.notifyAll(conn, ActionReconnect)
}

func (room *RoomRecorder) AddPlayer(conn *Connection) {
	room.notifyAll(conn, ActionConnectAsPlayer)
	room.se.PlayerEnter(conn)
}

func (room *RoomRecorder) AddObserver(conn *Connection) {
	room.notifyAll(conn, ActionConnectAsObserver)
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
	room.notifyAll(conn, ActionFlagСonflict)
}

func (room *RoomRecorder) notifyAll(conn *Connection, action int32) {
	room.s.Do(func() {
		pa := NewPlayerAction(conn.ID(), action)
		room.appendAction(pa)
		if !room.p.Empty() {
			room.se.Action(*pa, room.se.AllExceptThat(conn))
			room.l.Notify()
		}
	})
}

func (room *RoomRecorder) configure(info []models.Action) {
	room.setHistory(make([]*PlayerAction, 0))
	for _, actionDB := range info {
		action := room.mo.fromModelPlayerAction(actionDB)
		room.appendAction(action)
	}
}

/////////////////////////////// mutex

// history return '_history' field
func (room *RoomRecorder) history() []*PlayerAction {
	room.historyM.RLock()
	v := room._history
	room.historyM.RUnlock()
	return v
}

func (room *RoomRecorder) setHistory(history []*PlayerAction) {
	room.historyM.Lock()
	room._history = history
	room.historyM.Unlock()
}

// appendAction append action to action slice(history)
func (room *RoomRecorder) appendAction(action *PlayerAction) {
	room.historyM.Lock()
	defer room.historyM.Unlock()
	room._history = append(room._history, action)
}

// historyFree free action slice
func (room *RoomRecorder) historyFree() {
	room.historyM.Lock()
	room._history = nil
	room.historyM.Unlock()
}
