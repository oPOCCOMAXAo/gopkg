package worker

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"
)

const (
	recvTimeout = 2 * time.Second
	quietPeriod = 50 * time.Millisecond
)

type recordingHandler struct {
	mu      sync.Mutex
	records []slog.Record
}

func (h *recordingHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *recordingHandler) Handle(_ context.Context, record slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.records = append(h.records, record)

	return nil
}

func (h *recordingHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *recordingHandler) WithGroup(string) slog.Handler {
	return h
}

func recvCall(t *testing.T, calls <-chan struct{}) {
	t.Helper()

	select {
	case <-calls:
	case <-time.After(recvTimeout):
		t.Fatal("timed out waiting for function call")
	}
}

func assertNoCall(t *testing.T, calls <-chan struct{}) {
	t.Helper()

	select {
	case <-calls:
		t.Fatal("unexpected function call")
	case <-time.After(quietPeriod):
	}
}

func completedFunc(calls chan<- struct{}) WorkerFunc {
	return func(context.Context) (WorkerResult, error) {
		calls <- struct{}{}

		return WorkerResult{IsCompleted: true}, nil
	}
}

func serveAsync(w *Worker, ctx context.Context) chan error {
	done := make(chan error, 1)

	go func() {
		done <- w.Serve(ctx)
	}()

	return done
}
