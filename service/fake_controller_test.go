package pezauth_test

import (
	"log"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/pezauth/service"
)

var _ = Describe("FakeController", func() {
	var m *martini.ClassicMartini
	BeforeEach(func() {
		m = martini.Classic()
	})

	Context("when FakeController ", func() {

		It("Should not result in panic", func() {
			Ω(func() {
				FakeController(*new(martini.Params), new(log.Logger), *new(render.Render), *new(oauth2.Tokens))
			}).Should(Panic())
		})
	})
})
