package core

import "github.com/task-executor/pkg/api"

type StepRun struct {
	Step *api.Step
	Name     string
	Image    string
	Cmd      []string
	Args     []string
	CpuLimit int
	Memory   int
	BuildId  int64
}

type Stage struct {
	Name            string
	Image           string
	ImagePullPolicy string
	LimitMemory     int
	LimitCompute    int
	RequestMemory   int
	RequestCompute  int
	Command         []string
	Args            []string
	Volume          []InitVolume
	BuildId         int64
}

type InitContainer struct {
	Name            string
	Image           string
	ImagePullPolicy string
	Command         []string
	Secrets         []SecretSource
	Volume          []InitVolume
}

type SecretFromRef struct {
	Name string
	Key  string
}
type SecretSource struct {
	Name string
	From SecretFromRef
}
type InitVolume struct {
	Name      string
	MountPath string
}
