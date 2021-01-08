package staticdata

import (
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/events"
	"github.com/task-executor/pkg/api-server/services"
	"github.com/task-executor/pkg/core"
)

var buildStatusSvc = services.NewBuildStatusService()
var BuildStatusList = map[string]api.BuildStatus{}

var EventStream *events.DBEvents
var BuildStreamer *services.BuildStreamer
var LogStore core.LogStore

func Init() error {
	statusList, err := buildStatusSvc.List()
	if err != nil {
		return err
	}

	BuildStatusList = statusList
	return nil
}
