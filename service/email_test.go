package pezauth_test

import (
	"errors"
	"net/smtp"

	. "github.com/pivotalservices/pezauth/service"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Email", func() {
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
		// 		//emailServer = NewEmailServer("smtp.vchs.pivotal.io", 25, DefaultSMTPSendEmail)
		// 		auth = smtp.PlainAuth("", "test@gmail.com", "password", "smtp.gmail.com")
		// 		emailServer = NewEmailServer("smtp.gmail.com", 587, auth, DefaultSMTPSendEmail)
		// 	})
		//
		// 	It("Should send email without error", func() {
		// 		testEmail := &SMTPData{
		// 			From:    "from@gmail.com",
		// 			To:      "to@pivotal.io",
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
