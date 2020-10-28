package services

import (
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/dbstore"
	"time"
)

type BuildService struct {
}

func NewBuildService() BuildService {
	return BuildService{}
}

func (bs BuildService) Create(build *api.Build) (*api.Build, error) {
	db, err := dbstore.GetDb()
	if err != nil {
		return nil, err
	}
	defer dbstore.Release(db)

	insertSql := `INSERT INTO build (project_id,status) VALUES($1,$2) 
		RETURNING id, project_id, status, created_ts, updated_ts`
	row := db.QueryRow(insertSql, build.Project.Id, build.Status.Id)

	var (
		id                   int64
		projectId            int
		status               int
		createdTs, updatedTs time.Time
	)

	err = row.Scan(&id, &projectId, &status, &createdTs, &updatedTs)
	if err != nil {
		return nil, err
	}

	return &api.Build{
		Id:        id,
		Project:   api.Project{Id: projectId},
		Status:    api.BuildStatus{Id: status},
		CreatedTs: createdTs,
		UpdatedTs: updatedTs,
	}, nil
}
