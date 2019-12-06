package synced

// ObserverI represents observer pattern
type ObserverI interface {
	AddPreAction(a func()) ObserverI
	AddCheck(c func() bool) ObserverI
	AddPublisherCode(int32) ObserverI

	Start()
	Notify(msg Msg)
	Stop()
}

// ObserverBase ObserverI implementation
type ObserverBase struct {
	pairs     []Pair
	preAction func()
	check     func() bool
	publisher int32

	in   chan Msg
	stop chan struct{}
}

// NewObserver create new instance of ObserverBase
func NewObserver(pairs ...Pair) *ObserverBase {
	var o = &ObserverBase{}
	o.pairs = pairs
	o.in = make(chan Msg)
	o.stop = make(chan struct{})
	o.publisher = -1
	return o
}

// AddPreAction set the callback, that will be called on
//  any publisher event no matter what the event's code is
func (o *ObserverBase) AddPreAction(a func()) ObserverI {
	o.preAction = a
	return o
}

// AddCheck set the callback function to be called
// 	before starting the main action. If this function
// return false, the function will be aborted
func (o *ObserverBase) AddCheck(c func() bool) ObserverI {
	o.check = c
	return o
}

// AddPublisherCode set the publisher code checking
//  by default code is -1, it is mean no check. If code >= 0
// 	'publisher check' will turn on
func (o *ObserverBase) AddPublisherCode(publisher int32) ObserverI {
	o.publisher = publisher
	return o
}

// Notify send message to observer
func (o *ObserverBase) Notify(msg Msg) {
	o.in <- msg
}

// Stop listening for new messages from publisher
func (o *ObserverBase) Stop() {
	close(o.stop)
}

// Start listening for new messages from publisher
func (o *ObserverBase) Start() {
	go o.start()
}

// do - action performed when a message is received
//  from the publisher
func (o *ObserverBase) do(msg Msg) {
	// check callback
	if o.check != nil && !o.check() {
		return
	}

	// publisher code check
	if o.publisher > -1 && o.publisher != msg.Publisher {
		return
	}

	// preAction callback
	if o.preAction != nil {
		o.preAction()
	}
	for _, pair := range o.pairs {
		if pair.Code == msg.Action {
			pair.Do(msg)
			return
		}
	}
}

func (o *ObserverBase) start() {
	for {
		select {
		case msg := <-o.in:
			o.do(msg)
		case <-o.stop:
			close(o.in)
			return
		}
	}
}
