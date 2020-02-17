package synced

import (
	"sync"
	"time"
)

// SingleGoroutine struct that launch action only in one goroutine
// call Init for init me and dont forget to Close me!
type SingleGoroutine struct {
	ticker *time.Ticker
	single chan interface{}
	stop   chan interface{}

	stopedM  *sync.RWMutex
	_stopped bool

	action func()
}

func NewSingleGoroutine(
	d time.Duration,
	action func(),
) *SingleGoroutine {
	if d < 1 {
		d = 1
	}
	var sg = &SingleGoroutine{
		ticker:  time.NewTicker(d),
		single:  make(chan interface{}, 1),
		stop:    make(chan interface{}, 1),
		stopedM: &sync.RWMutex{},
		action:  action,
	}
	sg.single <- nil
	return sg
}

func (sg *SingleGoroutine) doIfNotStopped(f func()) {
	sg.stopedM.RLock()
	defer sg.stopedM.RUnlock()
	if sg._stopped {
		return
	}
	f()
}

func (sg *SingleGoroutine) setStopped() {
	sg.stopedM.Lock()
	defer sg.stopedM.Unlock()
	sg._stopped = true
}

func (sg *SingleGoroutine) Close() {
	sg.setStopped()
	sg.ticker.Stop()
	close(sg.single)
}

func (sg *SingleGoroutine) C() <-chan time.Time {
	return sg.ticker.C
}

func (sg *SingleGoroutine) Do() {
	_, ok := <-sg.single
	if !ok {
		return
	}
	sg.action()
	sg.doIfNotStopped(func() { sg.single <- nil })
}
