package user

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateJWT(user *User) (string, error) {
	claim := &JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    user.ID.String(),
			ExpiresAt: generateJWTExp(7),
			Subject:   "access_token",
			IssuedAt:  time.Now().Unix(),
		},
		User: *user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func generateJWTExp(days int) int64 {
	return time.Now().Add((time.Hour * 24) * time.Duration(days)).Unix()
}

func ConvertToUser(user *UserWithPassword) *User {
	userWithoutPassword := &User{
		Username:     user.Username,
		ProfileImage: user.ProfileImage,
		ID:           user.ID,
	}

	return userWithoutPassword
}

func ConvertToUserResponse(user *User, token string) UserResponse {
	userResponse := UserResponse{
		User:  *user,
		Token: token,
	}

	return userResponse
}
