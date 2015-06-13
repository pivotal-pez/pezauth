package integrations_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIntegrations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pez Auth Integrations Suite")
}
