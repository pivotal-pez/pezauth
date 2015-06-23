package pezauth

import (
	"fmt"
	"net/http"

	"github.com/martini-contrib/render"
)

//SUBJECT - Email subject for sandbox request
const SUBJECT = "Pez Request: Sandbox"

//BODY - Email body for sandbox request
const BODY = `Here is %s:
I would like to request a sandbox environment
My email address is: %s
Thanks
`

//NewSandBoxController - Create a Sandbox controller instance
func NewSandBoxController() Controller {
	return &sandBoxController{}
}

//SandBoxPostHandler Post Email send
type SandBoxPostHandler func(render.Render, *http.Request, Sender)

//Post - Post a sandbox request
func (e *sandBoxController) Post() interface{} {
	var handler SandBoxPostHandler = func(r render.Render, request *http.Request, emailServer Sender) {
		to := emailServer.GetSupportEmail()
		from, name := request.FormValue("from"), request.FormValue("name")
		emailData := &SMTPData{
			From:    from,
			To:      to,
			Subject: SUBJECT,
			Body:    fmt.Sprintf(BODY, name, from),
		}
		err := emailServer.SendEmail(emailData)
		genericResponseFormatter(r, "", map[string]interface{}{}, err)
	}
	return handler
}
