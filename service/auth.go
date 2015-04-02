package pezauth

import (
	"fmt"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/sessions"
	goauth2 "golang.org/x/oauth2"
)

const (
	ClientID      = "1083030294947-6g3bhhrgl3s7ul736jet625ajvp94f5p.apps.googleusercontent.com"
	ClientSecret  = "kfgM5mT3BqPQ84VeXsYokAK_"
	sessionName   = "randomSessionName"
	sessionSecret = "shhh.donttellanyone"
)

var (
	Scopes = []string{"https://www.googleapis.com/auth/plus.me"}
)

func getAppEnv() (appEnv *cfenv.App) {
	var (
		err error
	)

	switch os.Getenv("LOCAL") {
	case "true":
		appEnv = &cfenv.App{ApplicationURIs: []string{fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))}}

	default:
		if appEnv, err = cfenv.Current(); err != nil {
			panic(err)
		}
	}
	return
}

//InitAuth - initializes authentication middleware for controllers
func InitAuth(m *martini.ClassicMartini) {
	appEnv := getAppEnv()
	m.Use(sessions.Sessions(sessionName, sessions.NewCookieStore([]byte(sessionSecret))))
	m.Use(oauth2.Google(
		&goauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			Scopes:       Scopes,
			RedirectURL:  fmt.Sprintf("%s/oauth2callback", appEnv.ApplicationURIs[0]),
		},
	))
	m.Use(oauth2.LoginRequired)
}
