package service

import (
	"log"
	"time"

	"gorm.io/gorm"
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
	txManager repository.TransactionManager
	geoIP     *GeoIPService
}

// NewLinkService creates a new link service
func NewLinkService(linkRepo repository.LinkRepository, clickRepo repository.ClickRepository, txManager repository.TransactionManager, geoIP *GeoIPService) *LinkService {
	return &LinkService{
		linkRepo:  linkRepo,
		clickRepo: clickRepo,
		txManager: txManager,
		geoIP:     geoIP,
	}
}

// CreateLink creates a new shortened link with transaction support
// Uses SELECT FOR UPDATE to prevent race conditions on custom aliases
func (s *LinkService) CreateLink(originalURL string, customAlias *string, userID *uint, expiresAt *time.Time, shortCodeLength int) (*models.Link, error) {
	// Validate URL
	if !utils.ValidateURL(originalURL) {
		return nil, ErrInvalidURL
	}

	var shortCode string
	var err error
	var link *models.Link

	// Use custom alias if provided - needs transaction to prevent race condition
	if customAlias != nil && *customAlias != "" {
		if !utils.ValidateAlias(*customAlias) {
			return nil, ErrInvalidAlias
		}

		// Use transaction with row-level locking to prevent duplicate aliases
		err = s.txManager.ExecuteInTransaction(func(tx *gorm.DB) error {
			// Check if alias already exists with FOR UPDATE lock
			existing, _ := s.linkRepo.GetByShortCodeForUpdate(tx, *customAlias)
			if existing != nil {
				return ErrAliasAlreadyExists
			}

			link = &models.Link{
				UserID:      userID,
				ShortCode:   *customAlias,
				OriginalURL: originalURL,
				CustomAlias: customAlias,
				ExpiresAt:   expiresAt,
			}

			return s.linkRepo.CreateWithTx(tx, link)
		})

		if err != nil {
			return nil, err
		}
		return link, nil
	}

	// Generate random short code - retry on collision
	for i := 0; i < 5; i++ {
		shortCode, err = utils.GenerateShortCode(shortCodeLength)
		if err != nil {
			return nil, err
		}

		link = &models.Link{
			UserID:      userID,
			ShortCode:   shortCode,
			OriginalURL: originalURL,
			CustomAlias: customAlias,
			ExpiresAt:   expiresAt,
		}

		// Try to create - unique constraint will catch collisions
		if err := s.linkRepo.Create(link); err == nil {
			return link, nil
		}
		// If error is not unique violation, return it
		// Otherwise retry with new short code
	}

	return nil, ErrAliasAlreadyExists // All retries failed
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

// trackClick records a click event with transaction support
// Ensures click record and click_count are updated atomically
func (s *LinkService) trackClick(linkID uint, info *ClickInfo) {
	uaInfo := utils.ParseUserAgent(info.UserAgent)
	geoInfo, _ := s.geoIP.GetGeoIP(info.IPAddress)
	refInfo := utils.ParseReferer(info.Referer)

	click := &models.Click{
		LinkID:        linkID,
		IPAddress:     info.IPAddress,
		UserAgent:     info.UserAgent,
		Browser:       uaInfo.Browser,
		BrowserVer:    uaInfo.BrowserVer,
		OS:            uaInfo.OS,
		Device:        uaInfo.Device,
		Country:       geoInfo.Country,
		CountryCode:   geoInfo.CountryCode,
		City:          geoInfo.City,
		Referer:       info.Referer,
		RefererSource: refInfo.Source,
		RefererDomain: refInfo.Domain,
	}

	// Use transaction to ensure atomicity:
	// Both click record and click_count update succeed or both fail
	err := s.txManager.ExecuteInTransaction(func(tx *gorm.DB) error {
		if err := s.clickRepo.CreateWithTx(tx, click); err != nil {
			return err
		}
		return s.linkRepo.IncrementClickCountWithTx(tx, linkID)
	})

	if err != nil {
		log.Printf("Failed to track click for link %d: %v", linkID, err)
	}
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
