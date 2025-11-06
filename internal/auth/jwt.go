package auth

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type UsernameClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func CreateToken(username string, secret []byte, ttl time.Duration) (string, error) {
	claims := UsernameClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func VerifyToken(tokenString string, secret []byte) (string, error) {
	var claims UsernameClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	if claims.Username == "" {
		return "", errors.New("missing username claim")
	}
	return claims.Username, nil
}
