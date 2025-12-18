package dto

import "time"

// CreateLinkRequest represents a request to create a shortened link
type CreateLinkRequest struct {
	URL       string  `json:"url" binding:"required,url"`
	Alias     *string `json:"alias,omitempty"`
	ExpiresIn *int    `json:"expires_in,omitempty"` // hours
}

// LinkResponse represents a link in API responses
type LinkResponse struct {
	ID          uint       `json:"id"`
	ShortCode   string     `json:"short_code"`
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	ClickCount  int64      `json:"click_count"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ListLinksResponse represents a paginated list of links
type ListLinksResponse struct {
	Links   []LinkResponse `json:"links"`
	Total   int64          `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}

// LinkDetailResponse represents a link with analytics
type LinkDetailResponse struct {
	Link      LinkResponse      `json:"link"`
	Analytics *AnalyticsSummary `json:"analytics,omitempty"`
}
