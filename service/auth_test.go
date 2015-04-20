package pezauth_test

import (
	"fmt"
	"os"

	"github.com/go-martini/martini"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/pezauth/service"
)

var _ = Describe("Authentication", func() {

	Describe("InitAuth", func() {
		var (
			m *martini.ClassicMartini
		)

		BeforeEach(func() {
			setVcapApp("http://localhost:3000")
			setVcapServ()
			os.Setenv("PORT", "3000")
			m = martini.Classic()
		})

		Context("calling InitAuth with no enviornment variables set", func() {
			var (
				validDomain = "pivotal.io"
				validUser   = "testuser@pivotal.io"
			)

			setGetUserInfo(validDomain, validUser)

			BeforeEach(func() {
				os.Unsetenv("VCAP_APPLICATION")
				os.Unsetenv("VCAP_SERVICES")
			})

			It("Should panic", func() {
				Ω(func() {
					InitAuth(m, &mockRedisCreds{})
				}).Should(Panic())
			})
		})

		Context("when InitAuth is passed a classic martini", func() {
			It("Should setup the authentication middleware without panic", func() {
				Ω(func() {
					InitAuth(m, &mockRedisCreds{})
				}).ShouldNot(Panic())
			})
		})

		Context("calling DomainCheck with a valid domain", func() {
			var (
				validDomain = "pivotal.io"
				validUser   = "testuser@pivotal.io"
			)
			setGetUserInfo(validDomain, validUser)

			It("Should have a valid statuscode and body", func() {
				mock := new(mockResponseWriter)
				DomainChecker(mock, new(mockTokens))
				Ω(mock.StatusCode).ShouldNot(Equal(AuthFailStatus))
				Ω(mock.Body).ShouldNot(Equal(AuthFailureResponse))
			})

			Context("un-formatted domain", func() {
				BeforeEach(func() {
					setVcapApp("pivotal.io")
				})

				It("should format the domain in the config object", func() {
					control := fmt.Sprintf("https://%s/oauth2callback", validDomain)
					InitAuth(m, &mockRedisCreds{})
					Ω(OauthConfig.RedirectURL).Should(Equal(control))
				})
			})

			Context("version formatted domain", func() {
				BeforeEach(func() {
					setVcapApp("pivotal-1919241972nwdighsd921h192t23t.io")
				})

				It("should format the domain in the config object", func() {
					control := fmt.Sprintf("https://%s/oauth2callback", validDomain)
					InitAuth(m, &mockRedisCreds{})
					Ω(OauthConfig.RedirectURL).Should(Equal(control))
				})
			})

		})

		Context("calling DomainCheck with a in-valid domain", func() {
			var (
				inValidDomain = "google.com"
				validUser     = "testuser@pivotal.io"
			)
			setGetUserInfo(inValidDomain, validUser)

			It("Should return true", func() {
				mock := new(mockResponseWriter)
				DomainChecker(mock, new(mockTokens))
				Ω(mock.StatusCode).Should(Equal(AuthFailStatus))
			})
		})
	})
})
