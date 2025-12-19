package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"quocbui.dev/m/internal/dto"
	"quocbui.dev/m/internal/service"
	"quocbui.dev/m/pkg/utils"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtSecret   string
	jwtExpiry   int
}

func NewAuthHandler(authService *service.AuthService, jwtSecret string, jwtExpiry int) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
	}
}

// Register godoc
// @Summary      Register new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Register request"
// @Success      201 {object} dto.UserResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      409 {object} dto.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.ValidationError(c, err.Error())
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			dto.Error(c, http.StatusConflict, dto.ErrCodeEmailExists, "email already exists")
			return
		}
		dto.InternalServerError(c, "failed to register user")
		return
	}

	dto.Success(c, http.StatusCreated, dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login request"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.ValidationError(c, err.Error())
		return
	}

	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		dto.Error(c, http.StatusUnauthorized, dto.ErrCodeInvalidCredentials, "invalid credentials")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, h.jwtSecret, h.jwtExpiry)
	if err != nil {
		dto.InternalServerError(c, "failed to generate token")
		return
	}

	dto.Success(c, http.StatusOK, dto.LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(h.jwtExpiry) * time.Hour),
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
	})
}
