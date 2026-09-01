package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidateToken(t *testing.T) {
	userID := "user-123"
	username := "admin"
	email := "admin@infrapilot.io"
	role := "SuperAdmin"
	orgID := "0efc7e0f-27f9-4f1f-846a-7a3d2df9cfd6"

	// 1. Generate Token
	tokenStr, err := GenerateToken(userID, username, email, role, orgID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if tokenStr == "" {
		t.Fatalf("Generated token string is empty")
	}

	// 2. Validate Token
	claims, err := ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}
	if claims.Username != username {
		t.Errorf("Expected Username %s, got %s", username, claims.Username)
	}
	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}
	if claims.OrganizationID != orgID {
		t.Errorf("Expected OrganizationID %s, got %s", orgID, claims.OrganizationID)
	}
}

func TestValidateInvalidToken(t *testing.T) {
	// 1. Test completely bogus token string
	_, err := ValidateToken("invalid.jwt.string")
	if err == nil {
		t.Errorf("Expected error for invalid token string, got nil")
	}

	// 2. Test expired token
	claims := Claims{
		UserID:         "user-999",
		Username:       "expiredUser",
		Email:          "expired@infrapilot.io",
		Role:           "Viewer",
		OrganizationID: "org-1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredStr, _ := token.SignedString(getJWTSecret())

	_, err = ValidateToken(expiredStr)
	if err == nil {
		t.Errorf("Expected error for expired token, got nil")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	tokenStr, err := GenerateRefreshToken("user-123", "admin", "admin@infrapilot.io", "Admin", "org-1")
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	claims, err := ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("Failed to validate refresh token: %v", err)
	}

	if claims.Role != "Admin" {
		t.Errorf("Expected Role Admin, got %s", claims.Role)
	}
}
