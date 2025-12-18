package utils

import (
	"net/url"
	"regexp"
)

var aliasRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateURL checks if a string is a valid URL
func ValidateURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

// ValidateAlias checks if a custom alias is valid
func ValidateAlias(alias string) bool {
	if len(alias) < 3 || len(alias) > 20 {
		return false
	}
	return aliasRegex.MatchString(alias)
}
