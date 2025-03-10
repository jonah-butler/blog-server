package email

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PasswordResetRepository interface {
	CreatePasswordResetEntry(ctx context.Context, payload *PasswordResetMeta) error
}

type MongoPasswordResetRepository struct {
	collection *mongo.Collection
}

func NewPasswordResetRepository(db *mongo.Database) PasswordResetRepository {
	return &MongoPasswordResetRepository{
		collection: db.Collection("passowrdResetMeta"),
	}
}

func (r *MongoPasswordResetRepository) CreatePasswordResetEntry(ctx context.Context, payload *PasswordResetMeta) error {
	filter := bson.M{"user": payload.User}

	update := bson.M{
		"$set": bson.M{
			"createdAt": payload.CreatedAt,
			"hash":      payload.Hash,
			"user":      payload.User,
		},
	}

	opts := options.UpdateOne().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
