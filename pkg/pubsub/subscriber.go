package pubsub

import "sync"

type subscriber struct {
	sync.Mutex

	handler chan Message
	quit    chan struct{}
	done    bool
}

func (s *subscriber) publish(msg Message) {
	//select {
	s.handler <- msg
	//}
}

func (s *subscriber) close() {
	s.Lock()
	defer s.Unlock()
	if s.done == false {
		close(s.quit)
		s.done = true
	}
}
