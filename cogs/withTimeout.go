package cogs

import (
	"errors"
	"time"

	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//ErrTimeout is returned when a WithTimeout cog times out
var ErrTimeout = errors.New("timed out")

//WithTimeout will place a timeout on this specific cog.
//The cog will execute and return a success, a failur, a local failure or a global context error which ever comes first.
func WithTimeout(ctx context.Context, work func() error, timeout time.Duration) cogger.Cog {
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
			case <-time.After(timeout):
				out <- ErrTimeout
			}
		}()

		return out
	})
}
