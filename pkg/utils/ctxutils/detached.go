package ctxutils

import (
	"context"
)

//nolint:containedctx
type detachedContext struct {
	context.Context

	orig context.Context
}

func (c *detachedContext) Value(key any) any {
	return c.orig.Value(key)
}

func DetachedContext(ctx context.Context) context.Context {
	return &detachedContext{Context: context.Background(), orig: ctx}
}
