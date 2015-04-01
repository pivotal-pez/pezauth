package pezauth

import (
	"fmt"
	"log"

	"github.com/go-martini/martini"
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

func fakeController(params martini.Params, log *log.Logger, r render.Render) {
	r.JSON(200, map[string]interface{}{"hello": "world"})
}

//InitRoutes - initialize the mappings for controllers against valid routes
func InitRoutes(m *martini.ClassicMartini) {
	m.Use(martini.Static(StaticPath))

	m.Group("/", func(r martini.Router) {
		r.Get("info", func() (int, string) {
			return 200, "auth service"
		})
	})

	m.Group(URLAuthBaseV1, func(r martini.Router) {
		r.Put(APIKey, fakeController)
		r.Post(APIKey, fakeController)
		r.Get(APIKey, fakeController)
		r.Delete(APIKey, fakeController)

		r.Get(APIKeys, fakeController)
	})
}
