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

		Context("with a user that has no match in the system", func() {
			tokens := &mockTokens{}
			result := PivotOrg{
				Email:   fakeUser,
				OrgName: fakeOrg,
			}
			controlResponse := Response{ErrorMsg: ErrNoMatchInStore.Error()}
			var orgGet OrgGetHandler = NewOrgController(&mockPersistence{
				err:    ErrNoMatchInStore,
				result: result,
			}).Get().(OrgGetHandler)

			It("should return an error and a fail status", func() {
				orgGet(martini.Params{UserParam: fakeUser}, testLogger, render, tokens)
				Ω(render.StatusCode).Should(Equal(FailureStatus))
				Ω(render.ResponseObject).Should(Equal(controlResponse))
			})
		})

		Context("with a user that has an org", func() {
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

			It("should return a user object to the renderer", func() {
				orgGet(martini.Params{UserParam: fakeUser}, testLogger, render, tokens)
				Ω(render.StatusCode).Should(Equal(SuccessStatus))
				Ω(render.ResponseObject).Should(Equal(controlResponse))
			})
		})
	})
})
