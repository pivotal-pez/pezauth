package pezauth_test

import (
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/pezauth/service"
)

var _ = Describe("NewMeController", func() {
	Context("calling controller", func() {
		var (
			render     *mockRenderer
			testLogger = log.New(os.Stdout, "testLogger", 0)
		)
		setGetUserInfo("pivotal.io", "jcalabrese@pivotal.io")

		BeforeEach(func() {
			render = new(mockRenderer)
		})

		It("should return a user object to the renderer", func() {
			controlResponse := Response{}
			var meGet ValidateGetHandler = NewValidateV1().Get().(ValidateGetHandler)
			meGet(testLogger, render)
			Ω(render.StatusCode).Should(Equal(200))
			Ω(render.ResponseObject).Should(Equal(controlResponse))
		})
	})
})
