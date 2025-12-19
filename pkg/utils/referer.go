package utils

import (
	"net/url"
	"strings"
)

type RefererInfo struct {
	Source   string // Facebook, Google, Twitter, Direct, Other
	Domain   string // facebook.com, google.com
	Original string // Raw referer URL
}

var knownSources = map[string]string{
	"facebook.com":    "Facebook",
	"fb.com":          "Facebook",
	"m.facebook.com":  "Facebook",
	"l.facebook.com":  "Facebook",
	"lm.facebook.com": "Facebook",
	"google.com":      "Google",
	"google.co":       "Google",
	"twitter.com":     "Twitter",
	"t.co":            "Twitter",
	"x.com":           "Twitter",
	"instagram.com":   "Instagram",
	"l.instagram.com": "Instagram",
	"linkedin.com":    "LinkedIn",
	"lnkd.in":         "LinkedIn",
	"youtube.com":     "YouTube",
	"youtu.be":        "YouTube",
	"tiktok.com":      "TikTok",
	"reddit.com":      "Reddit",
	"pinterest.com":   "Pinterest",
	"telegram.org":    "Telegram",
	"t.me":            "Telegram",
	"zalo.me":         "Zalo",
	"messenger.com":   "Messenger",
}

func ParseReferer(referer string) RefererInfo {
	if referer == "" {
		return RefererInfo{Source: "Direct", Domain: "", Original: ""}
	}

	parsed, err := url.Parse(referer)
	if err != nil {
		return RefererInfo{Source: "Other", Domain: "", Original: referer}
	}

	host := strings.ToLower(parsed.Host)
	host = strings.TrimPrefix(host, "www.")

	for domain, source := range knownSources {
		if host == domain || strings.HasSuffix(host, "."+domain) {
			return RefererInfo{Source: source, Domain: host, Original: referer}
		}
	}

	return RefererInfo{Source: "Other", Domain: host, Original: referer}
}
