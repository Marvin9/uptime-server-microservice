package mailer

import (
	"fmt"
	"net/smtp"
	"os"
)

var from = "mayursiinh@gmail.com"
var password = os.Getenv("EMAIL_PASSWORD")
var smtpHost = "smtp.gmail.com"
var smtpPort = "587"

// Mail is used to send emails
func Mail(to []string, message string) error {
	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))

	if err != nil {
		fmt.Printf("\nError sending mail: %v", err)
		return err
	}

	return nil
}
