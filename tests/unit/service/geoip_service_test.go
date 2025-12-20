package service_test

import (
	"testing"

	"quocbui.dev/m/internal/service"
)

func TestGeoIPService_GetGeoIP_LocalhostIPv4(t *testing.T) {
	svc := service.NewGeoIPService()

	info, err := svc.GetGeoIP("127.0.0.1")
	if err != nil {
		t.Fatalf("GetGeoIP returned error: %v", err)
	}

	if info.Country != "Unknown" {
		t.Errorf("Country = %s, want Unknown", info.Country)
	}

	if info.CountryCode != "XX" {
		t.Errorf("CountryCode = %s, want XX", info.CountryCode)
	}

	if info.City != "Unknown" {
		t.Errorf("City = %s, want Unknown", info.City)
	}
}

func TestGeoIPService_GetGeoIP_LocalhostIPv6(t *testing.T) {
	svc := service.NewGeoIPService()

	info, err := svc.GetGeoIP("::1")
	if err != nil {
		t.Fatalf("GetGeoIP returned error: %v", err)
	}

	if info.Country != "Unknown" {
		t.Errorf("Country = %s, want Unknown", info.Country)
	}

	if info.CountryCode != "XX" {
		t.Errorf("CountryCode = %s, want XX", info.CountryCode)
	}
}

func TestGeoIPService_GetGeoIP_EmptyIP(t *testing.T) {
	svc := service.NewGeoIPService()

	info, err := svc.GetGeoIP("")
	if err != nil {
		t.Fatalf("GetGeoIP returned error: %v", err)
	}

	if info.Country != "Unknown" {
		t.Errorf("Country = %s, want Unknown", info.Country)
	}
}

func TestGeoIPService_GetGeoIP_PublicIP(t *testing.T) {
	// Skip in CI/CD or if no internet connection
	if testing.Short() {
		t.Skip("Skipping test that requires internet connection")
	}

	svc := service.NewGeoIPService()

	// Google DNS IP (should return US)
	info, err := svc.GetGeoIP("8.8.8.8")
	if err != nil {
		t.Fatalf("GetGeoIP returned error: %v", err)
	}

	// Should return some country (not Unknown)
	if info.Country == "" {
		t.Error("Country should not be empty for public IP")
	}

	if info.CountryCode == "" {
		t.Error("CountryCode should not be empty for public IP")
	}
}

func TestGeoIPService_GetGeoIP_InvalidIP(t *testing.T) {
	svc := service.NewGeoIPService()

	// Invalid IP should return Unknown gracefully
	info, err := svc.GetGeoIP("invalid-ip")
	if err != nil {
		t.Fatalf("GetGeoIP returned error: %v", err)
	}

	// Should handle gracefully and return Unknown
	if info == nil {
		t.Fatal("GetGeoIP returned nil info")
	}
}
