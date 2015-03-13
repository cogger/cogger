package wait

import (
	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//Resolve waits on a cog until all errors are returned
func Resolve(ctx context.Context, cog cogger.Cog) []error {
	errs := []error{}
	for err := range cog.Do(ctx) {
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
