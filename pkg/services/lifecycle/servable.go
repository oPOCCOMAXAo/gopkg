package lifecycle

import (
	"context"

	pkgerr "github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func ServeAll(
	ctx context.Context,
	servables ...Servable,
) error {
	var eg errgroup.Group

	for _, srv := range servables {
		eg.Go(func() error {
			return srv.Serve(ctx)
		})
	}

	return pkgerr.WithStack(eg.Wait())
}

type ServableFunc func(context.Context) error

func (s ServableFunc) Serve(ctx context.Context) error {
	return s(ctx)
}
