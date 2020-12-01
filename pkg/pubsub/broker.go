package pubsub

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

//https://eli.thegreenplace.net/2020/pubsub-using-channels-in-go/
type Broker struct {
	sync.RWMutex
	subscribers map[string]subscriberSet
}

type subscriberSet map[*subscriber]struct{}

func NewBroker() Pubsub {
	return &Broker{
		subscribers: make(map[string]subscriberSet),
	}
}

func (broker *Broker) NewTopic(topic string) {
	broker.subscribers[topic] = make(subscriberSet)
}

func (broker *Broker) Publish(ctx context.Context, topic string, msg Message) error {
	broker.RLock()
	defer broker.RUnlock()

	sub, ok := broker.subscribers[topic]
	if !ok {
		return errors.New(fmt.Sprintf("%s topic does not exist", topic))
	}

	for s := range sub {
		s.publish(msg)
	}

	return nil
}

func (broker *Broker) Subscribe(ctx context.Context, topic string) (chan Message, error) {
	broker.Lock()
	defer broker.Unlock()

	subscribers, ok := broker.subscribers[topic]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s topic does not exist", topic))
	}

	s := &subscriber{
		handler: make(chan Message),
		quit:    make(chan struct{}),
	}
	log.Println("About to subscriber")
	subscribers[s] = struct{}{}
	log.Println("Subscriber already set!!!")

	//Allow subscriber to cancel the subscription
	go func() {
		select {
		case <-ctx.Done():
			broker.Lock()
			defer broker.Unlock()
			delete(subscribers, s)
			s.close()
		}
	}()

	return s.handler, nil

}
