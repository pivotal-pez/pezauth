package pezauth_test

import (
	"os"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/pezauth/service"
)

var _ = Describe("Authentication", func() {

	Describe("InitAuth", func() {
		var (
			m              *martini.ClassicMartini
			oldGetUserInfo func(tokens oauth2.Tokens) map[string]interface{}
		)

		BeforeEach(func() {
			setVcapApp()
			setVcapServ()
			os.Setenv("PORT", "3000")
			m = martini.Classic()
		})

		Context("calling InitAuth with no enviornment variables set", func() {
			var validDomain = "pivotal.io"

			BeforeEach(func() {
				os.Unsetenv("VCAP_APPLICATION")
				os.Unsetenv("VCAP_SERVICES")
				oldGetUserInfo = GetUserInfo
				GetUserInfo = func(tokens oauth2.Tokens) map[string]interface{} {
					return map[string]interface{}{
						"domain": validDomain,
					}
				}
			})

			AfterEach(func() {
				GetUserInfo = oldGetUserInfo
			})

			It("Should panic", func() {
				Ω(func() {
					InitAuth(m)
				}).Should(Panic())
			})
		})

		Context("when InitAuth is passed a classic martini", func() {
			It("Should setup the authentication middleware without panic", func() {
				Ω(func() {
					InitAuth(m)
				}).ShouldNot(Panic())
			})
		})

		Context("calling DomainCheck with a valid domain", func() {
			var validDomain = "pivotal.io"

			BeforeEach(func() {
				oldGetUserInfo = GetUserInfo
				GetUserInfo = func(tokens oauth2.Tokens) map[string]interface{} {
					return map[string]interface{}{
						"domain": validDomain,
					}
				}
			})

			AfterEach(func() {
				GetUserInfo = oldGetUserInfo
			})

			It("Should return true", func() {
				mock := new(mockResponseWriter)
				DomainChecker(mock, new(mockTokens))
				Ω(mock.StatusCode).ShouldNot(Equal(AuthFailStatus))
			})
		})

		Context("calling DomainCheck with a in-valid domain", func() {
			var inValidDomain = "google.com"

			BeforeEach(func() {
				oldGetUserInfo = GetUserInfo
				GetUserInfo = func(tokens oauth2.Tokens) map[string]interface{} {
					return map[string]interface{}{
						"domain": inValidDomain,
					}
				}
			})

			AfterEach(func() {
				GetUserInfo = oldGetUserInfo
			})

			It("Should return true", func() {
				mock := new(mockResponseWriter)
				DomainChecker(mock, new(mockTokens))
				Ω(mock.StatusCode).Should(Equal(AuthFailStatus))
			})
		})
	})
})
