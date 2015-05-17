package cogs

import "github.com/cogger/cogger"

func ReturnErr(err error) cogger.Cog {
	return cogger.NewCog(func() chan error {
		out := make(chan error, 1)
		out <- err
		return out
	})
}
