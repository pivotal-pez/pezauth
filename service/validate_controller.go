package pezauth

import (
	"log"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

//ValidateGetHandler - a type of handler for validation get endpoints
type (
	ValidateGetHandler func(params martini.Params, log *log.Logger, r render.Render)
)

//NewValidateV1 - create a validation controller
func NewValidateV1() Controller {
	return new(validateV1)
}

type validateV1 struct {
	Controller
}

func (s *validateV1) Get() interface{} {
	var handler ValidateGetHandler = func(params martini.Params, log *log.Logger, r render.Render) {

	}
	return handler
}
