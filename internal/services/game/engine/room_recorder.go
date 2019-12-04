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

	Free()

	ConnectionSub() synced.SubscriberI
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

func (room *RoomRecorder) Init(builder RBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildInformation(&room.i)
	builder.BuildPeople(&room.p)
	builder.BuildField(&room.f)
	builder.BuildSender(&room.se)
	builder.BuildModelsAdapter(&room.mo)

	room.historyM = &sync.RWMutex{}
	room.setHistory(make([]*action_.PlayerAction, 0))
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

// history return '_history' field
func (room *RoomRecorder) history() []*action_.PlayerAction {
	room.historyM.RLock()
	v := room._history
	room.historyM.RUnlock()
	return v
}

func (room *RoomRecorder) setHistory(history []*action_.PlayerAction) {
	room.historyM.Lock()
	room._history = history
	room.historyM.Unlock()
}

// appendAction append action to action slice(history)
func (room *RoomRecorder) appendAction(action *action_.PlayerAction) {
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

///////////////////////// callbacks

func (room *RoomRecorder) ConnectionSub() synced.SubscriberI {
	return synced.NewSubscriber(room.connectionCallback)
}

func (room *RoomRecorder) connectionCallback(msg synced.Msg) {
	if msg.Code != room_.UpdateConnection {
		return
	}
	action, ok := msg.Content.(ConnectionMsg)
	if !ok {
		return
	}
	switch action.code {
	case action_.BackToLobby:
		isPlayer, ok := action.content.(bool)
		if !ok {
			return
		}
		room.Leave(action.connection, action.code, isPlayer)
	case action_.Restart:
		room.Restart(action.connection)
	}
}

func (room *RoomRecorder) PeopleSub() synced.SubscriberI {
	return synced.NewSubscriber(room.peopleCallback)
}

func (room *RoomRecorder) peopleCallback(msg synced.Msg) {
	if msg.Code != room_.UpdatePeople {
		return
	}
	action, ok := msg.Content.(ConnectionMsg)
	if !ok {
		return
	}
	switch action.code {
	case room_.ObserverEnter:
		needRecover, ok := action.content.(bool)
		if !ok {
			return
		}
		room.AddConnection(action.connection, false, needRecover)
	case room_.PlayerEnter:
		needRecover, ok := action.content.(bool)
		if !ok {
			return
		}
		room.AddConnection(action.connection, false, needRecover)
	}
}
