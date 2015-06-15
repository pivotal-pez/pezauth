package integrations_test

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/pezauth/integrations"
)

var _ = Describe("MyMongo", func() {
	var (
		appEnv *cfenv.App
		err    error
	)

	BeforeEach(func() {
		ip := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
		port := os.Getenv("MONGO_PORT_27017_TCP_PORT")
		setMongoEnv(ip, port)

		validEnv := []string{
			`VCAP_APPLICATION={}`,
			fmt.Sprintf("VCAP_SERVICES=%s", os.Getenv("VCAP_SERVICES")),
		}

		testEnv := cfenv.Env(validEnv)
		appEnv, err = cfenv.New(testEnv)
	})

	Context("Calling .New", func() {
		It("Should return a valid mongo session", func() {
			mngo := new(integrations.MyMongo).New(appEnv)
			Ω(mngo.Session.Ping()).Should(BeNil())
		})
	})

	Context("Calling .Remove", func() {
		It("Should not error", func() {
			mngo := new(integrations.MyMongo).New(appEnv)
			col := mngo.Collection()
			Ω(col.Remove(nil)).Should(BeNil())
		})
	})

	Context("Calling .Upsert", func() {
		It("Should not error", func() {
			mngo := new(integrations.MyMongo).New(appEnv)
			col := mngo.Collection()
			Ω(col.Upsert(nil, nil)).Should(BeNil())
		})
	})

	Context("Calling .FindOne", func() {
		It("Should not error", func() {
			mngo := new(integrations.MyMongo).New(appEnv)
			col := mngo.Collection()
			Ω(col.FindOne(nil, nil)).Should(BeNil())
		})
	})
})
