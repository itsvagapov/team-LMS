package jwt

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
	RoleAdmin   UserRole = "admin"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   UserRole `json:"role"`
	jwt.RegisteredClaims
}

func getSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return nil, errors.New("")
	}

	return []byte(secret), nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret, err := getSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
