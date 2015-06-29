package integrations_test

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-pez/pezauth/integrations"
)

var _ = Describe("MyRedis", func() {
	var (
		appEnv *cfenv.App
		err    error
		key    = "INTEGRATION_TEST_KEY"
		val    = "INTEGRATION_TEST_KEYS_VALUE"

		hash       = "INTEGRATION_TEST_HASH"
		hashSubkey = "INTEGRATION_TEST_HASH_SUBKEY"
		hashSubval = "INTEGRATION_TEST_HASH_SUBKEY_VALUE"
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
		rdis := new(integrations.MyRedis).New(appEnv)

		rdis.Conn.Do("DEL", key)
		rdis.Conn.Do("DEL", hash)
	})

	Context("Calling .New", func() {
		It("Should return a valid redis session", func() {
			rdis := new(integrations.MyRedis).New(appEnv)
			Ω(rdis.Conn.Err()).Should(BeNil())
		})
	})

	Context("Keys", func() {
		Context("Calling SET for a key/value", func() {
			It("Should not error", func() {
				rdis := new(integrations.MyRedis).New(appEnv)
				s, e := rdis.Conn.Do("SET", key, val)
				Ω(s).Should(Equal("OK"))
				Ω(e).Should(BeNil())
			})
		})

		Context("Calling GET on a key", func() {
			It("Should return the key's value", func() {
				rdis := new(integrations.MyRedis).New(appEnv)
				rdis.Conn.Do("SET", key, val)
				o, e := rdis.Conn.Do("GET", key)
				Ω(fmt.Sprintf("%s", o)).Should(Equal(val))
				Ω(e).Should(BeNil())
			})
		})

		Context("Calling DEL on a key", func() {
			It("Should delete the matching key", func() {
				rdis := new(integrations.MyRedis).New(appEnv)
				rdis.Conn.Do("SET", key, val)
				o, e := rdis.Conn.Do("GET", key)
				Ω(e).Should(BeNil())
				Ω(fmt.Sprintf("%s", o)).Should(Equal(val))
				_, e = rdis.Conn.Do("DEL", key)
				Ω(e).Should(BeNil())
				o, e = rdis.Conn.Do("GET", key)
				Ω(o).Should(BeNil())
				Ω(e).Should(BeNil())
			})
		})

		Context("Calling KEYS with patter", func() {
			It("Should find keys that match", func() {
				rdis := new(integrations.MyRedis).New(appEnv)
				rdis.Conn.Do("SET", key, val)
				o, e := rdis.Conn.Do("KEYS", key[0:len(key)-1]+"*")
				s, _ := o.([]interface{})
				Ω(len(s)).Should(Equal(1))
				Ω(fmt.Sprintf("%s", s[0])).Should(Equal(key))
				Ω(e).Should(BeNil())
			})
		})
	})

	Context("Hashes", func() {
		Context("Calling HMSET for a hash with key/value", func() {
			It("Should not error", func() {
				rdis := new(integrations.MyRedis).New(appEnv)
				s, e := rdis.Conn.Do("HMSET", hash, hashSubkey, hashSubval)
				Ω(s).Should(Equal("OK"))
				Ω(e).Should(BeNil())
			})
		})

		Context("Calling HMGET for a hash with a given key", func() {
			It("Should return stored value for that key", func() {
				rdis := new(integrations.MyRedis).New(appEnv)
				rdis.Conn.Do("HMSET", hash, hashSubkey, hashSubval)
				o, e := rdis.Conn.Do("HMGET", hash, hashSubkey)
				Ω(e).Should(BeNil())
				s, _ := o.([]interface{})
				Ω(fmt.Sprintf("%s", s[0])).Should(Equal(hashSubval))
			})
		})
	})
})
