package blog

import (
	r "blog-api/repositories/blog"
	"blog-api/s3"
	"context"

	"github.com/microcosm-cc/bluemonday"
)

type BlogService struct {
	blogRepo r.BlogRepository
}

func NewBlogService(repo r.BlogRepository) *BlogService {
	return &BlogService{blogRepo: repo}
}

func (s *BlogService) GetBlogIndex(ctx context.Context, q *r.BlogQuery) (r.BlogIndexResponse, error) {
	blogs, hasMore, err := s.blogRepo.GetBlogIndex(ctx, q)

	reponse := r.BlogIndexResponse{
		Blogs:   blogs,
		HasMore: hasMore,
	}

	return reponse, err
}

func (s *BlogService) GetBlogBySlug(ctx context.Context, slug string) (r.SingleBlogResponse, error) {
	var response r.SingleBlogResponse

	blog, err := s.blogRepo.GetBlogBySlug(ctx, slug)
	if err != nil {
		return response, err
	}

	if blog == nil {
		return response, nil
	}

	previousBlog, err := s.blogRepo.GetPreviousBlog(ctx, blog.ID)
	if err != nil {
		return response, err
	}

	nextBlog, err := s.blogRepo.GetNextBlog(ctx, blog.ID)
	if err != nil {
		return response, err
	}

	response.Post1 = blog
	response.NextPost = nextBlog
	response.PreviousPost = previousBlog

	s.blogRepo.IncrementViewCount(blog.Slug)

	return response, nil
}

func (s *BlogService) GetRandomBlog(ctx context.Context) (r.SingleBlogResponse, error) {
	var response r.SingleBlogResponse

	blogs, err := s.blogRepo.GetRandomBlog(ctx)
	if err != nil || blogs[0] == nil {
		return response, err
	}

	blog := blogs[0]

	previousBlog, err := s.blogRepo.GetPreviousBlog(ctx, blog.ID)
	if err != nil {
		return response, err
	}

	nextBlog, err := s.blogRepo.GetNextBlog(ctx, blog.ID)
	if err != nil {
		return response, err
	}

	response.Post1 = blog
	response.NextPost = nextBlog
	response.PreviousPost = previousBlog

	s.blogRepo.IncrementViewCount(blog.Slug)

	return response, nil
}

func (s *BlogService) GetBlogsByCategory(ctx context.Context, category string, q *r.BlogQuery) (r.BlogIndexResponse, error) {
	response := r.BlogIndexResponse{}
	blogs, hasMore, err := s.blogRepo.GetBlogsByCategory(ctx, category, q)
	if err != nil {
		return response, err
	}

	response.Blogs = blogs
	response.HasMore = hasMore

	return response, err
}

func (s *BlogService) GetDraftsByUser(ctx context.Context, q *r.BlogQuery) (r.BlogIndexResponse, error) {
	response := r.BlogIndexResponse{}

	blogs, hasMore, err := s.blogRepo.GetDraftsByUser(ctx, q)
	if err != nil {
		return response, err
	}

	response.HasMore = hasMore
	response.Blogs = blogs

	return response, nil
}

func (s *BlogService) GetBlogsBySearchQuery(ctx context.Context, searchQuery string, q *r.BlogQuery) (r.BlogIndexResponse, error) {
	response := r.BlogIndexResponse{}

	blogs, hasMore, err := s.blogRepo.GetBlogsBySearchQuery(ctx, searchQuery, q)
	if err != nil {
		return response, err
	}

	response.HasMore = hasMore
	response.Blogs = blogs

	return response, nil
}

func (s *BlogService) LikeBlog(ctx context.Context, id string) (r.BlogUpdateResponse, error) {
	var response r.BlogUpdateResponse

	blog, err := s.blogRepo.LikeBlog(ctx, id)
	if err != nil {
		return response, err
	}

	response.Blog = blog

	return response, nil
}

func (s *BlogService) UpdateBlog(ctx context.Context, input *r.UpdateBlogInput) (r.BlogUpdateResponse, error) {
	var response r.BlogUpdateResponse
	// if a file was included process first
	if input.Image != nil {
		url, err := s3.UploadToS3(input.Image, input.ImageBytes)
		if err != nil {
			return response, err
		}

		// set url and filename
		input.ImageLocation = url
		input.ImageKey = input.Image.Filename
	}

	// sanitize input text html
	if input.Text != "" {
		p := bluemonday.UGCPolicy()

		sanitized := p.Sanitize(input.Text)

		input.Text = sanitized
	}

	blog, err := s.blogRepo.UpdateBlog(ctx, input)
	if err != nil {
		return response, err
	}

	response.Blog = blog

	return response, nil
}

func (s *BlogService) ValidateSlug(ctx context.Context, slug string) (r.SlugValidationResponse, error) {
	var response r.SlugValidationResponse

	isAvailable, err := s.blogRepo.ValidateSlug(ctx, slug)
	if err != nil {
		return response, err
	}

	response.IsAvailable = isAvailable

	return response, err
}

func (s *BlogService) CreateBlog(ctx context.Context, input *r.CreateBlogInput) (r.BlogUpdateResponse, error) {
	var response r.BlogUpdateResponse
	// if a file was included process first
	if input.Image != nil {
		url, err := s3.UploadToS3(input.Image, input.ImageBytes)
		if err != nil {
			return response, err
		}

		// set url and filename
		input.ImageLocation = url
		input.ImageKey = input.Image.Filename
	}

	if input.GenerateSlug {
		isValid := false
		slug := generateSlug(input.Title)

		for !isValid {
			response, err := s.ValidateSlug(ctx, slug)
		}
	}

	// sanitize input text html
	if input.Text != "" {
		p := bluemonday.UGCPolicy()

		sanitized := p.Sanitize(input.Text)

		input.Text = sanitized
	}

	blog, err := s.blogRepo.CreateBlog(ctx, input)
	if err != nil {
		return response, err
	}

	response.Blog = blog

	return response, nil
}
