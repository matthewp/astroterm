package actors

type broker[T any] struct {
	stopCh    chan int
	publishCh chan T
	subCh     chan chan T
	unsubCh   chan chan T
}

func newBroker[T any]() *broker[T] {
	return &broker[T]{
		stopCh:    make(chan int),
		publishCh: make(chan T, 1),
		subCh:     make(chan chan T, 1),
		unsubCh:   make(chan chan T, 1),
	}
}

func (b *broker[T]) Start() {
	subs := map[chan T]struct{}{}
	for {
		select {
		case <-b.stopCh:
			return
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
		case msg := <-b.publishCh:
			for msgCh := range subs {
				// msgCh is buffered, use non-blocking send to protect the broker:
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

func (b *broker[T]) Stop() {
	close(b.stopCh)
}

func (b *broker[T]) Subscribe() chan T {
	msgCh := make(chan T, 5)
	b.subCh <- msgCh
	return msgCh
}

func (b *broker[T]) Unsubscribe(msgCh chan T) {
	b.unsubCh <- msgCh
}

func (b *broker[T]) Publish(msg T) {
	b.publishCh <- msg
}

func (b *broker[T]) Listen(cb func(T)) {
	go func() {
		ch := b.Subscribe()
		for {
			msg := <-ch
			cb(msg)
		}
	}()
}
