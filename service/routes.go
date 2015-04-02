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

	m.Group("/", func(r martini.Router) {
		r.Get("info", func() (int, string) {
			return 200, "auth service"
		})
	})

	m.Group(URLAuthBaseV1, func(r martini.Router) {
		r.Put(APIKey, fakeController)    //this will re-generate a new key for the user or return an error if one does not yet exist
		r.Post(APIKey, fakeController)   //this will generate a key for the user or return an error if one already exists
		r.Get(APIKey, fakeController)    //will return the key for the username (pivotal.io email) it is given... this needs to be locked in that only the current user or admin will receive a result
		r.Delete(APIKey, fakeController) //this will remove the key from the user

		r.Get(APIKeys, fakeController) //will return the list of current keys and its associated user
	})
}
