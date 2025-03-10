package email

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PasswordResetResponse struct {
	Message string `json:"message"`
}

type PasswordResetMeta struct {
	User      bson.ObjectID `bson:"user" json:"user"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	Hash      string        `bson:"hash" json:"hash"`
}
