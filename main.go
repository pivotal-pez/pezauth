package main

import (
	"fmt"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	pez "github.com/pivotalservices/pezauth/service"
)

type redisCreds struct {
	pass string
	uri  string
}

func (s *redisCreds) Pass() string {
	return s.pass
}

func (s *redisCreds) Uri() string {
	return s.uri
}

func main() {
	appEnv, _ := cfenv.Current()
	redisName := os.Getenv("REDIS_SERVICE_NAME")
	redisHost := os.Getenv("REDIS_HOSTNAME_NAME")
	redisPass := os.Getenv("REDIS_PASSWORD_NAME")
	redisPort := os.Getenv("REDIS_PORT_NAME")
	m := martini.Classic()
	name, _ := appEnv.Services.WithName(redisName)
	connectionURI := fmt.Sprintf("%s:%s", name.Credentials[redisHost], name.Credentials[redisPort])

	if c, err := redis.Dial("tcp", connectionURI); err == nil {

		if _, err := c.Do("AUTH", name.Credentials[redisPass]); err == nil {
			pez.InitSession(m, &redisCreds{
				pass: name.Credentials[redisPass],
				uri:  connectionURI,
			})
			pez.InitRoutes(m, c)
			m.Run()

		} else {
			fmt.Println(err)
		}
		c.Close()

	} else {
		fmt.Println(err)
	}
}
