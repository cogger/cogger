package limiter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLimiter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Limiter Suite")
}
