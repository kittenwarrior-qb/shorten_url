package dto

import "time"

// CreateLinkRequest represents a request to create a shortened link
type CreateLinkRequest struct {
	URL       string  `json:"url" binding:"required,url" example:"https://github.com"`
	Alias     *string `json:"alias,omitempty" example:"my-link"`
	ExpiresIn *int    `json:"expires_in,omitempty" example:"24"`
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

// PublicLinkResponse includes token for guest user
type PublicLinkResponse struct {
	Link  LinkResponse `json:"link"`
	Token string       `json:"token,omitempty"` // JWT token for guest user
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
