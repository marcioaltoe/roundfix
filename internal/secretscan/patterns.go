// Package secretscan reports secret locations without retaining their values.
package secretscan

import "regexp"

type pattern struct {
	kind       string
	expression *regexp.Regexp
}

// The table order decides which kind a line reports. All secret expressions
// live here; credential assignments additionally require a letter and a digit.
var patterns = []pattern{
	{"private-key-block", regexp.MustCompile(`-----BEGIN (?:[A-Z0-9]+ )*PRIVATE KEY-----`)},
	{"github-token", regexp.MustCompile(`(?:^|\W)(?:gh[pousr]_[A-Za-z0-9]{36}|github_pat_[A-Za-z0-9_]{82})(?:$|\W)`)},
	{"openrouter-key", regexp.MustCompile(`(?:^|\W)sk-or-v1-[a-f0-9]{64}(?:$|\W)`)},
	{"anthropic-key", regexp.MustCompile(`(?:^|\W)sk-ant-[A-Za-z0-9_-]{32,}(?:$|\W)`)},
	{"openai-key", regexp.MustCompile(`(?:^|\W)sk-(?:proj-)?[A-Za-z0-9_-]{32,}(?:$|\W)`)},
	{"aws-access-key-id", regexp.MustCompile(`(?:^|\W)(?:AKIA|ASIA)[A-Z0-9]{16}(?:$|\W)`)},
	{"slack-token", regexp.MustCompile(`(?:^|\W)xox[abprs]-[A-Za-z0-9-]{10,}(?:$|\W)`)},
	{"google-api-key", regexp.MustCompile(`(?:^|\W)AIza[A-Za-z0-9_-]{35}(?:$|\W)`)},
	{"stripe-secret-key", regexp.MustCompile(`(?:^|\W)(?:sk_live_|rk_live_)[A-Za-z0-9]{24,}(?:$|\W)`)},
	{"credential-assignment", regexp.MustCompile(`(?:^|\W)(?:[A-Z][A-Z0-9_]*)?(?:API_KEY|TOKEN|SECRET|PASSWORD)\s*(?:=|:=|\?=|:)\s*["']?([A-Za-z0-9_+/.-][A-Za-z0-9_+/.=-]{19,})`)},
}

func matchesKind(line []byte) string {
	for _, candidate := range patterns {
		if candidate.kind != "credential-assignment" {
			if candidate.expression.Match(line) {
				return candidate.kind
			}
			continue
		}
		for _, groups := range candidate.expression.FindAllSubmatch(line, -1) {
			letter, digit := false, false
			for _, char := range groups[1] {
				letter = letter || char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'
				digit = digit || char >= '0' && char <= '9'
			}
			if letter && digit {
				return candidate.kind
			}
		}
	}
	return ""
}

// Kinds returns a copy of the kind names in table order.
func Kinds() []string {
	kinds := make([]string, len(patterns))
	for index, candidate := range patterns {
		kinds[index] = candidate.kind
	}
	return kinds
}
