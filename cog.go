package cogger

import (
	"gopkg.in/cogger/cogger.v1/limiter"
	// "gopkg.in/cogger/cogger.v1/limiter"
	"golang.org/x/net/context"
)

//Cog is the interface for a cog
type Cog interface {
	Do(ctx context.Context) chan error
	SetLimit(limiter.Limit) Cog
}

type defaultCog struct {
	f     func() chan error
	limit limiter.Limit
}

func (cog *defaultCog) Do(ctx context.Context) chan error {
	if cog.limit != nil {
		out := make(chan error)
		go func() {
			defer cog.limit.Done(ctx)
			defer close(out)
			<-cog.limit.Next(ctx)
			out <- <-cog.f()
		}()
		return out
	}
	return cog.f()
}

func (cog *defaultCog) SetLimit(limit limiter.Limit) Cog {
	cog.limit = limit
	return cog
}

//NewCog creats a default cog from the provided function
//This is a lower level function and generally a predefined work cog should be used.
func NewCog(f func() chan error) Cog {
	return &defaultCog{
		f: f,
	}
}
