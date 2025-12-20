package utils

import (
	"regexp"
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 6", 6},
		{"length 8", 8},
		{"length 10", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := GenerateShortCode(tt.length)
			if err != nil {
				t.Fatalf("GenerateShortCode(%d) returned error: %v", tt.length, err)
			}

			if len(code) != tt.length {
				t.Errorf("GenerateShortCode(%d) returned code of length %d, want %d", tt.length, len(code), tt.length)
			}
		})
	}
}

func TestGenerateShortCode_CharsetOnly(t *testing.T) {
	// Verify generated code only contains valid charset characters
	validCharset := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	for i := 0; i < 100; i++ {
		code, err := GenerateShortCode(6)
		if err != nil {
			t.Fatalf("GenerateShortCode returned error: %v", err)
		}

		if !validCharset.MatchString(code) {
			t.Errorf("GenerateShortCode returned invalid characters: %s", code)
		}
	}
}

func TestGenerateShortCode_Uniqueness(t *testing.T) {
	// Generate multiple codes and verify they are unique
	codes := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		code, err := GenerateShortCode(6)
		if err != nil {
			t.Fatalf("GenerateShortCode returned error: %v", err)
		}

		if codes[code] {
			t.Logf("Duplicate code found after %d iterations: %s (this is statistically possible but rare)", i, code)
		}
		codes[code] = true
	}

	// With 62^6 = 56.8 billion combinations, duplicates in 1000 iterations should be extremely rare
	uniqueCount := len(codes)
	if uniqueCount < iterations-5 {
		t.Errorf("Too many duplicates: got %d unique codes out of %d iterations", uniqueCount, iterations)
	}
}
