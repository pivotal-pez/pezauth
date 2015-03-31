package pezauth

import (
	"fmt"

	"github.com/go-martini/martini"
)

//Constants to construct routes with
const (
	APIVersion1 = "v1"
	AuthGroup   = "auth"
)

//formatted strings based on constants, to be used in URLs
var (
	URLAuthBaseV1 = fmt.Sprintf("/%s/%s", APIVersion1, AuthGroup)
	APIKey        = "/api-key/:user"
	APIKeys       = "/api-keys"
)

func fakeController(params martini.Params) (int, string) {
	return 200, "success"
}

//InitRoutes - initialize the mappings for controllers against valid routes
func InitRoutes(m *martini.ClassicMartini) {
	m.Group("/", func(r martini.Router) {
		r.Get("info", func() string {
			return "auth service"
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
