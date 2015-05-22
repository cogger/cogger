package cogs

import "gopkg.in/cogger/cogger.v1"

//NoOp creates a cog that instantly returns
func NoOp() cogger.Cog {
	return cogger.NewCog(func() chan error {
		out := make(chan error)
		close(out)
		return out
	})
}
