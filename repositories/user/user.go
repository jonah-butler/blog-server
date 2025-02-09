package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	GetUser(ctx context.Context) error
	RegisterUser(ctx context.Context) error
	UserLogin(ctx context.Context, p UserLoginPost) error
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) GetUser(ctx context.Context) error { return nil }

func (r *MongoUserRepository) RegisterUser(ctx context.Context) error { return nil }

func (r *MongoUserRepository) UserLogin(ctx context.Context, payload UserLoginPost) error {
	var user *User

	filter := bson.M{"username": payload.Username}

	r.collection.FindOne(ctx, filter).Decode(&user)

	if user == nil {
		return errors.New("user not found")
	}

	isMatch := comparePasswords(user.Password, payload.Password)
	if !isMatch {
		return errors.New("invalid password")
	}

	return nil
}
