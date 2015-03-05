package order

import (
	"sync"

	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//Parallel will execute all functions is parallel that can be executed.
func Parallel(ctx context.Context, cogs ...cogger.Cog) cogger.Cog {
	return cogger.NewCog(func() chan error {
		out := make(chan error, len(cogs))
		wg := &sync.WaitGroup{}
		wg.Add(len(cogs))

		go func() {
			for _, cog := range cogs {
				go func(cog cogger.Cog) {
					defer wg.Done()
					select {
					case <-ctx.Done():
						out <- ctx.Err()
					case err := <-cog.Do(ctx):
						out <- err
					}
				}(cog)
			}
		}()

		go func() {
			defer close(out)
			wg.Wait()
		}()
		return out
	})
}
