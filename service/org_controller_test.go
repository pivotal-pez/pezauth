package pezauth_test

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/structs"
	"github.com/go-martini/martini"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/pezauth/service"
)

var _ = Describe("NewOrgController", func() {
	Context("calling controller", func() {
		var (
			fakeName   = "testuser"
			fakeUser   = fmt.Sprintf("%s@pivotal.io", fakeName)
			fakeOrg    = fmt.Sprintf("pivot-%s", fakeName)
			render     *mockRenderer
			testLogger = log.New(os.Stdout, "testLogger", 0)
		)
		setGetUserInfo("pivotal.io", "jcalabrese@pivotal.io")

		BeforeEach(func() {
			render = new(mockRenderer)
		})

		It("should return a user object to the renderer", func() {
			tokens := &mockTokens{}
			result := PivotOrg{
				Email:   fakeUser,
				OrgName: fakeOrg,
			}
			controlResponse := Response{Payload: structs.Map(result)}
			var orgGet OrgGetHandler = NewOrgController(&mockPersistence{
				err:    nil,
				result: result,
			}).Get().(OrgGetHandler)
			orgGet(martini.Params{UserParam: fakeUser}, testLogger, render, tokens)
			Ω(render.StatusCode).Should(Equal(200))
			Ω(render.ResponseObject).Should(Equal(controlResponse))
		})
	})
})
