package core

import (
	"context"
	"io"
)

// LogStore persists build output to storage.
type LogStore interface {
	// Find returns a log stream from the datastore.
	Find(ctx context.Context, stage int64) (io.ReadCloser, error)

	// Update writes copies the log stream from Reader r to the datastore.
	Upload(ctx context.Context, stepId int64, log io.Reader) error
}
