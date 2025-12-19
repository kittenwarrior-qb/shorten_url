package dto

import "time"

// AnalyticsSummary contains aggregated analytics data
type AnalyticsSummary struct {
	TotalClicks int64            `json:"total_clicks"`
	Browsers    map[string]int64 `json:"browsers,omitempty"`
	OS          map[string]int64 `json:"os,omitempty"`
	Devices     map[string]int64 `json:"devices,omitempty"`
	Countries   map[string]int64 `json:"countries,omitempty"`
	Referers    map[string]int64 `json:"referers,omitempty"`
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
