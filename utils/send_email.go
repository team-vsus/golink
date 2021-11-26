package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(from string, to []string, message []byte) {
	// Authentication.
	auth := smtp.PlainAuth("", from, os.Getenv("SMTP_PW"), "smtp.gmail.com")

	// Sending email.
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}

}
