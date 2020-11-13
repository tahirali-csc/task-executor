package services

import (
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/dbstore"
)

type BuildStatusService struct{}

func NewBuildStatusService() BuildStatusService {
	return BuildStatusService{}
}

func (bs BuildStatusService) List() (map[string]api.BuildStatus, error) {

	rows, err := dbstore.DataSource.Query("SELECT * FROM build_status")
	if err != nil {
		return nil, err
	}

	status := make(map[string]api.BuildStatus)
	for rows.Next() {
		bs := api.BuildStatus{}
		err = rows.Scan(&bs.Id, &bs.Name)
		if err != nil {
			return nil, err
		}
		status[bs.Name] = bs
	}

	return status, nil
}
