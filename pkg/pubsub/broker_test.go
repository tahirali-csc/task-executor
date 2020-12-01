package pubsub

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
	"sync"
	"testing"
)

func TestBroker(t *testing.T) {
	const (
		Build = "Build"
		Task  = "Task"
	)

	var wg sync.WaitGroup
	broker := NewBroker()
	broker.NewTopic(Build)

	count := 1

	receiver := func(idx int64) {
		//receiveCount := count
		ctx, _ := context.WithCancel(context.Background())
		msgChan, _ := broker.Subscribe(ctx, Build)

		go func() {
			var received []*api.Build
			select {
			case m := <-msgChan:
				received = append(received, m.(*api.Build))
				log.Println("Received::", idx, (m.(*api.Build)).Id)
				//receiveCount--

				//if receiveCount <= 0 {
				//	for _, b := range received {
				//		log.Println(b.Id)
				//	}
				ctx.Done()
				wg.Done()
				return
				//}
			}
			//}
		}()

		log.Println("Getting out....")

	}

	producer := func(idx int64) {
		//log.Println("Sending :::", idx)
		ctx, _ := context.WithCancel(context.Background())
		err := broker.Publish(ctx, Build, &api.Build{
			Id: idx,
		})
		if err != nil {
			log.Println(err)
		}
	}

	for i := 1; i <= count; i++ {
		wg.Add(1)
		receiver(int64(i))
		go producer(int64(i))

	}

	//time.Sleep(time.Second * 10)
	wg.Wait()
	log.Println("Done")

}
