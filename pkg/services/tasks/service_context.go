package tasks

import (
	"context"

	"github.com/opoccomaxao/gopkg/pkg/utils/stack"
)

func (s *Service) ContextWithSource(
	ctx context.Context,
	sourceName string,
) context.Context {
	return stack.PushString(ctx, sourceName, s.maxStack)
}
