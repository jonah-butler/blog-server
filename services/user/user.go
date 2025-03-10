package user

import (
	er "blog-api/repositories/email"
	r "blog-api/repositories/user"
	es "blog-api/services/email"
	"context"
	"errors"
	"time"
)

type UserService struct {
	userRepo             r.UserRepository
	passwordResetService es.PasswordResetService
}

func NewUserService(userRepo r.UserRepository, passwordResetService es.PasswordResetService) *UserService {
	return &UserService{
		userRepo:             userRepo,
		passwordResetService: passwordResetService,
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

func (s *UserService) UserResetPassword(ctx context.Context, payload r.UserResetPasswordPost) (*er.PasswordResetResponse, error) {
	response := &er.PasswordResetResponse{
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

	passwordResetMeta := &er.PasswordResetMeta{
		User:      user.ID,
		CreatedAt: time.Now(),
		Hash:      hash,
	}

	err = s.passwordResetService.CreatePasswordResetEntry(ctx, passwordResetMeta)
	if err != nil {
		return response, err
	}

	err = s.passwordResetService.SendEmail(user.Email, token)
	if err != nil {
		return response, err
	}

	return response, nil
}
