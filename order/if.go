package order

import (
	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//If will call a function to determine if it should run or not
func If(ctx context.Context, shouldRun func() bool, cog cogger.Cog) cogger.Cog {
	return cogger.NewCog(func() chan error {
		out := make(chan error, 1)

		if shouldRun() {
			go func(cog cogger.Cog) {
				defer close(out)
				select {
				case <-ctx.Done():
					out <- ctx.Err()
				case err := <-cog.Do(ctx):
					out <- err
				}
			}(cog)
		} else {
			close(out)
		}

		return out
	})
}
