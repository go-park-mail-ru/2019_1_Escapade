package engine

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// RoomNotifier notify actions to room history and to users
type RoomRecorder struct {
	//r  *Room
	s  SyncI
	i  *RoomInformation
	l  RoomLobbyCommunicationI
	p  *RoomPeople
	f  *RoomField
	se *RoomSender

	historyM *sync.RWMutex
	_history []*PlayerAction
}

func (room *RoomRecorder) Init(s SyncI, i *RoomInformation,
	l RoomLobbyCommunicationI, p *RoomPeople, f *RoomField, se *RoomSender) {
	room.s = s
	room.i = i
	room.l = l
	room.p = p
	room.f = f
	room.se = se
	room.historyM = &sync.RWMutex{}
	room.setHistory(make([]*PlayerAction, 0))
}

func (room *RoomRecorder) Free() {
	go room.historyFree()
}

// LeaveMeta update metainformation about user leaving room
func (room *RoomRecorder) Leave(conn *Connection, action int32) {
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
	cell := room.p.Players.m.Flag(conn.Index())
	cells := make([]Cell, 0)
	room.f.Field.saveCell(&cell.Cell, cells)
	go room.se.NewCells(cell.Cell)
}

func (room *RoomRecorder) Kill(conn *Connection,
	action int32, isDeathmatch bool) {
	room.notifyAll(conn, action)
	room.s.do(func() {
		if isDeathmatch {
			room.flag(conn)
		}
	})
}

func (room *RoomRecorder) ModelActions() []models.Action {
	history := room.history()
	actions := make([]models.Action, 0)
	room.s.do(func() {
		for _, actionHistory := range history {
			action := models.Action{
				PlayerID: actionHistory.Player,
				ActionID: actionHistory.Action,
				Date:     actionHistory.Time,
			}
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
	go room.se.PlayerEnter(conn)
}

func (room *RoomRecorder) AddObserver(conn *Connection) {
	room.notifyAll(conn, ActionConnectAsObserver)
	go room.se.ObserverEnter(conn)
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
	room.s.do(func() {
		pa := NewPlayerAction(conn.ID(), action)
		room.appendAction(pa)
		if !room.p.Empty() {
			room.se.Action(*pa, room.se.AllExceptThat(conn))
			go room.l.Notify()
		}
	})
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
