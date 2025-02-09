package user

import (
	r "blog-api/repositories/user"
	s "blog-api/services/user"
	"encoding/json"
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

func (h *UserHandler) UserLogin(w http.ResponseWriter, req *http.Request) {
	var login r.UserLoginPost

	err := json.NewDecoder(req.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if login.Password == "" || login.Username == "" {
		http.Error(w, "Invalid login payload", http.StatusInternalServerError)
		return
	}

	h.userService.UserLogin(req.Context(), login)
}
