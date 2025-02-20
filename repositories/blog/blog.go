package blog

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	ck "blog-api/contextkeys"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BlogRepository interface {
	GetBlogIndex(ctx context.Context, q *BlogQuery) ([]Blog, bool, error)
	GetBlogBySlug(ctx context.Context, slug string) (*Blog, error)
	GetPreviousBlog(ctx context.Context, id bson.ObjectID) (*Blog, error)
	GetNextBlog(ctx context.Context, id bson.ObjectID) (*Blog, error)
	GetRandomBlog(ctx context.Context) ([]*Blog, error)
	GetBlogsByCategory(ctx context.Context, category string, q *BlogQuery) ([]Blog, bool, error)
	GetBlogsBySearchQuery(ctx context.Context, searchQuery string, q *BlogQuery) ([]Blog, bool, error)
	GetDraftsByUser(ctx context.Context, q *BlogQuery) ([]Blog, bool, error)
	LikeBlog(ctx context.Context, id string) (*Blog, error)
	IncrementViewCount(slug string)
	UpdateBlog(ctx context.Context, input *BlogInput) error
}

type MongoBlogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepository(db *mongo.Database) BlogRepository {
	return &MongoBlogRepository{
		collection: db.Collection("blogposts"),
	}
}

/*
*

	Accepts: context, query

	Takes the provided offset and looks up 10 blogpost documents and returns the slice
	along with a bool indicating if there are any additional blogs available after the
	provided offset.
*/
func (r *MongoBlogRepository) GetBlogIndex(ctx context.Context, q *BlogQuery) ([]Blog, bool, error) {
	limit := 10
	var blogs []Blog

	filter := bson.M{"published": true}

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetLimit(int64(limit)).
		SetSkip(int64(q.Offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return blogs, false, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &blogs); err != nil {
		return blogs, false, err
	}

	totalDocuments, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return blogs, false, err
	}

	hasMore := q.Offset+limit < int(totalDocuments)

	return blogs, hasMore, nil
}

/*
*

	Accepts: context, id (document ObjectID)

	Takes a document's ID and looks up the document just before it. Returns nil if the provided ID
	of the document is the first in the collection.
*/
func (r *MongoBlogRepository) GetPreviousBlog(ctx context.Context, id bson.ObjectID) (*Blog, error) {
	var previousBlog *Blog

	previousFilter := bson.M{
		"$and": []bson.M{
			{"_id": bson.M{"$lt": id}},
			{"published": true},
		},
	}

	previousOpts := options.FindOne().SetSort(bson.M{"_id": -1})

	if err := r.collection.FindOne(ctx, previousFilter, previousOpts).Decode(&previousBlog); err != nil {
		if err != mongo.ErrNoDocuments {
			return previousBlog, err
		}
	}

	return previousBlog, nil
}

/*
*

	Accepts: context, id (document ObjectID)

	Takes a document's ID and looks up the document next to it. Returns nil if the provided ID
	of the document is the latest in the collection.
*/
func (r *MongoBlogRepository) GetNextBlog(ctx context.Context, id bson.ObjectID) (*Blog, error) {
	var nextBlog *Blog

	nextFilter := bson.M{
		"$and": []bson.M{
			{"_id": bson.M{"$gt": id}},
			{"published": true},
		},
	}

	opts := options.FindOne().SetSort(bson.M{"_id": 1})

	if err := r.collection.FindOne(ctx, nextFilter, opts).Decode(&nextBlog); err != nil {
		if err != mongo.ErrNoDocuments {
			return nextBlog, err
		}
	}

	return nextBlog, nil
}

/*
*

	Accepts: context, slug

	Looks up a blog by the provided slug
*/
func (r *MongoBlogRepository) GetBlogBySlug(ctx context.Context, slug string) (*Blog, error) {
	var blog *Blog

	filter := bson.M{"published": true, "slug": slug}

	err := r.collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		return blog, nil
	}

	return blog, nil
}

/*
*

	Accepts: context

	Uses Mongo's aggregate method with the sample key to retrieve
	a random document from the blogposts collection.
*/
func (r *MongoBlogRepository) GetRandomBlog(ctx context.Context) ([]*Blog, error) {
	var blogs []*Blog

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"published": true}}},
		{{Key: "$sample", Value: bson.M{"size": 1}}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return blogs, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &blogs); err != nil {
		return blogs, err
	}

	return blogs, nil
}

/*
*

	Accepts: slug

	Initializes a concurrent update to a blog document's view field.
	Utilizes a new context with a timeout so this can be called
	within other retrieval functions and not be reliant on the
	requests provided context.
*/
func (r *MongoBlogRepository) IncrementViewCount(slug string) {
	go func() {
		updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"slug": slug}
		update := bson.M{"$inc": bson.M{"views": 1}}

		r.collection.UpdateOne(updateCtx, filter, update)
	}()
}

/*
*

	Accepts: context, category, BlogQuery

	Parses comma seperated category value from url path
	and looks up blogs which contain all of the provided input.
	Returns the found blogs and if the collection contains more
	after the provided offset.
*/
func (r *MongoBlogRepository) GetBlogsByCategory(ctx context.Context, category string, q *BlogQuery) ([]Blog, bool, error) {
	limit := 10
	var blogs []Blog

	categorySlice := splitAndTrim(category)

	filter := bson.M{
		"categories": bson.M{"$all": categorySlice},
		"published":  true,
	}

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetLimit(int64(limit)).
		SetSkip(int64(q.Offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return blogs, false, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &blogs); err != nil {
		return blogs, false, err
	}

	totalDocuments, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return blogs, false, err
	}

	hasMore := q.Offset+limit < int(totalDocuments)

	return blogs, hasMore, nil
}

/*
*

	Accepts: context, BlogQuery

	Lookup drafts for a provided offset and userID
	which will be extracted from a verified token.
	A request with an invalid token or non-existent
	userID will not reach this repository method.
*/
func (r *MongoBlogRepository) GetDraftsByUser(ctx context.Context, q *BlogQuery) ([]Blog, bool, error) {
	// lots of repeated code - see what i can do about that...
	limit := 10
	var blogs []Blog

	// maybe store the userID in request context
	// as ObjectID without converting via .Hex()
	// then just convert it when needed since in
	// most cases in queries it will just be needed
	// in its original form
	userID, ok := ctx.Value(ck.UserIDKey).(string)
	if !ok {
		return blogs, false, errors.New("failed to access context values")
	}

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return blogs, false, err
	}

	filter := bson.M{
		"published": false,
		"author":    userObjectID,
	}

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetLimit(int64(limit)).
		SetSkip(int64(q.Offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return blogs, false, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &blogs); err != nil {
		return blogs, false, err
	}

	totalDocuments, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return blogs, false, err
	}

	hasMore := q.Offset+limit < int(totalDocuments)

	return blogs, hasMore, nil
}

func (r *MongoBlogRepository) GetBlogsBySearchQuery(ctx context.Context, searchQuery string, q *BlogQuery) ([]Blog, bool, error) {
	limit := 10
	var blogs []Blog

	escapedQuery := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(searchQuery))

	fmt.Println("query: ", escapedQuery)

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(q.Offset))

	filter := bson.M{
		"$and": []bson.M{
			{"published": true},
			{
				"$or": []bson.M{
					{"text": bson.M{"$regex": escapedQuery, "$options": "i"}},
					{"title": bson.M{"$regex": escapedQuery, "$options": "i"}},
					{"categories": bson.M{"$regex": escapedQuery, "$options": "i"}},
				},
			},
		},
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return blogs, false, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &blogs); err != nil {
		return blogs, false, err
	}

	totalDocuments, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return blogs, false, err
	}

	hasMore := q.Offset+limit < int(totalDocuments)

	return blogs, hasMore, nil
}

func (r *MongoBlogRepository) LikeBlog(ctx context.Context, id string) (*Blog, error) {
	var blog *Blog

	postID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return blog, err
	}

	filter := bson.M{"_id": postID}
	update := bson.M{"$inc": bson.M{"rating": 1}}

	err = r.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&blog)
	if err != nil {
		return blog, err
	}

	return blog, nil
}

func (r *MongoBlogRepository) UpdateBlog(ctx context.Context, input *BlogInput) error {
	fmt.Println("update: ", input)
	return nil
}
