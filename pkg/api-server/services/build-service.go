package services

import (
	"fmt"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/dbstore"
	"github.com/task-executor/pkg/api-server/querybuilder"
)

type BuildService struct {}

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

	rows, err := dbstore.DataSource.Query(fmt.Sprintf("SELECT * FROM build %s", filter))
	if err != nil {
		return nil, err
	}

	var builds []api.Build
	for rows.Next() {
		res := api.Build{}
		//err := rows.Scan(&res.Id, &res.Project.Id, &res.Status.Id, &res.StartTs, &res.FinishedTs, &res.CreatedTs, &res.UpdatedTs)

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
