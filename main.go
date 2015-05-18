package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/gorelic"
	pez "github.com/pivotalservices/pezauth/service"
	"github.com/xchapter7x/cloudcontroller-client"
	"gopkg.in/mgo.v2"
)

func main() {
	appEnv, _ := cfenv.Current()
	m := martini.Classic()
	newRelic := new(myNewRelic).New(appEnv)
	gorelic.InitNewrelicAgent(newRelic.Key, newRelic.App, true)
	m.Use(gorelic.Handler)
	rds := new(myRedis).New(appEnv)
	defer rds.Conn.Close()
	pez.InitSession(m, &redisCreds{
		pass: rds.Pass,
		uri:  rds.URI,
	})
	h := new(myHeritage).New(appEnv)
	heritageClient := &heritage{
		Client:   ccclient.New(h.LoginTarget, h.LoginUser, h.LoginPass, new(http.Client)),
		ccTarget: h.CCTarget,
	}
	mngo := new(myMongo).New(appEnv)
	defer mngo.Session.Close()

	if _, err := heritageClient.Login(); err == nil {
		pez.InitRoutes(m, rds.Conn, mngo.Col, heritageClient)
		m.Run()

	} else {
		panic(fmt.Sprintf("heritage client login error: %s", err.Error()))
	}
}

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

type (
	myRedis struct {
		URI  string
		Pass string
		Conn redis.Conn
	}
	myMongo struct {
		Col     *mgo.Collection
		Session *mgo.Session
	}
	myHeritage struct {
		LoginTarget string
		LoginUser   string
		LoginPass   string
		CCTarget    string
	}
	myNewRelic struct {
		Key string
		App string
	}
)

func (s *myHeritage) New(appEnv *cfenv.App) *myHeritage {
	heritageAdminServiceName := os.Getenv("UPS_PEZ_HERITAGE_ADMIN_NAME")
	heritageLoginTargetName := os.Getenv("HERITAGE_LOGIN_TARGET_NAME")
	heritageLoginUserName := os.Getenv("HERITAGE_LOGIN_USER_NAME")
	heritageLoginPassName := os.Getenv("HERITAGE_LOGIN_PASS_NAME")
	heritageCCTargetName := os.Getenv("HERITAGE_CC_TARGET_NAME")
	heritageAdminService, err := appEnv.Services.WithName(heritageAdminServiceName)

	if err != nil {
		panic(fmt.Sprintf("heritage service name error: %s", err.Error()))
	}
	s.LoginTarget = heritageAdminService.Credentials[heritageLoginTargetName]
	s.LoginUser = heritageAdminService.Credentials[heritageLoginUserName]
	s.LoginPass = heritageAdminService.Credentials[heritageLoginPassName]
	s.CCTarget = heritageAdminService.Credentials[heritageCCTargetName]
	return s
}

func (s *myNewRelic) New(appEnv *cfenv.App) *myNewRelic {
	serviceName := os.Getenv("NEWRELIC_SERVICE_NAME")
	keyName := os.Getenv("NEWRELIC_KEY_NAME")
	appName := os.Getenv("NEWRELIC_APP_NAME")
	service, err := appEnv.Services.WithName(serviceName)

	if err != nil {
		panic(fmt.Sprintf("new relic service name error: %s", err.Error()))
	}
	s.Key = service.Credentials[keyName]
	s.App = service.Credentials[appName]
	return s
}

func (s *myMongo) New(appEnv *cfenv.App) *myMongo {
	mongoServiceName := os.Getenv("MONGO_SERVICE_NAME")
	mongoURIName := os.Getenv("MONGO_URI_NAME")
	mongoCollName := os.Getenv("MONGO_COLLECTION_NAME")
	mongoService, err := appEnv.Services.WithName(mongoServiceName)

	if err != nil {
		panic(fmt.Sprintf("mongodb service name error: %s", err.Error()))
	}
	mongoConnectionURI := mongoService.Credentials[mongoURIName]
	parsedURI := strings.Split(mongoConnectionURI, "/")
	mongoDBName := parsedURI[len(parsedURI)-1]

	if s.Session, err = mgo.Dial(mongoConnectionURI); err != nil {
		panic(fmt.Sprintf("mongodb dial error: %s", err.Error()))
	}
	s.Session.SetMode(mgo.Monotonic, true)
	s.Col = s.Session.DB(mongoDBName).C(mongoCollName)
	return s
}

func (s *myRedis) New(appEnv *cfenv.App) *myRedis {
	redisName := os.Getenv("REDIS_SERVICE_NAME")
	redisHost := os.Getenv("REDIS_HOSTNAME_NAME")
	redisPass := os.Getenv("REDIS_PASSWORD_NAME")
	redisPort := os.Getenv("REDIS_PORT_NAME")
	redisService, err := appEnv.Services.WithName(redisName)
	s.Pass = redisService.Credentials[redisPass]
	s.URI = fmt.Sprintf("%s:%s", redisService.Credentials[redisHost], redisService.Credentials[redisPort])

	if err != nil {
		panic(fmt.Sprintf("redis service name error: %s", err.Error()))
	}

	if s.Conn, err = redis.Dial("tcp", s.URI); err != nil {
		panic(fmt.Sprintf("redis dial error: %s", err.Error()))
	}

	if _, err = s.Conn.Do("AUTH", s.Pass); err != nil {
		panic(fmt.Sprintf("redis auth error: %s", err.Error()))
	}
	return s
}
