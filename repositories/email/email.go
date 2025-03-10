package email

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PasswordResetRepository interface {
	CreatePasswordResetEntry(ctx context.Context, payload *PasswordResetMeta) error
	ValidatePasswordReset(ctx context.Context, hash string) (*PasswordResetMeta, error)
	DeletePasswordResetEntry(ctx context.Context, hash string, user bson.ObjectID) (bool, error)
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

func (r *MongoPasswordResetRepository) ValidatePasswordReset(ctx context.Context, hash string) (*PasswordResetMeta, error) {
	passwordResetMeta := new(PasswordResetMeta)

	filter := bson.M{
		"hash": hash,
	}

	result := r.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return passwordResetMeta, fmt.Errorf("invalid token, %s", err)
		}

		return passwordResetMeta, err
	}

	result.Decode(passwordResetMeta)

	return passwordResetMeta, nil
}

func (r *MongoPasswordResetRepository) DeletePasswordResetEntry(ctx context.Context, hash string, user bson.ObjectID) (bool, error) {
	filter := bson.M{
		"user": user,
		"hash": hash,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	return result.DeletedCount == 1, nil
}
