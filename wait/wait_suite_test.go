package wait_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWait(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wait Suite")
}
