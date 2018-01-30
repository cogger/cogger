package cogs

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
	"gopkg.in/cogger/cogger.v1/order"
	"gopkg.in/cogger/cogger.v1/wait"
)

var _ = Describe("DeferredCreate", func() {
	coggerInterface := reflect.TypeOf((*cogger.Cog)(nil)).Elem()
	It("should be a cog", func() {
		cog := DeferredCreate(func() cogger.Cog {
			return NoOp()
		})

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).ToNot(HaveOccurred())
	})

	It("should not equal the cog it is going to return", func() {
		cog := NoOp()
		dc := DeferredCreate(func() cogger.Cog {
			return cog
		})

		Expect(dc).ToNot(Equal(cog))
	})

	It("should wait to create a cog", func() {
		ctx := context.Background()
		count := 0
		errs := wait.Resolve(ctx,
			order.Series(ctx,
				DeferredCreate(func() cogger.Cog {
					Expect(count).To(Equal(0))
					return NoOp()
				}),
				Simple(ctx, func() error {
					count++
					return nil
				}),
				Simple(ctx, func() error {
					count++
					return nil
				}),
				Simple(ctx, func() error {
					count++
					return nil
				}),
				DeferredCreate(func() cogger.Cog {
					Expect(count).To(Equal(3))
					return NoOp()
				}),
			),
		)

		Expect(errs).To(HaveLen(0))
	})

	It("should implement SetLimit function", func() {
		cog := DeferredCreate(func() cogger.Cog {
			return NoOp()
		})

		limit := &mockLimit{}

		cog.SetLimit(limit)

		ctx := context.Background()
		Expect(<-cog.Do(ctx)).ToNot(HaveOccurred())
		Expect(limit.NextHits).To(Equal(1))
		Expect(limit.Completed).To(BeTrue())
	})
})

type mockLimit struct {
	Completed bool
	NextHits  int
}

func (limit *mockLimit) Next(ctx context.Context) chan struct{} {
	next := make(chan struct{})
	go func() {
		limit.NextHits++
		next <- struct{}{}
	}()
	return next
}

func (limit *mockLimit) Done(ctx context.Context) {
	limit.Completed = true
}
