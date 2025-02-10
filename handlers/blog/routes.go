package blog

import "net/http"

func (h *BlogHandler) RegisterBlogRoutes(prefix string, server *http.ServeMux) {
	// blog
	server.HandleFunc("GET "+prefix+"/", h.handleBlogIndex)
	server.HandleFunc("GET "+prefix+"/random", h.handleRandomBlog)
	server.HandleFunc("GET "+prefix+"/{slug}", h.handleBlogBySlug)
	// blog categories
	server.HandleFunc("GET "+prefix+"/category/{category}", h.handleBlogsByCategory)
	// protected
	server.HandleFunc("GET "+prefix+"/drafts/{userID}", h.handleDrafts)
}
