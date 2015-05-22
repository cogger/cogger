package wait_test

import (
	"time"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1/cogs"
	. "gopkg.in/cogger/cogger.v1/wait"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run", func() {
	It("should eventually finish", func() {
		ctx := context.Background()

		finished := false
		NoBlock(ctx, cogs.Simple(ctx, func() error {
			time.Sleep(100 * time.Millisecond)
			finished = true
			return nil
		}))

		Eventually(func() bool {
			return finished
		}).Should(BeTrue())
	})
})
