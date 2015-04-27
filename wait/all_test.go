package wait_test

import (
	"errors"

	. "github.com/cogger/cogger/wait"

	"time"

	"github.com/cogger/cogger/cogs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("All", func() {
	It("should return all cogs when none fail", func() {
		ctx := context.Background()

		order := []int{}
		cog := All(ctx,
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

	It("should return on first failure", func() {
		ctx := context.Background()

		fakeErr := errors.New("fake error")
		order := []int{}
		cog := All(ctx,
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
				return fakeErr
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

		Expect(<-cog.Do(ctx)).To(Equal(fakeErr))

		for i, o := range order {
			Expect(o).To(Equal(i))
		}
		Expect(order).To(HaveLen(3))
	})

	It("should exit if context is canceled before completion", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		order := []int{}
		cog := All(ctx,
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
