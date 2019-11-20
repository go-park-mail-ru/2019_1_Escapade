package synced

import (
	"time"
)

// SingleGoroutine struct that launch action only in one goroutine
// call Init for init me and dont forget to Close me!
type SingleGoroutine struct {
	ticker *time.Ticker
	single chan interface{}
	action func()
}

func (sg *SingleGoroutine) Init(d time.Duration, action func()) {
	if d < 1 {
		d = 1
	}
	sg.ticker = time.NewTicker(d)
	sg.single = make(chan interface{}, 1)
	sg.action = action
}

func (sg *SingleGoroutine) Close() {
	sg.ticker.Stop()
	close(sg.single)
}

func (sg *SingleGoroutine) C() <-chan time.Time {
	return sg.ticker.C
}

func (sg *SingleGoroutine) Do() {
	sg.single <- nil
	sg.action()
	<-sg.single
}
