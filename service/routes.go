package pezauth

import (
	"fmt"
	"log"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	goauth2 "golang.org/x/oauth2"
)

//Constants to construct routes with
const (
	APIVersion1 = "v1"
	AuthGroup   = "auth"
	APIKey      = "/api-key/:user"
	APIKeys     = "/api-keys"
	StaticPath  = "public"
)

//formatted strings based on constants, to be used in URLs
var (
	URLAuthBaseV1 = fmt.Sprintf("/%s/%s", APIVersion1, AuthGroup)
)

func fakeController(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
	statusCode := 200
	responseBody := map[string]interface{}{"hello": "world"}

	if tokens.Expired() {
		statusCode = 403
		responseBody = map[string]interface{}{"hello": "not logged in, or the access token is expired"}
	}
	r.JSON(statusCode, responseBody)
}

//InitRoutes - initialize the mappings for controllers against valid routes
func InitRoutes(m *martini.ClassicMartini) {
	m.Use(render.Renderer())
	m.Use(martini.Static(StaticPath))

	m.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte("secret123"))))
	m.Use(oauth2.Google(
		&goauth2.Config{
			ClientID:     "1083030294947-6g3bhhrgl3s7ul736jet625ajvp94f5p.apps.googleusercontent.com",
			ClientSecret: "kfgM5mT3BqPQ84VeXsYokAK_",
			Scopes:       []string{"https://www.googleapis.com/auth/plus.me"},
			RedirectURL:  "http://localhost:3000/oauth2callback",
		},
	))

	m.Group("/", func(r martini.Router) {
		r.Get("info", oauth2.LoginRequired, func() (int, string) {
			return 200, "auth service"
		})
	})

	m.Group(URLAuthBaseV1, func(r martini.Router) {
		r.Put(APIKey, oauth2.LoginRequired, fakeController)    //this will re-generate a new key for the user or return an error if one does not yet exist
		r.Post(APIKey, oauth2.LoginRequired, fakeController)   //this will generate a key for the user or return an error if one already exists
		r.Get(APIKey, oauth2.LoginRequired, fakeController)    //will return the key for the username (pivotal.io email) it is given... this needs to be locked in that only the current user or admin will receive a result
		r.Delete(APIKey, oauth2.LoginRequired, fakeController) //this will remove the key from the user

		r.Get(APIKeys, oauth2.LoginRequired, fakeController) //will return the list of current keys and its associated user
	})
}
