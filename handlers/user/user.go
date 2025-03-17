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

	u.WriteJSON(w, http.StatusOK, response)
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

func (h *UserHandler) validatePasswordReset(w http.ResponseWriter, req *http.Request) {
	var payload r.UserNewPasswordPost

	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		error := fmt.Errorf("failed to decode request: %s", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	if payload.Password != payload.PasswordVerification || payload.Password == "" {
		error := fmt.Errorf("passwords must match")
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	if payload.ResetToken == "" {
		error := fmt.Errorf("missing reset token data")
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	didUpdate, err := h.userService.ValidatePasswordReset(req.Context(), &payload)
	if err != nil {
		error := fmt.Errorf("failed to validate password update request: %s", err)
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	if !didUpdate {
		error := fmt.Errorf("password did not update")
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	response := &r.UserPasswordResetResponse{
		DidUpdate: didUpdate,
	}

	u.WriteJSON(w, http.StatusOK, response)
}

func (h *UserHandler) sendEmail(w http.ResponseWriter, req *http.Request) {
	toAddress := req.PathValue("emailAddress")
	if !isValidEmail(toAddress) {
		error := fmt.Errorf("the provided to address is not a valid email")
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	emailData := new(r.UserSendEmailPost)

	if err := json.NewDecoder(req.Body).Decode(emailData); err != nil {
		error := fmt.Errorf("failed to decode email payload: %s", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	if !isValidEmail(emailData.To) || !isValidEmail(emailData.From) {
		error := fmt.Errorf("the emails provided in the payload are invalid")
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	if emailData.To != toAddress {
		error := fmt.Errorf("both TO addresses must match")
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	err := h.userService.SendEmailToUser(req.Context(), emailData)
	if err != nil {
		error := fmt.Errorf("failed to send email %s", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
	}

	u.WriteJSON(w, http.StatusOK, u.EmptyResponse())
}
