package service_test

import (
	"testing"

	"quocbui.dev/m/internal/service"
	"quocbui.dev/m/pkg/utils"
	"quocbui.dev/m/tests/mocks"
)

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	authService := service.NewAuthService(mockRepo)

	user, err := authService.Register("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if user == nil {
		t.Fatal("Register returned nil user")
	}

	if user.Email != "test@example.com" {
		t.Errorf("user.Email = %s, want test@example.com", user.Email)
	}

	if user.Name != "Test User" {
		t.Errorf("user.Name = %s, want Test User", user.Name)
	}

	if user.PasswordHash == "password123" {
		t.Error("Password should be hashed, not stored as plain text")
	}

	if !utils.CheckPassword("password123", user.PasswordHash) {
		t.Error("Password hash should be verifiable")
	}
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	authService := service.NewAuthService(mockRepo)

	_, err := authService.Register("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("First Register returned error: %v", err)
	}

	_, err = authService.Register("test@example.com", "different123", "Another User")
	if err == nil {
		t.Error("Expected error for duplicate email")
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	authService := service.NewAuthService(mockRepo)

	_, err := authService.Register("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	user, err := authService.Login("test@example.com", "password123")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	if user == nil {
		t.Fatal("Login returned nil user")
	}

	if user.Email != "test@example.com" {
		t.Errorf("user.Email = %s, want test@example.com", user.Email)
	}
}

func TestAuthService_Login_InvalidEmail(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	authService := service.NewAuthService(mockRepo)

	_, err := authService.Login("nonexistent@example.com", "password123")
	if err == nil {
		t.Error("Expected error for non-existent email")
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	authService := service.NewAuthService(mockRepo)

	_, err := authService.Register("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	_, err = authService.Login("test@example.com", "wrongpassword")
	if err == nil {
		t.Error("Expected error for wrong password")
	}
}

func TestAuthService_Register_AssignsUserID(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	authService := service.NewAuthService(mockRepo)

	user, err := authService.Register("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if user.ID == 0 {
		t.Error("User ID should be assigned after registration")
	}
}
