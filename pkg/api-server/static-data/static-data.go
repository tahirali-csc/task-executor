package staticdata

import (
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/services"
)

var buildStatusSvc = services.NewBuildStatusService()

var BuildStatusList = map[string]api.BuildStatus{}

func Init() error {
	statusList, err := buildStatusSvc.List()
	if err != nil {
		return err
	}

	BuildStatusList = statusList
	return nil
}