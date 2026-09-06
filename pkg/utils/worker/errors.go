package worker

import "errors"

var (
	// ErrNoFunction is returned by [New] when the function is not set.
	ErrNoFunction = errors.New("worker function is not set")

	// ErrRunning is returned by [Worker.Serve] when the worker is already running.
	ErrRunning = errors.New("worker is already running")

	// ErrStopped is returned by [Worker.Serve] after [Worker.Stop].
	// A stopped worker cannot be started again.
	ErrStopped = errors.New("worker is stopped")
)
