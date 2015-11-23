package y

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestY(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Y Suite")
}
