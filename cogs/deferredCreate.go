package cogs

import (
	"sync"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
	"gopkg.in/cogger/cogger.v1/limiter"
)

func DeferredCreate(create func() cogger.Cog) cogger.Cog {
	return &deferredCog{
		create: create,
		once:   &sync.Once{},
		cog:    NoOp(),
	}
}

type deferredCog struct {
	create func() cogger.Cog
	once   *sync.Once
	cog    cogger.Cog
	limit  limiter.Limit
}

func (dc *deferredCog) Do(ctx context.Context) chan error {
	dc.once.Do(func() {
		dc.cog = dc.create()
		dc.cog.SetLimit(dc.limit)
	})
	return dc.cog.Do(ctx)
}

func (dc *deferredCog) SetLimit(limit limiter.Limit) cogger.Cog {
	dc.limit = limit
	dc.cog.SetLimit(limit)
	return dc
}
