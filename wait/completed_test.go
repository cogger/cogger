package wait_test

import (
	. "github.com/cogger/cogger/wait"

	"time"

	"github.com/cogger/cogger/cogs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("Completed", func() {
	It("should eventually finish the work", func() {
		ctx := context.Background()

		finished := false
		Completed(ctx, cogs.Simple(ctx, func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		}), func() {
			finished = true
		})

		Eventually(func() bool {
			return finished
		}).Should(BeTrue())
	})
})
