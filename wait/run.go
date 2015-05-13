package wait

import (
	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//NoBlock runs a cog and shallows all errors and blockers
func NoBlock(ctx context.Context, cog cogger.Cog) {
	go func() {
		for _ = range cog.Do(ctx) {
		}
	}()
}
