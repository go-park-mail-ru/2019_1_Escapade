package engine

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
)

// GarbageCollectorI handle deleting connections, when they are disconnected
// Strategy Pattern
type GarbageCollectorI interface {
	Run()
	Close()
	SingleGoroutine() synced.SingleGoroutine
}

// RoomGarbageCollector implements GarbageCollectorI
type RoomGarbageCollector struct {
	s synced.SyncI
	e EventsI
	p PeopleI
	c RClientI

	tPlayer   time.Duration
	tObserver time.Duration
	t         config.GameTimeouts
	sg        synced.SingleGoroutine
}

// Init configure dependencies with other components of the room
func (room *RoomGarbageCollector) Init(builder RBuilderI,
	interval time.Duration, timeouts config.GameTimeouts) {

	builder.BuildSync(&room.s)
	builder.BuildEvents(&room.e)
	builder.BuildPeople(&room.p)
	builder.BuildConnectionEvents(&room.c)

	room.t = timeouts

	room.sg = synced.SingleGoroutine{}
	room.sg.Init(interval, room.Run)
}

func (room *RoomGarbageCollector) updateTimeouts() {
	status := room.e.Status()
	if status == room_.StatusRecruitment {
		room.tPlayer = room.t.PeopleFinding.Duration
		room.tObserver = room.t.PeopleFinding.Duration
	} else if status == room_.StatusFinished {
		room.tPlayer = room.t.Finished.Duration
		room.tObserver = room.t.Finished.Duration
	} else {
		room.tPlayer = room.t.RunningPlayer.Duration
		room.tObserver = room.t.RunningObserver.Duration
	}
}

func (room *RoomGarbageCollector) SingleGoroutine() synced.SingleGoroutine {
	return room.sg
}

func (room *RoomGarbageCollector) Run() {
	room.s.Do(func() {
		room.updateTimeouts()
		room.checkPeople()
	})
}

func (room *RoomGarbageCollector) Close() {
	room.sg.Close()
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

func (room *RoomGarbageCollector) isExpired(conn *Connection, timeout time.Duration) bool {
	t := conn.Time()
	return conn.Disconnected() && time.Since(t) > timeout
}
