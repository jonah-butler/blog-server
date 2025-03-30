package blog

import (
	ck "blog-api/contextkeys"
	r "blog-api/repositories/blog"
	"blog-api/s3"
	"context"
	"fmt"
	"os"

	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	response.Blog = blog
	response.Next = nextBlog
	response.Previous = previousBlog

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

	response.Blog = blog
	response.Next = nextBlog
	response.Previous = previousBlog

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

func (s *BlogService) GetDraftByUser(ctx context.Context, slug string) (r.SingleBlogResponse, error) {
	response := r.SingleBlogResponse{}

	blog, err := s.blogRepo.GetDraftByUser(ctx, slug)
	if err != nil {
		return response, err
	}

	nextBlog, err := s.blogRepo.GetNextDraft(ctx, blog.ID)
	if err != nil {
		return response, err
	}

	previousBlog, err := s.blogRepo.GetPreviousDraft(ctx, blog.ID)
	if err != nil {
		return response, err
	}

	response.Blog = blog
	response.Next = nextBlog
	response.Previous = previousBlog

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
		authorID, ok := ctx.Value(ck.UserIDKey).(string)
		if !ok {
			return response, fmt.Errorf("user id missing in context")
		}

		url, err := s3.UploadToS3(input.Image, input.ImageBytes, authorID)
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
	fmt.Println("is available: ", isAvailable, err)
	if err != nil {
		return response, err
	}

	response.IsAvailable = isAvailable

	return response, err
}

func (s *BlogService) CreateBlog(ctx context.Context, input *r.CreateBlogInput) (r.BlogUpdateResponse, error) {
	safetyNet := 50
	var response r.BlogUpdateResponse
	// if a file was included process first
	if input.Image != nil {
		authorID, ok := ctx.Value(ck.UserIDKey).(string)
		if !ok {
			return response, fmt.Errorf("user id missing in context")
		}

		url, err := s3.UploadToS3(input.Image, input.ImageBytes, authorID)
		if err != nil {
			return response, err
		}

		// set url and filename
		input.ImageLocation = url
		input.ImageKey = input.Image.Filename
	}

	if input.GenerateSlug {
		originalSlug := generateSlug(input.Title)
		fmt.Println(originalSlug)
		slug := originalSlug
		i := 1

		for {
			validationResponse, err := s.ValidateSlug(ctx, slug)
			if err != nil {
				fmt.Println("the error in service", err)
				return response, err
			}

			if validationResponse.IsAvailable {
				input.Slug = slug
				break
			}

			// likelihood of this happening is slim, but just in case
			if i == safetyNet {
				return response, fmt.Errorf("could not generate a unique slug after %d attempts", safetyNet)
			}

			slug = fmt.Sprintf("%s-%d", originalSlug, i)
			i++
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

func (s *BlogService) DeleteImage(ctx context.Context, blogID string) (*r.GenericUpdateResponse, error) {
	response := new(r.GenericUpdateResponse)

	blogObjectID, err := bson.ObjectIDFromHex(blogID)
	if err != nil {
		return response, err
	}

	blog, err := s.blogRepo.GetBlogById(ctx, blogObjectID)
	if err != nil {
		return response, err
	}

	imageKey := blog.ImageKey

	if imageKey == "" {
		return response, fmt.Errorf("the provded blog ID contains no featured image key")
	}

	err = s3.DeleteFromS3(imageKey)
	if err != nil {
		return response, err
	}

	updates := bson.M{
		"featuredImageKey":      "",
		"featuredImageLocation": "",
	}

	additionalFilters := bson.M{
		"featuredImageKey": imageKey,
	}

	docsAffected, err := s.blogRepo.ClearBlogFields(ctx, updates, additionalFilters)
	if err != nil {
		return response, err
	}

	response.Affected = docsAffected
	return response, nil
}

func (s *BlogService) DeleteBlog(ctx context.Context, blogID string) (*r.GenericUpdateResponse, error) {
	response := new(r.GenericUpdateResponse)

	blogObjectID, err := bson.ObjectIDFromHex(blogID)
	if err != nil {
		return response, err
	}

	userID, ok := ctx.Value(ck.UserIDKey).(string)
	if !ok {
		return response, fmt.Errorf("failed to access context values")
	}

	authorObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return response, err
	}

	blog, err := s.blogRepo.GetBlogByIdAndAuthor(ctx, blogObjectID, authorObjectID)
	if err != nil {
		return response, err
	}

	resources := os.Getenv("AWS_BUCKET")
	imageSources, err := extraImageSourcesFromHTML(blog.Text, resources)
	if err != nil {
		return response, err
	}

	if len(imageSources) > 0 {
		// delete images from s3
		for _, src := range imageSources {
			key := extractKeyFromImageSource(src, resources+".s3.amazonaws.com/")

			if key != "" {
				err := s3.DeleteFromS3(key)
				if err != nil {
					return response, err
				}
			}
		}
	}
	affected, err := s.blogRepo.DeleteBlog(ctx, blogObjectID, authorObjectID)
	if err != nil {
		return response, err
	}

	response.Affected = affected

	return response, err
}
