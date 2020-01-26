package synced

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, timeout time.Duration, actions ...func() error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	for _, action := range actions {
		g.Go(action)
	}

	c := make(chan error)
	go func() {
		defer close(c)
		c <- g.Wait()
	}()
	select {
	case err := <-c:
		return err
	case <-time.After(timeout):
		return errors.New("timeout happened")
	}
}
