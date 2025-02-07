package user

import (
	s "blog-api/services/user"
	"net/http"
)

type UserHandler struct {
	userService *s.UserService
}

func NewUserHandler(service *s.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, req *http.Request) {}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, req *http.Request) {}

func (h *UserHandler) UserLogin(w http.ResponseWriter, req *http.Request) {}
