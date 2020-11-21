package runner

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/core"
	engine2 "github.com/task-executor/pkg/engine"
	"github.com/task-executor/pkg/engine/kube"
)

type Runner struct {
}

func NewRunner() *Runner {
	return nil
}

//func (runner *Runner) Run() {
//	stepClient := apiclient.NewStepClient(apiclient.Config{
//		Host: "localhost:8080",
//	})
//	filter := apiclient.NewStepFilterBuilder()
//	filter.WithStatus(1)
//	filter.ByCreatedTs(false)
//
//	for {
//		steps, err := stepClient.GetSteps(filter)
//		if err != nil {
//			log.Println(err)
//		} else {
//			log.Println(steps)
//		}
//
//		for i := 0; i < len(steps); i++ {
//			runner.runStep(&steps[i])
//		}
//
//		time.Sleep(time.Second * 5)
//	}
//}

func (runner *Runner) Run(step *core.StepRun) {
	engine, err := kube.NewFile("", "/Users/tahir/.kube/config", "")
	if err != nil {
		log.Println(err)
		return
	}

	err = engine.Start(context.Background(), &engine2.Spec{
		Image:   step.Image,
		Command: step.Cmd,
		Args:    step.Args,
		Metadata: engine2.Metadata{
			Namespace: "default",
			//TODO: Can add more randomization
			UID: fmt.Sprintf("te-step-%d", step.Step.Id),
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

}
