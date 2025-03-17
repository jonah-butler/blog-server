package blog

import (
	"blog-api/middlewares"
	"net/http"
)

func (h *BlogHandler) RegisterBlogRoutes(prefix string, server *http.ServeMux) {
	//
	// BLOG
	//

	// get latest blogs
	server.HandleFunc("GET "+prefix+"/", h.handleBlogIndex)
	// new blog
	server.HandleFunc("POST "+prefix+"/", middlewares.BearerAuthMiddleware(h.handleNewBlog))
	// get random blog
	server.HandleFunc("GET "+prefix+"/random", h.handleRandomBlog)
	// checks if the provided slug value is available
	server.HandleFunc("GET "+prefix+"/validate-slug/{slug}", middlewares.BearerAuthMiddleware(h.handleSlugValidation))
	// search blogs
	server.HandleFunc("GET "+prefix+"/search/{query}", h.handleBlogSearch)
	// lookup blog by slug
	server.HandleFunc("GET "+prefix+"/{slug}", h.handleBlogBySlug)
	// delete blog by id
	server.HandleFunc("DELETE "+prefix+"/{id}", middlewares.BearerAuthMiddleware(h.handleDeleteBlog))
	// update blog rating
	server.HandleFunc("POST "+prefix+"/{id}/like", h.handleBlogLike)
	// update blog
	server.HandleFunc("POST "+prefix+"/{id}/edit", middlewares.BearerAuthMiddleware(h.handleUpdatetBlog))
	// delete blog featured image
	server.HandleFunc("DELETE "+prefix+"/{id}/delete-image", middlewares.BearerAuthMiddleware(h.handleImageDelete))

	//
	// BLOG CATEGORIES
	//

	// lookup blogs by category
	server.HandleFunc("GET "+prefix+"/category/{category}", h.handleBlogsByCategory)

	//
	// BLOG DRAFTS
	//

	// get user drafts
	server.HandleFunc("GET "+prefix+"/drafts/{userID}", middlewares.BearerAuthMiddleware(h.handleDrafts))
}
