package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"quocbui.dev/m/internal/dto"
	"quocbui.dev/m/internal/middleware"
	"quocbui.dev/m/internal/repository"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetMe godoc
// @Summary      Get current user
// @Description  Get authenticated user information
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.UserResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		dto.Unauthorized(c, "unauthorized")
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		dto.NotFound(c, "user not found")
		return
	}

	dto.Success(c, http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	})
}
