package tasks

import (
	"context"

	"github.com/opoccomaxao/gopkg/pkg/services/tasks/stack"
)

func (s *Service) ContextWithSource(
	ctx context.Context,
	sourceName string,
) context.Context {
	return stack.Push(ctx, sourceName, s.maxStack)
}
