package wait_test

import (
	. "github.com/cogger/cogger/wait"

	"errors"

	"github.com/cogger/cogger/cogs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("Resolve", func() {
	It("should return an array of errors when finished", func() {
		ctx := context.Background()

		fakeErr := errors.New("fake error")

		errs := Resolve(ctx, cogs.Simple(ctx, func() error {
			return fakeErr
		}))

		Expect(errs).To(HaveLen(1))
		Expect(errs[0]).To(Equal(fakeErr))
	})
})
