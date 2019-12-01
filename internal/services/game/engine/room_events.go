package engine

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// EventsI handle events in the game process
// Strategy Pattern
type EventsI interface {
	OpenCell(conn *Connection, cell *Cell)

	IsActive() bool
	Status() int

	Restart(conn *Connection) error
	Next() *Room

	prepareOver()
	RecruitingOver()

	tryFinish()
	tryClose()

	configure(status int, date time.Time)

	Free()
	Run()

	playingTime() time.Duration
	recruitmentTime() time.Duration

	Date() time.Time
}

// RoomEvents implements EventsI
type RoomEvents struct {
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

	room.chanStatus = make(chan int)

	room.setStatus(StatusRecruitment)
}

// initTimers launch game timers. Call it when flag placement starts
func (room *RoomEvents) initTimers(first bool) {
	if first {
		utils.Debug(false, "first time ")
		room.prepare = time.NewTimer(time.Hour * 24)
		room.play = time.NewTimer(time.Hour * 24)
	} else {
		var prepare = time.Millisecond
		if room.isDeathmatch {
			prepare = room.TimeToPrepare
		}
		utils.Debug(false, "!!!times:", prepare, room.TimeToPlay)
		room.prepare.Reset(prepare)
		room.play.Reset(room.TimeToPlay + prepare)
	}
	return
}

func (room *RoomEvents) RecruitingOver() {
	room.initTimers(false)
	if room.updateStatus(StatusFlagPlacing) {
		if room.isDeathmatch {
			room.se.StatusToAll(room.se.All, StatusFlagPlacing, nil)
		}
	}
}

func (room *RoomEvents) prepareOver() {
	room.prepare.Stop()
	if room.updateStatus(StatusRunning) {
		room.se.StatusToAll(room.se.All, StatusRunning, nil)
	}
}

func (room *RoomEvents) tryFinish() {
	if room.p.AllKilled() || room.f.Field().IsCleared() {
		room.playingOver()
	}
}

func (room *RoomEvents) tryClose() {
	if room.p.Empty() {
		room.Close()
	}
}

func (room *RoomEvents) playingOver() {
	room.play.Stop()
	if room.updateStatus(StatusFinished) {
		room.se.StatusToAll(room.se.All, StatusFinished, nil)
	}
}

func (room *RoomEvents) updateStatus(newStatus int) bool {
	if room.Status() != newStatus {
		go func() { room.chanStatus <- newStatus }()
		return true
	}
	return false
}

// StartFlagPlacing prepare field, players and observers
func (room *RoomEvents) StartFlagPlacing() {
	room.s.Do(func() {
		utils.Debug(false, "StartFlagPlacing")
		room.setStatus(StatusFlagPlacing)

		room.p.ForEach(func(c *Connection, isPlayer bool) {
			room.se.Room(c)
			room.l.WaiterToPlayer(c)
		})

		room.p.Start()
		room.l.Start()

		room.se.StatusToAll(room.se.All, StatusFlagPlacing, nil)
		room.se.Field(room.se.All)
	})
}

func (room *RoomEvents) CancelGame() {
	room.s.Do(func() {
		room.setStatus(StatusFinished)
		room.met.Observe(room.l.metricsEnabled(), true)
		room.l.Finish()
	})
}

// StartGame start game
func (room *RoomEvents) StartGame() {
	room.s.Do(func() {
		utils.Debug(false, "StartGame")
		//s := room.r.Settings.Width * room.r.Settings.Height
		//open := float64(room.r.Settings.Mines) / float64(s) * float64(100)
		//utils.Debug(false, "opennn", open, room.Settings.Width*room.Settings.Height)

		cells := room.f.Field().OpenZero() //room.Field.OpenSave(int(open))
		room.se.NewCells(cells...)
		room.setStatus(StatusRunning)
		room.setDate(room.l.Date())
		room.se.StatusToAll(room.se.All, StatusRunning, nil)
		room.se.Text("Battle began! Destroy your enemy!", room.se.All)
	})
}

// FinishGame finish game
func (room *RoomEvents) FinishGame(timer bool) {
	room.s.Do(func() {
		room.setStatus(StatusFinished)

		// save Group
		saveAndSendGroup := &sync.WaitGroup{}

		cells := make([]Cell, 0)
		room.f.Field().OpenEverything(cells)

		saveAndSendGroup.Add(1)
		room.se.GameOver(timer, room.se.All, cells, saveAndSendGroup)

		saveAndSendGroup.Add(1)
		room.mo.Save(saveAndSendGroup)

		saveAndSendGroup.Add(1)
		room.p.Finish(saveAndSendGroup)
		saveAndSendGroup.Wait()

		room.met.Observe(room.l.metricsEnabled(), false)

		room.l.Finish()
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
		utils.Debug(false, "We closed room :С")

		go room.l.Close()

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
		utils.Debug(false, "status", status)
		if status == StatusFlagPlacing {
			room.f.SetFlag(conn, cell)
		} else if status == StatusRunning {
			room.f.OpenCell(conn, cell)
		}
	})
}

func (room *RoomEvents) Run() {
	room.s.Do(func() {

		defer room.g.Close()
		// зассунуть в room.g оттуда доставать
		var gc = room.g.SingleGoroutine()

		room.initTimers(true)
		defer func() {
			room.prepare.Stop()
			room.play.Stop()
		}()

		var beginGame, timeOut bool

		for {
			select {
			case <-gc.C():
				go gc.Do()
			case <-room.prepare.C:
				if beginGame {
					room.prepareOver()
				}
			case <-room.play.C:
				if beginGame {
					timeOut = true
					room.playingOver()
				}
			case newStatus := <-room.chanStatus:
				oldStatus := room.Status()

				if newStatus == oldStatus || newStatus > StatusFinished {
					continue
				}
				if oldStatus == StatusRecruitment {
					room.setRecruitmentTime()
				} else if newStatus == StatusFinished {
					room.setPlayingTime()
				}
				switch newStatus {
				case StatusFlagPlacing:
					beginGame = true
					room.StartFlagPlacing()
				case StatusRunning:
					room.StartGame()
				case StatusFinished:
					if oldStatus == StatusRecruitment {
						room.CancelGame()
					} else {
						room.FinishGame(timeOut)
					}
				//return
				case StatusAborted:
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
		result = status == StatusFlagPlacing || status == StatusRunning
	})
	return result
}

// Free clear all resources. Call it when no
//  observers and players inside
func (room *RoomEvents) Free() {

	room.s.Clear(func() {
		room.chanStatus <- StatusAborted

		room.setStatus(StatusFinished)
		go room.re.Free()
		go room.mes.Free()
		go room.p.Free()
		go room.f.Field().Free()

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
