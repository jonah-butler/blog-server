package user

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	GetUser(ctx context.Context) error
	RegisterUser(ctx context.Context) error
	UserLogin(ctx context.Context) error
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &MongoUserRepository{
		collection: db.Collection("user"),
	}
}

func (r *MongoUserRepository) GetUser(ctx context.Context) error { return nil }

func (r *MongoUserRepository) RegisterUser(ctx context.Context) error { return nil }

func (r *MongoUserRepository) UserLogin(ctx context.Context) error { return nil }
