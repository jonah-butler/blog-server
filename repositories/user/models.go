package user

import (
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserLoginPost struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResetPasswordPost struct {
	Email *string `json:"email"`
}

type User struct {
	ID           bson.ObjectID `bson:"_id" json:"_id"`
	Username     string        `bson:"username" json:"username"`
	ProfileImage string        `bson:"profileImageLocation" json:"profileImageLocation"`
}

type UserWithPassword struct {
	User     `bson:",inline"`
	Password string `bson:"password" json:"password"`
}

type UserWithoutPassword struct {
	Email string        `bson:"email" json:"email"`
	ID    bson.ObjectID `bson:"_id" json:"_id"`
}

type UserResponse struct {
	User
	Token string `json:"token"`
}

type JWTClaims struct {
	jwt.StandardClaims
	User User `json:"user"`
}
