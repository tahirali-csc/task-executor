package services

import (
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/dbstore"
)

type StepService struct{}

func NewStepService() StepService {
	return StepService{}
}

func (bs StepService) Create(step *api.Step) (*api.Step, error) {

	insertSql := `INSERT INTO step (build_id,name,status) VALUES($1,$2,$3) 
		RETURNING id, build_id, status, created_ts, updated_ts`
	row := dbstore.DataSource.QueryRow(insertSql, step.Build.Id, step.Name, step.Status.Id)

	res := &api.Step{}
	err := row.Scan(&res.Id, &res.Build.Id, &res.Status.Id, &res.CreatedTs, &res.UpdatedTs)

	return res, err
}
