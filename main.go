package main

import (
	"fmt"
	"os"
	"strings"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	pez "github.com/pivotalservices/pezauth/service"
	"gopkg.in/mgo.v2"
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
	mongoServiceName := os.Getenv("MONGO_SERVICE_NAME")
	mongoURIName := os.Getenv("MONGO_URI_NAME")
	mongoCollName := os.Getenv("MONGO_COLLECTION_NAME")

	m := martini.Classic()
	redisService, _ := appEnv.Services.WithName(redisName)
	mongoService, _ := appEnv.Services.WithName(mongoServiceName)
	mongoConnectionURI := mongoService.Credentials[mongoURIName]
	parsedURI := strings.Split(mongoConnectionURI, "/")
	mongoDBName := parsedURI[len(parsedURI)-1]
	connectionURI := fmt.Sprintf("%s:%s", redisService.Credentials[redisHost], redisService.Credentials[redisPort])

	if redisConn, err := redis.Dial("tcp", connectionURI); err == nil {
		defer redisConn.Close()

		if _, err := redisConn.Do("AUTH", redisService.Credentials[redisPass]); err == nil {
			pez.InitSession(m, &redisCreds{
				pass: redisService.Credentials[redisPass],
				uri:  connectionURI,
			})

			if session, err := mgo.Dial(mongoConnectionURI); err == nil {
				defer session.Close()
				session.SetMode(mgo.Monotonic, true)
				mongoConn := session.DB(mongoDBName).C(mongoCollName)
				pez.InitRoutes(m, redisConn, mongoConn)
				m.Run()

			} else {
				fmt.Println(err)
			}

		} else {
			fmt.Println(err)
		}

	} else {
		fmt.Println(err)
	}
}
