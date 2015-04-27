package cogs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCogs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cogs Suite")
}
