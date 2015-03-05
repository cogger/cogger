package limiter

import (
	"time"

	"golang.org/x/net/context"
)

type timeLimit struct {
	tick <-chan time.Time
}

func (limit timeLimit) Next(ctx context.Context) chan struct{} {
	out := make(chan struct{})

	go func() {
		select {
		case <-ctx.Done():
		case <-limit.tick:
			out <- struct{}{}
		}
	}()
	return out
}

func (limit timeLimit) Done(ctx context.Context) {}

//PerSecond creates a cog limiter that limits the number of working cogs to the provided number per second.
//This limiter does not ensure that the previous functions have completed before another is started but only that
//X functions will be started per second
func PerSecond(count int) (Limit, error) {
	if count <= 0 {
		return nil, ErrLimitConfig
	}
	return timeLimit{
		tick: time.Tick(1e9 / time.Duration(count)),
	}, nil
}
