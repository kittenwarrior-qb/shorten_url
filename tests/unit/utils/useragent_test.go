package utils_test

import (
	"testing"

	"quocbui.dev/m/pkg/utils"
)

func TestParseUserAgent_Chrome_Desktop(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	info := utils.ParseUserAgent(ua)

	if info.Browser != "Chrome" {
		t.Errorf("Browser = %s, want Chrome", info.Browser)
	}

	if info.Device != "Desktop" {
		t.Errorf("Device = %s, want Desktop", info.Device)
	}

	if info.OS == "" {
		t.Error("OS should not be empty")
	}
}

func TestParseUserAgent_Firefox_Desktop(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0"

	info := utils.ParseUserAgent(ua)

	if info.Browser != "Firefox" {
		t.Errorf("Browser = %s, want Firefox", info.Browser)
	}

	if info.Device != "Desktop" {
		t.Errorf("Device = %s, want Desktop", info.Device)
	}
}

func TestParseUserAgent_Safari_Desktop(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15"

	info := utils.ParseUserAgent(ua)

	if info.Browser != "Safari" {
		t.Errorf("Browser = %s, want Safari", info.Browser)
	}

	if info.Device != "Desktop" {
		t.Errorf("Device = %s, want Desktop", info.Device)
	}
}

func TestParseUserAgent_Edge_Desktop(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"

	info := utils.ParseUserAgent(ua)

	if info.Device != "Desktop" {
		t.Errorf("Device = %s, want Desktop", info.Device)
	}
}

func TestParseUserAgent_Chrome_Mobile(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; SM-S918B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"

	info := utils.ParseUserAgent(ua)

	if info.Device != "Mobile" {
		t.Errorf("Device = %s, want Mobile", info.Device)
	}

	if info.Browser == "" {
		t.Error("Browser should not be empty")
	}
}

func TestParseUserAgent_Safari_iPhone(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1"

	info := utils.ParseUserAgent(ua)

	if info.Device != "Mobile" {
		t.Errorf("Device = %s, want Mobile", info.Device)
	}

	if info.Browser != "Safari" {
		t.Errorf("Browser = %s, want Safari", info.Browser)
	}
}

func TestParseUserAgent_Safari_iPad(t *testing.T) {
	ua := "Mozilla/5.0 (iPad; CPU OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1"

	info := utils.ParseUserAgent(ua)

	if info.Device != "Mobile" {
		t.Errorf("Device = %s, want Mobile", info.Device)
	}
}

func TestParseUserAgent_Bot_Googlebot(t *testing.T) {
	ua := "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"

	info := utils.ParseUserAgent(ua)

	if info.Device != "Bot" {
		t.Errorf("Device = %s, want Bot", info.Device)
	}
}

func TestParseUserAgent_Bot_Bingbot(t *testing.T) {
	ua := "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)"

	info := utils.ParseUserAgent(ua)

	if info.Device != "Bot" {
		t.Errorf("Device = %s, want Bot", info.Device)
	}
}

func TestParseUserAgent_EmptyString(t *testing.T) {
	info := utils.ParseUserAgent("")

	if info == nil {
		t.Fatal("ParseUserAgent returned nil")
	}

	// Should return some default values
	if info.Device == "" {
		t.Error("Device should not be empty")
	}
}

func TestParseUserAgent_InvalidString(t *testing.T) {
	info := utils.ParseUserAgent("invalid-user-agent")

	if info == nil {
		t.Fatal("ParseUserAgent returned nil")
	}

	// Should handle gracefully
	if info.Device == "" {
		t.Error("Device should not be empty")
	}
}

func TestParseUserAgent_BrowserVersion(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	info := utils.ParseUserAgent(ua)

	if info.BrowserVer == "" {
		t.Error("BrowserVer should not be empty for Chrome")
	}
}

func TestParseUserAgent_OS_Windows(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	info := utils.ParseUserAgent(ua)

	if info.OS == "" {
		t.Error("OS should not be empty for Windows")
	}
}

func TestParseUserAgent_OS_Mac(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15"

	info := utils.ParseUserAgent(ua)

	if info.OS == "" {
		t.Error("OS should not be empty for Mac")
	}
}

func TestParseUserAgent_OS_Linux(t *testing.T) {
	ua := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	info := utils.ParseUserAgent(ua)

	if info.OS == "" {
		t.Error("OS should not be empty for Linux")
	}
}

func TestParseUserAgent_OS_Android(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; SM-S918B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"

	info := utils.ParseUserAgent(ua)

	if info.OS == "" {
		t.Error("OS should not be empty for Android")
	}
}
