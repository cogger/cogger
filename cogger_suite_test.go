package cogger

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cogger Suite")
}
