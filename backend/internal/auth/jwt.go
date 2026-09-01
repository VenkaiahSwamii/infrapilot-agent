package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey []byte

func getJWTSecret() []byte {
	if jwtSecretKey == nil {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "ChangeThisToASecretKey"
		}
		jwtSecretKey = []byte(secret)
	}
	return jwtSecretKey
}

type Claims struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`

	jwt.RegisteredClaims
}

// GenerateToken issues a signed JWT token valid for 24 hours including organization_id claim
func GenerateToken(userID, username, email, role, organizationID string) (string, error) {
	if organizationID == "" {
		organizationID = "default"
	}
	claims := Claims{
		UserID:         userID,
		Username:       username,
		Email:          email,
		Role:           role,
		OrganizationID: organizationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ValidateToken parses and validates a signed JWT token string
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claim")
	}

	return claims, nil
}

// GenerateRefreshToken issues a signed JWT token valid for 7 days including organization_id claim
func GenerateRefreshToken(userID, username, email, role, organizationID string) (string, error) {
	if organizationID == "" {
		organizationID = "default"
	}
	claims := Claims{
		UserID:         userID,
		Username:       username,
		Email:          email,
		Role:           role,
		OrganizationID: organizationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}
