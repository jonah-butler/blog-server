package handlers

import (
	r "blog-api/repositories/blog"
	s "blog-api/services/blog"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type BlogHandler struct {
	blogService *s.BlogService
}

func NewBlogHandler(service *s.BlogService) *BlogHandler {
	return &BlogHandler{blogService: service}
}

/*
/blog

	Blog index accepts the following query params:

	offset: 0 / 10 / 20 / 30 /  etc

	 Retruns array of blogs and hasMore boolean indicating more are available after
	 the set offset.
*/
func (h *BlogHandler) handleBlogIndex(w http.ResponseWriter, req *http.Request) {
	offset := req.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}

	parsedOffset, err := strconv.Atoi(offset)
	if err != nil {
		parsedOffset = 0
	}

	blogQuery := &r.BlogQuery{
		Offset: parsedOffset,
	}

	blogs, err := h.blogService.GetBlogIndex(req.Context(), blogQuery)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get blogs at offset %s: %v", offset, err), http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(blogs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal blogs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

/*
/blog/{slug}

	Lookup blog by slug accepts slug value by route parameter

	 Returns the blog in question and its two surrounding blogs if any otherwise those values are null
*/
func (h *BlogHandler) handleBlogBySlug(w http.ResponseWriter, req *http.Request) {
	slug := req.PathValue("slug")
	if slug == "" {
		http.Error(w, fmt.Sprintf("Not a valid slug: %s", slug), http.StatusBadRequest)
		return
	}

	blogs, err := h.blogService.GetBlogBySlug(req.Context(), slug)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup blog by slug %s: %v", slug, err), http.StatusBadRequest)
		return
	}

	if blogs.Post1 == nil {
		http.Error(w, fmt.Sprintf("Failed to lookup blog by slug: %s", slug), http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(blogs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal blogs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

/*
/blog/random

	Runs a mongo query which executes a random find query on the blogposts collection

	 Return's the random blog and its previous and next blog posts if any.
*/
func (h *BlogHandler) handleRandomBlog(w http.ResponseWriter, req *http.Request) {
	response, err := h.blogService.GetRandomBlog(req.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup a random blog: %v", err), http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal blogs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

/*
/blog/{category}/comma,seperated,categories

	Accepts the following query params:
	offset: 0 / 10 / 20 / 30...

	Queries blogs that contain each of the provided categories

	 Retruns array of blogs and hasMore boolean indicating more are available after
	 the set offset.
*/
func (c *BlogHandler) handleBlogsByCategory(w http.ResponseWriter, req *http.Request) {
	category := req.PathValue("category")
	if category == "" {
		http.Error(w, fmt.Sprintf("Not a valid category: %s", category), http.StatusBadRequest)
		return
	}

	blogQuery := new(r.BlogQuery)

	//maybe change handleBlogIndex to this
	offset := req.URL.Query().Get("offset")
	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err == nil {
			blogQuery.Offset = parsedOffset
		}
	}

	response, err := c.blogService.GetBlogsByCategory(req.Context(), category, blogQuery)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve blogs by category: %s", category), http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal blogs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
