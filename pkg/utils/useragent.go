package utils

import (
	"github.com/mssola/useragent"
)

// UserAgentInfo contains parsed user agent information
type UserAgentInfo struct {
	Browser    string
	BrowserVer string
	OS         string
	Device     string
}

// ParseUserAgent parses a user agent string and extracts information
func ParseUserAgent(uaString string) *UserAgentInfo {
	ua := useragent.New(uaString)

	device := "Desktop"
	if ua.Mobile() {
		device = "Mobile"
	} else if ua.Bot() {
		device = "Bot"
	}

	browserName, browserVer := ua.Browser()

	return &UserAgentInfo{
		Browser:    browserName,
		BrowserVer: browserVer,
		OS:         ua.OS(),
		Device:     device,
	}
}
