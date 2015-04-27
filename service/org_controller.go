package pezauth

import (
	"log"

	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
)

//MeGetHandler - a get control handler for me requests
type (
	OrgGetHandler func(log *log.Logger, r render.Render, tokens oauth2.Tokens)
	OrgPutHandler func(log *log.Logger, r render.Render, tokens oauth2.Tokens)
)

//NewMeController - a controller for me requests
func NewOrgController() Controller {
	return new(orgController)
}

type orgController struct {
	Controller
}

//Get - get a get handler for org management
func (s *orgController) Get() interface{} {
	var handler OrgGetHandler = func(log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		userInfo := GetUserInfo(tokens)
		log.Println("getting userInfo: ", userInfo)
		genericResponseFormatter(r, "", userInfo, nil)
	}
	return handler
}

//Put - get a get handler for org management
func (s *orgController) Put() interface{} {
	var handler OrgPutHandler = func(log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		userInfo := GetUserInfo(tokens)
		log.Println("getting userInfo: ", userInfo)
		genericResponseFormatter(r, "", userInfo, nil)
	}
	return handler
}
