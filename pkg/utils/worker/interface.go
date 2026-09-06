package worker

import (
	"context"
	"time"
)

// WorkerFunc is a function executed by [Worker].
// Exactly one call is in-flight at any time.
//
// Return value controls the next execution:
//   - error is non-nil — retried after [WithErrRetry] delay;
//   - [WorkerResult.IsCompleted] is true — not retried;
//   - otherwise — retried after [WorkerResult.RetryAfter] delay.
//
// The context may be canceled — the function should return promptly.
// Panics are not recovered and crash the process.
type WorkerFunc func(context.Context) (WorkerResult, error)

// WorkerResult is the result of a single [WorkerFunc] execution.
type WorkerResult struct {
	// IsCompleted reports whether the task is completed.
	// Completed task is not retried.
	IsCompleted bool

	// RetryAfter is the delay before an incomplete task is retried.
	// Zero or negative means immediate retry.
	RetryAfter time.Duration
}

// Interface is the worker contract implemented by [Worker].
type Interface interface {
	// Serve starts the worker and blocks
	// until the context is canceled or [Worker.Stop] is called.
	// Returns [ErrRunning] when already running,
	// [ErrStopped] when called after [Worker.Stop].
	Serve(context.Context) error

	// Stop signals the worker goroutine to shut down
	// without waiting for an in-flight function call to complete.
	// Stopped worker cannot be started again.
	Stop()

	// WakeUp schedules the function for execution.
	// If it is already running, it is executed again after completion.
	WakeUp()
}
