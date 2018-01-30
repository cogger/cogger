package limiter

import "golang.org/x/net/context"

type countLimit struct {
	ticks chan struct{}
}

func (limit countLimit) Next(ctx context.Context) chan struct{} {
	out := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
		case tick := <-limit.ticks:
			out <- tick
		}
	}()
	return out
}

func (limit countLimit) Done(ctx context.Context) {
	select {
	case <-ctx.Done():
	case limit.ticks <- struct{}{}:
	}
}

//ByCount creates a cog limiter that limits the number of working cogs to the ammount provided to this function.
func ByCount(alive int) (Limit, error) {
	if alive <= 0 {
		return nil, ErrLimitConfig
	}

	limit := countLimit{
		ticks: make(chan struct{}, alive),
	}

	for i := 0; i < alive; i++ {
		limit.ticks <- struct{}{}
	}

	return limit, nil

}
