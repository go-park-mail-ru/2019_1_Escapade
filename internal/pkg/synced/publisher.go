package synced

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

// Msg the message that the publisher sends to observers
type Msg struct {
	Publisher int32
	Action    int32
	Extra     interface{}
}

// NewMsg create new instance Message
func NewMsg(publisher, Action int32, Extra interface{}) Msg {
	return Msg{
		Publisher: publisher,
		Action:    Action,
		Extra:     Extra,
	}
}

// Pair of event code and callback func
type Pair struct {
	Code int32
	Do   func(msg Msg)
}

// NewPair create new pair with selected code and func
func NewPair(code int32, do func(msg Msg)) Pair {
	return Pair{
		Code: code,
		Do:   do,
	}
}

// NewPairNoArgs create new pair with selected code and func*
//	* func does not need to know the message value
func NewPairNoArgs(code int32, do func()) Pair {
	return Pair{
		Code: code,
		Do:   func(msg Msg) { do() },
	}
}

// PublisherI represents observer pattern
type PublisherI interface {
	Observe(o ObserverI)
}

// PublisherBase PublisherI implementation
type PublisherBase struct {
	subscribers []ObserverI
	chanAdd     chan ObserverI
	chanRemove  chan ObserverI
	chanMsg     chan Msg
	chanStop    chan struct{}
}

// NewPublisher create new Published
func NewPublisher() *PublisherBase {
	p := &PublisherBase{
		subscribers: make([]ObserverI, 0),
		chanAdd:     make(chan ObserverI),
		chanRemove:  make(chan ObserverI),
		chanMsg:     make(chan Msg),
		chanStop:    make(chan struct{}),
	}
	return p
}

// Observe subscribe to events of publisher
func (p *PublisherBase) Observe(o ObserverI) {
	p.chanAdd <- o
}

// StartPublish start sending messages to observers
func (p *PublisherBase) StartPublish() {
	go p.start()
}

// StopPublish stop sending messages to observers
func (p *PublisherBase) StopPublish() {
	close(p.chanStop)
}

//////////////////// internal

func (p *PublisherBase) start() {
	defer utils.CatchPanic("catch panic")
	for {
		select {
		case sub := <-p.chanAdd:
			p.addSub(sub)
		case sub := <-p.chanRemove:
			p.removeSub(sub)
		case msg := <-p.chanMsg:
			p.notify(msg)
		case <-p.chanStop:
			p.stop()
			return
		}
	}
}

func (p *PublisherBase) addSub(sub ObserverI) {
	p.subscribers = append(p.subscribers, sub)
	sub.Start()
}

func (p *PublisherBase) removeSub(sub ObserverI) {
	for i, s := range p.subscribers {
		if sub == s {
			p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
			s.Stop()
			break
		}
	}
}

func (p *PublisherBase) Notify(msg Msg) {
	p.chanMsg <- msg
}

func (p *PublisherBase) notify(msg Msg) {
	for _, sub := range p.subscribers {
		sub.Notify(msg)
	}
}

func (p *PublisherBase) stop() {
	for _, sub := range p.subscribers {
		sub.Stop()
	}

	close(p.chanAdd)
	close(p.chanRemove)
	close(p.chanMsg)

	return
}
