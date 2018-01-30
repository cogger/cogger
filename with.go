package cogger

import (
	"time"

	"golang.org/x/net/context"
)

//With creates a function wrapper that handles all cog tear down
func With(f func(context.Context)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	f(ctx)
}

//WithTimeout creates a function wrapper that handles all cog tear down and has a defined timeout.
func WithTimeout(f func(context.Context), duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	f(ctx)
}

//WithDeadline creats a function wrapper that handles all cog tear down and has defined deadlime.
func WithDeadline(f func(context.Context), deadline time.Time) {
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	f(ctx)
}
