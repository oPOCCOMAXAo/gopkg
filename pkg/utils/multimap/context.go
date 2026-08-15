package multimap

import (
	"context"
)

type ctxMMKey struct{}

func FromContext(ctx context.Context) (*MultiMap, bool) {
	res, ok := ctx.Value(ctxMMKey{}).(*MultiMap)

	return res, ok
}

func ContextWith(ctx context.Context, multiMap *MultiMap) context.Context {
	return context.WithValue(ctx, ctxMMKey{}, multiMap)
}
