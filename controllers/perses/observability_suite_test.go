package perses_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPersesController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Perses Suite")
}
