package spec

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReviewLiveness is what a recorded outcome establishes about an orphan
// Review Artifact's Pull Request. Undecidable remains live for downstream
// decisions.
type ReviewLiveness string

const (
	ReviewFinished    ReviewLiveness = "finished"
	ReviewLive        ReviewLiveness = "live"
	ReviewUndecidable ReviewLiveness = "undecidable"
)

type reviewOutcomeMetadata struct {
	PullRequestState string `yaml:"pull_request_state"`
	MergeCommit      string `yaml:"merge_commit"`
	RecordedAt       string `yaml:"recorded_at"`
}

// ClassifyReview reads the Review Artifact's recorded outcome. Repository Git
// state is deliberately not an input. An undecidable answer is not an error:
// callers retain that Review Artifact and can report reason.
func ClassifyReview(_ context.Context, _ string, reviewDir string) (ReviewLiveness, string, error) {
	outcomePath := filepath.Join(reviewDir, "outcome.md")
	content, err := os.ReadFile(outcomePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ReviewUndecidable, fmt.Sprintf("recorded outcome %q is missing", outcomePath), nil
		}
		return ReviewUndecidable, fmt.Sprintf("recorded outcome %q cannot be read: %v", outcomePath, err), nil
	}

	frontmatter, _, err := splitFrontmatter(content)
	if err != nil {
		return ReviewUndecidable, fmt.Sprintf("recorded outcome %q is malformed: %v", outcomePath, err), nil
	}
	var outcome reviewOutcomeMetadata
	if err := yaml.Unmarshal(frontmatter, &outcome); err != nil {
		return ReviewUndecidable, fmt.Sprintf("recorded outcome %q is malformed: %v", outcomePath, err), nil
	}

	state := strings.TrimSpace(outcome.PullRequestState)
	if state == "" {
		return ReviewUndecidable, fmt.Sprintf("recorded outcome %q is missing pull_request_state", outcomePath), nil
	}
	if strings.TrimSpace(outcome.RecordedAt) == "" {
		return ReviewUndecidable, fmt.Sprintf("recorded outcome %q is missing recorded_at", outcomePath), nil
	}

	switch state {
	case "open":
		return ReviewLive, "recorded outcome is open", nil
	case "closed":
		return ReviewFinished, "recorded outcome is closed", nil
	case "merged":
		mergeCommit := strings.TrimSpace(outcome.MergeCommit)
		if mergeCommit == "" {
			return ReviewUndecidable, fmt.Sprintf("recorded outcome %q is missing merge_commit for merged pull request", outcomePath), nil
		}
		if !validReviewCommit(mergeCommit) {
			return ReviewUndecidable, fmt.Sprintf("recorded outcome %q has invalid merge_commit %q", outcomePath, mergeCommit), nil
		}
		return ReviewFinished, fmt.Sprintf("recorded outcome is merged with merge_commit %s", mergeCommit), nil
	default:
		return ReviewUndecidable, fmt.Sprintf("recorded outcome %q has invalid pull_request_state %q", outcomePath, state), nil
	}
}

func validReviewCommit(commit string) bool {
	if len(commit) < 4 || len(commit) > 64 {
		return false
	}
	for _, char := range commit {
		if (char < '0' || char > '9') && (char < 'a' || char > 'f') && (char < 'A' || char > 'F') {
			return false
		}
	}
	return true
}
