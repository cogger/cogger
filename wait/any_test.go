package wait_test

import (
	"errors"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1/cogs"
	. "gopkg.in/cogger/cogger.v1/wait"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Any", func() {
	It("should return on first success", func() {
		ctx := context.Background()

		order := []int{}
		cog := Any(ctx,
			cogs.Simple(ctx, func() error {
				time.Sleep(500 * time.Millisecond)
				order = append(order, 4)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(500 * time.Millisecond)
				order = append(order, 3)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(500 * time.Millisecond)
				order = append(order, 2)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(500 * time.Millisecond)
				order = append(order, 1)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				order = append(order, 0)
				return nil
			}),
		)

		Expect(<-cog.Do(ctx)).ToNot(HaveOccurred())
		Expect(order).To(HaveLen(1))
	})

	It("should tell us when none pass", func() {
		ctx := context.Background()

		fakeErr := errors.New("fake error")
		cog := Any(ctx,
			cogs.Simple(ctx, func() error {
				return fakeErr
			}),
			cogs.Simple(ctx, func() error {

				return fakeErr
			}),
			cogs.Simple(ctx, func() error {

				return fakeErr
			}),
			cogs.Simple(ctx, func() error {

				return fakeErr
			}),
			cogs.Simple(ctx, func() error {

				return fakeErr
			}),
		)

		Expect(<-cog.Do(ctx)).To(Equal(ErrNonePassed))
	})

	It("should exit if context is canceled before completion", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		order := []int{}
		cog := Any(ctx,
			cogs.Simple(ctx, func() error {
				time.Sleep(500 * time.Millisecond)
				order = append(order, 4)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(500 * time.Millisecond)
				order = append(order, 3)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(500 * time.Millisecond)
				order = append(order, 2)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(500 * time.Millisecond)
				order = append(order, 1)
				return nil
			}),
			cogs.Simple(ctx, func() error {
				time.Sleep(500 * time.Millisecond)
				order = append(order, 0)
				return nil
			}),
		)

		Expect(<-cog.Do(ctx)).To(Equal(context.DeadlineExceeded))
		Expect(order).To(HaveLen(0))
	})
})
