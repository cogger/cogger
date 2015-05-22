package cogger

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("With", func() {
	It("should on exit it should cancel the context", func() {
		var c context.Context
		With(func(ctx context.Context) {
			c = ctx
		})
		Expect(c.Err()).To(Equal(context.Canceled))
	})
})

var _ = Describe("WithTimeout", func() {
	It("should on exit it should cancel the context", func() {
		var c context.Context
		WithTimeout(func(ctx context.Context) {
			c = ctx
		}, 1*time.Second)
		Expect(c.Err()).To(Equal(context.Canceled))
	})

	It("should timeout the context", func() {
		var c context.Context
		WithTimeout(func(ctx context.Context) {
			c = ctx
			time.Sleep(4 * time.Second)
		}, 1*time.Second)

		Expect(c.Err()).To(Equal(context.DeadlineExceeded))
	})
})

var _ = Describe("WithDeadline", func() {
	It("should on exit it should cancel the context", func() {
		var c context.Context
		WithDeadline(func(ctx context.Context) {
			c = ctx
		}, time.Now().Add(5*time.Second))
		Expect(c.Err()).To(Equal(context.Canceled))
	})

	It("should timeout the context", func() {
		var c context.Context
		WithDeadline(func(ctx context.Context) {
			c = ctx
			time.Sleep(4 * time.Second)
		}, time.Now().Add(2*time.Second))

		Expect(c.Err()).To(Equal(context.DeadlineExceeded))
	})
})
