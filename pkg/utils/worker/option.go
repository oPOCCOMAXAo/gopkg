package worker

import (
	"log/slog"
	"time"
)

// WorkerOption configures [Worker] on creation.
type WorkerOption func(*Worker)

// WithPeriodic sets the interval for periodic execution.
// Non-positive value (default) disables it.
func WithPeriodic(periodic time.Duration) WorkerOption {
	return func(w *Worker) {
		w.periodic = periodic
	}
}

// WithFunction sets the function executed by the worker.
// Required — [New] returns [ErrNoFunction] without it.
func WithFunction(function WorkerFunc) WorkerOption {
	return func(w *Worker) {
		w.function = function
	}
}

// WithLogger sets the logger used to report function errors.
// Nil value is ignored, default is slog.Default.
func WithLogger(logger *slog.Logger) WorkerOption {
	return func(w *Worker) {
		if logger == nil {
			return
		}

		w.logger = logger
	}
}

// WithErrRetry sets the delay before the function is retried after an error.
// Default is one hour.
// Non-positive value results in immediate retry in a tight loop.
func WithErrRetry(errRetry time.Duration) WorkerOption {
	return func(w *Worker) {
		w.errRetry = errRetry
	}
}
