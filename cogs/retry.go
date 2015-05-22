package cogs

import (
	"errors"
	"math"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
)

//ErrRetry is the error that should be returned when a Retry cog should retry
var ErrRetry = errors.New("retry")

//Retry will continue to retry a cog until it passes, a none ErrRetry error is returned or the context finishes.
func Retry(ctx context.Context, work func() error, max int) cogger.Cog {
	maxAttempts := float64(max)
	return cogger.NewCog(func() chan error {
		out := make(chan error)
		inner := make(chan error)

		go func() {
			defer close(inner)
			attempts := 0.0
			err := ErrRetry
			for err == ErrRetry && ctx.Err() == nil && attempts < maxAttempts {
				attempts++
				err = work()
				if err == ErrRetry {
					time.Sleep(time.Duration(math.Pow(2.0, attempts)) * time.Millisecond)
				}
			}
			inner <- err
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
