package wait

import (
	"sync"

	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//All will execute all cogs provided in parallel and will return when the first one fails or all succeed
func All(ctx context.Context, cogs ...cogger.Cog) cogger.Cog {
	return cogger.NewCog(func() chan error {
		out := make(chan error, 1)

		wg := &sync.WaitGroup{}
		wg.Add(len(cogs))

		for _, cog := range cogs {
			go func(cog cogger.Cog) {
				defer wg.Done()
				select {
				case <-ctx.Done():
					out <- ctx.Err()
				case err := <-cog.Do(ctx):
					if err != nil {
						out <- err
					}
				}
			}(cog)
		}

		go func() {
			defer close(out)
			wg.Wait()
			select {
			case <-ctx.Done():
			case out <- nil:
			}

		}()

		return out
	})
}
