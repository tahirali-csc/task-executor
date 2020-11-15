package api

import "time"

type RunConfig struct {
	Image   string
	Command []string
	//TODO: Will review. Object/Property
	BuildId int64
}

const PendingBuildStatus = "Pending"
const StartedBuildStatus = "Started"
const FinishedBuildStatus = "Finished"

type BuildStatus struct {
	Id   int
	Name string
}

type Build struct {
	Id         int64
	RepoBranch RepoBranch
	Status     BuildStatus
	StartTs    *time.Time
	FinishedTs *time.Time
	CreatedTs  time.Time
	UpdatedTs  time.Time
}

const BasicAuthType = "BasicAuth"
const SSHAuthType = "SSHAuth"

type AuthType struct {
	Id   int
	Name string
}

const KubernetesSecretType = "Kubernetes"

type SecretType struct {
	Id   int
	Name string
}

type Repo struct {
	Id         int64
	Namespace  string
	Name       string
	SSHUrl     string
	HttpUrl    string
	AuthType   AuthType
	SecretType SecretType
	SecretName string
	CreatedTs  time.Time
	UpdatedTs  time.Time
}

type RepoBranch struct {
	Id   int64
	Repo Repo
	Name string
}

type Step struct {
	Id         int64
	Build      Build
	Status     BuildStatus
	StartTs    *time.Time
	FinishedTs *time.Time
	CreatedTs  time.Time
	UpdatedTs  time.Time
}
