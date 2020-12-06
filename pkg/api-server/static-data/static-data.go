package staticdata

import (
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/events"
	"github.com/task-executor/pkg/api-server/services"
)

var buildStatusSvc = services.NewBuildStatusService()
var BuildStatusList = map[string]api.BuildStatus{}

var EventStream *events.DBEvents
var BuildStreamer *services.BuildStreamer

func Init() error {
	statusList, err := buildStatusSvc.List()
	if err != nil {
		return err
	}

	BuildStatusList = statusList
	return nil
}