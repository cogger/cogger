package cogger

import (
	"errors"
	"reflect"

	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewCog", func() {
	coggerInterface := reflect.TypeOf((*Cog)(nil)).Elem()

	It("should create a cog", func() {
		cog := NewCog(func() chan error {
			return make(chan error)
		})

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
	})

	Context("the cogger interface", func() {
		Context("when a executing a job", func() {
			It("should return nil where there are no errors", func() {
				cog := NewCog(func() chan error {
					done := make(chan error)
					go func() {
						done <- nil
					}()
					return done
				})

				ctx := context.Background()
				Expect(<-cog.Do(ctx)).ToNot(HaveOccurred())
			})

			It("should return an error when there is an error", func() {
				testErr := errors.New("test error")
				cog := NewCog(func() chan error {
					done := make(chan error)
					go func() {
						done <- testErr
					}()
					return done
				})

				ctx := context.Background()
				Expect(<-cog.Do(ctx)).To(Equal(testErr))
			})
		})

		It("should implement SetLimit function", func() {
			cog := NewCog(func() chan error {
				done := make(chan error)
				go func() {
					done <- nil
				}()
				return done
			})

			limit := &mockLimit{}

			cog.SetLimit(limit)

			ctx := context.Background()
			Expect(<-cog.Do(ctx)).ToNot(HaveOccurred())
			Expect(limit.NextHits).To(Equal(1))
			Expect(limit.Completed).To(BeTrue())
		})
	})
})

type mockLimit struct {
	Completed bool
	NextHits  int
}

func (limit *mockLimit) Next(ctx context.Context) chan struct{} {
	next := make(chan struct{})
	go func() {
		limit.NextHits++
		next <- struct{}{}
	}()
	return next

}

func (limit *mockLimit) Done(ctx context.Context) {
	limit.Completed = true
}
