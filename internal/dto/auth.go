package dto

import "time"

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse represents a success message response
type MessageResponse struct {
	Message string `json:"message"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@gmail.com"`
	Password string `json:"password" binding:"required" example:"123"`
	Name     string `json:"name" binding:"required" example:"user"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@gmail.com"`
	Password string `json:"password" binding:"required" example:"123"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
