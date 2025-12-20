package utils_test

import (
	"testing"

	"quocbui.dev/m/pkg/utils"
)

func TestParseReferer_Direct(t *testing.T) {
	info := utils.ParseReferer("")

	if info.Source != "Direct" {
		t.Errorf("Source = %s, want Direct", info.Source)
	}

	if info.Domain != "" {
		t.Errorf("Domain = %s, want empty", info.Domain)
	}
}

func TestParseReferer_Facebook(t *testing.T) {
	tests := []struct {
		name    string
		referer string
	}{
		{"facebook.com", "https://facebook.com/page"},
		{"www.facebook.com", "https://www.facebook.com/page"},
		{"m.facebook.com", "https://m.facebook.com/page"},
		{"l.facebook.com", "https://l.facebook.com/?u=https://example.com"},
		{"fb.com", "https://fb.com/page"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := utils.ParseReferer(tt.referer)

			if info.Source != "Facebook" {
				t.Errorf("Source = %s, want Facebook", info.Source)
			}

			if info.Domain == "" {
				t.Error("Domain should not be empty")
			}

			if info.Original != tt.referer {
				t.Errorf("Original = %s, want %s", info.Original, tt.referer)
			}
		})
	}
}

func TestParseReferer_Google(t *testing.T) {
	tests := []struct {
		name    string
		referer string
	}{
		{"google.com", "https://www.google.com/search?q=test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := utils.ParseReferer(tt.referer)

			if info.Source != "Google" {
				t.Errorf("Source = %s, want Google", info.Source)
			}
		})
	}
}

func TestParseReferer_Twitter(t *testing.T) {
	tests := []struct {
		name    string
		referer string
	}{
		{"twitter.com", "https://twitter.com/user/status/123"},
		{"t.co", "https://t.co/abc123"},
		{"x.com", "https://x.com/user/status/123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := utils.ParseReferer(tt.referer)

			if info.Source != "Twitter" {
				t.Errorf("Source = %s, want Twitter", info.Source)
			}
		})
	}
}

func TestParseReferer_Instagram(t *testing.T) {
	info := utils.ParseReferer("https://www.instagram.com/p/abc123/")

	if info.Source != "Instagram" {
		t.Errorf("Source = %s, want Instagram", info.Source)
	}
}

func TestParseReferer_LinkedIn(t *testing.T) {
	tests := []struct {
		name    string
		referer string
	}{
		{"linkedin.com", "https://www.linkedin.com/feed/"},
		{"lnkd.in", "https://lnkd.in/abc123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := utils.ParseReferer(tt.referer)

			if info.Source != "LinkedIn" {
				t.Errorf("Source = %s, want LinkedIn", info.Source)
			}
		})
	}
}

func TestParseReferer_YouTube(t *testing.T) {
	tests := []struct {
		name    string
		referer string
	}{
		{"youtube.com", "https://www.youtube.com/watch?v=abc123"},
		{"youtu.be", "https://youtu.be/abc123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := utils.ParseReferer(tt.referer)

			if info.Source != "YouTube" {
				t.Errorf("Source = %s, want YouTube", info.Source)
			}
		})
	}
}

func TestParseReferer_Other(t *testing.T) {
	info := utils.ParseReferer("https://example.com/page")

	if info.Source != "Other" {
		t.Errorf("Source = %s, want Other", info.Source)
	}

	if info.Domain != "example.com" {
		t.Errorf("Domain = %s, want example.com", info.Domain)
	}
}

func TestParseReferer_InvalidURL(t *testing.T) {
	info := utils.ParseReferer("not-a-valid-url")

	if info.Source != "Other" {
		t.Errorf("Source = %s, want Other", info.Source)
	}

	if info.Original != "not-a-valid-url" {
		t.Errorf("Original = %s, want not-a-valid-url", info.Original)
	}
}

func TestParseReferer_WithWWW(t *testing.T) {
	info := utils.ParseReferer("https://www.facebook.com/page")

	if info.Source != "Facebook" {
		t.Errorf("Source = %s, want Facebook (www should be stripped)", info.Source)
	}
}

func TestParseReferer_CaseInsensitive(t *testing.T) {
	info := utils.ParseReferer("https://WWW.FACEBOOK.COM/page")

	if info.Source != "Facebook" {
		t.Errorf("Source = %s, want Facebook (should be case insensitive)", info.Source)
	}
}
