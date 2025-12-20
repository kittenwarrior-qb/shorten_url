package utils_test

import (
	"testing"
	"time"

	"quocbui.dev/m/pkg/utils"
)

func TestGenerateToken(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	secret := "test-secret-key"
	expiryHours := 24

	token, err := utils.GenerateToken(userID, email, secret, expiryHours)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	if token == "" {
		t.Error("GenerateToken returned empty token")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	userID := uint(123)
	email := "test@example.com"
	secret := "test-secret-key"
	expiryHours := 24

	token, err := utils.GenerateToken(userID, email, secret, expiryHours)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := utils.ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("claims.UserID = %d, want %d", claims.UserID, userID)
	}

	if claims.Email != email {
		t.Errorf("claims.Email = %s, want %s", claims.Email, email)
	}
}

func TestValidateToken_InvalidSecret(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	secret := "correct-secret"
	wrongSecret := "wrong-secret"

	token, err := utils.GenerateToken(userID, email, secret, 24)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	_, err = utils.ValidateToken(token, wrongSecret)
	if err == nil {
		t.Error("ValidateToken should return error for invalid secret")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	secret := "test-secret"

	token, err := utils.GenerateToken(userID, email, secret, -1)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	_, err = utils.ValidateToken(token, secret)
	if err == nil {
		t.Error("ValidateToken should return error for expired token")
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	secret := "test-secret"

	invalidTokens := []string{
		"",
		"invalid-token",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
	}

	for _, token := range invalidTokens {
		_, err := utils.ValidateToken(token, secret)
		if err == nil {
			t.Errorf("ValidateToken(%q) should return error", token)
		}
	}
}
