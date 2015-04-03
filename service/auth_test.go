package pezauth_test

import (
	"os"

	"github.com/go-martini/martini"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/pezauth/service"
)

var _ = Describe("Authentication", func() {
	Describe("InitAuth", func() {
		var m *martini.ClassicMartini
		BeforeEach(func() {
			setVcapApp()
			setVcapServ()
			os.Setenv("PORT", "3000")
			m = martini.Classic()
		})

		Context("when InitAuth is passed a classic martini", func() {
			It("Should setup the authentication middleware without panic", func() {
				Ω(func() {
					InitAuth(m)
				}).ShouldNot(Panic())
			})
		})
	})
})
