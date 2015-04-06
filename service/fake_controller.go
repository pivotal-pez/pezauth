package pezauth

import (
	"encoding/json"
	"log"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
)

func FakeController(params martini.Params, log *log.Logger, r render.Render, tokens oauth2.Tokens) {
	var (
		statusCode int
		err        error
		res        Response
	)
	userInfo := GetUserInfo(tokens)

	if _, err = json.Marshal(userInfo); err != nil {
		statusCode = 403
		res = Response{
			ErrorMsg: err.Error(),
		}

	} else {
		res = Response{
			ApiKey: "12345",
			User:   userInfo,
		}
	}
	r.JSON(statusCode, res)
}
