package blog

import (
	"blog-api/middlewares"
	"net/http"
)

func (h *BlogHandler) RegisterBlogRoutes(prefix string, server *http.ServeMux) {
	// blog
	server.HandleFunc("GET "+prefix+"/", h.handleBlogIndex)
	server.HandleFunc("GET "+prefix+"/random", h.handleRandomBlog)
	server.HandleFunc("GET "+prefix+"/search/{query}", h.handleBlogSearch)
	server.HandleFunc("GET "+prefix+"/{slug}", h.handleBlogBySlug)
	server.HandleFunc("POST "+prefix+"/{id}/like", h.handleBlogLike)
	// blog categories
	server.HandleFunc("GET "+prefix+"/category/{category}", h.handleBlogsByCategory)
	// protected
	// set ID within request context - check chat gpt
	// ensure id from verified token matches the userID value in the req route path
	server.HandleFunc("GET "+prefix+"/drafts/{userID}", middlewares.BearerAuthMiddleware(h.handleDrafts))
}
