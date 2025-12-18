package service

import (
	"time"

	"quocbui.dev/m/internal/models"
	"quocbui.dev/m/internal/repository"
	"quocbui.dev/m/pkg/utils"
)

// ClickInfo contains information about a click event
type ClickInfo struct {
	IPAddress string
	UserAgent string
	Referer   string
}

// LinkService handles link-related business logic
type LinkService struct {
	linkRepo  repository.LinkRepository
	clickRepo repository.ClickRepository
	geoIP     *GeoIPService
}

// NewLinkService creates a new link service
func NewLinkService(linkRepo repository.LinkRepository, clickRepo repository.ClickRepository, geoIP *GeoIPService) *LinkService {
	return &LinkService{
		linkRepo:  linkRepo,
		clickRepo: clickRepo,
		geoIP:     geoIP,
	}
}

// CreateLink creates a new shortened link
func (s *LinkService) CreateLink(originalURL string, customAlias *string, userID *uint, expiresAt *time.Time, shortCodeLength int) (*models.Link, error) {
	// Validate URL
	if !utils.ValidateURL(originalURL) {
		return nil, ErrInvalidURL
	}

	var shortCode string
	var err error

	// Use custom alias if provided
	if customAlias != nil && *customAlias != "" {
		if !utils.ValidateAlias(*customAlias) {
			return nil, ErrInvalidAlias
		}
		// Check if alias already exists
		existing, _ := s.linkRepo.GetByShortCode(*customAlias)
		if existing != nil {
			return nil, ErrAliasAlreadyExists
		}
		shortCode = *customAlias
	} else {
		// Generate random short code
		for i := 0; i < 5; i++ {
			shortCode, err = utils.GenerateShortCode(shortCodeLength)
			if err != nil {
				return nil, err
			}
			// Check if code already exists
			existing, _ := s.linkRepo.GetByShortCode(shortCode)
			if existing == nil {
				break
			}
		}
	}

	link := &models.Link{
		UserID:      userID,
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CustomAlias: customAlias,
		ExpiresAt:   expiresAt,
	}

	if err := s.linkRepo.Create(link); err != nil {
		return nil, err
	}

	return link, nil
}

// Redirect gets the original URL and tracks the click
func (s *LinkService) Redirect(shortCode string, clickInfo *ClickInfo) (string, error) {
	link, err := s.linkRepo.GetByShortCode(shortCode)
	if err != nil {
		return "", ErrLinkNotFound
	}

	// Check if link has expired
	if link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now()) {
		return "", ErrLinkExpired
	}

	// Track click asynchronously
	go s.trackClick(link.ID, clickInfo)

	return link.OriginalURL, nil
}

// trackClick records a click event
func (s *LinkService) trackClick(linkID uint, info *ClickInfo) {
	// Parse user agent
	uaInfo := utils.ParseUserAgent(info.UserAgent)

	// Get geo info
	geoInfo, _ := s.geoIP.GetGeoIP(info.IPAddress)

	// Create click record
	click := &models.Click{
		LinkID:      linkID,
		IPAddress:   info.IPAddress,
		UserAgent:   info.UserAgent,
		Browser:     uaInfo.Browser,
		BrowserVer:  uaInfo.BrowserVer,
		OS:          uaInfo.OS,
		Device:      uaInfo.Device,
		Country:     geoInfo.Country,
		CountryCode: geoInfo.CountryCode,
		City:        geoInfo.City,
		Referer:     info.Referer,
	}

	s.clickRepo.Create(click)
	s.linkRepo.IncrementClickCount(linkID)
}

// GetUserLinks returns all links for a user with pagination
func (s *LinkService) GetUserLinks(userID uint, page, pageSize int) ([]*models.Link, int64, error) {
	return s.linkRepo.GetByUserID(userID, page, pageSize)
}

// GetLinkWithAnalytics returns a link with its analytics if the user owns it
func (s *LinkService) GetLinkWithAnalytics(shortCode string, userID uint) (*models.Link, error) {
	link, err := s.linkRepo.GetByShortCode(shortCode)
	if err != nil {
		return nil, ErrLinkNotFound
	}

	// Check ownership
	if link.UserID == nil || *link.UserID != userID {
		return nil, ErrUnauthorized
	}

	return link, nil
}

// DeleteLink deletes a link if the user owns it
func (s *LinkService) DeleteLink(shortCode string, userID uint) error {
	link, err := s.linkRepo.GetByShortCode(shortCode)
	if err != nil {
		return ErrLinkNotFound
	}

	// Check ownership
	if link.UserID == nil || *link.UserID != userID {
		return ErrUnauthorized
	}

	return s.linkRepo.Delete(link.ID)
}
