package cogs

import (
	"net/http"

	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

//Request is a special cog for executing http.Request
func Request(ctx context.Context, client *http.Client, req *http.Request, f func(*http.Response, error) error) cogger.Cog {
	return cogger.NewCog(func() chan error {
		out := make(chan error, 1)
		c := make(chan error, 1)

		go func() { c <- f(client.Do(req)) }()

		go func() {
			select {
			case <-ctx.Done():
				out <- ctx.Err()
			case err := <-c:
				out <- err
			}
		}()

		return out
	})
}
