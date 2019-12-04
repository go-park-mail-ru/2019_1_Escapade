package engine

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
)

// EventsI handle events in the game process
// Strategy Pattern
type EventsI interface {
	OpenCell(conn *Connection, cell *Cell)

	IsActive() bool
	Status() int

	Restart(conn *Connection) error
	Next() *Room

	RecruitingOver()
	PrepareOver()

	tryFinish()

	Timeout() bool

	configure(status int, date time.Time)

	Free()
	Run(realRoom *Room)

	playingTime() time.Duration
	recruitmentTime() time.Duration

	Date() time.Time

	PeopleSub() synced.SubscriberI
}

// RoomEvents implements EventsI
type RoomEvents struct {
	synced.PublisherBase

	s   synced.SyncI
	i   RoomInformationI
	l   LobbyProxyI
	p   PeopleI
	f   FieldProxyI
	g   GarbageCollectorI
	mo  RModelsI
	met MetricsI
	mes MessagesI
	re  ActionRecorderI
	se  RSendI

	play    *time.Timer
	prepare *time.Timer

	nextM *sync.RWMutex
	_next *Room

	dateM *sync.RWMutex
	_date time.Time

	recruitmentTimeM *sync.RWMutex
	_recruitmentTime time.Duration

	playingTimeM *sync.RWMutex
	_playingTime time.Duration

	statusM    *sync.RWMutex
	_status    int
	chanStatus chan int

	timeoutM *sync.RWMutex
	_timeout bool

	isDeathmatch  bool
	TimeToPlay    time.Duration
	TimeToPrepare time.Duration
}

// Init configure dependencies with other components of the room
func (room *RoomEvents) Init(builder RBuilderI, settings *models.RoomSettings) {
	builder.BuildSync(&room.s)
	builder.BuildInformation(&room.i)
	builder.BuildLobby(&room.l)
	builder.BuildPeople(&room.p)
	builder.BuildField(&room.f)
	builder.BuildGarbageCollector(&room.g)
	builder.BuildModelsAdapter(&room.mo)
	builder.BuildMetrics(&room.met)
	builder.BuildMessages(&room.mes)
	builder.BuildRecorder(&room.re)
	builder.BuildSender(&room.se)

	room.PublisherBase = *synced.NewPublisher()
	room.PublisherBase.Start(room.f.EventsSub(), room.l.EventsSub())

	room.isDeathmatch = settings.Deathmatch
	room.TimeToPlay = time.Second * time.Duration(settings.TimeToPlay)
	room.TimeToPrepare = time.Second * time.Duration(settings.TimeToPrepare)

	room.statusM = &sync.RWMutex{}

	room.recruitmentTimeM = &sync.RWMutex{}

	room.playingTimeM = &sync.RWMutex{}

	room.dateM = &sync.RWMutex{}
	room.setDate(room.l.Date())

	room.nextM = &sync.RWMutex{}
	room._next = nil

	room.timeoutM = &sync.RWMutex{}
	room._timeout = false

	room.chanStatus = make(chan int)

	room.setStatus(room_.StatusRecruitment)
}

func (room *RoomEvents) Timeout() bool {
	room.timeoutM.RLock()
	defer room.timeoutM.RUnlock()
	return room._timeout
}

func (room *RoomEvents) setTimeout() {
	room.timeoutM.RLock()
	defer room.timeoutM.RUnlock()
	room._timeout = true
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

func (room *RoomEvents) tryFinish() {
	if room.f.Field().IsCleared() {
		room.chanStatus <- room_.StatusFinished
	}
}

func (room *RoomEvents) RecruitingOver() {
	room.chanStatus <- room_.StatusFlagPlacing
}

func (room *RoomEvents) PrepareOver() {
	room.chanStatus <- room_.StatusRunning
}

func (room *RoomEvents) CancelGame() {
	room.s.Do(func() {
		room.met.Observe(room.l.metricsEnabled(), true)
	})
}

// StartGame start game
func (room *RoomEvents) StartGame() {
	room.s.Do(func() {
		room.setDate(room.l.Date())
		room.se.Text("Battle began! Destroy your enemy!", room.se.All)
	})
}

// FinishGame finish game
func (room *RoomEvents) FinishGame() {
	room.s.Do(func() {
		room.mo.Save()
		room.met.Observe(room.l.metricsEnabled(), false)
	})
}

// Close drives away players out of the room, free resources
// and inform lobby, that rooms closes
func (room *RoomEvents) Close() bool {
	var result bool
	room.s.Do(func() {
		if !room.l.closeEnabled() {
			return
		}
		utils.Debug(false, "Prepare to free!")
		go room.Free()
		utils.Debug(false, "We did it")
		result = true
	})
	return result
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

func (room *RoomEvents) Run(realRoom *Room) {
	room.s.Do(func() {
		defer room.g.Close()
		var gc = room.g.SingleGoroutine()

		//go room.Start(..)
		defer room.Stop()

		room.initTimers(true)
		defer func() {
			room.prepare.Stop()
			room.play.Stop()
		}()

		var beginGame bool

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
					room.setTimeout()
					room.chanStatus <- room_.StatusFinished
				}
			case newStatus := <-room.chanStatus:
				if newStatus == room_.StatusRunning {
					beginGame = true
				}
				room.updateStatus(newStatus)
				if newStatus == room_.StatusFinished || newStatus == room_.StatusAborted {
					return
				}
			}
		}
	})
}

// Restart fill in the room fields with the original values
func (room *RoomEvents) Restart(conn *Connection) error {

	if room.Next() == nil || room.Next().sync.IsCleared() {
		next, err := room.l.CreateAndAddToRoom(conn)
		if err != nil {
			return err
		}
		room.setNext(next)
	}
	return nil
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

// Free clear all resources. Call it when no
//  observers and players inside
func (room *RoomEvents) Free() {

	room.s.Clear(func() {
		room.chanStatus <- room_.StatusAborted

		go room.re.Free()
		go room.mes.Free()
		room.play.Stop()
		room.prepare.Stop()

		close(room.chanStatus)
	})
}

func (room *RoomEvents) configure(status int, date time.Time) {
	room.setStatus(status)
	room.setDate(date)
}

////////////////////////////////////////////////////////// mutex

func (room *RoomEvents) setStatus(status int) {
	room.statusM.Lock()
	room._status = status
	room.statusM.Unlock()
}

func (room *RoomEvents) trySetStatus(status int) bool {
	room.statusM.Lock()
	defer room.statusM.Unlock()
	if status == room._status {
		return false
	}
	room._status = status
	return true
}

// Status return room's current status
func (room *RoomEvents) Status() int {
	room.statusM.RLock()
	v := room._status
	room.statusM.RUnlock()
	return v
}

func (room *RoomEvents) playingTime() time.Duration {
	room.playingTimeM.RLock()
	v := room._playingTime
	room.playingTimeM.RUnlock()
	return v
}

func (room *RoomEvents) setPlayingTime() {
	room.dateM.RLock()
	v := room._date
	room.dateM.RUnlock()

	t := room.l.Date()

	room.playingTimeM.Lock()
	room._playingTime = t.Sub(v)
	room.playingTimeM.Unlock()

	room.dateM.Lock()
	room._date = t
	room.dateM.Unlock()
}

// Date return date, when room was created
func (room *RoomEvents) Date() time.Time {
	room.dateM.RLock()
	v := room._date
	room.dateM.RUnlock()
	return v
}

func (room *RoomEvents) setDate(date time.Time) {
	room.dateM.Lock()
	room._date = date
	room.dateM.Unlock()
}

func (room *RoomEvents) recruitmentTime() time.Duration {
	room.recruitmentTimeM.RLock()
	v := room._recruitmentTime
	room.recruitmentTimeM.RUnlock()
	return v
}

func (room *RoomEvents) setRecruitmentTime() {
	room.dateM.RLock()
	v := room._date
	room.dateM.RUnlock()

	t := room.l.Date()

	room.recruitmentTimeM.Lock()
	room._recruitmentTime = t.Sub(v)
	room.recruitmentTimeM.Unlock()

	room.dateM.Lock()
	room._date = t
	room.dateM.Unlock()
}

// Next return next room to whick players from this room will be
// sent in case of pressing the restart button
func (room *RoomEvents) Next() *Room {
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

func (room *RoomEvents) updateStatus(newStatus int) bool {
	oldStatus := room.Status()
	done := room.trySetStatus(newStatus)
	if !done {
		return false
	}
	room.Notify(synced.Msg{
		Code:    room_.UpdateStatus,
		Content: newStatus,
	})

	room.se.StatusToAll(room.se.All, newStatus, nil)
	switch newStatus {
	case room_.StatusFlagPlacing:
		room.initTimers(false)
		room.se.Field(room.se.All)
	case room_.StatusRunning:
		room.setRecruitmentTime()
		room.StartGame()
	case room_.StatusFinished:
		room.setPlayingTime()
		if oldStatus == room_.StatusRecruitment {
			room.CancelGame()
		} else {
			room.FinishGame()
		}
	}
	return true
}

func (room *RoomEvents) PeopleSub() synced.SubscriberI {
	return synced.NewSubscriber(room.peopleCallback)
}

///////////////////////////////// callbacks

func (room *RoomEvents) peopleCallback(msg synced.Msg) {
	if msg.Code != room_.UpdatePeople {
		return
	}
	code, ok := msg.Content.(int)
	if !ok {
		return
	}
	switch code {
	case room_.AllDied:
		room.chanStatus <- room_.StatusFinished
	case room_.AllExit:
		room.Close()
	}
}

// 494
