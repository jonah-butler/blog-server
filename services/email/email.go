package email

import (
	ur "blog-api/repositories/user"
	"fmt"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) PreparePasswordResetData(token, email string) (*SendgridPayload, error) {
	payload := new(SendgridPayload)
	url := "https://jonahbutler.dev/password-reset?resetToken=" + token

	daemon, ok := os.LookupEnv("DAEMON_ADDRESS")
	if !ok {
		return payload, fmt.Errorf("no to address in environment")
	}

	payload.Subject = "Password Reset"

	payload.To = mail.NewEmail("Password Reset Requester", email)
	payload.From = mail.NewEmail("Password Daemon", daemon)

	payload.PlainText = "A request to reset your password was submitted.\n\n" +
		"If you did not make this request, ignore this email as someone may have accidentally typed your email address.\n\n" +
		"To update your password please visit: \n\n" +
		url + "\n\n"

	payload.HTMLText = "<div><h3>A request to reset your password was submitted.</h3></div>" +
		"<div><strong>If you did not make this request, ignore this email as someone may have accidentally typed your email address.</strong></div>" +
		"<div><p>To update your password please visit:</p></div>" +
		"<div><a href='" + url + "'" + ">" + url + "</a></div"

	return payload, nil
}

func (s *EmailService) PrepareContactEmail(emailData *ur.UserSendEmailPost) (*SendgridPayload, error) {
	payload := new(SendgridPayload)

	daemon, ok := os.LookupEnv("DAEMON_ADDRESS")
	if !ok {
		return payload, fmt.Errorf("no to address in environment")
	}

	payload.Subject = emailData.Subject

	payload.To = mail.NewEmail("Contact Requester", emailData.To)
	payload.From = mail.NewEmail("Contact Daemon", daemon)

	payload.PlainText = "This message was delivered on behalf of: " + emailData.From + "\n\n" +
		"EMAIL IS AS FOLLOWS:\n\n" +
		"--------------------\n\n" +
		emailData.Message + "\n\n" +
		"Respond to: " + emailData.From

	payload.HTMLText = "<div><h3>This message was delivered on behalf of: " + emailData.From + "</h3></div>" +
		"<div><p>EMAIL IS AS FOLLOWS:</p></div>" +
		"<div>--------------------</div>" +
		"<div><p>" + emailData.Message + "</p></div>" +
		"</br>" +
		"<div><strong>" + "Respond to: " + emailData.From + "</strong></div>"

	return payload, nil
}

func (s *EmailService) EvaluatedElapsedTime(timestamp time.Time, hours int) bool {
	now := time.Now()

	elapsed := now.Sub(timestamp)

	hoursElapsed := elapsed.Hours()

	return (hoursElapsed > 0 && hoursElapsed < float64(hours))
}
