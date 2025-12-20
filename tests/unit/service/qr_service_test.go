package service_test

import (
	"strings"
	"testing"

	"quocbui.dev/m/internal/service"
)

func TestQRService_GenerateQRCodeBase64_Success(t *testing.T) {
	svc := service.NewQRService("assets/logo.png")

	url := "https://example.com/abc123"
	qrCode, err := svc.GenerateQRCodeBase64(url)
	if err != nil {
		t.Fatalf("GenerateQRCodeBase64 returned error: %v", err)
	}

	if qrCode == "" {
		t.Error("GenerateQRCodeBase64 returned empty string")
	}

	// Check if it's a valid base64 data URL
	if !strings.HasPrefix(qrCode, "data:image/png;base64,") {
		t.Error("QR code should start with data:image/png;base64,")
	}

	// Check if base64 part is not empty
	base64Part := strings.TrimPrefix(qrCode, "data:image/png;base64,")
	if len(base64Part) == 0 {
		t.Error("Base64 part should not be empty")
	}
}

func TestQRService_GenerateQRCodeBase64_EmptyURL(t *testing.T) {
	svc := service.NewQRService("assets/logo.png")

	qrCode, err := svc.GenerateQRCodeBase64("")
	if err != nil {
		t.Fatalf("GenerateQRCodeBase64 returned error: %v", err)
	}

	// Should still generate QR code for empty string
	if qrCode == "" {
		t.Error("GenerateQRCodeBase64 returned empty string")
	}
}

func TestQRService_GenerateQRCodeBase64_LongURL(t *testing.T) {
	svc := service.NewQRService("assets/logo.png")

	// Generate a very long URL
	longURL := "https://example.com/" + strings.Repeat("a", 1000)
	qrCode, err := svc.GenerateQRCodeBase64(longURL)
	if err != nil {
		t.Fatalf("GenerateQRCodeBase64 returned error: %v", err)
	}

	if qrCode == "" {
		t.Error("GenerateQRCodeBase64 should handle long URLs")
	}
}

func TestQRService_GenerateQRCodeBase64_SpecialCharacters(t *testing.T) {
	svc := service.NewQRService("assets/logo.png")

	url := "https://example.com/test?param=value&foo=bar#section"
	qrCode, err := svc.GenerateQRCodeBase64(url)
	if err != nil {
		t.Fatalf("GenerateQRCodeBase64 returned error: %v", err)
	}

	if qrCode == "" {
		t.Error("GenerateQRCodeBase64 should handle URLs with special characters")
	}
}

func TestQRService_GenerateQRCodeBase64_WithoutLogo(t *testing.T) {
	// Test with non-existent logo path
	svc := service.NewQRService("non-existent-logo.png")

	url := "https://example.com/abc123"
	qrCode, err := svc.GenerateQRCodeBase64(url)

	// Should still work even if logo doesn't exist
	if err != nil {
		t.Fatalf("GenerateQRCodeBase64 returned error: %v", err)
	}

	if qrCode == "" {
		t.Error("GenerateQRCodeBase64 should work without logo")
	}
}
