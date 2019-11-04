package engine

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

type SyncI interface {
	do(f func())
	done() bool
	doWithConn(conn *Connection, f func())
	Free()
}

type RoomSync struct {
	r      *Room
	wGroup *sync.WaitGroup

	doneM *sync.RWMutex
	_done bool
}

func (room *RoomSync) Init(r *Room) {
	room.r = r
	room.wGroup = &sync.WaitGroup{}
	room._done = false
	room.doneM = &sync.RWMutex{}
}

// do any action in room by calling this func!
func (room *RoomSync) do(f func()) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("catch panic")
		room.wGroup.Done()
	}()

	f()
}

func (room *RoomSync) doWithConn(conn *Connection, f func()) {
	if room.done() || conn.done() {
		return
	}

	defer utils.CatchPanic("doWithConn")

	room.wGroup.Add(1)
	defer room.wGroup.Done()

	conn.wGroup.Add(1)
	defer conn.wGroup.Done()

	f()
}

// Free clear all resources. Call it when no
//  observers and players inside
func (room *RoomSync) Free() {

	if room.checkAndSetCleared() {
		return
	}

	groupWaitRoom := 60 * time.Second // TODO в конфиг
	fieldWaitRoom := 40 * time.Second // TODO в конфиг
	utils.WaitWithTimeout(room.wGroup, groupWaitRoom)

	room.r.events.chanStatus <- StatusAborted

	room.r.events.setStatus(StatusFinished)
	go room.r.connEvents.notify.Free()
	go room.r.messages.Free()
	go room.r.people.Free()
	go room.r.field.Free(fieldWaitRoom)

	close(room.r.events.chanStatus)
}

/////////////////////////// mutex

// checkAndSetCleared checks if the cleanup function was called. This check is
// based on 'done'. If it is true, then the function has already been called.
// If not, set done to True and return false.
// IMPORTANT: this function must only be called in the cleanup function
func (room *RoomSync) checkAndSetCleared() bool {
	room.doneM.Lock()
	defer room.doneM.Unlock()
	if room._done {
		return true
	}
	room._done = true
	return false
}

// done return room readiness flag to free up resources
func (room *RoomSync) done() bool {
	room.doneM.RLock()
	v := room._done
	room.doneM.RUnlock()
	return v
}
