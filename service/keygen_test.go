package pezauth_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyGen", func() {
	var (
		username = "myfakeusername"
		guid     = "myfakeguid"
		err      error
		response string
	)
	Context("Get function", func() {
		Context("calling get with a valid user arg", func() {
			BeforeEach(func() {
				k := getKeygen(false, guid)
				response, err = k.Get(username)
			})

			It("Should return an api key for that user", func() {
				Ω(response).Should(Equal(guid))
			})

			It("Should return a nil error", func() {
				Ω(err).Should(BeNil())
			})
		})

		Context("calling get with a In-valid user arg", func() {
			BeforeEach(func() {
				k := getKeygen(true, guid)
				response, err = k.Get(username)
			})

			It("Should return an api key for that user", func() {
				Ω(response).ShouldNot(Equal(guid))
			})

			It("Should return a nil error", func() {
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(Equal(errDoerCallFailure))
			})
		})
	})

	Context("Create function", func() {
		Context("calling Create with a valid user arg", func() {
			BeforeEach(func() {
				k := getKeygen(false, guid)
				err = k.Create(username)
			})

			It("Should return a nil error", func() {
				Ω(err).Should(BeNil())
			})
		})

		Context("calling Create with a In-valid user arg", func() {
			BeforeEach(func() {
				k := getKeygen(true, guid)
				err = k.Create(username)
			})

			It("Should return a nil error", func() {
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(Equal(errDoerCallFailure))
			})
		})
	})

	Context("Delete function", func() {
		Context("calling Delete with a valid user arg", func() {
			BeforeEach(func() {
				k := getKeygen(false, guid)
				err = k.Delete(username)
			})

			It("Should return a nil error", func() {
				Ω(err).Should(BeNil())
			})
		})

		Context("calling Delete with a In-valid user arg", func() {
			BeforeEach(func() {
				k := getKeygen(true, guid)
				err = k.Create(username)
			})

			It("Should return a nil error", func() {
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(Equal(errDoerCallFailure))
			})
		})
	})

})
