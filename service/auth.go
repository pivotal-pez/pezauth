package pezauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	i "github.com/xchapter7x/goutil/itertools"
	goauth2 "golang.org/x/oauth2"
)

//Constants to construct my oauth calls
const (
	ClientID       = "1083030294947-6g3bhhrgl3s7ul736jet625ajvp94f5p.apps.googleusercontent.com"
	ClientSecret   = "kfgM5mT3BqPQ84VeXsYokAK_"
	sessionName    = "pivotalpezauthservicesession"
	sessionSecret  = "shhh.donttellanyone"
	AuthFailStatus = 403
)

//Vars for my oauth calls
var (
	Scopes         = []string{"https://www.googleapis.com/auth/plus.me", "https://www.googleapis.com/auth/userinfo.email"}
	allowedDomains = []string{
		"pivotal.io",
	}
	oauthConfig *goauth2.Config
)

func isAllowedDomain(domain string) bool {
	o := i.Find(allowedDomains, func(pair i.Pair) bool {
		val := pair.Second.(string)
		return val == domain
	})
	return o == *new(i.Pair)
}

func cleanVersionFromURI(uri string) string {
	var digitsRegexp = regexp.MustCompile(`-.*?\.`)
	match := digitsRegexp.FindStringSubmatch(uri)

	if len(match) > 0 {
		newS := strings.Replace(uri, match[0], ".", -1)
		fmt.Println(newS)
		uri = newS
	}

	if !strings.HasPrefix(uri, "http") {
		uri = fmt.Sprintf("https://%s", uri)
	}
	return uri
}

func getAppEnv() (appEnv *cfenv.App) {
	var err error

	if appEnv, err = cfenv.Current(); err != nil {
		panic(err)
	}
	return
}

func getAppURI() string {
	appEnv := getAppEnv()
	return cleanVersionFromURI(appEnv.ApplicationURIs[0])
}

func DomainChecker(res http.ResponseWriter, tokens oauth2.Tokens) {
	userInfo := GetUserInfo(tokens)

	if tokens.Expired() || isAllowedDomain(userInfo["domain"].(string)) {
		statusCode := AuthFailStatus
		responseBody := `{"error": "not logged in as a valid user, or the access token is expired"}`
		res.WriteHeader(statusCode)
		res.Write([]byte(responseBody))
	}
}

var domainCheck = func() martini.Handler {
	return DomainChecker
}()

var GetUserInfo = func(tokens oauth2.Tokens) (userObject map[string]interface{}) {
	url := "https://www.googleapis.com/plus/v1/people/me"
	token := &goauth2.Token{
		AccessToken:  tokens.Access(),
		TokenType:    "Bearer",
		RefreshToken: tokens.Refresh(),
		Expiry:       tokens.ExpiryTime(),
	}
	client := oauthConfig.Client(goauth2.NoContext, token)
	resp, _ := client.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &userObject)
	return
}

func setOauthConfig() {
	oauthConfig = &goauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Scopes:       Scopes,
		RedirectURL:  fmt.Sprintf("%s/oauth2callback", getAppURI()),
	}
}

//InitAuth - initializes authentication middleware for controllers
func InitAuth(m *martini.ClassicMartini) {
	setOauthConfig()
	m.Use(render.Renderer())
	m.Use(sessions.Sessions(sessionName, sessions.NewCookieStore([]byte(sessionSecret))))
	m.Use(oauth2.Google(oauthConfig))
	m.Use(oauth2.LoginRequired)
	m.Use(domainCheck)
}
