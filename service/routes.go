package pezauth

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/pivotal-pez/pezauth/integrations"
	"github.com/pivotal-pez/pezdispenser/service"
)

//Constants to construct routes with
const (
	UserParam          = "user"
	APIVersion1        = "v1"
	AuthGroup          = "auth"
	OrgGroup           = "org"
	APIKeys            = "/api-keys"
	ValidKeyCheck      = "/valid-key"
	StaticPath         = "public"
	InventoryItemParam = "invitem"
)

//formatted strings based on constants, to be used in URLs
var (
	APIKey        = fmt.Sprintf("/api-key/:%s", UserParam)
	OrgUser       = fmt.Sprintf("/user/:%s", UserParam)
	URLAuthBaseV1 = fmt.Sprintf("/%s/%s", APIVersion1, AuthGroup)
	URLOrgBaseV1  = fmt.Sprintf("/%s/%s", APIVersion1, OrgGroup)
	LeaseURL      = fmt.Sprintf("/pcfaas/inventory/:%s", InventoryItemParam)
)

var displayNewServices = strings.ToUpper(os.Getenv("DISPLAY_NEW_SERVICES")) == "YES"

//InitRoutes - initialize the mappings for controllers against valid routes
func InitRoutes(m *martini.ClassicMartini, redisConn func() Doer, mongoConn pezdispenser.MongoCollectionGetter, authClient AuthRequestCreator, invClient *integrations.MyInventoryClient) {
	setOauthConfig()
	keyGen := NewKeyGen(redisConn, &GUIDMake{})
	m.Use(render.Renderer())
	m.Use(martini.Static(StaticPath))
	m.Use(oauth2.Google(OauthConfig))
	authKey := NewAuthKeyV1(keyGen)

	m.Get("/info", authKey.Get())
	m.Get(ValidKeyCheck, NewValidateV1(keyGen).Get())

	m.Get("/me", oauth2.LoginRequired, DomainCheck, NewMeController().Get())
	m.Get("/pcfaas/inventory", oauth2.LoginRequired, DomainCheck, NewPcfaasController(invClient).Get())
	m.Post(LeaseURL, oauth2.LoginRequired, DomainCheck, NewPcfaasController(invClient).Post())

	m.Get("/", oauth2.LoginRequired, DomainCheck, func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		userInfo := GetUserInfo(tokens)
		if displayNewServices {
			r.HTML(SuccessStatus, "index_newservices", userInfo)
		} else {
			r.HTML(SuccessStatus, "index", userInfo)
		}
	})

	m.Post("/sandbox", oauth2.LoginRequired, DomainCheck, NewSandBoxController().Post())

	m.Group(URLAuthBaseV1, func(r martini.Router) {
		r.Put(APIKey, authKey.Put())
		r.Get(APIKey, authKey.Get())
		r.Delete(APIKey, authKey.Delete())
	}, oauth2.LoginRequired, DomainCheck)

	m.Group(URLOrgBaseV1, func(r martini.Router) {
		pcfOrg := NewOrgController(mongoConn, authClient)
		r.Put(OrgUser, pcfOrg.Put())
		r.Get(OrgUser, pcfOrg.Get())
	}, oauth2.LoginRequired, DomainCheck)
}
