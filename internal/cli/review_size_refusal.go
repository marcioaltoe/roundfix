package cli

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"roundfix/internal/preflight"
)

type reviewDiffFile struct {
	Path  string
	Bytes int
}

func reviewDiffArgs(base, head string, omitted []reviewOmittedPath) []string {
	args := []string{"diff", "--no-ext-diff", "--no-textconv", "--no-color", base, head, "--"}
	for _, path := range omitted {
		args = append(args, ":(exclude,literal)"+path.Path)
	}
	return args
}

func reviewLargestDiffFiles(ctx context.Context, git preflight.GitRunner, root, base, head string, omitted []reviewOmittedPath, diff string) ([]reviewDiffFile, error) {
	// Match the measured patch's flags and pathspecs, including rename detection.
	args := append([]string{"diff", "--name-only", "-z"}, reviewDiffArgs(base, head, omitted)[1:]...)
	names, err := git.RunGit(ctx, root, args...)
	if err != nil {
		return nil, fmt.Errorf("list review diff files: %w", err)
	}
	var paths []string
	if names != "" {
		paths = strings.Split(strings.TrimSuffix(names, "\x00"), "\x00")
	}
	var sizes []int
	// SplitAfter retains every byte, including the final newline, without a
	// Scanner's line-length limit. Headers supply boundaries, never paths.
	for line := range strings.SplitAfterSeq(diff, "\n") {
		if strings.HasPrefix(line, "diff --git ") {
			sizes = append(sizes, 0)
		}
		if len(sizes) > 0 {
			sizes[len(sizes)-1] += len(line)
		}
	}
	if len(paths) != len(sizes) {
		return nil, fmt.Errorf("review diff has %d file sections but Git lists %d paths", len(sizes), len(paths))
	}
	files := make([]reviewDiffFile, len(paths))
	for i, path := range paths {
		files[i] = reviewDiffFile{Path: path, Bytes: sizes[i]}
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].Bytes != files[j].Bytes {
			return files[i].Bytes > files[j].Bytes
		}
		return files[i].Path < files[j].Path
	})
	return files[:min(3, len(files))], nil
}

func reviewSizeRefusal(reason string, files []reviewDiffFile) error {
	descriptions := make([]string, 0, min(3, len(files)))
	for _, file := range files[:min(3, len(files))] {
		descriptions = append(descriptions, fmt.Sprintf("%s (%d bytes)", file.Path, file.Bytes))
	}
	return reviewPrePromptError{errors.New(reason + "; largest files in the diff: " + strings.Join(descriptions, ", ") + "; pre_pr_review.omit_paths in the base commit's Project Config leaves generated files out of the review")}
}

func reviewSizeRefusalForDiff(ctx context.Context, git preflight.GitRunner, root, base, head string, omitted []reviewOmittedPath, diff, reason string) error {
	files, err := reviewLargestDiffFiles(ctx, git, root, base, head, omitted, diff)
	if err != nil {
		return reviewPrePromptError{err}
	}
	return reviewSizeRefusal(reason, files)
}
