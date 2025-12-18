package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"quocbui.dev/m/internal/dto"
	"quocbui.dev/m/internal/middleware"
	"quocbui.dev/m/internal/models"
	"quocbui.dev/m/internal/service"
	"quocbui.dev/m/pkg/utils"
)

type LinkHandler struct {
	linkService      *service.LinkService
	analyticsService *service.AnalyticsService
	authService      *service.AuthService
	domain           string
	shortCodeLength  int
	jwtSecret        string
	jwtExpiry        int
}

func NewLinkHandler(
	linkService *service.LinkService,
	analyticsService *service.AnalyticsService,
	authService *service.AuthService,
	domain string,
	shortCodeLength int,
	jwtSecret string,
	jwtExpiry int,
) *LinkHandler {
	return &LinkHandler{
		linkService:      linkService,
		analyticsService: analyticsService,
		authService:      authService,
		domain:           domain,
		shortCodeLength:  shortCodeLength,
		jwtSecret:        jwtSecret,
		jwtExpiry:        jwtExpiry,
	}
}

// Shorten godoc
// @Summary      Shorten URL
// @Description  Create shortened link. If token provided, link belongs to that user. Otherwise creates guest account.
// @Tags         links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateLinkRequest true "Create link request"
// @Success      201 {object} dto.PublicLinkResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      409 {object} dto.ErrorResponse
// @Router       /shorten [post]
func (h *LinkHandler) Shorten(c *gin.Context) {
	var req dto.CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresIn != nil {
		t := time.Now().Add(time.Duration(*req.ExpiresIn) * time.Hour)
		expiresAt = &t
	}

	var userID *uint
	var token string

	// Check if user is authenticated
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
		claims, err := utils.ValidateToken(tokenString, h.jwtSecret)
		if err == nil {
			userID = &claims.UserID
		}
	}

	// If no valid token, create guest user
	if userID == nil {
		guestEmail := fmt.Sprintf("guest_%d@temp.local", time.Now().UnixNano())
		guestPass := fmt.Sprintf("guest_%d", time.Now().UnixNano())
		guestUser, err := h.authService.Register(guestEmail, guestPass, "Guest")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create guest account"})
			return
		}
		userID = &guestUser.ID
		token, _ = utils.GenerateToken(guestUser.ID, guestUser.Email, h.jwtSecret, h.jwtExpiry)
	}

	link, err := h.linkService.CreateLink(req.URL, req.Alias, userID, expiresAt, h.shortCodeLength)
	if err != nil {
		h.handleLinkError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.PublicLinkResponse{
		Link:  h.toLinkResponse(link),
		Token: token,
	})
}

// Redirect godoc
// @Summary      Redirect to original URL
// @Description  Redirect short URL to original URL and track click
// @Tags         redirect
// @Param        code path string true "Short code"
// @Success      301 "Redirect to original URL"
// @Failure      404 {object} dto.ErrorResponse
// @Failure      410 {object} dto.ErrorResponse
// @Router       /{code} [get]
func (h *LinkHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	clickInfo := &service.ClickInfo{
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Referer:   c.GetHeader("Referer"),
	}
	originalURL, err := h.linkService.Redirect(code, clickInfo)
	if err != nil {
		if err == service.ErrLinkNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
			return
		}
		if err == service.ErrLinkExpired {
			c.JSON(http.StatusGone, gin.H{"error": "link has expired"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// GetMyLinks godoc
// @Summary      Get my links
// @Description  Get all links for authenticated user with pagination
// @Tags         links
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Items per page" default(10)
// @Success      200 {object} dto.ListLinksResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /me/links [get]
func (h *LinkHandler) GetMyLinks(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	links, total, err := h.linkService.GetUserLinks(userID, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch links"})
		return
	}
	linkResponses := make([]dto.LinkResponse, len(links))
	for i, link := range links {
		linkResponses[i] = h.toLinkResponse(link)
	}
	c.JSON(http.StatusOK, dto.ListLinksResponse{
		Links:   linkResponses,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	})
}

// GetMyLinkDetail godoc
// @Summary      Get link detail
// @Description  Get link detail with analytics
// @Tags         links
// @Produce      json
// @Security     BearerAuth
// @Param        code path string true "Short code"
// @Success      200 {object} dto.LinkDetailResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /me/links/{code} [get]
func (h *LinkHandler) GetMyLinkDetail(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	code := c.Param("code")
	link, err := h.linkService.GetLinkWithAnalytics(code, userID)
	if err != nil {
		if err == service.ErrLinkNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
			return
		}
		if err == service.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "you don't own this link"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	analytics, _ := h.analyticsService.GetAnalyticsSummary(link.ID, userID)
	c.JSON(http.StatusOK, dto.LinkDetailResponse{
		Link:      h.toLinkResponse(link),
		Analytics: analytics,
	})
}

// DeleteMyLink godoc
// @Summary      Delete link
// @Description  Delete a link owned by authenticated user
// @Tags         links
// @Produce      json
// @Security     BearerAuth
// @Param        code path string true "Short code"
// @Success      200 {object} dto.MessageResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /me/links/{code} [delete]
func (h *LinkHandler) DeleteMyLink(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	code := c.Param("code")
	err := h.linkService.DeleteLink(code, userID)
	if err != nil {
		if err == service.ErrLinkNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
			return
		}
		if err == service.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "you don't own this link"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "link deleted successfully"})
}

func (h *LinkHandler) toLinkResponse(link *models.Link) dto.LinkResponse {
	return dto.LinkResponse{
		ID:          link.ID,
		ShortCode:   link.ShortCode,
		ShortURL:    fmt.Sprintf("http://%s/%s", h.domain, link.ShortCode),
		OriginalURL: link.OriginalURL,
		ClickCount:  link.ClickCount,
		ExpiresAt:   link.ExpiresAt,
		CreatedAt:   link.CreatedAt,
	}
}

func (h *LinkHandler) handleLinkError(c *gin.Context, err error) {
	switch err {
	case service.ErrInvalidURL:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})
	case service.ErrInvalidAlias:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alias (3-20 alphanumeric characters)"})
	case service.ErrAliasAlreadyExists:
		c.JSON(http.StatusConflict, gin.H{"error": "alias already exists"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create link"})
	}
}
