package service

import (
	"quocbui.dev/m/internal/dto"
	"quocbui.dev/m/internal/models"
	"quocbui.dev/m/internal/repository"
)

// AnalyticsService handles analytics-related operations
type AnalyticsService struct {
	clickRepo repository.ClickRepository
	linkRepo  repository.LinkRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(clickRepo repository.ClickRepository, linkRepo repository.LinkRepository) *AnalyticsService {
	return &AnalyticsService{
		clickRepo: clickRepo,
		linkRepo:  linkRepo,
	}
}

// GetClicksByLinkID returns clicks for a specific link with pagination
func (s *AnalyticsService) GetClicksByLinkID(linkID uint, userID uint, page, pageSize int) ([]*models.Click, int64, error) {
	// Verify ownership
	link, err := s.linkRepo.GetByID(linkID)
	if err != nil {
		return nil, 0, ErrLinkNotFound
	}

	if link.UserID == nil || *link.UserID != userID {
		return nil, 0, ErrUnauthorized
	}

	return s.clickRepo.GetByLinkID(linkID, page, pageSize)
}

// GetAnalyticsSummary returns aggregated analytics for a link
func (s *AnalyticsService) GetAnalyticsSummary(linkID uint, userID uint) (*dto.AnalyticsSummary, error) {
	// Verify ownership
	link, err := s.linkRepo.GetByID(linkID)
	if err != nil {
		return nil, ErrLinkNotFound
	}

	if link.UserID == nil || *link.UserID != userID {
		return nil, ErrUnauthorized
	}

	return s.clickRepo.GetAnalytics(linkID)
}
