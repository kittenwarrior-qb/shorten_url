package utils_test

import (
	"strings"
	"testing"

	"quocbui.dev/m/pkg/utils"
)

func TestGenerateQRCode_Success(t *testing.T) {
	content := "https://example.com/test"

	qrBytes, err := utils.GenerateQRCode(content, "")
	if err != nil {
		t.Fatalf("GenerateQRCode returned error: %v", err)
	}

	if len(qrBytes) == 0 {
		t.Error("GenerateQRCode returned empty bytes")
	}

	// Check PNG header
	if len(qrBytes) < 8 {
		t.Fatal("QR code bytes too short to be valid PNG")
	}

	// PNG magic number: 89 50 4E 47 0D 0A 1A 0A
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	for i := 0; i < 8; i++ {
		if qrBytes[i] != pngHeader[i] {
			t.Errorf("Invalid PNG header at byte %d: got %x, want %x", i, qrBytes[i], pngHeader[i])
		}
	}
}

func TestGenerateQRCode_EmptyContent(t *testing.T) {
	qrBytes, err := utils.GenerateQRCode("", "")
	if err != nil {
		t.Fatalf("GenerateQRCode returned error: %v", err)
	}

	if len(qrBytes) == 0 {
		t.Error("GenerateQRCode should generate QR code even for empty content")
	}
}

func TestGenerateQRCode_LongContent(t *testing.T) {
	// Generate very long URL
	longContent := "https://example.com/" + strings.Repeat("a", 1000)

	qrBytes, err := utils.GenerateQRCode(longContent, "")
	if err != nil {
		t.Fatalf("GenerateQRCode returned error: %v", err)
	}

	if len(qrBytes) == 0 {
		t.Error("GenerateQRCode should handle long content")
	}
}

func TestGenerateQRCode_SpecialCharacters(t *testing.T) {
	content := "https://example.com/test?param=value&foo=bar#section"

	qrBytes, err := utils.GenerateQRCode(content, "")
	if err != nil {
		t.Fatalf("GenerateQRCode returned error: %v", err)
	}

	if len(qrBytes) == 0 {
		t.Error("GenerateQRCode should handle special characters")
	}
}

func TestGenerateQRCode_WithNonExistentLogo(t *testing.T) {
	content := "https://example.com/test"

	// Should work even if logo doesn't exist
	qrBytes, err := utils.GenerateQRCode(content, "non-existent-logo.png")
	if err != nil {
		t.Fatalf("GenerateQRCode returned error: %v", err)
	}

	if len(qrBytes) == 0 {
		t.Error("GenerateQRCode should work without logo")
	}
}

func TestGenerateQRCode_WithLogo(t *testing.T) {
	// Skip if logo doesn't exist
	content := "https://example.com/test"

	qrBytes, err := utils.GenerateQRCode(content, "assets/logo.png")
	if err != nil {
		t.Fatalf("GenerateQRCode returned error: %v", err)
	}

	if len(qrBytes) == 0 {
		t.Error("GenerateQRCode returned empty bytes")
	}
}

func TestGenerateQRCode_Unicode(t *testing.T) {
	content := "https://example.com/测试/テスト/тест"

	qrBytes, err := utils.GenerateQRCode(content, "")
	if err != nil {
		t.Fatalf("GenerateQRCode returned error: %v", err)
	}

	if len(qrBytes) == 0 {
		t.Error("GenerateQRCode should handle Unicode characters")
	}
}

func TestGenerateQRCode_MultipleGeneration(t *testing.T) {
	// Test generating multiple QR codes in sequence
	contents := []string{
		"https://example.com/1",
		"https://example.com/2",
		"https://example.com/3",
	}

	for i, content := range contents {
		qrBytes, err := utils.GenerateQRCode(content, "")
		if err != nil {
			t.Fatalf("GenerateQRCode #%d returned error: %v", i+1, err)
		}

		if len(qrBytes) == 0 {
			t.Errorf("GenerateQRCode #%d returned empty bytes", i+1)
		}
	}
}
