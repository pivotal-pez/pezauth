package pezauth

import (
	"bytes"
	"log"
	"net/smtp"
	"strconv"
	"text/template"
)

//SMTPTemplate template to generate smtp data
const SMTPTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}
`

//SMTPData data typr for smtp email info
type SMTPData struct {
	From    string
	To      string
	Subject string
	Body    string
}

//EmailServer - email server pez auth use to send email
type EmailServer struct {
	host         string
	port         int
	auth         smtp.Auth
	sendMailFunc SendMailFunc
}

//DefaultSMTPSendEmail - This is the default SMTP server send email behavior
func DefaultSMTPSendEmail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return smtp.SendMail(addr, a, from, to, msg)
}

//SendMailFunc - Function to wrap the smtp SendMail behavior
type SendMailFunc func(string, smtp.Auth, string, []string, []byte) error

//NewEmailServer - Create an email server
func NewEmailServer(host string, port int, auth smtp.Auth, sendMailFunc SendMailFunc) *EmailServer {
	return &EmailServer{
		host:         host,
		port:         port,
		auth:         auth,
		sendMailFunc: sendMailFunc,
	}
}

//SendEmail - send email
func (emailServer *EmailServer) SendEmail(data *SMTPData) error {
	var doc bytes.Buffer
	t := template.New("emailTemplate")
	t, err := t.Parse(SMTPTemplate)
	if err != nil {
		log.Fatal("error trying to parse mail template", err)
		return err
	}
	err = t.Execute(&doc, data)
	if err != nil {
		log.Fatal("error tring to map data to the smtp email template", err)
		return err
	}
	err = emailServer.sendMailFunc(emailServer.host+":"+strconv.Itoa(emailServer.port),
		emailServer.auth,
		data.From,
		[]string{data.To},
		doc.Bytes())
	return err
}
