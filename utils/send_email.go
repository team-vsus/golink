package utils

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func SendEmail(e *email.Email) error {
	//auth := smtp.PlainAuth("", e.From, os.Getenv("SMTP_PW"), "smtp-mail.outlook.com")
	auth := LoginAuth("muazahmed766@outlook.com", os.Getenv("SMTP_PW"))

	fmt.Println("Authenticating")
	return e.Send("smtp-mail.outlook.com:587", auth)
	//return smtp.SendMail("smtp-mail.outlook.com:587", auth, e.From, e.To, e.Text)

}
