package passwordreset

import (
	prr "blog-api/repositories/passwordreset"
	"context"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PasswordResetService struct {
	passwordResetRepo prr.PasswordResetRepository
}

func NewPasswordResetService(repo prr.PasswordResetRepository) *PasswordResetService {
	return &PasswordResetService{passwordResetRepo: repo}
}

func (s *PasswordResetService) CreatePasswordResetEntry(ctx context.Context, payload *prr.PasswordResetMeta) error {
	err := s.passwordResetRepo.CreatePasswordResetEntry(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *PasswordResetService) ValidatePasswordReset(ctx context.Context, hash string) (*prr.PasswordResetMeta, error) {
	meta, err := s.passwordResetRepo.ValidatePasswordReset(ctx, hash)
	if err != nil {
		return meta, err
	}

	return meta, nil
}

func (s *PasswordResetService) DeletePasswordResetEntry(ctx context.Context, hash string, user bson.ObjectID) (bool, error) {
	return s.passwordResetRepo.DeletePasswordResetEntry(ctx, hash, user)
}

func (s *PasswordResetService) SendEmail(payload *prr.SendgridPayload) error {
	message := mail.NewSingleEmail(payload.From, payload.Subject, payload.To, payload.PlainText, payload.HTMLText)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
