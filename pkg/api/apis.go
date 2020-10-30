package api

import "time"

type PipelineConfig struct {
}

type RunConfig struct {
	Image   string
	Command []string
	//TODO: Will review. Object/Property
	BuildId int64
}

type Project struct {
	Id int
}

type BuildStatus struct {
	Id   int
	Name string
}

type Build struct {
	Id         int64
	Project    Project
	Status     BuildStatus
	StartTs    *time.Time
	FinishedTs *time.Time
	CreatedTs  time.Time
	UpdatedTs  time.Time
}
