package services

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/events"
	"github.com/task-executor/pkg/pubsub"
	"sync"
	"time"
)

var buildService = NewBuildService()

type Subscriber struct {
	BuildChannel chan *api.Build
	StepChannel  chan *api.Step
}

type SubscriberSet map[*Subscriber]struct{}

type BuildStreamer struct {
	dbEvents    *events.DBEvents
	subscribers map[int64]SubscriberSet
	sync.RWMutex

	buildDBChannel chan pubsub.Message
	stepDBChannel  chan pubsub.Message
}

func NewBuildStreamer(dbEvents *events.DBEvents) *BuildStreamer {
	ctx, _ := context.WithCancel(context.Background())

	buildMessageChan, _ := dbEvents.Subscribe(ctx, "build")
	stepMessageChan, _ := dbEvents.Subscribe(ctx, "step")

	return &BuildStreamer{
		dbEvents:       dbEvents,
		subscribers:    make(map[int64]SubscriberSet),
		buildDBChannel: buildMessageChan,
		stepDBChannel:  stepMessageChan,
	}
}

func (bs *BuildStreamer) Subscribe(ctx context.Context, buildId int64) *Subscriber {
	defer bs.Unlock()
	bs.Lock()

	// Add new subscriber for build id
	newSubscriber := &Subscriber{
		BuildChannel: make(chan *api.Build),
		StepChannel:  make(chan *api.Step),
	}

	_, ok := bs.subscribers[buildId]
	if !ok {
		bs.subscribers[buildId] = make(map[*Subscriber]struct{})
	}

	//TODO: review it
	buildSubscribers := bs.subscribers[buildId]
	buildSubscribers[newSubscriber] = struct{}{}

	// Cleanup the subscriber
	go func() {
		select {
		case <-ctx.Done():
			bs.Lock()
			delete(buildSubscribers, newSubscriber)

			//If there are no subscriber for this build.
			if len(bs.subscribers[buildId]) == 0 {
				delete(bs.subscribers, buildId)
			}
			log.Debug("Deleting subscriber:::", len(buildSubscribers), bs.subscribers)
			bs.Unlock()
		}
	}()
	return newSubscriber
}

func (bs *BuildStreamer) Start() {
	timer := time.NewTicker(time.Second * 15)

	publish := func(build *api.Build) {
		defer bs.RUnlock()
		bs.RLock()
		for subscriber, _ := range bs.subscribers[build.Id] {
			go func() {
				subscriber.BuildChannel <- build
			}()
		}
	}

	//Consolidate build event from build and step table. It consolidates from:
	// 1. Using polling
	// 2. build table
	// 3. step table
	for {
		select {
		//Poll DB for any updates
		case <-timer.C:
			var currentBuilds []int64
			bs.RLock()
			for k, _ := range bs.subscribers {
				currentBuilds = append(currentBuilds, k)
			}
			bs.RUnlock()

			log.Debug("Getting up to date info from DB.....", currentBuilds)
			builds, _ := buildService.DeepFetch(currentBuilds)
			for i := 0; i < len(builds); i++ {
				publish(builds[i])
			}

		//Event received from build table
		case buildObj := <-bs.buildDBChannel:
			record := buildObj.(map[string]interface{})
			log.Debug("Build Event:::", record)
			buildId := int64(record["id"].(float64))

			bs.RLock()
			for subscriber, _ := range bs.subscribers[buildId] {
				go func() {
					subscriber.BuildChannel <- &api.Build{
						Id: buildId,
						RepoBranch: api.RepoBranch{
							Id: int64(record["repo_branch"].(float64)),
						},
						Status: api.BuildStatus{
							Id:   int(record["status"].(float64)),
							Name: record["status_name"].(string),
						},
						//TODO: Handle more fields
					}
				}()
			}
			bs.RUnlock()

		//Event received from step table
		case stepObj := <-bs.stepDBChannel:
			record := stepObj.(map[string]interface{})
			log.Println("Step Event:::", record)
			buildId := int64(record["build_id"].(float64))
			stepId := int64(record["id"].(float64))

			bs.RLock()
			for subscriber, _ := range bs.subscribers[buildId] {
				go func() {
					subscriber.StepChannel <- &api.Step{
						Id: stepId,
						Build: api.Build{
							Id: buildId,
						},
						Name: record["name"].(string),
						Status: api.BuildStatus{
							Id:   int(record["status"].(float64)),
							Name: record["status_name"].(string),
						},
						//TODO: Handle more fields
					}
				}()
			}
			bs.RUnlock()
		}
	}
}
