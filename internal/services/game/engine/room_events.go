package engine

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
)

// RunnableI the class that starts (allocates resources) at the
//  beginning of a game and stops (clears resources) after it
//  finishes
type RunnableI interface {
	start()
	stop()
}

// SubscribeRunnable subscribe Runnable to listen for start
//  and end game events
func (room *RoomEvents) SubscribeRunnable(runnable RunnableI) {
	observer := synced.NewObserver(
		synced.NewPairNoArgs(room_.StatusRecruitment, runnable.start),
		synced.NewPairNoArgs(room_.StatusAborted, runnable.stop))
	room.Observe(observer.AddPublisherCode(room_.UpdateStatus))
}

// EventsI handle events in the game process
// Strategy Pattern
type EventsI interface {
	synced.PublisherI

	SubscribeRunnable(RunnableI)

	OpenCell(conn *Connection, cell *Cell)

	IsActive() bool
	Status() int32

	Restart(conn *Connection) error

	UpdateStatus(status int32)

	Run(realRoom *Room)
}

// RoomEvents implements EventsI
type RoomEvents struct {
	synced.PublisherBase

	s  synced.SyncI
	i  RoomInformationI
	l  LobbyProxyI
	p  PeopleI
	f  FieldProxyI
	g  GarbageCollectorI
	se RSendI

	play    *time.Timer
	prepare *time.Timer

	nextM *sync.RWMutex
	_next *Room

	statusM    *sync.RWMutex
	_status    int32
	chanStatus chan int32

	isDeathmatch  bool
	canClose      bool
	TimeToPlay    time.Duration
	TimeToPrepare time.Duration
}

// Init configure dependencies with other components of the room
func (room *RoomEvents) Init(builder RBuilderI, settings *models.RoomSettings, canClose bool) {
	room.init(settings, canClose)
	room.build(builder)
	room.subscribe()
}

// UpdateStatus update room status
func (room *RoomEvents) UpdateStatus(newStatus int32) {
	room.chanStatus <- newStatus
}

// IsActive check if game is started and results not known
func (room *RoomEvents) IsActive() bool {
	var result = false
	room.s.Do(func() {
		status := room.Status()
		result = status == room_.StatusFlagPlacing || status == room_.StatusRunning
	})
	return result
}

// Restart fill in the room fields with the original values
func (room *RoomEvents) Restart(conn *Connection) error {
	if room.Status() != room_.StatusFinished {
		return re.ErrorWrongStatus()
	}
	if room.next() == nil || room.next().sync.IsCleared() {
		next, err := room.l.CreateAndAddToRoom(conn)
		if err != nil {
			return err
		}
		room.setNext(next)
	}
	room.next().people.Enter(conn, true, false)
	return nil
}

func (room *RoomEvents) Run(realRoom *Room) {
	room.s.Do(func() {
		defer room.g.Close()
		var gc = room.g.SingleGoroutine()

		room.StartPublish()
		defer room.StopPublish()

		room.initTimers(true)
		defer func() {
			room.prepare.Stop()
			room.play.Stop()
		}()

		var (
			beginGame bool
			timeout   bool
		)

		room.notify(room_.StatusRecruitment, nil)
		for {
			select {
			case <-gc.C():
				go gc.Do()
			case <-room.prepare.C:
				if !beginGame {
					room.chanStatus <- room_.StatusRunning
				}
			case <-room.play.C:
				if beginGame {
					timeout = true
					room.chanStatus <- room_.StatusFinished
				}
			case newStatus := <-room.chanStatus:
				if newStatus == room_.StatusRunning {
					beginGame = true
				}
				room.updateStatus(newStatus, timeout)
				if newStatus == room_.StatusFinished || newStatus == room_.StatusAborted {
					return
				}
			}
		}
	})
}

// init struct's values
func (room *RoomEvents) init(settings *models.RoomSettings, canClose bool) {
	room.isDeathmatch = settings.Deathmatch
	room.TimeToPlay = time.Second * time.Duration(settings.TimeToPlay)
	room.TimeToPrepare = time.Second * time.Duration(settings.TimeToPrepare)

	room.statusM = &sync.RWMutex{}
	room.setStatus(room_.StatusRecruitment)

	room.canClose = canClose

	room.nextM = &sync.RWMutex{}
	room.setNext(nil)

	room.chanStatus = make(chan int32)

	room.PublisherBase = *synced.NewPublisher()
}

// build components
func (room *RoomEvents) build(builder RBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildInformation(&room.i)
	builder.BuildLobby(&room.l)
	builder.BuildPeople(&room.p)
	builder.BuildField(&room.f)
	builder.BuildGarbageCollector(&room.g)
	builder.BuildSender(&room.se)
}

// subscribe to room events
func (room *RoomEvents) subscribe() {
	room.peopleSubscribe(room.p)
}

// initTimers launch game timers. Call it when flag placement starts
func (room *RoomEvents) initTimers(first bool) {
	if first {
		room.prepare = time.NewTimer(time.Hour * 24)
		room.play = time.NewTimer(time.Hour * 24)
	} else {
		var prepare = time.Millisecond
		if room.isDeathmatch {
			prepare = room.TimeToPrepare
		}
		room.prepare.Reset(prepare)
		room.play.Reset(room.TimeToPlay + prepare)
	}
	return
}

// close handles players out of the room, free resources
func (room *RoomEvents) close() bool {
	var result bool
	room.s.Do(func() {
		if !room.canClose {
			return
		}
		utils.Debug(false, "Prepare to free!")
		go room.free()
		utils.Debug(false, "We did it")
		result = true
	})
	return result
}

// notify observers about new event
func (room *RoomEvents) notify(code int32, extra interface{}) {
	room.Notify(synced.Msg{
		Publisher: room_.UpdateStatus,
		Action:    code,
		Extra:     extra,
	})
}

func (room *RoomEvents) OpenCell(conn *Connection, cell *Cell) {
	room.s.DoWithOther(conn, func() {
		if !room.p.isAlive(conn) {
			return
		}
		status := room.Status()
		if status == room_.StatusFlagPlacing {
			room.f.SetFlag(conn, cell)
		} else if status == room_.StatusRunning {
			room.f.OpenCell(conn, cell)
		}
	})
}

// free clear all resources. Call it when no
//  observers and players inside
func (room *RoomEvents) free() {

	room.s.Clear(func() {
		room.chanStatus <- room_.StatusAborted

		room.play.Stop()
		room.prepare.Stop()

		close(room.chanStatus)
	})
}

// updateStatus sets a new room status if it is different from the current one.
//   Depending on the status applies the appropriate action and informs observers
//   of the status change
func (room *RoomEvents) updateStatus(newStatus int32, timeout bool) bool {
	oldStatus := room.Status()
	done := room.trySetStatus(newStatus)
	if !done {
		return false
	}
	if newStatus != room_.StatusFinished {
		room.notify(newStatus, nil)
	} else {
		room.notify(newStatus, room_.FinishResults{
			Cancel:  oldStatus == room_.StatusRecruitment,
			Timeout: timeout,
		})
	}

	room.se.StatusToAll(room.se.All, newStatus, nil)
	switch newStatus {
	case room_.StatusFlagPlacing:
		room.initTimers(false)
		room.se.Field(room.se.All)
	}
	return true
}

////////////////////////////////////////////////////////// mutex

// setStatus set the status of the room
func (room *RoomEvents) setStatus(status int32) {
	room.statusM.Lock()
	room._status = status
	room.statusM.Unlock()
}

// trySetStatus set status of the room, if it is different from
//   the current one. Return true if new status set, otherwise false
func (room *RoomEvents) trySetStatus(status int32) bool {
	room.statusM.Lock()
	defer room.statusM.Unlock()
	if status == room._status {
		return false
	}
	room._status = status
	return true
}

// Status return room's current status
func (room *RoomEvents) Status() int32 {
	room.statusM.RLock()
	v := room._status
	room.statusM.RUnlock()
	return v
}

// Next return next room to whick players from this room will be
// sent in case of pressing the restart button
func (room *RoomEvents) next() *Room {
	room.nextM.RLock()
	v := room._next
	room.nextM.RUnlock()
	return v
}

func (room *RoomEvents) setNext(next *Room) {
	room.nextM.Lock()
	room._next = next
	room.nextM.Unlock()
}

///////////////////////////////// subscribe

// peopleSubscribe subscibe to events associated with room's members
func (room *RoomEvents) peopleSubscribe(p PeopleI) {
	observer := synced.NewObserver(
		synced.Pair{
			Code: room_.AllDied,
			Do: func(synced.Msg) {
				room.chanStatus <- room_.StatusFinished
			}},
		synced.Pair{
			Code: room_.AllExit,
			Do: func(synced.Msg) {
				room.close()
			}})
	p.Observe(observer)
}

// 494
