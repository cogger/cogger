package wait

import (
	"sync"

	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//Settle will execute all cogs in parallel and return in cog order all the states of the code when all are finished.
func Settle(ctx context.Context, cogs ...cogger.Cog) cogger.Cog {
	return cogger.NewCog(func() chan error {
		out := make(chan error, len(cogs))
		results := make([]error, len(cogs))

		wg := &sync.WaitGroup{}
		wg.Add(len(cogs))

		for i, cog := range cogs {
			go func(cog cogger.Cog, i int) {
				defer wg.Done()
				select {
				case <-ctx.Done():
					out <- ctx.Err()
				case err := <-cog.Do(ctx):
					results[i] = err
				}
			}(cog, i)
		}

		go func() {
			defer close(out)
			wg.Wait()

			for _, result := range results {
				select {
				case <-ctx.Done():
				case out <- result:
				}
			}

		}()

		return out
	})
}
