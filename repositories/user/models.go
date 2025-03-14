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
	ID           bson.ObjectID `bson:"_id" json:"_id"`
	Email        string        `bson:"email" json:"email"`
	Username     string        `bson:"username" json:"username"`
	ProfileImage string        `bson:"profileImageLocation" json:"profileImageLocation"`
}

// User with Password field
type UserWithPassword struct {
	User     `bson:",inline"`
	Password string `bson:"password" json:"password"`
}

// User Response
type UserResponse struct {
	User
	Token string `json:"token"`
}

type UserPasswordResetResponse struct {
	DidUpdate bool `json:"didUpdate"`
}

// JWT
type JWTClaims struct {
	jwt.StandardClaims
	User User `json:"user"`
}
