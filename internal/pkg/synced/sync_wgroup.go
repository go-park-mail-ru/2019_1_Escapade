package synced

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

type SyncUserI interface {
	GetSync() SyncI
}

// SyncI controls the launch of funcs and resource cleanup
// Strategy Pattern
type SyncI interface {
	Do(f func())
	DoWithOther(user SyncUserI, f func())
	DoWithOthers(f func(), users ...SyncUserI)

	Clear(clear func())
	IsCleared() bool
}

// SyncWgroup implements SyncI
type SyncWgroup struct {
	wGroup *sync.WaitGroup

	doneM *sync.RWMutex
	_done bool
	wait  time.Duration
}

func (synced *SyncWgroup) Init(wait time.Duration) {
	synced.wGroup = &sync.WaitGroup{}
	synced._done = false
	synced.doneM = &sync.RWMutex{}
	synced.wait = wait
}

// do any action in room by calling this func!
func (synced *SyncWgroup) Do(f func()) {
	if synced.IsCleared() {
		return
	}
	synced.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("catch panic")
		synced.wGroup.Done()
	}()
	f()
}

func (synced *SyncWgroup) DoWithOther(other SyncUserI, f func()) {
	synced.Do(func() {
		other.GetSync().Do(func() {
			f()
		})
	})
}

func (synced *SyncWgroup) DoWithOthers(f func(), next ...SyncUserI) {
	synced.Do(func() {
		if len(next) > 0 {
			next[0].GetSync().DoWithOthers(f, next[1:]...)
			return
		}
		f()
	})
}

func (synced *SyncWgroup) Clear(clear func()) {

	if synced.checkAndSetCleared() {
		return
	}

	synced.waitWithTimeout()
	clear()
}

func (synced *SyncWgroup) waitWithTimeout() bool {
	// without timeout mode
	if synced.wait.Seconds() < 1 {
		synced.wGroup.Wait()
		return false
	}

	c := make(chan struct{})
	go func() {
		defer close(c)
		synced.wGroup.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(synced.wait):
		return true // timed out
	}
}

/////////////////////////// mutex

// checkAndSetCleared checks if the cleanup function was called. This check is
// based on 'done'. If it is true, then the function has already been called.
// If not, set done to True and return false.
// IMPORTANT: this function must only be called in the cleanup function
func (synced *SyncWgroup) checkAndSetCleared() bool {
	synced.doneM.Lock()
	defer synced.doneM.Unlock()
	if synced._done {
		return true
	}
	synced._done = true
	return false
}

// IsCleared check were resourse cleared
func (synced *SyncWgroup) IsCleared() bool {
	synced.doneM.RLock()
	v := synced._done
	synced.doneM.RUnlock()
	return v
}
