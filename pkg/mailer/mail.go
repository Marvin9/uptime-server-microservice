package mailer

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sendgrid/sendgrid-go"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	fromName  = "Server Monitor"
	fromEmail = "test@example.com"
	subject   = "Report"
)

func generateBody(forURL string, status int, at time.Time) string {
	body := fmt.Sprintf("For <b>%v</b>,<br /> Status at: %v <br />was <b>%v</b>", forURL, at, status)
	return body
}

// Mail is used to send emails
func Mail(to string, forURL string, status int, at time.Time) error {
	if strings.HasSuffix(os.Args[0], ".test") {
		return nil
	}
	log.Printf("Sending mail to %v\n\n", to)
	from := mail.NewEmail(fromName, fromEmail)
	toS := mail.NewEmail("", to)
	body := generateBody(forURL, status, at)
	message := mail.NewSingleEmail(from, "Report", toS, "Report.\n", body)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)

	if err != nil {
		log.Printf("Error sending mail %v\n\n.", err)
	}

	return err
}
