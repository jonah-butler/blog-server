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
	Blog     *Blog `json:"blog"`
	Previous *Blog `json:"previous"`
	Next     *Blog `json:"next"`
}

type BlogUpdateResponse struct {
	Blog *Blog `json:"blog"`
}

type GenericUpdateResponse struct {
	Affected int `json:"affected"`
}

type SlugValidationResponse struct {
	IsAvailable bool `json:"isAvailable"`
}

type Blog struct {
	Categories    []string      `bson:"categories" json:"categories"`
	Rating        int           `bson:"rating" json:"rating"`
	Views         int           `bson:"views" json:"views"`
	ID            bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Author        bson.ObjectID `bson:"author" json:"author"`
	Title         string        `bson:"title" json:"title"`
	ImageLocation string        `bson:"featuredImageLocation" json:"featuredImageLocation"`
	ImageTag      string        `bson:"featuredImageTag" json:"featuredImageTag"`
	ImageKey      string        `bson:"featuredImageKey" json:"featuredImageKey"`
	Text          string        `bson:"text" json:"text"`
	Published     bool          `bson:"published" json:"published"`
	Slug          string        `bson:"slug" json:"slug"`
	// SanitizedHTML string        `bson:"sanitizedHTML" json:"sanitizedHTML"` not using atm
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type BaseBlogInput struct {
	Categories    []string              `bson:"categories" form:"categories"`
	Text          string                `bson:"text" form:"text"`
	Published     bool                  `bson:"published" form:"published"`
	Title         string                `bson:"title" form:"title"`
	Image         *multipart.FileHeader `bson:"-" form:"image"`     // ignored bson -> ignored in the mongo upsert
	ImageBytes    []byte                `bson:"-" form:"imageData"` // ignored bson -> ignored in the mongo upsert
	ImageLocation string                `bson:"featuredImageLocation"`
	ImageKey      string                `bson:"featuredImageKey"`
	Slug          string                `bson:"slug" form:"slug"`
}

type UpdateBlogInput struct {
	BaseBlogInput `bson:",inline"`
	ID            string `bson:"_id" form:"id"`
}

type CreateBlogInput struct {
	BaseBlogInput `bson:",inline"`
	GenerateSlug  bool `bson:"generateSlug" form:"generateSlug"`
}
