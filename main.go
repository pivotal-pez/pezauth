package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

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
	gorelic.InitNewrelicAgent(newRelic.Key, newRelic.App, false)
	m.Use(gorelic.Handler)
	oauth2Client := new(myOAuth2).New(appEnv)
	pez.ClientID = oauth2Client.ID
	pez.ClientSecret = oauth2Client.Secret
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
	myOAuth2 struct {
		ID     string
		Secret string
	}
	myRedis struct {
		URI  string
		Pass string
		Conn redis.Conn
	}
	myMongo struct {
		Col                *mgo.Collection
		Session            *mgo.Session
		mongoConnectionURI string
		mongoDBName        string
		mongoCollName      string
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

func (s *myOAuth2) New(appEnv *cfenv.App) *myOAuth2 {
	oauth2ServiceName := os.Getenv("OAUTH2_SERVICE_NAME")
	clientIdName := os.Getenv("OAUTH2_CLIENT_ID")
	clientSecretName := os.Getenv("OAUTH2_CLIENT_SECRET")
	oauth2Service, err := appEnv.Services.WithName(oauth2ServiceName)

	if err != nil {
		panic(fmt.Sprintf("oauth2 client service name error: %s", err.Error()))
	}
	s.ID = oauth2Service.Credentials[clientIdName]
	s.Secret = oauth2Service.Credentials[clientSecretName]
	return s
}

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
	s.mongoCollName = os.Getenv("MONGO_COLLECTION_NAME")
	mongoService, err := appEnv.Services.WithName(mongoServiceName)

	if err != nil {
		panic(fmt.Sprintf("mongodb service name error: %s", err.Error()))
	}
	s.mongoConnectionURI = mongoService.Credentials[mongoURIName]
	parsedURI := strings.Split(s.mongoConnectionURI, "/")
	s.mongoDBName = parsedURI[len(parsedURI)-1]
	s.connect()
	defer func() { go s.autoReconnect() }()
	return s
}

func (s *myMongo) connect() {
	var err error

	if s.Session, err = mgo.Dial(s.mongoConnectionURI); err != nil {
		panic(fmt.Sprintf("mongodb dial error: %s", err.Error()))
	}
	s.Session.SetMode(mgo.Monotonic, true)
	s.Col = s.Session.DB(s.mongoDBName).C(s.mongoCollName)
}

func (s *myMongo) autoReconnect() {
	for {

		if err := s.Session.Ping(); err != nil {
			fmt.Printf("mongodb connection lost... attempting to reconnect")
			s.Session.Close()
			s.connect()
		}
		time.Sleep(5000 * time.Millisecond)
	}
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
