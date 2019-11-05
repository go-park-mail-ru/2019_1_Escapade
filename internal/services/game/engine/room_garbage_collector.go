package engine

import (
	"time"
)

type RoomGarbageCollectorI interface {
	Init(s SyncI, e *RoomEvents, p *RoomPeople, c *RoomConnectionEvents,
		timeouts Timeouts)
	Run()
}

type Timeouts struct {
	timeoutPeopleFinding   float64
	timeoutRunningPlayer   float64
	timeoutRunningObserver float64
	timeoutFinished        float64
}

type RoomGarbageCollector struct {
	s SyncI
	e *RoomEvents
	p *RoomPeople
	c *RoomConnectionEvents

	tPlayer   float64
	tObserver float64
	t         Timeouts
}

func (room *RoomGarbageCollector) Init(s SyncI, e *RoomEvents,
	p *RoomPeople, c *RoomConnectionEvents, timeouts Timeouts) {
	room.s = s
	room.e = e
	room.p = p
	room.c = c
	room.t = timeouts
}

func (room *RoomGarbageCollector) updateTimeouts() {
	status := room.e.Status()
	if status == StatusRecruitment {
		room.tPlayer = room.t.timeoutPeopleFinding
		room.tObserver = room.t.timeoutPeopleFinding
	} else if status == StatusFinished {
		room.tPlayer = room.t.timeoutFinished
		room.tObserver = room.t.timeoutFinished
	} else {
		room.tPlayer = room.t.timeoutRunningPlayer
		room.tObserver = room.t.timeoutRunningObserver
	}
}

func (room *RoomGarbageCollector) Run() {
	room.s.do(func() {
		room.updateTimeouts()
		room.checkPeople()
	})
}

func (room *RoomGarbageCollector) checkPeople() {
	room.p.ForEach(func(c *Connection, isPlayer bool) {
		if isPlayer {
			if room.isExpired(c, room.tPlayer) {
				room.c.Timeout(c)
			}
		} else {
			if room.isExpired(c, room.tObserver) {
				room.c.Timeout(c)
			}
		}
	})
}

func (room *RoomGarbageCollector) isExpired(conn *Connection, timeout float64) bool {
	t := conn.Time()
	return conn.Disconnected() && time.Since(t).Seconds() > timeout
}
