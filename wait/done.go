package wait

import (
	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//Done runs a cog and shallows all errors and blockers
func Done(ctx context.Context, cog cogger.Cog,done func()) {
	go func() {
		for range cog.Do(ctx) {}
		done()

	}()
}
