package tasks

import (
	"context"
	"time"
)

type Task func(context.Context) error

type TaskDescription struct {
	// Name for logging and tracing.
	Name string

	Task Task

	// Schedule optional.
	// If specified, the task will be run periodically according to the schedule.
	//
	// Format see here: https://github.com/robfig/cron
	Schedule string

	// Timeout optional.
	// If specified, the task will be canceled if it runs longer than the timeout.
	Timeout time.Duration

	// WithInit optional.
	// If specified, the task will be run on startup before any scheduled tasks.
	WithInit bool

	// WithInitAsync optional.
	// If specified, the task will be run on startup asynchronously without waiting for its completion.
	// If used together with WithInit, WithInitAsync takes precedence.
	WithInitAsync bool

	// After optional.
	// Task is automatically triggered after successful execution
	// of any of specified tasks.
	After []string
}

type ServiceWithTasks interface {
	GetTasks(ctx context.Context) []TaskDescription
}
