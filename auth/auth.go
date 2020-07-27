package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenClaims is the package to request a new token
type TokenClaims struct {
	ID    string
	Email string
}

// Auth defines this struct
type Auth struct {
	secret string
}

// New returns an Auth pointer
func New(secret string) *Auth {
	return &Auth{
		secret: secret,
	}
}

// GetToken returns a new auth token
func (a *Auth) GetToken(tc *TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    tc.ID,
		"email": tc.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString([]byte(a.secret))

}

// IsValid checks if a token is valid
func (a *Auth) IsValid(authorization string) (bool, error) {
	token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Error on validating auth token")
		}
		return []byte(a.secret), nil
	})
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}
