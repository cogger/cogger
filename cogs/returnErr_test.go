package cogs_test

import (
	"errors"
	"reflect"

	"github.com/cogger/cogger"
	. "github.com/cogger/cogger/cogs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("ReturnErr", func() {
	coggerInterface := reflect.TypeOf((*cogger.Cog)(nil)).Elem()

	It("should create a cog that returns what ever error is passed to it", func() {
		err := errors.New("random error")
		cog := ReturnErr(err)

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).To(HaveOccurred())
		Expect(err).To(Equal(err))
	})
})
