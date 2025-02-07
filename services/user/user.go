package user

import (
	r "blog-api/repositories/user"
	"context"
)

type UserService struct {
	userRepo r.UserRepository
}

func NewUserService(repo r.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) GetUser(ctx context.Context) error { return nil }

func (s *UserService) RegisterUser(ctx context.Context) error { return nil }

func (s *UserService) UserLogin(ctx context.Context) error { return nil }
