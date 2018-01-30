package wait

import (
	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
)

//NoBlock runs a cog and shallows all errors and blockers
func NoBlock(ctx context.Context, cog cogger.Cog) {
	go func() {
		for _ = range cog.Do(ctx) {
		}
	}()
}
