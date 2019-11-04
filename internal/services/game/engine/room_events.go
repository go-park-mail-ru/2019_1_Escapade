package engine

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

type RoomEvents struct {
	r       *Room
	s       SyncI
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
}

func (room *RoomEvents) Init(r *Room, s SyncI) {
	room.r = r
	room.s = s
	room.statusM = &sync.RWMutex{}

	room.recruitmentTimeM = &sync.RWMutex{}

	room.playingTimeM = &sync.RWMutex{}

	room.dateM = &sync.RWMutex{}
	room.setDate(time.Now().In(room.r.lobby.location()))

	room.nextM = &sync.RWMutex{}
	room._next = nil

	room.chanStatus = make(chan int)

	room.setStatus(StatusRecruitment)
}

// initTimers launch game timers. Call it when flag placement starts
func (room *RoomEvents) initTimers(first bool) {
	if first {
		room.prepare = time.NewTimer(time.Millisecond)
		room.play = time.NewTimer(time.Millisecond)
	} else {
		if room.r.Settings.Deathmatch {
			room.prepare.Reset(time.Second *
				time.Duration(room.r.Settings.TimeToPrepare))
		} else {
			room.prepare.Reset(time.Millisecond)
		}
		room.play.Reset(time.Second *
			time.Duration(room.r.Settings.TimeToPlay))
	}
	return
}

func (room *RoomEvents) RecruitingOver() {
	room.initTimers(false)
	if room.updateStatus(StatusFlagPlacing) {
		if room.r.Settings.Deathmatch {
			go room.r.send.StatusToAll(room.r.All, StatusFlagPlacing, nil)
		}
	}
}

func (room *RoomEvents) prepareOver() {
	room.prepare.Stop()
	if room.updateStatus(StatusRunning) {
		go room.r.send.StatusToAll(room.r.All, StatusRunning, nil)
	}
}

func (room *RoomEvents) tryFinish() {
	if room.r.people.AllKilled() {
		room.playingOver()
	}
}

func (room *RoomEvents) tryClose() {
	if room.r.people.Empty() {
		room.Close()
	}
}

func (room *RoomEvents) playingOver() {
	room.play.Stop()
	if room.updateStatus(StatusFinished) {
		go room.r.send.StatusToAll(room.r.All, StatusFinished, nil)
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
	room.s.do(func() {
		room.setStatus(StatusFlagPlacing)

		room.r.people.ForEach(func(c *Connection, isPlayer bool) {
			go room.r.send.Room(c)
			room.r.lobby.waiterToPlayer(c, room.r)
		})

		room.r.field.Fill(room.r.people.Players.m.Flags())
		room.r.people.Players.Init()

		room.r.lobby.RoomStart(room.r)

		go room.r.send.StatusToAll(room.r.All, StatusFlagPlacing, nil)
		go room.r.send.Field(room.r.All)
	})
}

func (room *RoomEvents) CancelGame() {
	room.s.do(func() {
		room.setStatus(StatusFinished)
		go room.r.metrics.Observe(room.r.lobby.config().Metrics, true)
		room.r.lobby.roomFinish(room.r)
	})
}

// StartGame start game
func (room *RoomEvents) StartGame() {
	room.s.do(func() {
		//s := room.r.Settings.Width * room.r.Settings.Height
		//open := float64(room.r.Settings.Mines) / float64(s) * float64(100)
		//utils.Debug(false, "opennn", open, room.Settings.Width*room.Settings.Height)

		cells := room.r.field.Field.OpenZero() //room.Field.OpenSave(int(open))
		go room.r.send.NewCells(room.r.All, cells...)
		room.setStatus(StatusRunning)
		room.setDate(time.Now().In(room.r.lobby.location()))
		go room.r.send.StatusToAll(room.r.All, StatusRunning, nil)
		go room.r.send.Message("Battle began! Destroy your enemy!", room.r.All)
	})
}

// FinishGame finish game
func (room *RoomEvents) FinishGame(timer bool) {
	room.s.do(func() {
		room.setStatus(StatusFinished)

		// save Group
		saveAndSendGroup := &sync.WaitGroup{}

		cells := make([]Cell, 0)
		room.r.field.Field.OpenEverything(cells)

		saveAndSendGroup.Add(1)
		go room.r.send.GameOver(timer, room.r.All, cells, saveAndSendGroup)

		saveAndSendGroup.Add(1)
		go room.r.models.Save(saveAndSendGroup)

		saveAndSendGroup.Add(1)
		go room.r.people.Players.m.Finish(saveAndSendGroup)
		saveAndSendGroup.Wait()

		go room.r.metrics.Observe(room.r.lobby.config().Metrics, false)

		room.r.lobby.roomFinish(room.r)
	})
}

// Close drives away players out of the room, free resources
// and inform lobby, that rooms closes
func (room *RoomEvents) Close() bool {
	var result bool
	room.s.do(func() {
		if !room.r.lobby.config().CanClose {
			return
		}
		utils.Debug(false, "We closed room :С")

		go room.r.lobby.CloseRoom(room.r)

		utils.Debug(false, "Prepare to free!")
		go room.s.Free()
		utils.Debug(false, "We did it")
		result = true
	})
	return result
}

func (room *RoomEvents) Run() {
	room.s.do(func() {
		// все в конфиг
		ticker := time.NewTicker(time.Second * 10)

		room.initTimers(true)
		defer func() {
			ticker.Stop()
			room.prepare.Stop()
			room.play.Stop()
		}()

		var beginGame, timeOut bool

		for {
			select {
			case <-ticker.C:
				go room.r.garbageCollector.Run()
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
					ticker.Stop()
					return
				}
			}
		}
	})
}

// Restart fill in the room fields with the original values
func (room *RoomEvents) Restart(conn *Connection) error {

	if room.Next() == nil || room.Next().sync.done() {
		next, err := room.r.lobby.CreateAndAddToRoom(room.r.Settings, conn)
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
	room.s.do(func() {
		status := room.Status()
		result = status == StatusFlagPlacing || status == StatusRunning
	})
	return result
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

	t := time.Now().In(room.r.lobby.location())

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

	t := time.Now().In(room.r.lobby.location())

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
