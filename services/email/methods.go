package email

import (
	er "blog-api/repositories/email"
	"fmt"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func (s *PasswordResetService) PreparePasswordResetData(token, email string) (*er.SendgridPayload, error) {
	payload := new(er.SendgridPayload)
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

func (s *PasswordResetService) EvaluatedElapsedTime(timestamp time.Time, hours int) bool {
	now := time.Now()

	elapsed := now.Sub(timestamp)

	hoursElapsed := elapsed.Hours()

	return (hoursElapsed > 0 && hoursElapsed < float64(hours))
}
