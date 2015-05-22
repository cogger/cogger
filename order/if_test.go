package order_test

import (
	"time"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1/cogs"
	. "gopkg.in/cogger/cogger.v1/order"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("If", func() {
	It("should execute the cog if the function returns true", func() {
		ctx := context.Background()
		ran := false
		Expect(<-If(ctx, func() bool { return true }, cogs.Simple(ctx, func() error {
			ran = true
			return nil
		})).Do(ctx)).To(BeNil())
		Expect(ran).To(BeTrue())
	})
	It("should not execute the cog if the function returns false", func() {
		ctx := context.Background()
		ran := false
		Expect(<-If(ctx, func() bool { return false }, cogs.Simple(ctx, func() error {
			ran = true
			return nil
		})).Do(ctx)).To(BeNil())
		Expect(ran).To(BeFalse())
	})

	It("should exit if context is canceled before completion", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		ran := false
		Expect(<-If(ctx, func() bool { return true }, cogs.Simple(ctx, func() error {
			time.Sleep(1 * time.Second)
			ran = true
			return nil
		})).Do(ctx)).To(Equal(context.DeadlineExceeded))
		Expect(ran).To(BeFalse())
	})

	It("should not execute the check until the cog is run", func() {
		ctx := context.Background()
		check := false
		Series(ctx,
			cogs.Simple(ctx, func() error {
				check = true
				return nil
			}),
			If(ctx,
				func() bool {
					Expect(check).To(BeTrue())
					check = false
					return true
				},
				cogs.NoOp(),
			),
			If(ctx,
				func() bool {
					Expect(check).To(BeFalse())
					return true
				},
				cogs.NoOp(),
			),
		)
	})
})
