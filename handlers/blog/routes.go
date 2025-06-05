package blog

import (
	authmiddleware "blog-api/middlewares/auth"
	"net/http"
)

func (h *BlogHandler) RegisterBlogRoutes(prefix string, server *http.ServeMux) {
	//
	// BLOG
	//

	// get latest blogs
	server.HandleFunc("GET "+prefix, h.handleBlogIndex)
	// new blog
	server.HandleFunc("POST "+prefix, authmiddleware.BearerAuthMiddleware(h.handleNewBlog))
	// get random blog
	server.HandleFunc("GET "+prefix+"/random", h.handleRandomBlog)
	// checks if the provided slug value is available
	server.HandleFunc("GET "+prefix+"/validate-slug/{slug}", authmiddleware.BearerAuthMiddleware(h.handleSlugValidation))
	// search blogs
	server.HandleFunc("GET "+prefix+"/search/{query}", h.handleBlogSearch)
	// lookup blog by slug
	server.HandleFunc("GET "+prefix+"/{slug}", h.handleBlogBySlug)
	// get published blogs by user
	server.HandleFunc("GET "+prefix+"/user/{userID}", h.handleBlogsByUser)
	// delete blog by id
	server.HandleFunc("DELETE "+prefix+"/{id}", authmiddleware.BearerAuthMiddleware(h.handleDeleteBlog))
	// update blog rating
	server.HandleFunc("POST "+prefix+"/{id}/like", h.handleBlogLike)
	// update blog
	server.HandleFunc("PUT "+prefix+"/{id}/edit", authmiddleware.BearerAuthMiddleware(h.handleUpdatetBlog))
	// delete blog featured image
	server.HandleFunc("DELETE "+prefix+"/featured-image/{id}", authmiddleware.BearerAuthMiddleware(h.handleImageDelete))

	//
	// BLOG CATEGORIES
	//

	// lookup blogs by category
	server.HandleFunc("GET "+prefix+"/category/{category}", h.handleBlogsByCategory)

	//
	// BLOG DRAFTS
	//

	// get user drafts
	server.HandleFunc("GET "+prefix+"/drafts", authmiddleware.BearerAuthMiddleware(h.handleDrafts))

	// get single draft from author
	server.HandleFunc("GET "+prefix+"/drafts/{slug}", authmiddleware.BearerAuthMiddleware(h.handleDraft))
}
