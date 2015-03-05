package limiter

import (
	"errors"

	"golang.org/x/net/context"
)

//Limit is an interface that defines how cog limiters can be accessed.
//A cog with a limit will wait until Next() allows forward movement.
//Done will be called when a cog is finished working.
//Using a Limit on a cog might create deadlocks in the system.
type Limit interface {
	Next(ctx context.Context) chan struct{}
	Done(ctx context.Context)
}

//ErrLimitConfig is is the error for a limti beign configured incorrectly.
var ErrLimitConfig = errors.New("limit was not configured correctly")
