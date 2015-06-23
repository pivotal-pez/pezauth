package pezauth

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"text/template"

	"github.com/cloudfoundry-community/go-cfenv"
)

//NewEmailServerFromService - construct email server from vCap Service
func NewEmailServerFromService(appEnv *cfenv.App) *EmailServer {
	serviceName := os.Getenv("SMTP_SERVICE_NAME")
	hostName := os.Getenv("SMTP_HOST")
	portName := os.Getenv("SMTP_PORT")
	userName := os.Getenv("SMTP_USERNAME")
	passName := os.Getenv("SMTP_PASSNAME")
	supportEmail := os.Getenv("SUPPORT_EMAIL")
	service, err := appEnv.Services.WithName(serviceName)
	if err != nil {
		panic(fmt.Sprintf("email service name error: %s", err.Error()))
	}
	auth := smtp.PlainAuth("", service.Credentials[userName], service.Credentials[passName], service.Credentials[hostName])
	port, err := strconv.Atoi(service.Credentials[portName])
	if err != nil {
		panic(fmt.Sprintf("The port for email server is not a valid integer %s", err.Error()))
	}
	return &EmailServer{
		host:         service.Credentials[hostName],
		port:         port,
		auth:         auth,
		sendMailFunc: DefaultSMTPSendEmail,
		supportEmail: service.Credentials[supportEmail],
	}
}

//DefaultSMTPSendEmail - This is the default SMTP server send email behavior
//There are some issue with the smtp ssl certificate
//Reimplementing the http://golang.org/src/net/smtp/smtp.go?s=7610:7688#L263
//Will switch back to the default smtp.SendMail function
func DefaultSMTPSendEmail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()

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

//GetSupportEmail - retrieve the support email address
func (emailServer *EmailServer) GetSupportEmail() string {
	return emailServer.supportEmail
}
