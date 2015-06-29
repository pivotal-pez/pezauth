package pezauth_test

import (
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezauth/service"
)

var _ = Describe("create org", func() {
	Context("calling .SafeCreate() w/ successfull rest calls", func() {
		var (
			username                     = "testuser@email.com"
			log      *log.Logger         = log.New(os.Stdout, "logger: ", log.Lshortfile)
			store    *mockOrgPersistence = &mockOrgPersistence{
				CalledUpsert: 0,
			}
			authClient AuthRequestCreator = &mockAuthRequestCreator{
				MyMockDoer: NewMockDoer(),
			}
		)

		It("should write a record to mongo and not error", func() {
			org := NewOrg(username, log, nil, store, authClient)
			_, err := org.SafeCreate()
			Ω(err).Should(BeNil())
			Ω(store.CalledUpsert).Should(Equal(1))
		})
	})
})
