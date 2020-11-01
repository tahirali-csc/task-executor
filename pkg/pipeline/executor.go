package pipeline

import (
	dockerrunner "github.com/task-executor/pkg/docker-runner"
	"github.com/task-executor/pkg/docker-runner/api"
)

//func getPipelineCode() string {
//	return `
//package main
//
//func main(){
//	log.Println("I am running a pipeline")
//}
//`
//}

func Run() {
	// 1. scm clone
	//

	//
	mounts := []api.VolumeMount{
		{
			Target: "/brave",
			Source: "/Users/tahir/workspace/rnd-projects/task-executor/pkg/server",
		},
	}

	dockerrunner.Run(api.ContainerConfig{
		Image:   "golang:latest",
		Command: []string{"/bin/sh", "-c", "go get github.com/tahirali-csc/client-api && go run /brave/pipeline.go"},
		Volumes: mounts,
		Env:     []string{"BUILD_ID=1"},
	})

}
