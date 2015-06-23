package pezauth_test

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/pezauth/service"
)

var data *SMTPData

type FakeSender struct {
}

var sendEmailSucceed = true

var request = &http.Request{}

func (sender *FakeSender) SendEmail(smtpData *SMTPData) error {
	if !sendEmailSucceed {
		return errors.New("")
	}
	data = smtpData
	return nil
}

func (sender *FakeSender) GetSupportEmail() string {
	return "to@pivotal.io"
}

var _ = Describe("Request Sandbox", func() {
	Describe("Through email", func() {
		var (
			fromEmail              = "from@pivotal.io"
			toEmail                = "to@pivotal.io"
			name                   = "First, Last"
			fakeSender *FakeSender = &FakeSender{}
			render     *mockRenderer
			fn         SandBoxPostHandler = NewSandBoxController().Post().(SandBoxPostHandler)
		)
		Context("Send email failed", func() {
			BeforeEach(func() {
				request.Form = map[string][]string{"From": {fromEmail}, "Name": {name}}
				sendEmailSucceed = false
				render = new(mockRenderer)
			})

			It("Should request the sandbox with error", func() {
				fn(render, request, fakeSender)
				Ω(render.StatusCode).Should(Equal(FailureStatus))
			})

		})

		Context("Send email succeed", func() {
			BeforeEach(func() {
				request.Form = map[string][]string{"from": {fromEmail}, "name": {name}}
				sendEmailSucceed = true
				render = new(mockRenderer)
				os.Setenv("SANDBOX_SUPPORT_EMAIL", toEmail)
			})

			It("Should request the sandbox with success", func() {
				fn(render, request, fakeSender)
				Ω(render.StatusCode).Should(Equal(SuccessStatus))
			})
			It("Should send correct email", func() {
				fn(render, request, fakeSender)
				Ω(data.From).Should(Equal(fromEmail))
				Ω(data.To).Should(Equal(toEmail))
				Ω(data.Subject).Should(Equal(SUBJECT))
				Ω(data.Body).Should(Equal(fmt.Sprintf(BODY, name, fromEmail)))
			})

		})
	})
})
