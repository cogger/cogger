package order

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOrder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Order Suite")
}
