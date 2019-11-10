package engine

/*
type SyncUserI interface {
	GetSync() SyncI
}

// SyncI controls the launch of funcs and resource cleanup
// Strategy Pattern
type SyncI interface {
	do(f func())
	doAndFree(clear func())

	done() bool

	doWithOther(f func(), users ...SyncUserI)
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
func (synced *SyncWgroup) do(f func()) {
	fmt.Println("----1")
	if synced.done() {
		return
	}
	fmt.Println("----2")
	synced.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("catch panic")
		synced.wGroup.Done()
	}()
	fmt.Println("----3")
	f()
}

func (synced *SyncWgroup) doWithConn(conn *Connection, f func()) {
	synced.do(func() {
		conn.s.do(func() {
			f()
		})
	})
}

func (synced *SyncWgroup) doAndFree(clear func()) {

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

// done return room readiness flag to free up resources
func (synced *SyncWgroup) done() bool {
	synced.doneM.RLock()
	v := synced._done
	synced.doneM.RUnlock()
	return v
}
*/
