package blog

import (
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BlogQuery struct {
	Offset int
}

type BlogIndexResponse struct {
	Blogs   []Blog `json:"blogs"`
	HasMore bool   `json:"hasMore"`
}

// update this since it's not atually a SingleBlogResponse anymore
type SingleBlogResponse struct {
	Post1        *Blog `json:"post1"`
	PreviousPost *Blog `json:"previousPost"`
	NextPost     *Blog `json:"nextPost"`
}

type BlogUpdateResponse struct {
	Blog *Blog `json:"blog"`
}

type SlugValidationResponse struct {
	IsAvailable bool `json:"isAvailable"`
}

type Blog struct {
	Categories    []string      `bson:"categories" json:"categories"`
	Rating        int           `bson:"rating" json:"rating"`
	Views         int           `bson:"views" json:"views"`
	ID            bson.ObjectID `bson:"_id" json:"_id"`
	Author        bson.ObjectID `bson:"author" json:"author"`
	Title         string        `bson:"title" json:"title"`
	ImageLocation string        `bson:"featuredImageLocation" json:"featuredImageLocation"`
	ImageTag      string        `bson:"featuredImageTag" json:"featuredImageTag"`
	ImageKey      string        `bson:"featuredImageKey" json:"featuredImageKey"`
	Text          string        `bson:"text" json:"text"`
	Published     bool          `bson:"published" json:"published"`
	Slug          string        `bson:"slug" json:"slug"`
	SanitizedHTML string        `bson:"sanitizedHTML" json:"sanitizedHTML"`
	CreatedAt     time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type UpdateBlogInput struct {
	CreateBlogInput `bson:",inline"`
	ID              string `bson:"_id" param:"id"`
}

type CreateBlogInput struct {
	Categories    []string              `bson:"categories" param:"categories"`
	Text          string                `bson:"text" param:"title"`
	Published     bool                  `bson:"published" param:"published"`
	Title         string                `bson:"title" param:"title"`
	Image         *multipart.FileHeader `bson:"-" param:"image"`     // ignored bson -> ignored in the mongo upsert
	ImageBytes    []byte                `bson:"-" param:"imageData"` // ignored bson -> ignored in the mongo upsert
	ImageLocation string                `bson:"featuredImageLocation"`
	ImageKey      string                `bson:"featuredImageKey"`
	Slug          string                `bson:"slug" param:"slug"`
}
