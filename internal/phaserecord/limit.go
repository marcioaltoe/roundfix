package phaserecord

import (
	"regexp"
	"strings"
	"unicode"
)

// LimitKind distinguishes a model limit from shared and other limits.
type LimitKind string

// UsageLimit holds only the first matching, bounded line of a harness error.
type UsageLimit struct {
	Kind LimitKind
	Name string
	Line string
}

var usageHit = regexp.MustCompile(`(?i)you['’]ve\s+hit your (.+?) limit\b`)
var usageReached = regexp.MustCompile(`(?i)^(.+?) limit reached\b`)

// ReadUsageLimit classifies harness text without knowing any model names.
func ReadUsageLimit(message string) UsageLimit {
	for _, line := range strings.Split(message, "\n") {
		line = strings.TrimSpace(strings.Map(func(r rune) rune {
			if unicode.IsControl(r) {
				return ' '
			}
			return r
		}, line))
		match := usageHit.FindStringSubmatch(line)
		if match == nil {
			match = usageReached.FindStringSubmatch(line)
		}
		if match == nil {
			continue
		}
		name := strings.TrimSpace(match[1])
		kind := LimitKind("model")
		if len(strings.Fields(name)) > 1 {
			kind = "other"
		} else if strings.EqualFold(name, "session") || strings.EqualFold(name, "weekly") {
			kind = "shared"
		}
		runes := []rune(line)
		if len(runes) > 200 {
			line = string(runes[:200])
		}
		return UsageLimit{Kind: kind, Name: name, Line: line}
	}
	return UsageLimit{Kind: "none"}
}
