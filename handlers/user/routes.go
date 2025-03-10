package user

import (
	"net/http"
)

func (h *UserHandler) RegisterUserRoutes(prefix string, server *http.ServeMux) {
	// user
	server.HandleFunc("GET "+prefix+"/{userID}", h.getUser)
	server.HandleFunc("POST "+prefix+"/register", h.registerUser)
	server.HandleFunc("POST "+prefix+"/login", h.userLogin)
	server.HandleFunc("POST "+prefix+"/reset-password", h.resetPassword)
}
