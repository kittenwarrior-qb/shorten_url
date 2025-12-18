package dto

import "time"

// AnalyticsSummary contains aggregated analytics data
type AnalyticsSummary struct {
	TotalClicks  int64            `json:"total_clicks"`
	BrowserStats map[string]int64 `json:"browser_stats"`
	OSStats      map[string]int64 `json:"os_stats"`
	DeviceStats  map[string]int64 `json:"device_stats"`
	CountryStats map[string]int64 `json:"country_stats"`
	RefererStats map[string]int64 `json:"referer_stats"`
}

// ClickResponse represents a single click event
type ClickResponse struct {
	ID          uint      `json:"id"`
	IPAddress   string    `json:"ip_address"`
	Browser     string    `json:"browser"`
	BrowserVer  string    `json:"browser_ver"`
	OS          string    `json:"os"`
	Device      string    `json:"device"`
	Country     string    `json:"country"`
	CountryCode string    `json:"country_code"`
	City        string    `json:"city"`
	Referer     string    `json:"referer"`
	ClickedAt   time.Time `json:"clicked_at"`
}

// AnalyticsResponse represents analytics data for a link
type AnalyticsResponse struct {
	Summary *AnalyticsSummary `json:"summary"`
	Clicks  []ClickResponse   `json:"clicks,omitempty"`
	Total   int64             `json:"total"`
	Page    int               `json:"page"`
	PerPage int               `json:"per_page"`
}
