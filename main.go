package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	pez "github.com/pivotalservices/pezauth/service"
	"github.com/xchapter7x/cloudcontroller-client"
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

type heritage struct {
	*ccclient.Client
	ccTarget string
}

func (s *heritage) CCTarget() string {
	return s.ccTarget
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
	heritageAdminServiceName := os.Getenv("UPS_PEZ_HERITAGE_ADMIN_NAME")
	heritageLoginTargetName := os.Getenv("HERITAGE_LOGIN_TARGET_NAME")
	heritageLoginUserName := os.Getenv("HERITAGE_LOGIN_USER_NAME")
	heritageLoginPassName := os.Getenv("HERITAGE_LOGIN_PASS_NAME")
	heritageCCTargetName := os.Getenv("HERITAGE_CC_TARGET_NAME")

	m := martini.Classic()
	redisService, err := appEnv.Services.WithName(redisName)

	if err != nil {
		panic(err.Error())
	}
	mongoService, err := appEnv.Services.WithName(mongoServiceName)

	if err != nil {
		panic(err.Error())
	}
	heritageAdminService, err := appEnv.Services.WithName(heritageAdminServiceName)

	if err != nil {
		panic(err.Error())
	}
	heritageLoginTarget := heritageAdminService.Credentials[heritageLoginTargetName]
	heritageLoginUser := heritageAdminService.Credentials[heritageLoginUserName]
	heritageLoginPass := heritageAdminService.Credentials[heritageLoginPassName]
	heritageCCTarget := heritageAdminService.Credentials[heritageCCTargetName]
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
				heritageClient := &heritage{
					Client:   ccclient.New(heritageLoginTarget, heritageLoginUser, heritageLoginPass, new(http.Client)),
					ccTarget: heritageCCTarget,
				}

				if _, err := heritageClient.Login(); err == nil {
					pez.InitRoutes(m, redisConn, mongoConn, heritageClient)
					m.Run()

				} else {
					fmt.Println("heritage client login error: ", err)
				}

			} else {
				fmt.Println("mongodb dial error: ", err)
			}

		} else {
			fmt.Println("redis auth error: ", err)
		}

	} else {
		fmt.Println("redis dial error: ", err)
	}
}
