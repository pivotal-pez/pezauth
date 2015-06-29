package integrations_test

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-pez/pezauth/integrations"
	"github.com/pivotal-pez/pezauth/service"
	"github.com/pivotal-pez/pezdispenser/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("MyMongo", func() {
	var (
		appEnv       *cfenv.App
		err          error
		col          pezdispenser.Persistence
		controlField = bson.M{"fieldname": "test"}
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
		mngo := new(integrations.MyMongo).New(appEnv)
		col = mngo.Collection()
		col.Remove(controlField)
	})

	AfterEach(func() {
		col.Remove(controlField)
	})

	Context("Calling .New", func() {
		It("Should return a valid mongo session", func() {
			mngo := new(integrations.MyMongo).New(appEnv)
			Ω(mngo.Session.Ping()).Should(BeNil())
		})
	})

	Context("Calling .Remove on non-existing record", func() {
		It("Should return error", func() {
			Ω(col.Remove(controlField)).Should(Equal(mgo.ErrNotFound))
		})
	})

	Context("Calling .Remove on valid record", func() {
		It("Should return error", func() {
			col.Upsert(controlField, controlField)
			err := col.Remove(controlField)
			Ω(err).Should(BeNil())
			Ω(err).ShouldNot(Equal(mgo.ErrNotFound))
		})
	})

	Context("Calling .Upsert", func() {
		It("Should not error", func() {
			mngo := new(integrations.MyMongo).New(appEnv)
			col := mngo.Collection()
			Ω(col.Upsert(controlField, controlField)).Should(BeNil())
		})
	})

	Context("Calling .FindOne on noexistent record", func() {
		It("Should return error", func() {
			mngo := new(integrations.MyMongo).New(appEnv)
			col := mngo.Collection()
			Ω(col.FindOne(controlField, nil)).Should(Equal(pezauth.ErrNoMatchInStore))
		})
	})

	Context("Calling .FindOne", func() {
		It("Should not error", func() {
			mngo := new(integrations.MyMongo).New(appEnv)
			col := mngo.Collection()
			col.Upsert(controlField, controlField)
			Ω(col.FindOne(controlField, nil)).ShouldNot(Equal(pezauth.ErrNoMatchInStore))
			Ω(col.FindOne(controlField, nil)).Should(BeNil())
		})
	})
})
