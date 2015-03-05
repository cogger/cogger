package cogs

import (
	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//Simple will create a cog that follows the basic context cancellation pattern
func Simple(ctx context.Context, work func() error) cogger.Cog {
	return cogger.NewCog(func() chan error {
		out := make(chan error)

		inner := make(chan error)

		go func() {
			defer close(inner)
			inner <- work()
		}()

		go func() {
			defer close(out)
			select {
			case <-ctx.Done():
				out <- ctx.Err()
			case i := <-inner:
				out <- i
			}
		}()

		return out
	})
}
