package cogs

import "github.com/cogger/cogger"

//NoOp creates a cog that instantly returns
func NoOp() cogger.Cog {
	return cogger.NewCog(func() chan error {
		return nil
	})
}
