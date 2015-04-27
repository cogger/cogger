package limiter_test

import (
	"reflect"
	"time"

	"github.com/cogger/cogger/cogs"
	. "github.com/cogger/cogger/limiter"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PerSecond", func() {
	limitInterface := reflect.TypeOf((*Limit)(nil)).Elem()

	It("should create a limiter", func() {
		limiter, err := PerSecond(10)
		Expect(err).ToNot(HaveOccurred())
		Expect(reflect.TypeOf(limiter).Implements(limitInterface)).To(BeTrue())
	})

	It("should return a config error when 0 or less is provided", func() {
		limiter, err := PerSecond(0)
		Expect(err).To(Equal(ErrLimitConfig))
		Expect(limiter).To(BeNil())
	})

	It("should limit cogs", func() {
		ctx, cancel := context.WithCancel(context.Background())

		count := 0
		perSecond := 2
		limiter, err := PerSecond(perSecond)
		Expect(err).ToNot(HaveOccurred())

		for i := 0; i < (perSecond*2)+1; i++ {
			cog := cogs.Simple(ctx, func() error {
				count++
				time.Sleep(1 * time.Second)
				return nil
			})
			cog.SetLimit(limiter)
			go func() {
				<-cog.Do(ctx)
			}()
		}

		time.Sleep(time.Duration(perSecond) * time.Second)
		cancel()

		Expect(count).To(Equal(perSecond + 1))
	})
})
