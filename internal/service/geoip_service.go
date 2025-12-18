package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GeoIPInfo contains geolocation information
type GeoIPInfo struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	City        string `json:"city"`
}

// GeoIPService handles IP geolocation
type GeoIPService struct {
	client *http.Client
}

// NewGeoIPService creates a new GeoIP service
func NewGeoIPService() *GeoIPService {
	return &GeoIPService{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetGeoIP fetches geolocation information for an IP address
func (s *GeoIPService) GetGeoIP(ip string) (*GeoIPInfo, error) {
	if ip == "" || ip == "::1" || ip == "127.0.0.1" {
		return &GeoIPInfo{
			Country:     "Unknown",
			CountryCode: "XX",
			City:        "Unknown",
		}, nil
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=country,countryCode,city", ip)
	resp, err := s.client.Get(url)
	if err != nil {
		return &GeoIPInfo{
			Country:     "Unknown",
			CountryCode: "XX",
			City:        "Unknown",
		}, nil
	}
	defer resp.Body.Close()

	var info GeoIPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return &GeoIPInfo{
			Country:     "Unknown",
			CountryCode: "XX",
			City:        "Unknown",
		}, nil
	}

	return &info, nil
}
