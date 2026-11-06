package config

import (
	"path"
	"strings"
)

func validateReviewOmitPaths(paths []string) error {
	return validateDerivedPathPatterns("pre_pr_review.omit_paths", paths)
}

func validateReviewDataPaths(paths []string) error {
	return validateDerivedPathPatterns("pre_pr_review.data_paths", paths)
}

// Omits covers an exact file, a directory prefix, or a whole-path Go pattern.
func (review PrePRReview) Omits(name string) bool {
	return matchesReviewPath(review.OmitPaths, name)
}

// ReviewsApart reports whether name is a Review Data Path (ADR-0368).
func (review PrePRReview) ReviewsApart(name string) bool {
	return matchesReviewPath(review.DataPaths, name)
}

func matchesReviewPath(paths []string, name string) bool {
	for _, entry := range paths {
		if entry == name || (strings.HasSuffix(entry, "/") && strings.HasPrefix(name, entry)) {
			return true
		}
		if matched, _ := path.Match(entry, name); matched {
			return true
		}
	}
	return false
}
