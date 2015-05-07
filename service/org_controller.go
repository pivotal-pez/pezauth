package pezauth

import (
	"errors"
	"log"

	"github.com/fatih/structs"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
)

const (
	EmailFieldName = "email"
)

var (
	ErrNoMatchInStore      = errors.New("Could not find a matching user org or connection failure")
	ErrCanNotCreateOrg     = errors.New("Could not create a new org")
	ErrCanNotAddOrgRec     = errors.New("Could not add a new org record")
	ErrCantCallAcrossUsers = errors.New("user calling another users endpoint")
)

type (
	//OrgGetHandler - func signature of org get handler
	OrgGetHandler func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
	//OrgPutHandler - func signature of org put handler
	OrgPutHandler func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
)

type (
	persistence interface {
		FindOne(query interface{}, result interface{}) (err error)
		Upsert(selector interface{}, update interface{}) (err error)
	}
	//PivotOrg - struct for pivot org record
	PivotOrg struct {
		Email   string
		OrgName string
		OrgGuid string
	}
	orgController struct {
		Controller
		store      persistence
		authClient authRequestCreator
	}
	//APIResponse - cc http response object
	APIResponse struct {
		Metadata APIMetadata            `json:"metadata"`
		Entity   map[string]interface{} `json:"entity"`
	}
	//APIMetadata = cc http response metadata
	APIMetadata struct {
		Guid      string `json:"guid"`
		URL       string `json:"url"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)

//NewMeController - a controller for me requests
func NewOrgController(c persistence, authClient authRequestCreator) Controller {
	return &orgController{
		store:      c,
		authClient: authClient,
	}
}

//Get - get a get handler for org management
func (s *orgController) Get() interface{} {
	var handler OrgGetHandler = func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		username := params[UserParam]
		org := newOrg(username, log, tokens, s.store, s.authClient)
		result, err := org.Show()
		genericResponseFormatter(r, "", structs.Map(result), err)
	}
	return handler
}

//Put - get a get handler for org management
func (s *orgController) Put() interface{} {
	var handler OrgPutHandler = func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		var (
			err             error
			payload         *PivotOrg
			responsePayload map[string]interface{}
		)
		username := params[UserParam]
		org := newOrg(username, log, tokens, s.store, s.authClient)

		if _, err = org.Show(); err == ErrNoMatchInStore {
			payload, err = org.Create()
			responsePayload = structs.Map(payload)

		} else {
			err = ErrCanNotCreateOrg
		}
		genericResponseFormatter(r, "", responsePayload, err)
	}
	return handler
}
