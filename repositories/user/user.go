package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	GetUser(ctx context.Context) error
	RegisterUser(ctx context.Context) error
	FindUser(ctx context.Context, payload UserLoginPost) (*UserWithPassword, error)
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
