package pezauth

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/sessions"
	goauth2 "golang.org/x/oauth2"
)

//Constants to construct my oauth calls
const (
	ClientID      = "1083030294947-6g3bhhrgl3s7ul736jet625ajvp94f5p.apps.googleusercontent.com"
	ClientSecret  = "kfgM5mT3BqPQ84VeXsYokAK_"
	sessionName   = "randomSessionName"
	sessionSecret = "shhh.donttellanyone"
)

//Vars for my oauth calls
var (
	Scopes = []string{"https://www.googleapis.com/auth/plus.me"}
)

func cleanVersionFromURI(uri string) string {
	var digitsRegexp = regexp.MustCompile(`-.*?\.`)
	match := digitsRegexp.FindStringSubmatch(uri)

	if len(match) > 0 {
		newS := strings.Replace(uri, match[0], ".", -1)
		fmt.Println(newS)
		uri = newS
	}
	return uri
}

func getAppEnv() (appEnv *cfenv.App) {
	var (
		err error
	)

	switch os.Getenv("LOCAL") {
	case "true":
		appEnv = &cfenv.App{
			ApplicationURIs: []string{
				fmt.Sprintf("http://localhost-lkashdgaskhdglaskdhgasd:%s", os.Getenv("PORT")),
				fmt.Sprintf("http://localhost:%s", os.Getenv("PORT")),
			},
		}

	default:
		if appEnv, err = cfenv.Current(); err != nil {
			panic(err)
		}
	}
	return
}

func getAppURI() string {
	appEnv := getAppEnv()
	return cleanVersionFromURI(appEnv.ApplicationURIs[0])
}

//InitAuth - initializes authentication middleware for controllers
func InitAuth(m *martini.ClassicMartini) {
	m.Use(sessions.Sessions(sessionName, sessions.NewCookieStore([]byte(sessionSecret))))
	m.Use(oauth2.Google(
		&goauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			Scopes:       Scopes,
			RedirectURL:  fmt.Sprintf("%s/oauth2callback", getAppURI()),
		},
	))
	m.Use(oauth2.LoginRequired)
}
