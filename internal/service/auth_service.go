package service

import (
	"quocbui.dev/m/internal/models"
	"quocbui.dev/m/internal/repository"
	"quocbui.dev/m/pkg/utils"
)

// AuthService handles authentication logic
type AuthService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register creates a new user account
func (s *AuthService) Register(email, password, name string) (*models.User, error) {
	// Check if email already exists
	existing, _ := s.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Name:         name,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns the user if successful
func (s *AuthService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
