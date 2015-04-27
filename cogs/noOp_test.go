package cogs_test

import (
	"reflect"

	"github.com/cogger/cogger"
	. "github.com/cogger/cogger/cogs"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NoOp", func() {
	coggerInterface := reflect.TypeOf((*cogger.Cog)(nil)).Elem()

	It("should create a noOp Cog", func() {
		cog := NoOp()

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).ToNot(HaveOccurred())
	})
})
