package events

import (
	"context"
	"encoding/json"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api-server/dbstore"
	"github.com/task-executor/pkg/pubsub"
	"time"
)

type DBEvents struct {
	pubsub pubsub.Pubsub
}

func NewDBEvents() *DBEvents {
	broker := pubsub.NewBroker()
	return &DBEvents{
		pubsub: broker,
	}
}

func (dbEvents *DBEvents) NewTopic(topic string) {
	dbEvents.pubsub.NewTopic(topic)
}

func (dbEvents *DBEvents) Subscribe(ctx context.Context, topic string) (chan pubsub.Message, error) {
	return dbEvents.pubsub.Subscribe(ctx, topic)
}

func (dbEvents *DBEvents) Start() error {

	eventListener := pq.NewListener(dbstore.ConnString, 10*time.Second, time.Minute, nil)
	if err := eventListener.Listen("events"); err != nil {
		return err
	}

	for {
		select {
		case eve := <-eventListener.Notify:
			//log.Debug(eve)


			m := make(map[string]interface{})
			json.Unmarshal([]byte(eve.Extra), &m)
			//log.Println(m)

			table := m["table"].(string)

			err := dbEvents.pubsub.Publish(context.Background(), table, m["data"])
			if err != nil {
				log.Println(err)
			}
			//eve.Extra
		}
	}
}
