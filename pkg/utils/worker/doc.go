// Package worker provides a single-goroutine background worker
// that executes a function with retry and periodic scheduling support.
//
// The function is executed once immediately after [Worker.Serve] starts,
// then re-executed according to its [WorkerResult] (see [WorkerFunc]).
// It can also be woken up externally via [Worker.WakeUp]
// and scheduled periodically via [WithPeriodic].
//
// Serve blocks until the context is canceled or [Worker.Stop] is called.
// Panics in the function are not recovered.
package worker
