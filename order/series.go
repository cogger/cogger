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
			defer close(out)
			for _, cog := range cogs {
				if ctx.Err() == nil {
					for err := range cog.Do(ctx) {
						if ctx.Err() != nil {
							out <- err
							break
						}
						if err != nil {
							out <- err
						}
					}
				}
			}
		}()
		return out
	})
}
