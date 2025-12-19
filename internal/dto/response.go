package dto

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Standard API Response wrapper
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// APIError represents error details
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse for Swagger documentation
type ErrorResponse struct {
	Success   bool      `json:"success" example:"false"`
	Error     APIError  `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

// MessageResponse for simple success messages
type MessageResponse struct {
	Success   bool      `json:"success" example:"true"`
	Data      Message   `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

type Message struct {
	Message string `json:"message"`
}

// Error codes
const (
	ErrCodeBadRequest         = "BAD_REQUEST"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeConflict           = "CONFLICT"
	ErrCodeGone               = "GONE"
	ErrCodeValidation         = "VALIDATION_ERROR"
	ErrCodeInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrCodeInvalidURL         = "INVALID_URL"
	ErrCodeInvalidAlias       = "INVALID_ALIAS"
	ErrCodeAliasExists        = "ALIAS_EXISTS"
	ErrCodeLinkNotFound       = "LINK_NOT_FOUND"
	ErrCodeLinkExpired        = "LINK_EXPIRED"
	ErrCodeEmailExists        = "EMAIL_EXISTS"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeRateLimitExceeded  = "RATE_LIMIT_EXCEEDED"
)

// Response helpers
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func Error(c *gin.Context, statusCode int, code string, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
		Timestamp: time.Now(),
	})
}

// Common error responses
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, ErrCodeBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, ErrCodeUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, ErrCodeForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, ErrCodeNotFound, message)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, ErrCodeConflict, message)
}

func Gone(c *gin.Context, message string) {
	Error(c, http.StatusGone, ErrCodeGone, message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, ErrCodeInternalServer, message)
}

func ValidationError(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, ErrCodeValidation, message)
}
