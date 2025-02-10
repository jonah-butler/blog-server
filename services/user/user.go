package user

import (
	r "blog-api/repositories/user"
	"context"
	"errors"
)

type UserService struct {
	userRepo r.UserRepository
}

func NewUserService(repo r.UserRepository) *UserService {
	return &UserService{userRepo: repo}
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
