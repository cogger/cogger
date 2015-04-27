package wait

import (
	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//Completed runs a cog and shallows all errors and blockers
func Completed(ctx context.Context, cog cogger.Cog, completed func()) {
	go func() {
		for range cog.Do(ctx) {
		}
		completed()

	}()
}
