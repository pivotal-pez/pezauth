package pezauth

import (
	"log"

	"github.com/martini-contrib/render"
)

//ValidateGetHandler - a type of handler for validation get endpoints
type (
	ValidateGetHandler func(log *log.Logger, r render.Render)
)

//NewValidateV1 - create a validation controller
func NewValidateV1() Controller {
	return new(validateV1)
}

type validateV1 struct {
	Controller
}

func (s *validateV1) Get() interface{} {
	var handler ValidateGetHandler = func(log *log.Logger, r render.Render) {
		r.JSON(200, Response{})
	}
	return handler
}
