package integrations_test

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/pezauth/integrations"
)

var _ = Describe("MyRedis", func() {
	var (
		appEnv *cfenv.App
		err    error
	)

	BeforeEach(func() {
		ip := os.Getenv("REDIS_PORT_6379_TCP_ADDR")
		port := os.Getenv("REDIS_PORT_6379_TCP_PORT")
		setRedisEnv(ip, port)

		validEnv := []string{
			`VCAP_APPLICATION={}`,
			fmt.Sprintf("VCAP_SERVICES=%s", os.Getenv("VCAP_SERVICES")),
		}

		testEnv := cfenv.Env(validEnv)
		appEnv, err = cfenv.New(testEnv)
	})

	Context("Calling .New", func() {
		It("Should return a valid redis session", func() {
			rdis := new(integrations.MyRedis).New(appEnv)
			Ω(rdis.Conn.Err()).Should(BeNil())
		})
	})
})
