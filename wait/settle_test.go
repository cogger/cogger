package wait

import (
	"time"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1/cogs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Settle", func() {
	It("should settle all cogs", func() {
		ctx := context.Background()

		order := []int{}
		cog := Settle(ctx,
			cogs.Simple(ctx, func() error {
				time.Sleep(100 * time.Millisecond)
				order = append(order, 4)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(75 * time.Millisecond)
				order = append(order, 3)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(50 * time.Millisecond)
				order = append(order, 2)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(25 * time.Millisecond)
				order = append(order, 1)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				order = append(order, 0)
				return nil
			}),
		)

		for err := range cog.Do(ctx) {
			Expect(err).To(BeNil())
		}

		for i, o := range order {
			Expect(o).To(Equal(i))
		}
		Expect(order).To(HaveLen(5))
	})

	It("should exit if context is canceled before completion", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		cog := Settle(ctx,
			cogs.Simple(ctx, func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				return nil
			}),
			cogs.Simple(ctx, func() error {
				return nil
			}),
		)

		count := 0
		for err := range cog.Do(ctx) {
			if count < 3 {
				Expect(err).To(Equal(context.DeadlineExceeded))
			} else {
				Expect(err).To(BeNil())
			}
			count++
		}
	})
})
