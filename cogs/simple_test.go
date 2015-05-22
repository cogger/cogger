package cogs_test

import (
	"reflect"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
	. "gopkg.in/cogger/cogger.v1/cogs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Simple", func() {
	coggerInterface := reflect.TypeOf((*cogger.Cog)(nil)).Elem()

	It("should create a simple Cog", func() {
		ctx := context.Background()
		cog := Simple(ctx, func() error {
			return nil
		})

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).ToNot(HaveOccurred())
	})

	It("should cancel when the context cancels", func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		cog := Simple(ctx, func() error {
			time.Sleep(2 * time.Second)
			return nil
		})

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(ctx)).To(Equal(context.DeadlineExceeded))
	})
})
