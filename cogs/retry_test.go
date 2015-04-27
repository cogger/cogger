package cogs_test

import (
	"reflect"
	"time"

	"github.com/cogger/cogger"
	. "github.com/cogger/cogger/cogs"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Retry", func() {
	coggerInterface := reflect.TypeOf((*cogger.Cog)(nil)).Elem()

	It("should create a retry Cog", func() {
		ctx := context.Background()
		cog := Retry(ctx, func() error { return nil }, 1)
		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).ToNot(HaveOccurred())
	})

	It("should retry when ErrRetry is returned from the worker", func() {
		ctx := context.Background()

		count := 0
		first := true

		cog := Retry(ctx, func() error {
			count++
			if first {
				first = false
				return ErrRetry
			}
			return nil
		}, 10)

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).ToNot(HaveOccurred())
		Expect(count).To(Equal(2))
	})

	It("should only retry for the max number of attempts", func() {
		ctx := context.Background()

		count := 0
		max := 10
		cog := Retry(ctx, func() error {
			count++
			return ErrRetry
		}, max)

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).To(Equal(ErrRetry))
		Expect(count).To(Equal(max))
	})

	It("should exit if context is canceled before completion", func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		count := 0
		max := 10
		cog := Retry(ctx, func() error {
			count++
			return ErrRetry
		}, max)

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).To(Equal(context.DeadlineExceeded))
		Expect(count).ToNot(Equal(max))

	})
})
