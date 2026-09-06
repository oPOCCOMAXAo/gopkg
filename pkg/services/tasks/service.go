package tasks

import (
	"cmp"
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/alitto/pond/v2"
	"github.com/opoccomaxao/gopkg/pkg/errors"
	"github.com/opoccomaxao/gopkg/pkg/services/lifecycle"
	pkgerr "github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
)

var (
	_ lifecycle.Servable     = (*Service)(nil)
	_ lifecycle.Shutdownable = (*Service)(nil)
)

type Service struct {
	logger   *slog.Logger
	pool     pond.Pool
	cron     *cron.Cron
	timeout  time.Duration
	maxStack int

	runLock  NameLock
	registry map[string]*registeredTask

	served   bool
	services []ServiceWithTasks
}

//nolint:mnd
func NewService(
	config Config,
	logger *slog.Logger,
	cronLogger *slog.Logger,
) *Service {
	res := Service{
		logger:   logger,
		timeout:  config.Timeout,
		maxStack: lo.Clamp(config.MaxStack, 1, 1000),
		registry: make(map[string]*registeredTask),
	}

	if config.QueueSize == 0 {
		config.QueueSize = pond.Unbounded
	}

	res.pool = pond.NewPool(
		config.WorkerCount,
		pond.WithNonBlocking(false),
		pond.WithQueueSize(config.QueueSize),
	)

	res.cron = cron.New(
		cron.WithLocation(time.UTC),
		cron.WithLogger(WrapCronLogger(cronLogger)),
	)

	return &res
}

func (s *Service) Register(
	svc ...ServiceWithTasks,
) error {
	if s.served {
		return pkgerr.WithStack(errRegisterAfterServe)
	}

	s.services = append(s.services, svc...)

	return nil
}

func (s *Service) Serve(ctx context.Context) error {
	s.served = true

	initial := make([]TaskDescription, 0)
	periodic := make([]TaskDescription, 0)

	dependents := map[string][]string{}

	for _, svc := range s.services {
		tasks := svc.GetTasks(ctx)

		for _, task := range tasks {
			if _, exists := s.registry[task.Name]; exists {
				return pkgerr.Wrapf(errors.ErrDuplicate, "task %s", task.Name)
			}

			s.registry[task.Name] = &registeredTask{
				Name:    task.Name,
				Task:    task.Task,
				Timeout: cmp.Or(task.Timeout, s.timeout),
			}

			for _, dep := range task.After {
				dependents[dep] = append(dependents[dep], task.Name)
			}

			if task.WithInit || task.WithInitAsync {
				initial = append(initial, task)
			}

			if task.Schedule != "" {
				periodic = append(periodic, task)
			}
		}
	}

	for _, task := range s.registry {
		task.Dependents = dependents[task.Name]
	}

	for _, task := range initial {
		if task.WithInitAsync {
			s.RunIgnore(ctx, task.Name)

			continue
		}

		err := s.RunWait(ctx, task.Name)
		if err != nil {
			return err
		}
	}

	s.cron.Start()

	for _, task := range periodic {
		//nolint:contextcheck
		err := s.registerPeriodic(task)
		if err != nil {
			return err
		}
	}

	s.services = nil

	return nil
}

func (s *Service) Shutdown(
	ctx context.Context,
) error {
	done := make(chan struct{})

	go func() {
		s.pool.StopAndWait()

		close(done)
	}()

	select {
	case <-ctx.Done():
		// if we failed to shutdown gracefully within the timeout,
		// we just exit and let the OS clean up our resources.
		return nil
	case <-done:
		return nil
	}
}

func (s *Service) RunIgnore(
	ctx context.Context,
	name string,
) {
	task, err := s.getTaskByName(ctx, name)
	if err != nil {
		return
	}

	s.runIgnore(ctx, task)
}

func (s *Service) RunWait(
	ctx context.Context,
	name string,
) error {
	task, err := s.getTaskByName(ctx, name)
	if err != nil {
		return err
	}

	return s.runWait(ctx, task)
}

func (s *Service) RunAfter(
	ctx context.Context,
	name string,
	duration time.Duration,
) {
	time.AfterFunc(duration, func() {
		s.RunIgnore(ctx, name)
	})
}

func (s *Service) ListTasks() []string {
	res := make([]string, 0, len(s.registry))
	for name := range s.registry {
		res = append(res, name)
	}

	slices.Sort(res)

	return res
}
