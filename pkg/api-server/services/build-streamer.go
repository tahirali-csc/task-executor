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

type Streamer struct {
	BuildChannel chan *api.Build
	StepChannel  chan *api.Step
}

type StreamerSet map[*Streamer]struct{}

type BuildStreamer struct {
	dbEvents          *events.DBEvents
	buildIdSet        map[int64]struct{}
	buildEventChannel map[int64]StreamerSet
	sync.RWMutex

	buildDBChannel chan pubsub.Message
	stepDBChannel  chan pubsub.Message
}

func NewBuildStreamer(dbEvents *events.DBEvents) *BuildStreamer {
	ctx, _ := context.WithCancel(context.Background())

	buildMessageChan, _ := dbEvents.Subscribe(ctx, "build")
	stepMessageChan, _ := dbEvents.Subscribe(ctx, "step")

	return &BuildStreamer{
		dbEvents:          dbEvents,
		buildIdSet:        make(map[int64]struct{}),
		buildEventChannel: make(map[int64]StreamerSet),
		buildDBChannel:    buildMessageChan,
		stepDBChannel:     stepMessageChan,
	}
}

func (bs *BuildStreamer) Subscribe(ctx context.Context, buildId int64) *Streamer {
	defer bs.Unlock()
	bs.Lock()

	streamer := &Streamer{
		BuildChannel: make(chan *api.Build),
		StepChannel:  make(chan *api.Step),
	}

	_, ok := bs.buildEventChannel[buildId]
	if !ok {
		bs.buildEventChannel[buildId] = make(map[*Streamer]struct{})
	}

	subscribers := bs.buildEventChannel[buildId]

	//??
	subscribers[streamer] = struct{}{}
	//subscribers = append(subscribers, streamer)
	//???
	//bs.buildEventChannel[buildId] = subscribers
	go func() {
		select {
		case <-ctx.Done():
			bs.Lock()
			delete(subscribers, streamer)
			log.Println("Deleting subsciber:::", len(subscribers))
			bs.Unlock()
		}
	}()
	return streamer
}

func (bs *BuildStreamer) Start() {
	timer := time.NewTicker(time.Second * 20)

	for {
		select {
		case <-timer.C:
			var temp []int64
			for k, _ := range bs.buildIdSet {
				temp = append(temp, k)
			}

			builds, _ := buildService.DeepFetch(temp)
			for i := 0; i < len(builds); i++ {
				bid := builds[i].Id

				bs.RLock()
				for subscriber, _ := range bs.buildEventChannel[bid] {
					go func() {
						subscriber.BuildChannel <- builds[i]
					}()
				}
				bs.RUnlock()

			}

		case buildObj := <-bs.buildDBChannel:
			record := buildObj.(map[string]interface{})
			log.Println("Build Change:::", record)
			buildId := int64(record["id"].(float64))

			bs.RLock()
			for subscriber, _ := range bs.buildEventChannel[buildId] {
				go func() {
					subscriber.BuildChannel <- &api.Build{
						Id: buildId,
						RepoBranch: api.RepoBranch{
							Id: int64(record["repo_branch"].(float64)),
						},
						Status: api.BuildStatus{
							Id: int(record["status"].(float64)),
						},
						//TODO: Handle more fields
					}
				}()
			}
			bs.RUnlock()

		case stepObj := <-bs.stepDBChannel:
			record := stepObj.(map[string]interface{})
			log.Println("Step Obj:::", record)
			buildId := int64(record["build_id"].(float64))
			stepId := int64(record["id"].(float64))

			bs.RLock()
			for subscriber, _ := range bs.buildEventChannel[buildId] {
				go func() {
					subscriber.StepChannel <- &api.Step{
						Id: stepId,
						Build: api.Build{
							Id: buildId,
						},
						Name: record["name"].(string),
						Status: api.BuildStatus{
							Id: int(record["status"].(float64)),
						},
						//TODO: Handle more fields
					}
				}()
			}
			bs.RUnlock()
		}
	}
}
