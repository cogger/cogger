package order

import (
	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//Series will execute the workers in order.
//It will wait for the previous to finish before starting the next.
func Series(ctx context.Context, cogs ...cogger.Cog) cogger.Cog {
	return cogger.NewCog(func() chan error {
		out := make(chan error, len(cogs))
		go func() {
			for _, cog := range cogs {
				select {
				case <-ctx.Done():
					out <- ctx.Err()
				case err := <-cog.Do(ctx):
					out <- err
				}
			}
			close(out)
		}()
		return out
	})
}
