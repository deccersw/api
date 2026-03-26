// pkg/jwtutil/jwtutil.go
package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID string
	Email  string
}

func GenerateToken(secret string, claims Claims, ttl time.Duration) (string, error) {
	jwtClaims := jwt.MapClaims{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"exp":     time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(secret))
}

func ParseToken(secret, tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID, ok := mapClaims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	email, ok := mapClaims["email"].(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &Claims{
		UserID: userID,
		Email:  email,
	}, nil
}
