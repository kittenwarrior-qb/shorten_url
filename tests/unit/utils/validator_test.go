package utils_test

import (
	"testing"

	"quocbui.dev/m/pkg/utils"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		// Valid URLs
		{"valid http url", "http://example.com", true},
		{"valid https url", "https://example.com", true},
		{"valid url with path", "https://example.com/path/to/resource", true},
		{"valid url with query params", "https://example.com?param1=value1&param2=value2", true},
		{"valid url with port", "https://example.com:8080/path", true},
		{"valid url with subdomain", "https://sub.example.com", true},

		// Invalid URLs
		{"empty string", "", false},
		{"whitespace only", "   ", false},
		{"no scheme", "example.com", false},
		{"ftp scheme", "ftp://example.com", false},
		{"javascript scheme", "javascript:alert(1)", false},
		{"no host", "https://", false},
		{"invalid url", "not-a-url", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidateURL(tt.url)
			if result != tt.expected {
				t.Errorf("ValidateURL(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestValidateURL_MaxLength(t *testing.T) {
	// Create URL with exactly 2048 characters
	longPath := make([]byte, 2030)
	for i := range longPath {
		longPath[i] = 'a'
	}
	validLongURL := "https://example.com/" + string(longPath)

	if len(validLongURL) > 2048 {
		validLongURL = validLongURL[:2048]
	}

	if !utils.ValidateURL(validLongURL) {
		t.Error("URL at max length should be valid")
	}

	tooLongURL := validLongURL + "extra"
	if utils.ValidateURL(tooLongURL) {
		t.Error("URL exceeding 2048 characters should be invalid")
	}
}

func TestValidateAlias(t *testing.T) {
	tests := []struct {
		name     string
		alias    string
		expected bool
	}{
		// Valid aliases
		{"valid lowercase", "mylink", true},
		{"valid uppercase", "MYLINK", true},
		{"valid mixed case", "MyLink", true},
		{"valid with numbers", "link123", true},
		{"valid with underscore", "my_link", true},
		{"valid with hyphen", "my-link", true},
		{"valid min length (3)", "abc", true},
		{"valid max length (20)", "abcdefghij1234567890", true},

		// Invalid aliases
		{"empty string", "", false},
		{"too short (2)", "ab", false},
		{"too long (21)", "abcdefghij12345678901", false},
		{"contains space", "my link", false},
		{"contains special char", "my@link", false},
		{"contains dot", "my.link", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidateAlias(tt.alias)
			if result != tt.expected {
				t.Errorf("ValidateAlias(%q) = %v, want %v", tt.alias, result, tt.expected)
			}
		})
	}
}
