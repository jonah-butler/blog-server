package user

import "go.mongodb.org/mongo-driver/v2/bson"

type UserLoginPost struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       bson.ObjectID `bson:"_id" json:"_id"`
	Username string        `bson:"username" json:"username"`
	Password string        `bson:"password" json:"password"`
}
