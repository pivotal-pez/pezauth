package pezauth

import (
	"log"
	"net/http"
	"net/smtp"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/pivotal-pez/pezdispenser/cloudfoundryclient"
	"github.com/pivotal-pez/pezdispenser/service"
	"github.com/xchapter7x/cloudcontroller-client"
)

type (
	//AuthPutHandler - auth control handler for put calls
	AuthPutHandler func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
	//AuthPostHandler - auth control handler for post calls
	AuthPostHandler func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
	//AuthGetHandler - auth control handler for get calls
	AuthGetHandler func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
	//AuthDeleteHandler - auth control handler for delete calls
	AuthDeleteHandler func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
	//Controller - interface of a base controller
	Controller interface {
		Put() interface{}
		Post() interface{}
		Get() interface{}
		Delete() interface{}
	}
	//AuthRequestCreator - interface to an object which can decorate a request with auth tokens
	AuthRequestCreator interface {
		CreateAuthRequest(verb, requestURL, path string, args interface{}) (*http.Request, error)
		CCTarget() string
		HttpClient() ccclient.ClientDoer
		Login() (*ccclient.Client, error)
	}
	//GUIDMaker - interface for a guid maker
	GUIDMaker interface {
		Create() string
	}
	//GUIDMake - struct for making guids
	GUIDMake struct {
	}
	//KeyGen - and implementation of the KeyGenerator interface
	KeyGen struct {
		store     Doer
		guidMaker GUIDMaker
	}
	//KeyGenerator - interface to work with apikeys
	KeyGenerator interface {
		Get(user string) (string, error)
		GetByKey(key string) (hash string, val interface{}, err error)
		Create(user, details string) error
		Delete(user string) error
	}

	//Doer - interface to make a call to pezdispenser.Persistence store
	Doer interface {
		Do(commandName string, args ...interface{}) (reply interface{}, err error)
	}

	//MeGetHandler - a get control handler for me requests
	MeGetHandler func(log *log.Logger, r render.Render, tokens oauth2.Tokens)

	meController struct {
		Controller
	}

	//OrgGetHandler - func signature of org get handler
	OrgGetHandler func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
	//OrgPutHandler - func signature of org put handler
	OrgPutHandler func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)

	//PivotOrg - struct for pivot org record
	PivotOrg struct {
		Email   string
		OrgName string
		OrgGUID string
	}
	orgController struct {
		Controller
		store      func() pezdispenser.Persistence
		authClient AuthRequestCreator
	}

	sandBoxController struct {
		Controller
	}

	//Response - generic response object
	Response struct {
		Payload  interface{}
		APIKey   string
		ErrorMsg string
	}

	//UserMatch - an object used to check if a user is updating the records on a user key they are able to access
	UserMatch struct {
		userInfo    map[string]interface{}
		username    string
		successFunc func()
		failFunc    func()
	}

	//ValidateGetHandler - a type of handler for validation get endpoints
	ValidateGetHandler func(log *log.Logger, r render.Render, req *http.Request)

	validateV1 struct {
		Controller
		keyGenerator KeyGenerator
	}

	redisCreds interface {
		Pass() string
		Uri() string
	}

	authKeyV1 struct {
		Controller
		keyGen KeyGenerator
	}
	//OrgManager - interface to the org creation functionality
	OrgManager interface {
		Show() (result *PivotOrg, err error)
		SafeCreate() (record *PivotOrg, err error)
	}
	orgManager struct {
		username string
		userGUID string
		log      *log.Logger
		tokens   oauth2.Tokens
		store    pezdispenser.Persistence
		cfClient cloudfoundryclient.CloudFoundryClient
		apiInfo  map[string]interface{}
	}

	//SMTPData data typr for smtp email info
	SMTPData struct {
		From    string
		To      string
		Subject string
		Body    string
	}

	//EmailServer - email server pez auth use to send email
	EmailServer struct {
		host         string
		port         int
		auth         smtp.Auth
		sendMailFunc SendMailFunc
		supportEmail string //TODO maybe make this as an independent environment variable
	}

	//Sender - the interface that can send email
	Sender interface {
		SendEmail(data *SMTPData) error
		GetSupportEmail() string
	}
)
