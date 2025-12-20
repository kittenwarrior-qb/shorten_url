package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"quocbui.dev/m/internal/dto"
	"quocbui.dev/m/internal/middleware"
	"quocbui.dev/m/internal/models"
	"quocbui.dev/m/internal/service"
)

type LinkHandler struct {
	linkService      *service.LinkService
	analyticsService *service.AnalyticsService
	qrService        *service.QRService
	domain           string
	shortCodeLength  int
}

func NewLinkHandler(
	linkService *service.LinkService,
	analyticsService *service.AnalyticsService,
	qrService *service.QRService,
	domain string,
	shortCodeLength int,
) *LinkHandler {
	return &LinkHandler{
		linkService:      linkService,
		analyticsService: analyticsService,
		qrService:        qrService,
		domain:           domain,
		shortCodeLength:  shortCodeLength,
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
		dto.ValidationError(c, err.Error())
		return
	}

	var expiresAt *time.Time
	if req.ExpiresIn != nil {
		t := time.Now().Add(time.Duration(*req.ExpiresIn) * time.Hour)
		expiresAt = &t
	}

	// Get authorization header
	authHeader := c.GetHeader("Authorization")

	// Service handles authentication and guest user creation
	link, token, err := h.linkService.CreateLinkWithAuth(
		req.URL,
		req.Alias,
		expiresAt,
		authHeader,
		h.shortCodeLength,
	)
	if err != nil {
		h.handleLinkError(c, err)
		return
	}

	dto.Success(c, http.StatusCreated, dto.PublicLinkResponse{
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
			dto.Error(c, http.StatusNotFound, dto.ErrCodeLinkNotFound, "link not found")
			return
		}
		if err == service.ErrLinkExpired {
			dto.Error(c, http.StatusGone, dto.ErrCodeLinkExpired, "link has expired")
			return
		}
		dto.InternalServerError(c, "internal server error")
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
		dto.Unauthorized(c, "unauthorized")
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
		dto.InternalServerError(c, "failed to fetch links")
		return
	}
	linkResponses := make([]dto.LinkResponse, len(links))
	for i, link := range links {
		linkResponses[i] = h.toLinkResponse(link)
	}
	dto.Success(c, http.StatusOK, dto.ListLinksResponse{
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
		dto.Unauthorized(c, "unauthorized")
		return
	}
	code := c.Param("code")
	link, err := h.linkService.GetLinkWithAnalytics(code, userID)
	if err != nil {
		if err == service.ErrLinkNotFound {
			dto.Error(c, http.StatusNotFound, dto.ErrCodeLinkNotFound, "link not found")
			return
		}
		if err == service.ErrUnauthorized {
			dto.Forbidden(c, "you don't own this link")
			return
		}
		dto.InternalServerError(c, "internal server error")
		return
	}
	analytics, _ := h.analyticsService.GetAnalyticsSummary(link.ID, userID)
	dto.Success(c, http.StatusOK, dto.LinkDetailResponse{
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
		dto.Unauthorized(c, "unauthorized")
		return
	}
	code := c.Param("code")
	err := h.linkService.DeleteLink(code, userID)
	if err != nil {
		if err == service.ErrLinkNotFound {
			dto.Error(c, http.StatusNotFound, dto.ErrCodeLinkNotFound, "link not found")
			return
		}
		if err == service.ErrUnauthorized {
			dto.Forbidden(c, "you don't own this link")
			return
		}
		dto.InternalServerError(c, "internal server error")
		return
	}
	dto.Success(c, http.StatusOK, dto.Message{Message: "link deleted successfully"})
}

func (h *LinkHandler) toLinkResponse(link *models.Link) dto.LinkResponse {
	shortURL := fmt.Sprintf("https://%s/%s", h.domain, link.ShortCode)

	// Generate QR code using QR service
	qrCode, _ := h.qrService.GenerateQRCodeBase64(shortURL)

	return dto.LinkResponse{
		ID:          link.ID,
		ShortCode:   link.ShortCode,
		ShortURL:    shortURL,
		OriginalURL: link.OriginalURL,
		ClickCount:  link.ClickCount,
		QRCode:      qrCode,
		ExpiresAt:   link.ExpiresAt,
		CreatedAt:   link.CreatedAt,
	}
}

func (h *LinkHandler) handleLinkError(c *gin.Context, err error) {
	switch err {
	case service.ErrInvalidURL:
		dto.Error(c, http.StatusBadRequest, dto.ErrCodeInvalidURL, "invalid URL")
	case service.ErrInvalidAlias:
		dto.Error(c, http.StatusBadRequest, dto.ErrCodeInvalidAlias, "invalid alias (3-20 alphanumeric characters)")
	case service.ErrAliasAlreadyExists:
		dto.Error(c, http.StatusConflict, dto.ErrCodeAliasExists, "alias already exists")
	default:
		dto.InternalServerError(c, "failed to create link")
	}
}
