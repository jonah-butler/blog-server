package email

import (
	e "blog-api/repositories/email"
	"context"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type PasswordResetService struct {
	passwordResetRepo e.PasswordResetRepository
}

func NewPasswordResetService(repo e.PasswordResetRepository) *PasswordResetService {
	return &PasswordResetService{passwordResetRepo: repo}
}

func (s *PasswordResetService) CreatePasswordResetEntry(ctx context.Context, payload *e.PasswordResetMeta) error {
	err := s.passwordResetRepo.CreatePasswordResetEntry(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *PasswordResetService) SendEmail(email, token string) error {
	from := mail.NewEmail("The Daemon", "daemon@jonahbutler.dev")
	to := mail.NewEmail("Password Reset Requester", email)
	subject := "Blog Password Reset"

	url := "https://jonahbutler.dev/password-reset?token=" + token

	plainTextContent := "A request to reset your password was submitted.\n\n" +
		"If you did not make this request, ignore this email as someone may have accidentally typed your email address.\n\n" +
		"To update your password please visit: \n\n" +
		url + "\n\n"

	htmlContent := "<div><h3>A request to reset your password was submitted.</h3></div>" +
		"<div><strong>If you did not make this request, ignore this email as someone may have accidentally typed your email address.</strong></div>" +
		"<div><p>To update your password please visit:</p></div>" +
		"<div><a href='" + url + "'" + ">" + url + "</a></div"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
