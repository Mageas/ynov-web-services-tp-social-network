package auth

import (
	"errors"
	"time"

	"ynov-social-api/internal/pkg/apperrors"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secret []byte
	ttl    time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService(secret []byte, ttl time.Duration) *JWTService {
	return &JWTService{
		secret: secret,
		ttl:    ttl,
	}
}

// GenerateToken generates a new JWT token for the given email
func (s *JWTService) GenerateToken(email string) (string, error) {
	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", apperrors.Wrap(err, 500, "failed to sign token")
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the email
func (s *JWTService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return "", apperrors.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Email, nil
	}

	return "", apperrors.ErrInvalidToken
}
