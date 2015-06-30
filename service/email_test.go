package pezauth_test

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/pivotal-pez/pezauth/service"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Email", func() {
	Describe("NewEmailServerFromService", func() {

		Context("when given valid smtp credentials", func() {
			os.Setenv("SMTP_SERVICE_NAME", "email-server-service")
			os.Setenv("SMTP_HOST", "smtp-host")
			os.Setenv("SMTP_PORT", "smtp-port")
			os.Setenv("SUPPORT_EMAIL", "support-email")
			validEnv := []string{
				`VCAP_APPLICATION={}`,
				fmt.Sprintf("VCAP_SERVICES=%s", `{    
          "user-provided": [   
            {
                "name": "email-server-service",
                "label": "user-provided",
                "tags": [],
                "credentials": {
                  "smtp-host": "smtp.test.com",
                  "smtp-port": "25",
                  "support-email": "someone@gmail.com"
                },
                "syslog_drain_url": ""
              }
            ]
          }`),
			}

			testEnv := cfenv.Env(validEnv)
			appEnv, _ := cfenv.New(testEnv)

			It("should not panic", func() {
				Ω(func() { NewEmailServerFromService(appEnv) }).ShouldNot(Panic())
			})
		})
	})

	Describe("SendEmail", func() {
		var (
			address, from string
			to            []string
			msg           []byte
			auth          smtp.Auth
			hasErr        bool
			smtpData      *SMTPData = &SMTPData{
				From:    "test@pivotal.io",
				To:      "to@pivotal.io",
				Subject: "This is an email",
				Body:    "This is the body",
			}
			argumentCatcherFunc SendMailFunc = func(addr string, a smtp.Auth, f string, t []string, m []byte) (err error) {
				address = addr
				auth = a
				from = f
				to = t
				msg = m
				if hasErr {
					return errors.New("")
				}
				return
			}
			emailServer *EmailServer
		)

		Context("Valid email server", func() {

			BeforeEach(func() {
				hasErr = false
				emailServer = NewEmailServer("localhost", 25, nil, argumentCatcherFunc)
			})

			It("Should send email without error", func() {
				err := emailServer.SendEmail(smtpData)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Should get correct from, to email address and body", func() {
				emailServer.SendEmail(smtpData)
				Ω(address).Should(Equal("localhost:25"))
				Ω(from).Should(Equal("test@pivotal.io"))
				Ω(to[0]).Should(Equal("to@pivotal.io"))
				expectString := `From: test@pivotal.io
To: to@pivotal.io
Subject: This is an email

This is the body
`
				Ω(string(msg)).Should(Equal(expectString))
			})
		})

		Context("Invalid email server", func() {
			BeforeEach(func() {
				hasErr = true
				emailServer = NewEmailServer("localhost", 25, nil, argumentCatcherFunc)
			})
			It("Should send email with error", func() {
				err := emailServer.SendEmail(smtpData)
				Expect(err).To(HaveOccurred())
			})

		})

		// Context("Integration with real email server", func() {
		// 	BeforeEach(func() {
		// 		hasErr = true
		// 		auth = smtp.PlainAuth("", "", "", "smtp.gmail.com")
		// 		emailServer = NewEmailServer("smtp.vchs.gopivotal.com", 25, auth, DefaultSMTPSendEmail)
		// 		//emailServer = NewEmailServer("smtp.gmail.com", 587, auth, DefaultSMTPSendEmail)
		// 	})
		//
		// 	It("Should send email without error", func() {
		// 		testEmail := &SMTPData{
		// 			From:    "sding@pivotal.io",
		// 			To:      "dsz0111@gmail.com",
		// 			Subject: "This is a test email",
		// 			Body:    "This is the body",
		// 		}
		// 		err := emailServer.SendEmail(testEmail)
		// 		Expect(err).NotTo(HaveOccurred())
		// 	})
		//
		// })

	})
})
