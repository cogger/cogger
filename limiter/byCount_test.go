package limiter

import (
	"reflect"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1/cogs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ByCount", func() {
	limitInterface := reflect.TypeOf((*Limit)(nil)).Elem()

	It("should create a limiter", func() {
		limiter, err := ByCount(10)
		Expect(err).ToNot(HaveOccurred())
		Expect(reflect.TypeOf(limiter).Implements(limitInterface)).To(BeTrue())
	})

	It("should return a config error when 0 or less is provided", func() {
		limiter, err := ByCount(0)
		Expect(err).To(Equal(ErrLimitConfig))
		Expect(limiter).To(BeNil())
	})

	It("should limit cogs", func() {
		ctx, cancel := context.WithCancel(context.Background())

		count := 0
		countLimit := 2
		limiter, err := ByCount(countLimit)
		Expect(err).ToNot(HaveOccurred())

		for i := 0; i < (countLimit*2)+1; i++ {
			cog := cogs.Simple(ctx, func() error {
				count++
				time.Sleep(100 * time.Millisecond)
				return nil
			})
			cog.SetLimit(limiter)
			go func() {
				<-cog.Do(ctx)
			}()
		}

		time.Sleep(200 * time.Millisecond)
		cancel()

		Expect(count).To(Equal(countLimit * 2))
	})
})
