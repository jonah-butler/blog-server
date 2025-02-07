package user

import "net/http"

func (h *UserHandler) RegisterUserRoutes(prefix string, server *http.ServeMux) {
	// user
	server.HandleFunc("GET"+prefix+"/{userID}", h.GetUser)
	server.HandleFunc("POST"+prefix+"/register", h.RegisterUser)
	server.HandleFunc("POST"+prefix+"/login", h.UserLogin)
}
