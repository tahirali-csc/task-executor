package staticdata

import (
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/services"
	"sync"
)

var buildStatusSvc = services.NewBuildStatusService()
var BuildStatusList = map[string]api.BuildStatus{}

//var StepChannel = make(chan *api.Step)
//var BuildChannel = make(chan *api.Build)
var EventBroker = Broker{}

func Init() error {
	statusList, err := buildStatusSvc.List()
	if err != nil {
		return err
	}

	BuildStatusList = statusList
	return nil
}

type StepCallback func(step *api.Step)

type Broker struct {
	lock sync.Mutex
	subscribers map[int64][]StepCallback
}

func (br *Broker) Publish(buildId int64, step *api.Step) {
	for _, fn := range br.subscribers[buildId] {
		log.Println("Publishing...", buildId, step)
		fn(step)
	}
}

func (br *Broker) Subscribe(buildId int64, callBack StepCallback) {
	br.lock.Lock()
	list := br.subscribers[buildId]
	list = append(list, callBack)
	br.lock.Unlock()
}
