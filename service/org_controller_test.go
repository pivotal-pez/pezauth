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

	Describe("Put Control handler", func() {
		var (
			fakeName   = "testuser"
			fakeUser   = fmt.Sprintf("%s@pivotal.io", fakeName)
			fakeOrg    = fmt.Sprintf("pivot-%s", fakeName)
			render     *mockRenderer
			testLogger = log.New(os.Stdout, "testLogger", 0)
		)
		setGetUserInfo("pivotal.io", fakeUser)

		BeforeEach(func() {
			render = new(mockRenderer)
		})

		Context("when email is not in the system", func() {
			tokens := &mockTokens{}
			result := PivotOrg{
				Email:   fakeUser,
				OrgName: fakeOrg,
			}
			controlResponse := Response{Payload: structs.Map(result)}
			var orgPut OrgPutHandler = NewOrgController(&mockPersistence{
				err:    ErrNoMatchInStore,
				result: nil,
			}, new(mockHeritageClient)).Put().(OrgPutHandler)

			It("should create a new org record", func() {
				orgPut(martini.Params{UserParam: fakeUser}, testLogger, render, tokens)
				Ω(render.StatusCode).Should(Equal(SuccessStatus))
				Ω(render.ResponseObject).Should(Equal(controlResponse))
			})
		})
	})

	Describe("Get Control handler", func() {
		Context("calling controller with a bad user token combo", func() {
			var (
				fakeName   = "testuser"
				fakeUser   = fmt.Sprintf("%s@pivotal.io", fakeName)
				badUser    = fmt.Sprintf("%s@pivotal.io", "baduser")
				render     *mockRenderer
				testLogger = log.New(os.Stdout, "testLogger", 0)
			)
			setGetUserInfo("pivotal.io", badUser)

			BeforeEach(func() {
				render = new(mockRenderer)
			})

			Context("with a user that has no match in the system", func() {
				tokens := &mockTokens{}
				result := new(PivotOrg)
				controlResponse := Response{ErrorMsg: ErrCantCallAcrossUsers.Error()}
				var orgGet OrgGetHandler = NewOrgController(&mockPersistence{
					err:    ErrCantCallAcrossUsers,
					result: result,
				}, new(mockHeritageClient)).Get().(OrgGetHandler)

				It("should return an error and a fail status", func() {
					orgGet(martini.Params{UserParam: fakeUser}, testLogger, render, tokens)
					Ω(render.StatusCode).Should(Equal(FailureStatus))
					Ω(render.ResponseObject).Should(Equal(controlResponse))
				})
			})

			Context("with a user that has an org", func() {
				tokens := &mockTokens{}
				result := new(PivotOrg)
				controlResponse := Response{ErrorMsg: ErrCantCallAcrossUsers.Error()}
				var orgGet OrgGetHandler = NewOrgController(&mockPersistence{
					err:    ErrCantCallAcrossUsers,
					result: result,
				}, new(mockHeritageClient)).Get().(OrgGetHandler)

				It("should return a error object to the renderer", func() {
					orgGet(martini.Params{UserParam: fakeUser}, testLogger, render, tokens)
					Ω(render.StatusCode).Should(Equal(FailureStatus))
					Ω(render.ResponseObject).Should(Equal(controlResponse))
				})
			})
		})

		Context("calling controller", func() {
			var (
				fakeName   = "testuser"
				fakeUser   = fmt.Sprintf("%s@pivotal.io", fakeName)
				fakeOrg    = fmt.Sprintf("pivot-%s", fakeName)
				render     *mockRenderer
				testLogger = log.New(os.Stdout, "testLogger", 0)
			)
			setGetUserInfo("pivotal.io", fakeUser)

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
				}, new(mockHeritageClient)).Get().(OrgGetHandler)

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
				}, new(mockHeritageClient)).Get().(OrgGetHandler)

				It("should return a user object to the renderer", func() {
					orgGet(martini.Params{UserParam: fakeUser}, testLogger, render, tokens)
					Ω(render.StatusCode).Should(Equal(SuccessStatus))
					Ω(render.ResponseObject).Should(Equal(controlResponse))
				})
			})
		})
	})
})
