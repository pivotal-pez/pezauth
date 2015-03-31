package pezauth_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPezAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pez Auth Suite")
}
