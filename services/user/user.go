package user

import (
	prr "blog-api/repositories/passwordreset"
	r "blog-api/repositories/user"
	es "blog-api/services/email"
	prs "blog-api/services/passwordreset"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo             r.UserRepository
	passwordResetService prs.PasswordResetService
	emailService         es.EmailService
}

func NewUserService(userRepo r.UserRepository, passwordResetService prs.PasswordResetService, emailService es.EmailService) *UserService {
	return &UserService{
		userRepo:             userRepo,
		passwordResetService: passwordResetService,
		emailService:         emailService,
	}
}

func (s *UserService) GetUser(ctx context.Context) error { return nil }

func (s *UserService) RegisterUser(ctx context.Context) error { return nil }

func (s *UserService) UserLogin(ctx context.Context, payload r.UserLoginPost) (r.UserResponse, error) {
	var userResponse r.UserResponse

	userWithPassword, err := s.userRepo.FindUser(ctx, payload)
	if err != nil || userWithPassword == nil {
		return userResponse, err
	}

	isMatch := r.ComparePasswords(userWithPassword.Password, payload.Password)
	if !isMatch {
		return userResponse, errors.New("invalid password")
	}

	user := r.ConvertToUser(userWithPassword)

	token, err := r.GenerateJWT(user)

	userResponse = r.ConvertToUserResponse(user, token)

	return userResponse, err
}

func (s *UserService) UserResetPassword(ctx context.Context, payload r.UserResetPasswordPost) (*prr.PasswordResetResponse, error) {
	response := &prr.PasswordResetResponse{
		Message: "If the provided email address exists in our system, you should receive an email soon!",
	}

	token, hash, err := generateToken()
	if err != nil {
		return response, err
	}

	user, err := s.userRepo.GetUserByEmail(ctx, *payload.Email)
	if err != nil {
		return response, err
	}

	passwordResetMeta := &prr.PasswordResetMeta{
		User:      user.ID,
		CreatedAt: time.Now(),
		Hash:      hash,
	}

	err = s.passwordResetService.CreatePasswordResetEntry(ctx, passwordResetMeta)
	if err != nil {
		return response, err
	}

	message, err := s.emailService.PreparePasswordResetData(token, user.Email)
	if err != nil {
		return response, err
	}

	err = s.passwordResetService.SendEmail(message)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s *UserService) ValidatePasswordReset(ctx context.Context, payload *r.UserNewPasswordPost) (bool, error) {
	hash := computeHash(payload.ResetToken)

	meta, err := s.passwordResetService.ValidatePasswordReset(ctx, hash)
	if err != nil || meta.Hash == "" {
		return false, err
	}

	// evaluate createAt and revoke stale token
	isFresh := s.emailService.EvaluatedElapsedTime(meta.CreatedAt, 1)
	if !isFresh {
		_, err := s.passwordResetService.DeletePasswordResetEntry(ctx, meta.Hash, meta.User)
		if err != nil {
			return false, fmt.Errorf("failed to delete password reset meta")
		}
		return false, fmt.Errorf("please resubmit your request to reset your password")
	}

	didUpdate, err := s.UpdateUserPassword(ctx, payload.Password, meta.User)
	if err != nil {
		return false, err
	}

	if didUpdate {
		_, err = s.passwordResetService.DeletePasswordResetEntry(ctx, meta.Hash, meta.User)
		if err != nil {
			return didUpdate, err
		}
	}

	return didUpdate, nil
}

func (s *UserService) UpdateUserPassword(ctx context.Context, password string, user bson.ObjectID) (bool, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	didUpdate, err := s.userRepo.UpdateUserPassword(ctx, string(hashedPassword), user)
	if err != nil {
		return false, err
	}

	return didUpdate, nil
}

func (s *UserService) SendEmailToUser(ctx context.Context, emailData *r.UserSendEmailPost) error {
	user, err := s.userRepo.GetUserByEmail(ctx, emailData.To)
	if err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("the provided email address is not a supported account")
	}

	email, err := s.emailService.PrepareContactEmail(emailData)
	if err != nil {
		return err
	}

	err = s.passwordResetService.SendEmail(email)
	if err != nil {
		return err
	}

	return nil
}
