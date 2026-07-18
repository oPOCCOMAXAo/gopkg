package tasks

import (
	"context"
	"log/slog"
	"time"

	"github.com/opoccomaxao/gopkg/pkg/errors"
	"github.com/opoccomaxao/gopkg/pkg/services/tasks/stack"
	"github.com/opoccomaxao/gopkg/pkg/utils/ctxutils"
	pkgerr "github.com/pkg/errors"
)

func (s *Service) registerPeriodic(task TaskDescription) error {
	_, err := s.cron.AddFunc(task.Schedule, func() {
		ctx := context.Background()
		ctx = stack.Push(ctx, "cron", s.maxStack)

		s.RunIgnore(ctx, task.Name)
	})
	if err != nil {
		return pkgerr.WithStack(err)
	}

	return nil
}

func (s *Service) runTask(
	ctx context.Context,
	task *registeredTask,
) error {
	logger := s.logger.With(
		slog.String("name", task.Name),
	)

	{
		prevStack := stack.Get(ctx)
		if len(prevStack) > 0 {
			logger = logger.With(slog.Any("stack", prevStack))
		}
	}

	if !s.runLock.TryLock(task.Name) {
		logger.InfoContext(ctx, "task skipped (already running)")

		return nil
	}
	defer s.runLock.Unlock(task.Name)

	logger.InfoContext(ctx, "task start")

	start := time.Now()

	ctx = stack.Push(ctx, task.Name, s.maxStack)

	ctx, cancel := context.WithTimeout(ctx, task.Timeout)
	defer cancel()

	err := task.Task(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "task error",
			slog.Any("error", err),
		)

		return pkgerr.WithStack(err)
	}

	duration := time.Since(start)

	logger.InfoContext(ctx, "task complete",
		slog.Duration("duration", duration),
	)

	childCtx := ctxutils.DetachedContext(ctx)

	for _, depName := range task.Dependents {
		depTask := s.registry[depName]
		s.runIgnore(childCtx, depTask)
	}

	return nil
}

func (s *Service) runWait(
	ctx context.Context,
	task *registeredTask,
) error {
	err := s.pool.
		SubmitErr(func() error {
			return s.runTask(ctxutils.DetachedContext(ctx), task)
		}).
		Wait()
	if err != nil {
		return pkgerr.WithStack(err)
	}

	return nil
}

func (s *Service) runIgnore(
	ctx context.Context,
	task *registeredTask,
) {
	s.pool.SubmitErr(func() error {
		return s.runTask(ctxutils.DetachedContext(ctx), task)
	})
}

func (s *Service) getTaskByName(
	ctx context.Context,
	name string,
) (*registeredTask, error) {
	task, ok := s.registry[name]
	if !ok {
		s.logger.ErrorContext(ctx, "task not found",
			slog.String("name", name),
		)

		return nil, pkgerr.Wrapf(errors.ErrNotFound, "task %s", name)
	}

	return task, nil
}
