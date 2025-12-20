package utils_test

import (
	"testing"

	"quocbui.dev/m/pkg/utils"
)

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty hash")
	}

	if hash == password {
		t.Error("HashPassword should not return plain password")
	}
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "samePassword"

	hash1, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	hash2, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	// bcrypt generates different hashes due to salt
	if hash1 == hash2 {
		t.Error("HashPassword should generate different hashes for same password")
	}
}

func TestCheckPassword_Valid(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if !utils.CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for correct password")
	}
}

func TestCheckPassword_Invalid(t *testing.T) {
	password := "correctPassword"
	wrongPassword := "wrongPassword"

	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if utils.CheckPassword(wrongPassword, hash) {
		t.Error("CheckPassword should return false for incorrect password")
	}
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	password := "somePassword"

	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if utils.CheckPassword("", hash) {
		t.Error("CheckPassword should return false for empty password")
	}
}
