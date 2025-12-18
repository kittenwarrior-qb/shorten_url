package utils

import (
	"net/url"
	"regexp"
	"strings"
)

var aliasRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateURL checks if a string is a valid URL
// Rules:
// - Must have http or https scheme
// - Must have valid host
// - Max length 2048 characters
func ValidateURL(urlStr string) bool {
	// Check empty
	if strings.TrimSpace(urlStr) == "" {
		return false
	}

	// Check max length
	if len(urlStr) > 2048 {
		return false
	}

	// Parse URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Must have http or https scheme
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// Must have host
	if u.Host == "" {
		return false
	}

	return true
}

// ValidateAlias checks if a custom alias is valid
// Rules:
// - 3-20 characters
// - Only alphanumeric, underscore, hyphen
func ValidateAlias(alias string) bool {
	if len(alias) < 3 || len(alias) > 20 {
		return false
	}
	return aliasRegex.MatchString(alias)
}
