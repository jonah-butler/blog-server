package user

import (
	r "blog-api/repositories/user"
	s "blog-api/services/user"
	u "blog-api/utilities"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHandler struct {
	userService *s.UserService
}

func NewUserHandler(service *s.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

func (h *UserHandler) getUser(w http.ResponseWriter, req *http.Request) {}

func (h *UserHandler) registerUser(w http.ResponseWriter, req *http.Request) {}

func (h *UserHandler) userLogin(w http.ResponseWriter, req *http.Request) {
	var login r.UserLoginPost

	err := json.NewDecoder(req.Body).Decode(&login)
	if err != nil {
		error := fmt.Errorf("failed to decode request body: %s", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	if login.Password == "" || login.Username == "" {
		error := fmt.Errorf("invalid login payload")
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	response, err := h.userService.UserLogin(req.Context(), login)
	if err != nil {
		error := fmt.Errorf("failed to login user: %s", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		error := fmt.Errorf("failed to marshall json: %s", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	u.WriteJSON(w, http.StatusOK, jsonData)
}

func (h *UserHandler) resetPassword(w http.ResponseWriter, req *http.Request) {
	var passwordReset r.UserResetPasswordPost

	err := json.NewDecoder(req.Body).Decode(&passwordReset)
	if err != nil {
		error := fmt.Errorf("failed to decode password reset payload: %s", err)
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	if passwordReset.Email == nil {
		error := fmt.Errorf("invalid reset password payload: %s", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	response, err := h.userService.UserResetPassword(req.Context(), passwordReset)
	if err != nil {
		error := fmt.Errorf("encountered an error during the password reset flow: %s", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	u.WriteJSON(w, http.StatusOK, response)
}
