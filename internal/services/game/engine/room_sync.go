package engine

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// SyncI controls the launch of funcs and resource cleanup
// Strategy Pattern
type SyncI interface {
	do(f func())
	doAndFree(clear func())
	doWithConn(conn *Connection, f func())

	done() bool
}

// SyncWgroup implements SyncI
type SyncWgroup struct {
	wGroup *sync.WaitGroup

	doneM *sync.RWMutex
	_done bool
}

func (room *SyncWgroup) Init() {
	room.wGroup = &sync.WaitGroup{}
	room._done = false
	room.doneM = &sync.RWMutex{}
}

// do any action in room by calling this func!
func (room *SyncWgroup) do(f func()) {
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

func (room *SyncWgroup) doWithConn(conn *Connection, f func()) {
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

func (room *SyncWgroup) doAndFree(clear func()) {

	if room.checkAndSetCleared() {
		return
	}

	groupWaitRoom := 60 * time.Second // TODO в конфиг
	utils.WaitWithTimeout(room.wGroup, groupWaitRoom)
	clear()
}

/////////////////////////// mutex

// checkAndSetCleared checks if the cleanup function was called. This check is
// based on 'done'. If it is true, then the function has already been called.
// If not, set done to True and return false.
// IMPORTANT: this function must only be called in the cleanup function
func (room *SyncWgroup) checkAndSetCleared() bool {
	room.doneM.Lock()
	defer room.doneM.Unlock()
	if room._done {
		return true
	}
	room._done = true
	return false
}

// done return room readiness flag to free up resources
func (room *SyncWgroup) done() bool {
	room.doneM.RLock()
	v := room._done
	room.doneM.RUnlock()
	return v
}
