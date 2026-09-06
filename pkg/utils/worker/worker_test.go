package worker

import (
	"context"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("function is required", func(t *testing.T) {
		t.Parallel()

		_, err := New()

		require.ErrorIs(t, err, ErrNoFunction)
	})

	t.Run("defaults", func(t *testing.T) {
		t.Parallel()

		calls := make(chan struct{}, 1)

		w, err := New(WithFunction(completedFunc(calls)))

		require.NoError(t, err)
		require.Equal(t, time.Hour, w.errRetry)
		require.Equal(t, slog.Default(), w.logger)
	})

	t.Run("options applied", func(t *testing.T) {
		t.Parallel()

		calls := make(chan struct{}, 1)
		logger := slog.New(slog.DiscardHandler)

		w, err := New(
			WithFunction(completedFunc(calls)),
			WithPeriodic(time.Minute),
			WithErrRetry(time.Second),
			WithLogger(logger),
		)

		require.NoError(t, err)
		require.Equal(t, time.Minute, w.periodic)
		require.Equal(t, time.Second, w.errRetry)
		require.Equal(t, logger, w.logger)
	})

	t.Run("nil logger ignored", func(t *testing.T) {
		t.Parallel()

		calls := make(chan struct{}, 1)

		w, err := New(
			WithFunction(completedFunc(calls)),
			WithLogger(nil),
		)

		require.NoError(t, err)
		require.Equal(t, slog.Default(), w.logger)
	})
}

func TestServeCompleted(t *testing.T) {
	t.Parallel()

	calls := make(chan struct{}, 8)

	w, err := New(WithFunction(completedFunc(calls)))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := serveAsync(w, ctx)

	recvCall(t, calls)

	assertNoCall(t, calls)

	cancel()

	require.NoError(t, <-done)
}

func TestServeRetryImmediate(t *testing.T) {
	t.Parallel()

	calls := make(chan struct{}, 8)

	var count atomic.Int32

	w, err := New(WithFunction(func(context.Context) (WorkerResult, error) {
		calls <- struct{}{}

		switch count.Add(1) {
		case 1:
			// negative RetryAfter is clamped to zero, retry is immediate.
			return WorkerResult{RetryAfter: -time.Second}, nil
		case 2:
			return WorkerResult{}, nil
		default:
			return WorkerResult{IsCompleted: true}, nil
		}
	}))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := serveAsync(w, ctx)

	recvCall(t, calls)
	recvCall(t, calls)
	recvCall(t, calls)

	assertNoCall(t, calls)

	cancel()

	require.NoError(t, <-done)
	require.Equal(t, int32(3), count.Load())
}

func TestServeRetryAfter(t *testing.T) {
	t.Parallel()

	calls := make(chan struct{}, 8)

	var count atomic.Int32

	w, err := New(WithFunction(func(context.Context) (WorkerResult, error) {
		calls <- struct{}{}

		if count.Add(1) == 1 {
			return WorkerResult{RetryAfter: 150 * time.Millisecond}, nil
		}

		return WorkerResult{IsCompleted: true}, nil
	}))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := serveAsync(w, ctx)

	recvCall(t, calls)

	// retry is not due yet.
	assertNoCall(t, calls)

	recvCall(t, calls)

	assertNoCall(t, calls)

	cancel()

	require.NoError(t, <-done)
}

func TestServeErrorRetry(t *testing.T) {
	t.Parallel()

	handler := &recordingHandler{}

	calls := make(chan struct{}, 8)

	var count atomic.Int32

	w, err := New(
		WithFunction(func(context.Context) (WorkerResult, error) {
			calls <- struct{}{}

			if count.Add(1) == 1 {
				return WorkerResult{}, assert.AnError
			}

			return WorkerResult{IsCompleted: true}, nil
		}),
		WithErrRetry(150*time.Millisecond),
		WithLogger(slog.New(handler)),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := serveAsync(w, ctx)

	recvCall(t, calls)

	// error retry is not due yet.
	assertNoCall(t, calls)

	recvCall(t, calls)

	cancel()

	require.NoError(t, <-done)
	require.Equal(t, int32(2), count.Load())

	handler.mu.Lock()
	defer handler.mu.Unlock()

	require.Len(t, handler.records, 1)
	require.Equal(t, slog.LevelError, handler.records[0].Level)
}

func TestServePeriodic(t *testing.T) {
	t.Parallel()

	calls := make(chan struct{}, 8)

	w, err := New(
		WithFunction(completedFunc(calls)),
		WithPeriodic(40*time.Millisecond),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := serveAsync(w, ctx)

	// initial wake plus two periodic ticks.
	recvCall(t, calls)
	recvCall(t, calls)
	recvCall(t, calls)

	cancel()

	require.NoError(t, <-done)
}

func TestWakeUp(t *testing.T) {
	t.Parallel()

	calls := make(chan struct{}, 8)

	w, err := New(WithFunction(completedFunc(calls)))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := serveAsync(w, ctx)

	// initial wake.
	recvCall(t, calls)

	assertNoCall(t, calls)

	w.WakeUp()

	recvCall(t, calls)

	assertNoCall(t, calls)

	w.WakeUp()

	recvCall(t, calls)

	cancel()

	require.NoError(t, <-done)
}

func TestWakeUpCoalesced(t *testing.T) {
	t.Parallel()

	calls := make(chan struct{}, 8)
	release := make(chan struct{})

	w, err := New(WithFunction(func(context.Context) (WorkerResult, error) {
		calls <- struct{}{}

		<-release

		return WorkerResult{IsCompleted: true}, nil
	}))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := serveAsync(w, ctx)

	// first call is in-flight and blocked.
	recvCall(t, calls)

	w.WakeUp()
	w.WakeUp()

	close(release)

	// both wakes are coalesced into a single extra run.
	recvCall(t, calls)

	assertNoCall(t, calls)

	cancel()

	require.NoError(t, <-done)
}

func TestServeState(t *testing.T) {
	t.Parallel()

	calls := make(chan struct{}, 8)
	release := make(chan struct{})

	fn := func(context.Context) (WorkerResult, error) {
		calls <- struct{}{}

		<-release

		return WorkerResult{IsCompleted: true}, nil
	}

	w, err := New(WithFunction(fn))
	require.NoError(t, err)

	ctx := t.Context()

	done := serveAsync(w, ctx)

	// call is in-flight, Serve is running.
	recvCall(t, calls)

	require.ErrorIs(t, w.Serve(ctx), ErrRunning)

	w.Stop()

	require.ErrorIs(t, w.Serve(ctx), ErrStopped)

	// Stop is graceful, in-flight call completes first.
	close(release)

	require.NoError(t, <-done)

	require.ErrorIs(t, w.Serve(ctx), ErrStopped)
}

func TestServeStateStopBeforeServe(t *testing.T) {
	t.Parallel()

	calls := make(chan struct{}, 1)

	w, err := New(WithFunction(completedFunc(calls)))
	require.NoError(t, err)

	w.Stop()

	require.ErrorIs(t, w.Serve(t.Context()), ErrStopped)
}

func TestServeStateStopWhileIdle(t *testing.T) {
	t.Parallel()

	calls := make(chan struct{}, 8)

	w, err := New(WithFunction(completedFunc(calls)))
	require.NoError(t, err)

	ctx := t.Context()

	done := serveAsync(w, ctx)

	recvCall(t, calls)

	w.Stop()

	require.NoError(t, <-done)
}

func TestServeRestart(t *testing.T) {
	t.Parallel()

	calls := make(chan struct{}, 8)

	w, err := New(WithFunction(completedFunc(calls)))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	done := serveAsync(w, ctx)

	recvCall(t, calls)

	cancel()

	require.NoError(t, <-done)

	// worker can be served again after context cancel exit.
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	done = serveAsync(w, ctx)

	recvCall(t, calls)

	cancel()

	require.NoError(t, <-done)
}
