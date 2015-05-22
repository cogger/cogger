package cogger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"time"

	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var handlerInterface = reflect.TypeOf((*Handler)(nil)).Elem()

var _ = Describe("NewHandler", func() {
	It("should create a handler that implements cogger Handler interface", func() {
		handler := NewHandler()
		Expect(reflect.TypeOf(handler).Implements(handlerInterface)).To(BeTrue())
	})

	It("should return a 500 when no operation is set", func() {
		handler := NewHandler()

		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", &bytes.Buffer{})
		Expect(err).ToNot(HaveOccurred())

		handler.ServeHTTP(recorder, req)
		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
	})

	It("should return the status code returns from the handler if it is above or equal to 400", func() {
		for _, status := range []int{
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusInternalServerError,
			http.StatusForbidden,
		} {
			handler := NewHandler()

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/", &bytes.Buffer{})
			Expect(err).ToNot(HaveOccurred())
			handler.SetHandler(func(ctx context.Context, writer http.ResponseWriter, req *http.Request) int {
				return status
			})
			handler.ServeHTTP(recorder, req)
			Expect(recorder.Code).To(Equal(status))
		}
	})

	It("should timeout after the provided timeout", func() {
		handler := NewHandler()

		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", &bytes.Buffer{})
		Expect(err).ToNot(HaveOccurred())

		duration := time.Second * 1
		handler.SetTimeout(duration)
		handler.SetHandler(func(ctx context.Context, writer http.ResponseWriter, req *http.Request) int {
			time.Sleep(4 * time.Second)

			return http.StatusOK
		})
		handler.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
	})

	It("should add context before execution", func() {
		handler := NewHandler()

		counts := []int{}
		max := 10

		for i := 0; i < max; i++ {
			func(i int) {
				handler.AddContext(func(ctx context.Context, req *http.Request) context.Context {
					counts = append(counts, i)
					return ctx
				})
			}(i)
		}

		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", &bytes.Buffer{})
		Expect(err).ToNot(HaveOccurred())
		handler.ServeHTTP(recorder, req)

		Expect(counts).To(HaveLen(max))
		for i := 0; i < max; i++ {
			c := counts[i]
			Expect(c).To(Equal(i))
		}
	})

	It("should execute a handler by test method", func() {
		handler := NewHandler()
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", &bytes.Buffer{})
		Expect(err).ToNot(HaveOccurred())

		executed := false
		handler.SetHandler(func(ctx context.Context, writer http.ResponseWriter, req *http.Request) int {
			executed = true
			return http.StatusOK
		})

		handler.Test(context.Background(), recorder, req)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(executed).To(BeTrue())
	})
})

var _ = Describe("WithHandler", func() {
	It("should returns a handler that implements cogger Handler interface", func() {
		handler := WithHandler(func(ctx context.Context, writer http.ResponseWriter, req *http.Request) int {
			return http.StatusOK
		})
		Expect(reflect.TypeOf(handler).Implements(handlerInterface)).To(BeTrue())
	})

	It("should return the status code returns from the handler if it is above or equal to 400", func() {
		for _, status := range []int{
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusInternalServerError,
			http.StatusForbidden,
		} {
			handler := WithHandler(func(ctx context.Context, writer http.ResponseWriter, req *http.Request) int {
				return status
			})

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/", &bytes.Buffer{})
			Expect(err).ToNot(HaveOccurred())

			handler.ServeHTTP(recorder, req)
			Expect(recorder.Code).To(Equal(status))
		}
	})

	It("should timeout after the provided timeout", func() {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", &bytes.Buffer{})
		Expect(err).ToNot(HaveOccurred())

		duration := time.Second * 1

		handler := WithHandler(func(ctx context.Context, writer http.ResponseWriter, req *http.Request) int {
			time.Sleep(4 * time.Second)

			return http.StatusOK
		})

		handler.SetTimeout(duration)

		handler.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
	})

	It("should add context before execution", func() {
		handler := WithHandler(func(ctx context.Context, writer http.ResponseWriter, req *http.Request) int {
			return http.StatusOK
		})

		counts := []int{}
		max := 10

		for i := 0; i < max; i++ {
			func(i int) {
				handler.AddContext(func(ctx context.Context, req *http.Request) context.Context {
					counts = append(counts, i)
					return ctx
				})
			}(i)
		}

		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", &bytes.Buffer{})
		Expect(err).ToNot(HaveOccurred())
		handler.ServeHTTP(recorder, req)

		Expect(counts).To(HaveLen(max))
		for i := 0; i < max; i++ {
			c := counts[i]
			Expect(c).To(Equal(i))
		}
	})

	It("should execute a handler by test method", func() {
		executed := false
		handler := WithHandler(func(ctx context.Context, writer http.ResponseWriter, req *http.Request) int {
			executed = true
			return http.StatusOK
		})
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", &bytes.Buffer{})
		Expect(err).ToNot(HaveOccurred())

		handler.Test(context.Background(), recorder, req)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(executed).To(BeTrue())
	})
})
