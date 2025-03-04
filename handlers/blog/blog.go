package blog

import (
	r "blog-api/repositories/blog"
	s "blog-api/services/blog"
	u "blog-api/utilities"
	"fmt"
	"net/http"
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
	blogQuery := new(r.BlogQuery)
	queryValues := req.URL.Query()

	u.ParseBlogQueryParams(blogQuery, queryValues)

	blogs, err := h.blogService.GetBlogIndex(req.Context(), blogQuery)
	if err != nil {
		message := fmt.Errorf("failed to get blogs at offset %d: %v", blogQuery.Offset, err)
		u.WriteJSONErr(w, http.StatusBadRequest, message)
		return
	}

	u.WriteJSON(w, http.StatusOK, blogs)
}

/*
/blog/{slug}

	Lookup blog by slug accepts slug value by route parameter

	 Returns the blog in question and its two surrounding blogs if any otherwise those values are null
*/
func (h *BlogHandler) handleBlogBySlug(w http.ResponseWriter, req *http.Request) {
	slug := req.PathValue("slug")
	if slug == "" {
		error := fmt.Errorf("not a valid slug: %s", slug)
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	blogs, err := h.blogService.GetBlogBySlug(req.Context(), slug)
	if err != nil {
		error := fmt.Errorf("failed to lookup blog by slug %s: %v", slug, err)
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	if blogs.Post1 == nil {
		error := fmt.Errorf("failed to lookup blog by slug: %s", slug)
		u.WriteJSONErr(w, http.StatusNotFound, error)
		return
	}

	u.WriteJSON(w, http.StatusOK, blogs)
}

/*
/blog/random

	Runs a mongo query which executes a random find query on the blogposts collection

	 Return's the random blog and its previous and next blog posts if any.
*/
func (h *BlogHandler) handleRandomBlog(w http.ResponseWriter, req *http.Request) {
	response, err := h.blogService.GetRandomBlog(req.Context())
	if err != nil {
		error := fmt.Errorf("failed to lookup a random blog: %v", err)
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	u.WriteJSON(w, http.StatusOK, response)
}

/*
/blog/{category}/comma,seperated,categories

	Accepts the following query params:
	offset: 0 / 10 / 20 / 30...

	Queries blogs that contain each of the provided categories

	 Retruns array of blogs and hasMore boolean indicating more are available after
	 the set offset.
*/
func (h *BlogHandler) handleBlogsByCategory(w http.ResponseWriter, req *http.Request) {
	category := req.PathValue("category")
	if category == "" {
		error := fmt.Errorf("not a valid category: %s", category)
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	blogQuery := new(r.BlogQuery)

	u.ParseBlogQueryParams(blogQuery, req.URL.Query())

	response, err := h.blogService.GetBlogsByCategory(req.Context(), category, blogQuery)
	if err != nil {
		error := fmt.Errorf("failed to retrieve blogs by category: %s", category)
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	u.WriteJSON(w, http.StatusOK, response)
}

/*
/blog/drafts/{userID]}

	Accepts the following query params:
	offset: 0 / 10 / 20 / 30...

	Queries drafts for the provided user with offset

	Protected endpoint requiring authorized token

	 Retruns array of blogs and hasMore boolean indicating more are available after
	 the set offset.
*/
func (h *BlogHandler) handleDrafts(w http.ResponseWriter, req *http.Request) {
	blogQuery := new(r.BlogQuery)

	u.ParseBlogQueryParams(blogQuery, req.URL.Query())

	response, err := h.blogService.GetDraftsByUser(req.Context(), blogQuery)
	if err != nil {
		error := fmt.Errorf("error getting drafts: %s", err)
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	u.WriteJSON(w, http.StatusOK, response)
}

/*
/blog/search/{query}

	Accepts the following query params:
	offset: 0 / 10 / 20 / 30...

	Queries published blogs by search query where:
	- text can contain search query
	- title can contain search query
	- categories can contain search query

	 Retruns array of blogs and hasMore boolean indicating more are available after
	 the set offset.
*/
func (h *BlogHandler) handleBlogSearch(w http.ResponseWriter, req *http.Request) {
	blogQuery := new(r.BlogQuery)

	u.ParseBlogQueryParams(blogQuery, req.URL.Query())

	searchQuery := req.PathValue("query")
	if searchQuery == "" {
		error := fmt.Errorf("search query is empty")
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	response, err := h.blogService.GetBlogsBySearchQuery(req.Context(), searchQuery, blogQuery)
	if err != nil {
		error := fmt.Errorf("error searching blogs: %s", err)
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	u.WriteJSON(w, http.StatusOK, response)
}

/*
POST
/blog/{id}/like

Validates the provided object id path value and increments
the blog's rating property by one

	Returns the updated document.
*/
func (h *BlogHandler) handleBlogLike(w http.ResponseWriter, req *http.Request) {
	blogID := req.PathValue("id")
	if blogID == "" {
		error := fmt.Errorf("not a valid post ID")
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	response, err := h.blogService.LikeBlog(req.Context(), blogID)
	if err != nil {
		error := fmt.Errorf("failed to update blog rating: %v", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	u.WriteJSON(w, http.StatusOK, response)
}

func (h *BlogHandler) handleUpdatetBlog(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, 32<<20+512)

	isValidMime := ValidateRequestMime(req.Header.Get("Content-Type"), "multipart/form-data")
	if !isValidMime {
		error := fmt.Errorf("invalid content type")
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	reader, err := req.MultipartReader()
	if err != nil {
		error := fmt.Errorf("error reading mutlipart form: %v", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	input, err := ParseMultiPartFormBlogUpdate(reader)
	if err != nil {
		error := fmt.Errorf("error parsing mutlipart form: %v", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	if input.ID == "" || input.ID != req.PathValue("id") {
		error := fmt.Errorf("blog id missing in payload or mismatched")
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	response, err := h.blogService.UpdateBlog(req.Context(), input)
	if err != nil {
		error := fmt.Errorf("error updating blog: %v", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	u.WriteJSON(w, http.StatusOK, response)
}

func (h *BlogHandler) handleSlugValidation(w http.ResponseWriter, req *http.Request) {
	slug := req.PathValue("slug")
	if slug == "" {
		error := fmt.Errorf("no slug provided")
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	response, err := h.blogService.ValidateSlug(req.Context(), slug)
	if err != nil {
		error := fmt.Errorf("error validating slug: %v", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	u.WriteJSON(w, http.StatusOK, response)
}

// use required fields on a struct to ensure the
// new document has all the fields it requires
func (h *BlogHandler) handleNewBlog(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, 32<<20+512)

	isValidMime := ValidateRequestMime(req.Header.Get("Content-Type"), "multipart/form-data")
	if !isValidMime {
		error := fmt.Errorf("invalid content type")
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	reader, err := req.MultipartReader()
	if err != nil {
		error := fmt.Errorf("error reading mutlipart form: %v", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	input, err := ParseMultiPartFormBlogCreate(reader)
	if err != nil {
		error := fmt.Errorf("error parsing mutlipart form: %v", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	if input.Slug == "" {
		error := fmt.Errorf("missing required form value: slug")
		u.WriteJSONErr(w, http.StatusBadRequest, error)
		return
	}

	response, err := h.blogService.CreateBlog(req.Context(), input)
	if err != nil {
		error := fmt.Errorf("error updating blog: %v", err)
		u.WriteJSONErr(w, http.StatusInternalServerError, error)
		return
	}

	u.WriteJSON(w, http.StatusOK, response)

}
