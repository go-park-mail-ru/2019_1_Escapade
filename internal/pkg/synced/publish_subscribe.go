package synced

type SubscriberI interface {
	Start()
	Notify(msg Msg)
	Close()
}

type PublisherI interface {
	Start(subs ...SubscriberI)

	start()
	addSubscriber(sub SubscriberI)
	removeSubscriber(sub SubscriberI)
	Notify(msg Msg)

	Stop()
}

type PublisherBase struct {
	subscribers []SubscriberI
	addSubCh    chan SubscriberI
	removeSubCh chan SubscriberI
	msgCh       chan Msg
	stopCh      chan struct{}
}

func (p *PublisherBase) addSubscriber(sub SubscriberI) {
	p.addSubCh <- sub
}
func (p *PublisherBase) removeSubscribe(sub SubscriberI) {
	p.removeSubCh <- sub
}
func (p *PublisherBase) Notify(msg Msg) {
	p.msgCh <- msg
}
func (p *PublisherBase) Stop() {
	close(p.stopCh)
}

func (p *PublisherBase) start() {
	for {
		select {
		case sub := <-p.addSubCh:
			p.addSub(sub)
		case sub := <-p.removeSubCh:
			p.removeSub(sub)
		case msg := <-p.msgCh:
			p.notify(msg)
		case <-p.stopCh:
			p.stop()
		}
	}
}

func (p *PublisherBase) addSub(sub SubscriberI) {
	p.subscribers = append(p.subscribers, sub)
}

func (p *PublisherBase) removeSub(sub SubscriberI) {
	for i, s := range p.subscribers {
		if sub == s {
			p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
			s.Close()
			break
		}
	}
}

func (p *PublisherBase) notify(msg Msg) {
	for _, sub := range p.subscribers {
		sub.Notify(msg)
	}
}

func (p *PublisherBase) stop() {
	for _, sub := range p.subscribers {
		sub.Close()
	}

	close(p.addSubCh)
	close(p.removeSubCh)
	close(p.msgCh)

	return
}

type SubscriberBase struct {
	in      chan Msg
	stop    chan struct{}
	process func(msg Msg)
}

func (s SubscriberBase) Notify(msg Msg) {
	s.in <- msg
}
func (s SubscriberBase) Close() {
	close(s.stop)
}
func (s SubscriberBase) Start() {
	for {
		select {
		case msg := <-s.in:
			s.process(msg)
		case <-s.stop:
			close(s.in)
			return
		}
	}
}

func NewSubscriber(call func(msg Msg)) SubscriberBase {
	sub := SubscriberBase{
		in:      make(chan Msg),
		stop:    make(chan struct{}),
		process: call,
	}
	return sub
}

func (p PublisherBase) Start(subs ...SubscriberI) {
	p.start()
	for _, sub := range subs {
		sub.Start()
		p.addSubscriber(sub)
	}
}

func NewPublisher() *PublisherBase {
	p := PublisherBase{
		addSubCh:    make(chan SubscriberI),
		removeSubCh: make(chan SubscriberI),
		msgCh:       make(chan Msg),
		stopCh:      make(chan struct{}),
	}
	return &p
}

type Msg struct {
	Code    int32
	Content interface{}
}
