package cogs

import "gopkg.in/cogger/cogger.v1"

//ReturnErr creates a cog that instantly returns the error given to it
func ReturnErr(err error) cogger.Cog {
	return cogger.NewCog(func() chan error {
		out := make(chan error, 1)
		out <- err
		return out
	})
}
