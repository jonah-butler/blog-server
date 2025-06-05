package user

import (
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// User POST payloads
type UserLoginPost struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResetPasswordPost struct {
	Email *string `json:"email"`
}

type UserNewPasswordPost struct {
	ResetToken           string `json:"resetToken"`
	Password             string `json:"password"`
	PasswordVerification string `json:"passwordVerification"`
}

type UserSendEmailPost struct {
	From    string `json:"from"`
	Message string `json:"message"`
	Subject string `json:"subject"`
	To      string `json:"to"`
}

// User Base
type User struct {
	Email        string `bson:"email" json:"email"`
	Username     string `bson:"username" json:"username"`
	ProfileImage string `bson:"profileImageLocation" json:"profileImageLocation"`
}

// User base with ID field used in scenarios where ID is not necessary
type UserWithID struct {
	ID   bson.ObjectID `bson:"_id" json:"_id"`
	User `bson:",inline"`
}

// User with Password field
type UserWithPassword struct {
	UserWithID `bson:",inline"`
	Password   string `bson:"password" json:"password"`
}

// User Response
type UserResponse struct {
	User  UserWithID `json:"user"`
	Token string     `json:"token"`
}

type UserPasswordResetResponse struct {
	DidUpdate bool `json:"didUpdate"`
}

// JWT
type JWTClaims struct {
	jwt.StandardClaims
	User UserWithID `json:"user"`
}
