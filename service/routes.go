package pezauth

import (
	"fmt"
	"log"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
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

type Response struct {
	User     map[string]interface{}
	ApiKey   string
	ErrorMsg string
}

//InitRoutes - initialize the mappings for controllers against valid routes
func InitRoutes(m *martini.ClassicMartini) {
	m.Use(render.Renderer())
	m.Use(martini.Static(StaticPath))

	m.Get("/info", FakeController)
	m.Get("/", func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		userInfo := GetUserInfo(tokens)
		r.HTML(200, "index", userInfo)
	})

	m.Group(URLAuthBaseV1, func(r martini.Router) {
		r.Put(APIKey, FakeController)    //this will re-generate a new key for the user or return an error if one does not yet exist
		r.Post(APIKey, FakeController)   //this will generate a key for the user or return an error if one already exists
		r.Get(APIKey, FakeController)    //will return the key for the username (pivotal.io email) it is given... this needs to be locked in that only the current user or admin will receive a result
		r.Delete(APIKey, FakeController) //this will remove the key from the user

		r.Get(APIKeys, FakeController) //will return the list of current keys and its associated user
	})
}
