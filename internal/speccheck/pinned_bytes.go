package speccheck

import (
	"fmt"
	"regexp"
	"strings"

	"roundfix/internal/spec"
)

// CodeQABytesPinned identifies a QA requirement that pins a byte invariant to a commit.
const CodeQABytesPinned = "SC-QA-BYTES-PINNED"

var (
	byteIdentityPhrasePattern  = regexp.MustCompile(`(?i)byte-identical|byte identical|byte identity|byte-for-byte|byte-equal|identical bytes|same bytes|unchanged bytes`)
	byteInvariantTokenPattern  = regexp.MustCompile(`[0-9A-Za-z_./-]+`)
	byteInvariantCommitPattern = regexp.MustCompile(`^[0-9a-f]{7,40}$`)
)

// PinnedByteInvariants reports each requirement of a pending QA Task that
// names a byte-identity phrase and a commit-shaped token together. It reads
// only parsed requirement items, with no commit horizon or Git history lookup.
func PinnedByteInvariants(task spec.Task, isQA bool) []Finding {
	if task.Status == spec.StatusCompleted || !isQA && task.Type != spec.TaskTypeQA {
		return nil
	}
	var findings []Finding
	for _, requirement := range task.Requirements {
		if !byteIdentityPhrasePattern.MatchString(requirement.Text) {
			continue
		}
		token := pinnedByteCommitToken(requirement.Text)
		if token == "" {
			continue
		}
		findings = append(findings, Finding{
			Code:     CodeQABytesPinned,
			Severity: SeverityGap,
			Summary:  fmt.Sprintf("%s line %d pins a byte invariant to commit %s", task.File, requirement.Line, token),
			Where:    []Location{{Path: task.File, Line: requirement.Line}},
			Fix:      "Compare with the Delivery Base, the merge base of the audited head with the default branch, and keep the commit where the plan measured out of the requirement, following the write-tasks rule \"Byte invariants name the delivery base\".",
		})
	}
	return findings
}

func pinnedByteCommitToken(text string) string {
	for _, token := range byteInvariantTokenPattern.FindAllString(text, -1) {
		for _, part := range strings.Split(strings.TrimRight(token, "."), "..") {
			if byteInvariantCommitPattern.MatchString(part) && strings.ContainsAny(part, "0123456789") && strings.ContainsAny(part, "abcdef") {
				return part
			}
		}
	}
	return ""
}
