package services

import (
	"fmt"
	"strings"

	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/dbstore"
	"github.com/task-executor/pkg/api-server/querybuilder"
)

type BuildService struct{}

func NewBuildService() BuildService {
	return BuildService{}
}

func (bs BuildService) Create(build *api.Build) (*api.Build, error) {

	insertSql := `INSERT INTO build (repo_branch,status) VALUES($1,$2) 
		RETURNING id, repo_branch, status, created_ts, updated_ts`
	row := dbstore.DataSource.QueryRow(insertSql, build.RepoBranch.Id, build.Status.Id)

	res := &api.Build{}
	err := row.Scan(&res.Id, &res.RepoBranch.Id, &res.Status.Id, &res.CreatedTs, &res.UpdatedTs)

	return res, err
}

func (bs BuildService) Filter(values map[string][]string) ([]api.Build, error) {
	fieldMapping := bs.getFieldMapping()

	filter, err := querybuilder.GetFilterClause(values, fieldMapping)
	if err != nil {
		return nil, err
	}

	rows, err := dbstore.DataSource.Query(fmt.Sprintf("SELECT b.*, bs.name FROM build b INNER JOIN build_status bs ON b.status=bs.id %s", filter))
	if err != nil {
		return nil, err
	}

	var builds []api.Build
	for rows.Next() {
		res := api.Build{}
		err := rows.Scan(&res.Id, &res.RepoBranch.Id, &res.Status.Id, &res.StartTs, &res.FinishedTs, &res.CreatedTs, &res.UpdatedTs, &res.Status.Name)

		if err != nil {
			return nil, err
		}

		builds = append(builds, res)
	}

	return builds, nil
}

//TODO: Set only on bootstrapping
func (bs BuildService) getFieldMapping() map[string]querybuilder.Column {
	fieldMap := make(map[string]querybuilder.Column)
	fieldMap["id"] = querybuilder.NewColumn("id", querybuilder.NumberType)
	fieldMap["projectId"] = querybuilder.NewColumn("project_id", querybuilder.NumberType)
	fieldMap["status"] = querybuilder.NewColumn("status", querybuilder.NumberType)
	fieldMap["createdTs"] = querybuilder.NewColumn("created_ts", querybuilder.TimestampType)
	fieldMap["updatedTs"] = querybuilder.NewColumn("updated_ts", querybuilder.TimestampType)
	return fieldMap
}

func (bs BuildService) GetStatus(buildId int64) (*api.BuildStatus, error) {

	selectStmt := `SELECT b.status, bs.name
	FROM build b
	INNER JOIN build_status bs ON b.status = bs.id
	WHERE b.id=$1`

	//selectStmt := `SELECT s.status, bs.name FROM step s
	//INNER JOIN build_status bs ON s.status = bs.id
	//WHERE s.id=$1`
	row := dbstore.DataSource.QueryRow(selectStmt, buildId)

	status := &api.BuildStatus{}
	err := row.Scan(&status.Id, &status.Name)

	return status, err
}

func (bs BuildService) UpdateStatus(stepId int64, statusId int) error {
	updateStmt := `UPDATE build SET status=$1 WHERE id=$2`
	_, err := dbstore.DataSource.Exec(updateStmt, statusId, stepId)
	return err
}

func (bs BuildService) DeepFetch(buildId []int64) ([]*api.Build, error) {
	var strArr []string
	for i := 0; i < len(buildId); i++ {
		strArr = append(strArr, fmt.Sprintf("%d", buildId[i]))
	}

	sid := strings.Join(strArr, ",")

	selectStmt := `SELECT b.id, b.repo_branch, b.status, s.name as status_name, b.start_ts, b.finished_ts , b.created_ts , b.updated_ts
FROM build b
INNER JOIN build_status s ON b.status = s.id
WHERE b.id IN `

	selectStmt = selectStmt + "(" + sid + ")"

	rows, err := dbstore.DataSource.Query(selectStmt)
	if err != nil {
		return nil, err
	}

	var builds []*api.Build
	buildMap := make(map[int64]*api.Build)
	for rows.Next() {
		b := &api.Build{Steps: []*api.Step{}}
		if err := rows.Scan(&b.Id, &b.RepoBranch.Id, &b.Status.Id, &b.Status.Name, &b.StartTs, &b.FinishedTs, &b.CreatedTs, &b.UpdatedTs); err != nil {
			return nil, err
		}
		buildMap[b.Id] = b
		builds = append(builds, b)
	}

	selectStmt = `SELECT st.id, st.build_id, st.name, st.status, st.start_ts, st.finished_ts, st.created_ts, st.updated_ts, s.name
FROM step st
INNER JOIN build_status s ON st.status = s.id
WHERE build_id IN`

	selectStmt = selectStmt + "(" + sid + ")"

	rows, err = dbstore.DataSource.Query(selectStmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		s := &api.Step{}
		if err := rows.Scan(&s.Id, &s.Build.Id, &s.Name, &s.Status.Id, &s.StartTs, &s.FinishedTs, &s.CreatedTs,
			&s.UpdatedTs, &s.Status.Name); err != nil {
			return nil, err
		}
		buildMap[s.Build.Id].Steps = append(buildMap[s.Build.Id].Steps, s)
	}

	return builds, nil
}
