package config

import (
	"path"
	"strings"
)

func validateReviewOmitPaths(paths []string) error {
	return validateDerivedPathPatterns("pre_pr_review.omit_paths", paths)
}

// Omits covers an exact file, a directory prefix, or a whole-path Go pattern.
func (review PrePRReview) Omits(name string) bool {
	for _, entry := range review.OmitPaths {
		if entry == name || (strings.HasSuffix(entry, "/") && strings.HasPrefix(name, entry)) {
			return true
		}
		if matched, _ := path.Match(entry, name); matched {
			return true
		}
	}
	return false
}
