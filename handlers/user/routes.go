package user

import (
	authmiddleware "blog-api/middlewares/auth"
	"net/http"
)

func (h *UserHandler) RegisterUserRoutes(prefix string, server *http.ServeMux) {
	// get public user data
	server.HandleFunc("GET "+prefix+"/{username}", h.getUserPublic)

	// register a new user - not working atm
	server.HandleFunc("POST "+prefix+"/register", h.registerUser)

	// authenticate a user and generate an access token
	server.HandleFunc("POST "+prefix+"/login", h.userLogin)

	// initiative the password reset flow
	server.HandleFunc("POST "+prefix+"/reset-password", h.resetPassword)

	// password reset validation handler
	server.HandleFunc("POST "+prefix+"/validate-password-reset", h.validatePasswordReset)

	// PRIVATE: updates a user's associated profile data
	server.HandleFunc("POST "+prefix+"/{userID}", authmiddleware.BearerAuthMiddleware(h.updateUser))

	// looks up the provided user in the user's table and sends the email if found
	server.HandleFunc("POST "+prefix+"/send-email/{emailAddress}", h.sendEmail)
}
