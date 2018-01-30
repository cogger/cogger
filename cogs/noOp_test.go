package cogs

import (
	"reflect"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"

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
