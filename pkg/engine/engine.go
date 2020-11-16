package engine

import "context"

// Engine defines a runtime engine for pipeline execution.
type Engine interface {
	// Setup the pipeline environment.
	Setup(context.Context, *Spec) error

	// Start the pipeline step.
	Start(context.Context, *Spec) error
}
