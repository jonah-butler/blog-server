package user

import (
	"errors"
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
	return &User{
		Username:     user.Username,
		ProfileImage: user.ProfileImage,
		ID:           user.ID,
		Email:        user.Email,
	}
}

func ConvertToUserResponse(user *User, token string) UserResponse {
	userResponse := UserResponse{
		User:  *user,
		Token: token,
	}

	return userResponse
}

func VerifyJWT(token string) (string, error) {
	claims := new(JWTClaims)

	if len(token) == 0 {
		return "", errors.New("invalid token length")
	}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing signature")
		}

		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if !parsedToken.Valid {
		return "", errors.New("token is invalid")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return "", errors.New("token is expired")
	}

	return claims.User.ID.Hex(), nil
}
