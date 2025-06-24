package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository interface {
	GetUserByID(ctx context.Context) error
	GetUserByEmail(ctx context.Context, email string) (*UserWithID, error)
	RegisterUser(ctx context.Context) error
	FindUser(ctx context.Context, payload UserLoginPost) (*UserWithPassword, error)
	UpdateUserPassword(ctx context.Context, password string, user bson.ObjectID) (bool, error)
	UpdateUser(ctx context.Context, authorId string, input *UserUpdatePost) (*User, error)
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) GetUserByID(ctx context.Context) error { return nil }

func (r *MongoUserRepository) RegisterUser(ctx context.Context) error { return nil }

func (r *MongoUserRepository) GetUserByEmail(ctx context.Context, email string) (*UserWithID, error) {
	var user *UserWithID

	filter := bson.M{"email": email}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}

func (r *MongoUserRepository) FindUser(ctx context.Context, payload UserLoginPost) (*UserWithPassword, error) {
	var user UserWithPassword
	filter := bson.M{"username": payload.Username}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserRepository) UpdateUserPassword(ctx context.Context, password string, user bson.ObjectID) (bool, error) {
	filter := bson.M{
		"_id": user,
	}

	updates := bson.M{
		"$set": bson.M{
			"password": password,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, updates)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount == 1, nil
}

func (r *MongoUserRepository) UpdateUser(ctx context.Context, authorId string, input *UserUpdatePost) (*User, error) {
	var user *User

	hexAuthorID, err := bson.ObjectIDFromHex(authorId)
	if err != nil {
		return user, err
	}

	filter := bson.M{
		"_id": hexAuthorID,
	}

	updateFields := bson.M{}

	if input.ImageKey != "" {
		updateFields["profileImageKey"] = input.ImageKey
	}

	if input.ImageLocation != "" {
		updateFields["profileImageLocation"] = input.ImageLocation
	}

	update := bson.M{"$set": updateFields}

	err = r.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}
