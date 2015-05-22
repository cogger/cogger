package wait

import (
	"errors"
	"sync"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
)

//ErrNonePassed is the error for when no cogs finish without an error in an Any call
var ErrNonePassed = errors.New("no worker finished without an error")

//Any returns the first cog to finish successfully or ErrNonePassed when all fail
func Any(ctx context.Context, cogs ...cogger.Cog) cogger.Cog {
	return cogger.NewCog(func() chan error {
		out := make(chan error, 1)
		first := make(chan struct{})

		wg := &sync.WaitGroup{}
		wg.Add(len(cogs))

		for _, cog := range cogs {
			go func(cog cogger.Cog) {
				defer wg.Done()
				select {
				case <-ctx.Done():
					out <- ctx.Err()
				case err := <-cog.Do(ctx):
					if err == nil && ctx.Err() == nil {
						out <- nil
						close(first)
					}
				case <-first:
				}
			}(cog)
		}

		go func() {
			defer close(out)
			wg.Wait()
			select {
			case <-first:
			default:
				select {
				case <-ctx.Done():
				case out <- ErrNonePassed:
				}
			}
		}()

		return out
	})
}
