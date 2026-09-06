package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	pkgerr "github.com/pkg/errors"
)

// Ensure [Worker] implements [Interface].
var _ Interface = (*Worker)(nil)

type workerState int

const (
	stateCreated workerState = iota
	stateRunning
	stateStopped
)

type Worker struct {
	mu     sync.Mutex
	state  workerState
	stopCh chan struct{}

	retryTimer *time.Timer
	retryCh    <-chan time.Time

	// periodic is the interval for periodic execution.
	// Non-positive value (default) disables it.
	periodic time.Duration

	// function is the [WorkerFunc] executed by the worker.
	function WorkerFunc

	logger *slog.Logger

	wakeCh chan struct{}

	// errRetry is the delay before the function is retried after an error.
	errRetry time.Duration
}

// New creates a new [Worker] with provided options.
// Returns [ErrNoFunction] when [WithFunction] is not provided or is nil.
func New(options ...WorkerOption) (*Worker, error) {
	res := &Worker{
		wakeCh: make(chan struct{}, 1),
		stopCh: make(chan struct{}),

		logger: slog.Default(),

		errRetry: time.Hour,
	}

	for _, option := range options {
		option(res)
	}

	if res.function == nil {
		return nil, pkgerr.WithStack(ErrNoFunction)
	}

	return res, nil
}

// Serve starts the worker and blocks
// until the context is canceled or [Worker.Stop] is called.
// The function is executed once immediately after start;
// an in-flight call is not interrupted by shutdown.
// Returns [ErrRunning] when already running,
// [ErrStopped] when called after [Worker.Stop].
func (w *Worker) Serve(ctx context.Context) error {
	err := w.startServe()
	if err != nil {
		return err
	}

	defer w.finishServe()

	var periodic <-chan time.Time

	if w.periodic > 0 {
		ticker := time.NewTicker(w.periodic)
		defer ticker.Stop()

		periodic = ticker.C
	}

	w.WakeUp()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-w.stopCh:
			return nil

		case <-periodic:
			w.WakeUp()

		case <-w.wakeCh:
			w.run(ctx)

		case <-w.retryCh:
			w.run(ctx)
		}
	}
}

// WakeUp schedules the function for execution.
// If it is already running, it is executed again after completion.
func (w *Worker) WakeUp() {
	select {
	case w.wakeCh <- struct{}{}:
	default:
	}
}

// Stop shuts down the worker goroutine
// without waiting for an in-flight function call to complete.
func (w *Worker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state == stateStopped {
		return
	}

	w.state = stateStopped

	close(w.stopCh)
}

// run executes the function once and schedules a retry according to the result.
func (w *Worker) run(ctx context.Context) {
	w.stopRetry()

	res, err := w.function(ctx)
	if err != nil {
		w.logger.ErrorContext(ctx, "worker run failed",
			slog.Any("error", err),
		)

		w.startRetry(w.errRetry)

		return
	}

	if res.IsCompleted {
		return
	}

	w.startRetry(max(res.RetryAfter, 0))
}

// startServe transitions the worker to the running state.
func (w *Worker) startServe() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state == stateStopped {
		return pkgerr.WithStack(ErrStopped)
	}

	if w.state == stateRunning {
		return pkgerr.WithStack(ErrRunning)
	}

	w.state = stateRunning

	return nil
}

// finishServe stops the pending retry and transitions the worker
// back to the created state unless it is stopped.
func (w *Worker) finishServe() {
	w.stopRetry()

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state != stateStopped {
		w.state = stateCreated
	}
}

// stopRetry stops the pending retry timer, if any.
func (w *Worker) stopRetry() {
	if w.retryTimer == nil {
		return
	}

	w.retryTimer.Stop()

	select {
	case <-w.retryTimer.C:
	default:
	}

	w.retryTimer = nil
	w.retryCh = nil
}

// startRetry schedules a retry after the provided duration.
func (w *Worker) startRetry(after time.Duration) {
	w.retryTimer = time.NewTimer(after)
	w.retryCh = w.retryTimer.C
}
