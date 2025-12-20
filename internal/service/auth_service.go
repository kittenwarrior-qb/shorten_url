package service

import (
	"fmt"
	"time"

	"quocbui.dev/m/internal/models"
	"quocbui.dev/m/internal/repository"
	"quocbui.dev/m/pkg/utils"
)

// AuthService handles authentication logic
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry int
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry int) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
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

// LoginWithToken authenticates a user and returns user with JWT token
func (s *AuthService) LoginWithToken(email, password string) (*models.User, string, error) {
	user, err := s.Login(email, password)
	if err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// ValidateToken validates JWT token and returns user ID
func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	claims, err := utils.ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		return 0, ErrInvalidToken
	}
	return claims.UserID, nil
}

// CreateGuestUser creates a temporary guest user account
func (s *AuthService) CreateGuestUser() (*models.User, string, error) {
	guestEmail := fmt.Sprintf("guest_%d@temp.local", time.Now().UnixNano())
	guestPass := fmt.Sprintf("guest_%d", time.Now().UnixNano())

	user, err := s.Register(guestEmail, guestPass, "Guest")
	if err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetUserFromToken extracts and validates user from authorization header
func (s *AuthService) GetUserFromToken(authHeader string) (*uint, error) {
	if authHeader == "" {
		return nil, nil
	}

	// Remove "Bearer " prefix if present
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	userID, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}
