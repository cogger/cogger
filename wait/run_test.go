package wait_test

import (
	"time"

	"github.com/cogger/cogger/cogs"
	. "github.com/cogger/cogger/wait"
	"golang.org/x/net/context"

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
