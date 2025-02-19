package blog

import (
	"blog-api/middlewares"
	"net/http"
)

func (h *BlogHandler) RegisterBlogRoutes(prefix string, server *http.ServeMux) {
	//
	// blog
	//

	// get latest blogs
	server.HandleFunc("GET "+prefix+"/", h.handleBlogIndex)
	// get random blog
	server.HandleFunc("GET "+prefix+"/random", h.handleRandomBlog)
	// search blogs
	server.HandleFunc("GET "+prefix+"/search/{query}", h.handleBlogSearch)
	// lookup blog by slug
	server.HandleFunc("GET "+prefix+"/{slug}", h.handleBlogBySlug)
	// update blog rating
	server.HandleFunc("POST "+prefix+"/{id}/like", h.handleBlogLike)
	// update blog
	server.HandleFunc("POST "+prefix+"/{id}/edit/{userID}", middlewares.BearerAuthMiddleware(h.handleEditBlog))

	//
	// blog categories
	//

	// lookup blogs by category
	server.HandleFunc("GET "+prefix+"/category/{category}", h.handleBlogsByCategory)

	//
	// blog drafts
	//

	// get user drafts
	server.HandleFunc("GET "+prefix+"/drafts/{userID}", middlewares.BearerAuthMiddleware(h.handleDrafts))
}
