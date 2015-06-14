package main

import (
	"fmt"
	"net/http"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/gorelic"
	"github.com/pivotalservices/pezauth/integrations"
	pez "github.com/pivotalservices/pezauth/service"
	"github.com/xchapter7x/cloudcontroller-client"
)

func main() {
	appEnv, _ := cfenv.Current()
	m := martini.Classic()
	newRelic := new(integrations.MyNewRelic).New(appEnv)
	gorelic.InitNewrelicAgent(newRelic.Key, newRelic.App, false)
	m.Use(gorelic.Handler)
	oauth2Client := new(integrations.MyOAuth2).New(appEnv)
	pez.ClientID = oauth2Client.ID
	pez.ClientSecret = oauth2Client.Secret
	rds := new(integrations.MyRedis).New(appEnv)
	defer rds.Conn.Close()
	pez.InitSession(m, &redisCreds{
		pass: rds.Pass,
		uri:  rds.URI,
	})
	h := new(integrations.MyHeritage).New(appEnv)
	heritageClient := &heritage{
		Client:   ccclient.New(h.LoginTarget, h.LoginUser, h.LoginPass, new(http.Client)),
		ccTarget: h.CCTarget,
	}
	mngo := new(integrations.MyMongo).New(appEnv)
	defer mngo.Session.Close()

	if _, err := heritageClient.Login(); err == nil {
		pez.InitRoutes(m, rds.Conn, pez.NewMongoCollectionWrapper(mngo.Col), heritageClient)
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
