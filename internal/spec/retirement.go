package spec

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var legacyADRStatusPattern = regexp.MustCompile(`(?im)^\s*(?:\*\*)?status(?:\*\*)?:\s*(proposed|accepted|rejected|deprecated|superseded)\b`)

// Retirement reports whether a decision or intent document has retired.
// A proposed decision record is pending, never retired.
type Retirement struct {
	Retired bool
	Reason  string
}

// ClassifyADR reports whether content is a retired decision record.
func ClassifyADR(content []byte) Retirement {
	status, body := retirementStatus(content)
	if status == "" {
		match := legacyADRStatusPattern.FindSubmatch(body)
		if len(match) == 2 {
			status = strings.ToLower(string(match[1]))
		}
	}

	switch status {
	case "rejected", "deprecated", "superseded":
		return Retirement{Retired: true, Reason: status}
	default:
		return Retirement{}
	}
}

// ClassifyBacklogEntry reports whether content is a retired typed intent entry.
func ClassifyBacklogEntry(content []byte) Retirement {
	status, _ := retirementStatus(content)
	if status == "declined" {
		return Retirement{Retired: true, Reason: status}
	}
	return Retirement{}
}

func retirementStatus(content []byte) (string, []byte) {
	frontmatter, body, err := splitFrontmatter(content)
	if err != nil {
		return "", content
	}
	var document struct {
		Status string `yaml:"status"`
	}
	if err := yaml.Unmarshal(frontmatter, &document); err != nil {
		return "", body
	}
	return strings.ToLower(strings.TrimSpace(document.Status)), body
}
