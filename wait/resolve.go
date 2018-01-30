package wait

import (
	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
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
