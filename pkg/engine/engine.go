package engine

import (
	"context"
	"io"
)

// Engine defines a runtime engine for pipeline execution.
type Engine interface {
	// Setup the pipeline environment.
	Setup(context.Context, *Spec) error

	// Start the pipeline step.
	Start(context.Context, *Spec) error

	Tail(ctx context.Context, spec *Spec) (io.ReadCloser, error)

	Wait(ctx context.Context, spec *Spec) error
}
