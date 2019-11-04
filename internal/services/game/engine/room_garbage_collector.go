package engine

import (
	"time"
)

type RoomGarbageCollectorI interface {
	Init(r *Room, s SyncI, timeouts Timeouts)
	Run()
}

type Timeouts struct {
	timeoutPeopleFinding   float64
	timeoutRunningPlayer   float64
	timeoutRunningObserver float64
	timeoutFinished        float64
}

type RoomGarbageCollector struct {
	r         *Room
	s         SyncI
	tPlayer   float64
	tObserver float64
	t         Timeouts
}

func (room *RoomGarbageCollector) Init(r *Room, s SyncI, timeouts Timeouts) {
	room.r = r
	room.s = s
	room.t = timeouts
}

func (room *RoomGarbageCollector) updateTimeouts() {
	status := room.r.events.Status()
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
	room.r.people.ForEach(func(c *Connection, isPlayer bool) {
		if isPlayer {
			if room.isExpired(c, room.tPlayer) {
				room.r.connEvents.Timeout(c, isPlayer)
			}
		} else {
			if room.isExpired(c, room.tObserver) {
				room.r.connEvents.Timeout(c, isPlayer)
			}
		}
	})
}

func (room *RoomGarbageCollector) isExpired(conn *Connection, timeout float64) bool {
	t := conn.Time()
	return conn.Disconnected() && time.Since(t).Seconds() > timeout
}
