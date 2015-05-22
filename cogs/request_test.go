package cogs_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
	. "gopkg.in/cogger/cogger.v1/cogs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Request", func() {
	coggerInterface := reflect.TypeOf((*cogger.Cog)(nil)).Elem()

	It("should create a request Cog", func() {
		ctx := context.Background()
		cog := Simple(ctx, func() error {
			return nil
		})

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).ToNot(HaveOccurred())
	})

	It("should execute the request", func() {
		ctx := context.Background()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		transport := &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(server.URL)
			},
		}

		httpClient := &http.Client{Transport: transport}

		req, err := http.NewRequest("GET", "http://localhost/", &bytes.Buffer{})
		Expect(err).ToNot(HaveOccurred())

		var status int
		cog := Request(ctx, httpClient, req, func(resp *http.Response, err error) error {
			Expect(err).ToNot(HaveOccurred())
			status = resp.StatusCode
			return nil
		})

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(ctx)).ToNot(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
	})

	It("should cancel when the context cancels", func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		transport := &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(server.URL)
			},
		}

		httpClient := &http.Client{Transport: transport}

		req, err := http.NewRequest("GET", "http://localhost/", &bytes.Buffer{})
		Expect(err).ToNot(HaveOccurred())

		var status int = 0
		cog := Request(ctx, httpClient, req, func(resp *http.Response, err error) error {
			Expect(err).ToNot(HaveOccurred())
			status = resp.StatusCode
			return nil
		})

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(ctx)).To(Equal(context.DeadlineExceeded))
		Expect(status).To(Equal(0))
	})
})
