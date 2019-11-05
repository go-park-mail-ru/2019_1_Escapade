package engine

import (
	"time"
)

// GarbageCollectorI handle deleting connections, when they are disconnected
// Strategy Pattern
type GarbageCollectorI interface {
	Run()
}

// Timeouts contains the timeouts required for the garbage collector to run
type Timeouts struct {
	timeoutPeopleFinding   float64
	timeoutRunningPlayer   float64
	timeoutRunningObserver float64
	timeoutFinished        float64
}

// RoomGarbageCollector implements GarbageCollectorI
type RoomGarbageCollector struct {
	s SyncI
	e EventsI
	p PeopleI
	c ConnectionEventsI

	tPlayer   float64
	tObserver float64
	t         Timeouts
}

// Init configure dependencies with other components of the room
func (room *RoomGarbageCollector) Init(builder ComponentBuilderI, timeouts Timeouts) {
	builder.BuildSync(&room.s)
	builder.BuildEvents(&room.e)
	builder.BuildPeople(&room.p)
	builder.BuildRoomConnectionEvents(&room.c)

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
