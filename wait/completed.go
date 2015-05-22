package wait

import (
	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
)

//Completed runs a cog and shallows all errors and blockers
func Completed(ctx context.Context, cog cogger.Cog, completed func()) {
	go func() {
		defer completed()
		Resolve(ctx, cog)
	}()
}
