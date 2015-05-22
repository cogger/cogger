package wait

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1/cogs"
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
