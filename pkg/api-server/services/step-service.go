package services

import (
	"fmt"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/dbstore"
	"github.com/task-executor/pkg/api-server/querybuilder"
)

type StepService struct{}

func NewStepService() StepService {
	return StepService{}
}

func (ss StepService) Create(step *api.Step) (*api.Step, error) {

	insertSql := `INSERT INTO step (build_id, name, status, start_ts) VALUES($1, $2, $3, $4) 
		RETURNING id, build_id, status, created_ts, updated_ts`
	row := dbstore.DataSource.QueryRow(insertSql, step.Build.Id, step.Name, step.Status.Id, step.StartTs)

	res := &api.Step{}
	err := row.Scan(&res.Id, &res.Build.Id, &res.Status.Id, &res.CreatedTs, &res.UpdatedTs)

	return res, err
}

func (ss StepService) GetSteps(buildId int64) ([]*api.Step, error) {
	selectStmt := `SELECT * FROM step WHERE build_id=$1`
	rows, err := dbstore.DataSource.Query(selectStmt, buildId)
	if err != nil {
		return nil, err
	}

	var steps []*api.Step
	for rows.Next() {
		step := &api.Step{}
		err := rows.Scan(&step.Id, &step.Build.Id, &step.Name, &step.Status.Id, &step.StartTs, &step.FinishedTs, &step.CreatedTs, &step.FinishedTs)
		if err != nil {
			return nil, err
		}

		steps = append(steps, step)
	}

	return steps, nil
}

func (ss StepService) Filter(values map[string][]string) ([]*api.Step, error) {
	fieldMapping := ss.getFieldMapping()

	filter, err := querybuilder.GetFilterClause(values, fieldMapping)
	if err != nil {
		return nil, err
	}

	rows, err := dbstore.DataSource.Query(fmt.Sprintf("SELECT * FROM step %s", filter))
	if err != nil {
		return nil, err
	}

	var steps []*api.Step
	for rows.Next() {
		step := &api.Step{}
		rows.Scan(&step.Id, &step.Build.Id, &step.Name, &step.Status.Id, &step.StartTs, &step.FinishedTs,
			&step.CreatedTs, &step.FinishedTs)
		steps = append(steps, step)
	}

	return steps, nil
}

//TODO: Set only on bootstrapping
func (ss StepService) getFieldMapping() map[string]querybuilder.Column {
	fieldMap := make(map[string]querybuilder.Column)
	fieldMap["id"] = querybuilder.NewColumn("id", querybuilder.NumberType)
	fieldMap["buildId"] = querybuilder.NewColumn("build_id", querybuilder.NumberType)
	fieldMap["status"] = querybuilder.NewColumn("status", querybuilder.NumberType)
	fieldMap["createdTs"] = querybuilder.NewColumn("created_ts", querybuilder.TimestampType)
	fieldMap["updatedTs"] = querybuilder.NewColumn("updated_ts", querybuilder.TimestampType)
	return fieldMap
}
