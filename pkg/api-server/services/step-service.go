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
	selectStmt := `SELECT st.*, s.name as step_status FROM step st
	INNER JOIN build_status s ON st.status = s.id
	WHERE build_id=$1 ORDER BY id asc`
	rows, err := dbstore.DataSource.Query(selectStmt, buildId)
	if err != nil {
		return nil, err
	}

	var steps []*api.Step
	for rows.Next() {
		step := &api.Step{}
		err := rows.Scan(&step.Id, &step.Build.Id, &step.Name, &step.Status.Id, &step.StartTs, &step.FinishedTs,
			&step.CreatedTs, &step.FinishedTs, &step.Status.Name)
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

func (ss StepService) UpdateStatus(stepId int64, statusId int) error {
	updateStmt := `UPDATE step SET status=$1 WHERE id=$2`
	_, err := dbstore.DataSource.Exec(updateStmt, statusId, stepId)
	return err
}

func (ss StepService) GetStatus(stepId int64) (*api.BuildStatus, error) {

	selectStmt := `SELECT s.status, bs.name FROM step s
	INNER JOIN build_status bs ON s.status = bs.id
	WHERE s.id=$1`
	row := dbstore.DataSource.QueryRow(selectStmt, stepId)

	status := &api.BuildStatus{}
	err := row.Scan(&status.Id, &status.Name)

	return status, err
}

func (ss StepService) UploadLogs(stepId int64, log []byte) error {
	insertStmt := `INSERT INTO logs(step_id, log_data) VALUES($1, $2)
ON conflict(step_ID) do update set log_data = logs.log_data || $2
	`
	_, err := dbstore.DataSource.Exec(insertStmt, stepId, log)
	return err
}

func (ss StepService) GetLogs(stepId int64) ([]byte, error) {

	var logs []byte
	row := dbstore.DataSource.QueryRow(`SELECT log_data FROM logs WHERE step_id=$1`, stepId)
	err := row.Scan(&logs)

	if err != nil {
		return nil, err
	}

	return logs, nil
}
