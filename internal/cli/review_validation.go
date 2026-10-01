package cli

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type reviewFindingAnchor struct {
	Path      string `json:"path"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
}

type reviewFindingValidation struct {
	Status string `json:"status"`
	Rule   string `json:"rule,omitempty"`
	Reason string `json:"reason"`
}

type reviewValidation struct {
	Conventions string `json:"conventions"`
	Validator   string `json:"validator"`
	Reason      string `json:"reason,omitempty"`
}

var reviewAnchorPattern = regexp.MustCompile("^`?([^\\s:`]+):([0-9]+)(?:-([0-9]+))?`?")
var reviewQuotedDiffHeaderPattern = regexp.MustCompile(`^"(?:[^"\\]|\\.)*" `)

var reviewHunkPattern = regexp.MustCompile(`^@@ -[0-9]+(?:,[0-9]+)? \+([0-9]+)(?:,([0-9]+))? @@`)

func parseReviewFindingAnchor(text string) (reviewFindingAnchor, bool) {
	match := reviewAnchorPattern.FindStringSubmatch(text)
	if match == nil || strings.IndexFunc(match[1], unicode.IsSpace) >= 0 {
		return reviewFindingAnchor{}, false
	}
	if len(match[0]) < len(text) {
		next, _ := utf8.DecodeRuneInString(text[len(match[0]):])
		if !unicode.IsSpace(next) && next != ':' && next != ',' && next != ')' {
			return reviewFindingAnchor{}, false
		}
	}
	start, err := strconv.Atoi(match[2])
	if err != nil || start < 0 {
		return reviewFindingAnchor{}, false
	}
	end := start
	if match[3] != "" {
		end, err = strconv.Atoi(match[3])
		if err != nil || end < start {
			return reviewFindingAnchor{}, false
		}
	}
	return reviewFindingAnchor{Path: match[1], StartLine: start, EndLine: end}, true
}

type reviewDiffFile struct {
	whole  bool
	ranges []reviewFindingAnchor
}
type reviewDiffIndex map[string]reviewDiffFile

// Git quotes paths containing unusual bytes using Go-compatible C escapes.
func reviewDiffPath(value string) string {
	value = strings.TrimSuffix(value, "\t")
	if strings.HasPrefix(value, `"`) {
		if decoded, err := strconv.Unquote(value); err == nil {
			value = decoded
		}
	}
	if strings.HasPrefix(value, "a/") || strings.HasPrefix(value, "b/") {
		value = value[2:]
	}
	return value
}

func indexReviewDiff(diff string) reviewDiffIndex {
	index := make(reviewDiffIndex)
	var path, oldPath string
	inHunk := false
	file := reviewDiffFile{}
	flush := func() {
		if path != "" && path != "/dev/null" {
			if len(file.ranges) == 0 {
				file.whole = true
			}
			index[path] = file
		}
	}
	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "diff --git "):
			flush()
			path, oldPath, file = "", "", reviewDiffFile{}
			inHunk = false
			// The b/ path is authoritative for files without ---/+++ headers.
			header := strings.TrimPrefix(line, "diff --git ")
			if strings.HasPrefix(header, `"`) {
				if quoted := reviewQuotedDiffHeaderPattern.FindStringIndex(header); quoted != nil {
					path = reviewDiffPath(header[quoted[1]:])
				}
			} else if midpoint := len(header) / 2; len(header)%2 == 1 && header[midpoint] == ' ' && strings.HasPrefix(header, "a/") && strings.HasPrefix(header[midpoint+1:], "b/") && header[2:midpoint] == header[midpoint+3:] {
				path = reviewDiffPath(header[midpoint+1:])
			} else if at := strings.LastIndex(header, " b/"); at >= 0 {
				path = reviewDiffPath(header[at+1:])
			}
		case !inHunk && strings.HasPrefix(line, "--- "):
			oldPath = reviewDiffPath(strings.TrimPrefix(line, "--- "))
		case !inHunk && strings.HasPrefix(line, "+++ "):
			path = reviewDiffPath(strings.TrimPrefix(line, "+++ "))
			if path == "/dev/null" {
				path = oldPath
				file.whole = true
			}
		case !inHunk && strings.HasPrefix(line, "rename to "):
			path = strings.TrimPrefix(line, "rename to ")
			if strings.HasPrefix(path, `"`) {
				if decoded, err := strconv.Unquote(path); err == nil {
					path = decoded
				}
			}
		case strings.HasPrefix(line, "@@ "):
			inHunk = true
			match := reviewHunkPattern.FindStringSubmatch(line)
			if match == nil {
				continue
			}
			start, err := strconv.Atoi(match[1])
			if err != nil {
				continue
			}
			count := 1
			if match[2] != "" {
				count, err = strconv.Atoi(match[2])
				if err != nil {
					continue
				}
			}
			end := start
			if count > 0 {
				if count-1 > int(^uint(0)>>1)-start {
					continue
				}
				end = start + count - 1
			}
			file.ranges = append(file.ranges, reviewFindingAnchor{StartLine: start, EndLine: end})
		}
	}
	flush()
	return index
}

func (index reviewDiffIndex) holds(anchor reviewFindingAnchor) bool {
	file, ok := index[anchor.Path]
	if !ok || anchor.StartLine < 0 || anchor.EndLine < anchor.StartLine {
		return false
	}
	if file.whole {
		return true
	}
	for _, span := range file.ranges {
		if anchor.StartLine <= span.EndLine && anchor.EndLine >= span.StartLine {
			return true
		}
	}
	return false
}

func reviewFindingHasFailureClause(text string) bool {
	lower := strings.ToLower(text)
	at := strings.Index(lower, "failure:")
	return at >= 0 && strings.TrimFunc(lower[at+len("failure:"):], unicode.IsSpace) != ""
}

func reviewFindingDismissedByValidation(finding reviewFinding) bool {
	return finding.Validation != nil && finding.Validation.Status == "dismissed-by-validation"
}

func reviewFindingEvidenceDismissed(record reviewRecord, finding reviewFinding) bool {
	matches := 0
	dismissed := false
	for _, disposition := range record.Dispositions {
		if disposition.Repository == record.Repository && disposition.HeadCommit == record.HeadCommit && disposition.Finding == finding.ID && disposition.Text == finding.Text {
			matches++
			dismissed = disposition.Disposition == "dismissed" && strings.TrimSpace(disposition.Evidence) != "" && disposition.FixedBy == ""
		}
	}
	return matches == 1 && dismissed
}

func validateReviewFindingAnchors(record reviewRecord, diff string) (reviewRecord, int) {
	record.Validation = &reviewValidation{Conventions: deliveryConventionsVersion, Validator: "not-needed"}
	if record.Outcome != reviewOutcomeFindings {
		return record, exitOK
	}
	index := indexReviewDiff(diff)
	readable, standing := 0, 0
	for i := range record.FindingItems {
		finding := &record.FindingItems[i]
		anchor, ok := parseReviewFindingAnchor(finding.Text)
		validation := &reviewFindingValidation{Status: "dismissed-by-validation", Rule: "unanchored", Reason: "finding has no path:line anchor"}
		if ok {
			readable++
			finding.Anchor = &anchor
			validation.Reason = "anchor names no line of the candidate diff"
			if index.holds(anchor) {
				standing++
				validation = &reviewFindingValidation{Status: "stands", Reason: "anchored in the candidate diff"}
			}
		}
		finding.Validation = validation
	}
	if readable == 0 {
		record.Outcome = reviewOutcomeBlocked
		record.Reason = "findings name no file and line"
		record.Findings = ""
		record.FindingItems = nil
		return record, exitPreflight
	}
	if standing == 0 {
		record.Outcome = reviewOutcomeFindingsDismissed
		return record, exitOK
	}
	return record, exitRunFailed
}
