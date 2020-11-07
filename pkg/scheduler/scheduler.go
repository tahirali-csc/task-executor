package scheduler

import (
	"context"
	"github.com/task-executor/pkg/core"
)

type Scheduler interface {
	Schedule(context context.Context, stage *core.Stage, initContainers []core.InitContainer) error
}
