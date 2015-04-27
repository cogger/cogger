package order_test

import (
	"errors"
	"time"

	"github.com/cogger/cogger/cogs"
	. "github.com/cogger/cogger/order"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Series", func() {
	It("should execute cogs in series", func() {
		ctx := context.Background()

		order := []int{}
		cog := Series(ctx,
			cogs.Simple(ctx, func() error {
				time.Sleep(100 * time.Millisecond)
				order = append(order, 0)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(75 * time.Millisecond)
				order = append(order, 1)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(50 * time.Millisecond)
				order = append(order, 2)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(25 * time.Millisecond)
				order = append(order, 3)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				order = append(order, 4)
				return nil
			}),
		)

		Expect(<-cog.Do(ctx)).ToNot(HaveOccurred())
		for i, o := range order {
			Expect(o).To(Equal(i))
		}
		Expect(order).To(HaveLen(5))
	})

	It("should return all errors", func() {
		ctx := context.Background()

		ErrFake := errors.New("fake error")

		order := []int{}
		cog := Series(ctx,
			cogs.Simple(ctx, func() error {
				time.Sleep(100 * time.Millisecond)
				order = append(order, 0)
				return ErrFake
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(75 * time.Millisecond)
				order = append(order, 1)
				return ErrFake
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(50 * time.Millisecond)
				order = append(order, 2)
				return ErrFake
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(25 * time.Millisecond)
				order = append(order, 3)
				return ErrFake
			}),
			cogs.Simple(ctx, func() error {
				order = append(order, 4)
				return ErrFake
			}),
		)

		for err := range cog.Do(ctx) {
			Expect(err).To(Equal(ErrFake))
		}

		for i, o := range order {
			Expect(o).To(Equal(i))
		}
		Expect(order).To(HaveLen(5))
	})

	It("should exit if context is canceled before completion", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		order := []int{}
		cog := Series(ctx,
			cogs.Simple(ctx, func() error {
				order = append(order, 0)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(1 * time.Second)
				order = append(order, 1)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				order = append(order, 2)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				order = append(order, 3)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				order = append(order, 4)
				return nil
			}),
		)

		for range cog.Do(ctx) {
		}

		for i, o := range order {
			Expect(o).To(Equal(i))
		}
		Expect(order).To(HaveLen(1))
	})
})
