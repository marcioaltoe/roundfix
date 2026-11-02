package spec

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	EvidenceLinkIgnored = "ignored"
	EvidenceLinkMissing = "missing"
)

// UnresolvedEvidenceLink is a local link of a QA Report the committed tree
// would not hold.
type UnresolvedEvidenceLink struct {
	Line        int
	Destination string
	Target      string
	Reason      string // EvidenceLinkIgnored or EvidenceLinkMissing
}

// UnresolvedEvidenceLinks classifies the report's local links in document order.
// Git lists tracked and addable targets once; outside Git only disk existence
// decides resolution. Read and Git failures are infrastructure errors.
func UnresolvedEvidenceLinks(ctx context.Context, reportPath string) ([]UnresolvedEvidenceLink, error) {
	content, err := os.ReadFile(reportPath)
	if err != nil {
		return nil, fmt.Errorf("read QA Report evidence links: %w", err)
	}
	// Git returns a physical root; use the same namespace for report targets.
	path, err := filepath.EvalSymlinks(reportPath)
	if err != nil {
		return nil, fmt.Errorf("resolve QA Report path: %w", err)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve QA Report path: %w", err)
	}
	rootOutput, gitErr := archiveGitRead(ctx, filepath.Dir(path), "rev-parse", "--show-toplevel")
	root := strings.TrimSpace(string(rootOutput))
	if gitErr != nil {
		// A broken worktree is not a plain directory. Preserve Git's error if
		// any ancestor has Git metadata, or the failure is not "outside Git".
		for dir := filepath.Dir(path); ; dir = filepath.Dir(dir) {
			_, statErr := os.Lstat(filepath.Join(dir, ".git"))
			if !errors.Is(statErr, os.ErrNotExist) {
				return nil, fmt.Errorf("resolve QA Report Git root: %w", errors.Join(gitErr, statErr))
			}
			if filepath.Dir(dir) == dir {
				break
			}
		}
		if ctx.Err() != nil || !strings.Contains(gitErr.Error(), "not a git repository") {
			return nil, fmt.Errorf("resolve QA Report Git root: %w", gitErr)
		}
	}
	links := RelativeLinks(path, content)
	relative := make([]string, len(links))
	exists := make([]bool, len(links))
	directories := make([]bool, len(links))
	args := []string{"ls-files", "-z", "--cached", "--others", "--exclude-standard", "--"}
	pathspecStart := len(args)
	for i, link := range links {
		if root != "" {
			rel, err := filepath.Rel(root, filepath.FromSlash(link.Target))
			if err != nil {
				return nil, fmt.Errorf("resolve QA evidence target %s: %w", link.Destination, err)
			}
			if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
				continue
			}
			relative[i] = filepath.ToSlash(rel)
		}
		info, err := os.Stat(filepath.FromSlash(link.Target))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("stat QA evidence target %s: %w", link.Destination, err)
		}
		exists[i], directories[i] = true, info.IsDir()
		if root != "" {
			args = append(args, ":(literal)"+relative[i])
		}
	}
	var listed []string
	if root != "" && len(args) > pathspecStart {
		output, err := archiveGitRead(ctx, root, args...)
		if err != nil {
			return nil, fmt.Errorf("list QA evidence targets: %w", err)
		}
		listed = strings.Split(string(output), "\x00")
	}
	var unresolved []UnresolvedEvidenceLink
	for i, link := range links {
		if link.Target == filepath.ToSlash(path) {
			continue
		}
		reason := EvidenceLinkMissing
		if exists[i] {
			if root == "" {
				continue
			}
			resolved := false
			for _, candidate := range listed {
				if candidate == relative[i] || (directories[i] && (relative[i] == "." && candidate != "" || strings.HasPrefix(candidate, relative[i]+"/"))) {
					resolved = true
					break
				}
			}
			if resolved {
				continue
			}
			reason = EvidenceLinkIgnored
		}
		unresolved = append(unresolved, UnresolvedEvidenceLink{Line: link.Line, Destination: link.Destination, Target: link.Target, Reason: reason})
	}
	return unresolved, nil
}

// EvidenceLinksError is the shared refusal of unresolved QA evidence links.
type EvidenceLinksError struct{ Links []UnresolvedEvidenceLink }

func (err EvidenceLinksError) Error() string {
	parts := make([]string, 0, len(err.Links))
	for _, link := range err.Links {
		parts = append(parts, fmt.Sprintf("%s (line %d, %s)", link.Destination, link.Line, link.Reason))
	}
	return "QA Report evidence links do not resolve in the committed tree: " + strings.Join(parts, "; ")
}

// RequireResolvedEvidenceLinks returns a shared refusal or an infrastructure
// error, or nil when every local link resolves.
func RequireResolvedEvidenceLinks(ctx context.Context, reportPath string) error {
	links, err := UnresolvedEvidenceLinks(ctx, reportPath)
	if err != nil {
		return err
	}
	if len(links) > 0 {
		return EvidenceLinksError{Links: links}
	}
	return nil
}
