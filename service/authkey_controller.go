package pezauth

import (
	"encoding/json"
	"log"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
)

//Authentication Handler function type definitions
type (
	AuthPutHandler    func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
	AuthPostHandler   func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
	AuthGetHandler    func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
	AuthDeleteHandler func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens)
)

//Controller - interface of a base controller
type Controller interface {
	Put() interface{}
	Post() interface{}
	Get() interface{}
	Delete() interface{}
}

//NewAuthKeyV1 - get an instance of a V1 authkey controller
func NewAuthKeyV1(kg KeyGenerator) Controller {
	return &authKeyV1{
		keyGen: kg,
	}
}

type authKeyV1 struct {
	keyGen KeyGenerator
}

//Put - get a put handler for authkeyv1
func (s *authKeyV1) Put() interface{} {
	var handler AuthPutHandler = func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		username := params[UserParam]
		userInfo := GetUserInfo(tokens)
		s.keyGen.Delete(username)
		s.keyGen.Create(username)
		apikey, err := s.keyGen.Get(username)
		genericResponseFormatter(r, apikey, userInfo, err)
	}
	return handler
}

//Post - get a post handler for authkeyv1
func (s *authKeyV1) Post() interface{} {
	var handler AuthPostHandler = func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		username := params[UserParam]
		userInfo := GetUserInfo(tokens)
		s.keyGen.Delete(username)
		s.keyGen.Create(username)
		apikey, err := s.keyGen.Get(username)
		genericResponseFormatter(r, apikey, userInfo, err)
	}
	return handler
}

//Get - get a get handler for authkeyv1
func (s *authKeyV1) Get() interface{} {
	var handler AuthGetHandler = func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		username := params[UserParam]
		userInfo := GetUserInfo(tokens)
		apikey, err := s.keyGen.Get(username)
		genericResponseFormatter(r, apikey, userInfo, err)
	}
	return handler
}

//Delete - get a delete handler for authkeyv1
func (s *authKeyV1) Delete() interface{} {
	var handler AuthDeleteHandler = func(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
		username := params[UserParam]
		userInfo := GetUserInfo(tokens)
		err := s.keyGen.Delete(username)
		genericResponseFormatter(r, "", userInfo, err)
	}
	return handler
}

func genericResponseFormatter(r render.Render, apikey string, userInfo map[string]interface{}, extErr error) {
	var (
		statusCode int
		err        error
		res        Response
	)

	if extErr != nil {
		statusCode = 403
		res = Response{
			ErrorMsg: extErr.Error(),
		}

	} else {

		if _, err = json.Marshal(userInfo); err != nil {
			statusCode = 403
			res = Response{
				ErrorMsg: err.Error(),
			}

		} else {
			statusCode = 200
			res = Response{
				ApiKey: apikey,
				User:   userInfo,
			}
		}
	}
	r.JSON(statusCode, res)
}
