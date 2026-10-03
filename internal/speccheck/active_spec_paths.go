package speccheck

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"roundfix/internal/spec"
)

// CodeSpecPathPinned identifies a file that depends on an active Spec directory.
const CodeSpecPathPinned = "SC-SPEC-PATH-PINNED"

// SpecPathPin locates a repository-relative, slash-separated path dependency.
type SpecPathPin struct {
	Path string
	Line int
}

func activeSpecArchiveRoot(repoRoot, specsRoot string) string {
	return spec.ArchiveSpecRoot(specsRoot, filepath.Clean(specsRoot) == filepath.Join(filepath.Clean(repoRoot), "docs", "specs"))
}

// ActiveSpecPathPins lists non-Markdown files outside the Spec Root, its
// archive root and docs/history that name the active Spec's directory.
func ActiveSpecPathPins(repoRoot, specsRoot, slug string) ([]SpecPathPin, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve repository root: %w", err)
	}
	specsRoot, err = filepath.Abs(specsRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve Spec Root: %w", err)
	}
	relative, err := filepath.Rel(repoRoot, specsRoot)
	if err != nil {
		return nil, fmt.Errorf("relativize Spec Root: %w", err)
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return nil, nil
	}
	needle := filepath.ToSlash(filepath.Join(relative, slug))
	pattern := regexp.MustCompile(regexp.QuoteMeta(needle) + `(?:/|\b|$)`)
	excluded := []string{specsRoot, activeSpecArchiveRoot(repoRoot, specsRoot), filepath.Join(repoRoot, "docs", "history"), filepath.Join(repoRoot, ".git")}
	skip := func(path string) bool {
		for _, root := range excluded {
			if path == root || strings.HasPrefix(path, root+string(filepath.Separator)) {
				return true
			}
		}
		return strings.EqualFold(filepath.Ext(path), ".md")
	}
	var pins []SpecPathPin
	probe := exec.Command("git", "-c", "core.fsmonitor=false", "rev-parse", "--is-inside-work-tree")
	probe.Dir = repoRoot
	output, probeErr := probe.Output()
	if probeErr == nil && strings.TrimSpace(string(output)) == "true" {
		args := []string{"grep", "-n", "-z", "-I", "-F", "--untracked", "-e", needle, "--", ".", ":(exclude,icase)*.md"}
		for _, root := range excluded {
			rel, relErr := filepath.Rel(repoRoot, root)
			if relErr != nil {
				return nil, fmt.Errorf("relativize excluded root: %w", relErr)
			}
			args = append(args, ":(exclude,literal)"+filepath.ToSlash(rel))
		}
		command := exec.Command("git", append([]string{"-c", "core.fsmonitor=false"}, args...)...)
		command.Dir = repoRoot
		var stderr bytes.Buffer
		command.Stderr = &stderr
		output, err = command.Output()
		if err != nil {
			var exit *exec.ExitError
			if !errors.As(err, &exit) || exit.ExitCode() != 1 {
				return nil, fmt.Errorf("search active Spec paths: %w: %s", err, stderr.String())
			}
		}
		for len(output) > 0 {
			path, rest, ok := bytes.Cut(output, []byte{0})
			if !ok {
				return nil, fmt.Errorf("parse active Spec path search: missing path separator")
			}
			line, rest, ok := bytes.Cut(rest, []byte{0})
			if !ok {
				return nil, fmt.Errorf("parse active Spec path search: missing line separator")
			}
			content, next, _ := bytes.Cut(rest, []byte{'\n'})
			output = next
			number, parseErr := strconv.Atoi(string(line))
			if parseErr != nil {
				return nil, fmt.Errorf("parse active Spec path line: %w", parseErr)
			}
			if pattern.Match(content) {
				pins = append(pins, SpecPathPin{Path: string(path), Line: number})
			}
		}
	} else {
		if probeErr != nil {
			var exit *exec.ExitError
			if !errors.As(probeErr, &exit) {
				return nil, fmt.Errorf("inspect Git work tree: %w", probeErr)
			}
		}
		err = filepath.WalkDir(repoRoot, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if skip(path) {
				if entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if !entry.Type().IsRegular() {
				return nil
			}
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			if bytes.ContainsRune(content, 0) {
				return nil
			}
			rel, relErr := filepath.Rel(repoRoot, path)
			if relErr != nil {
				return relErr
			}
			for index, line := range bytes.Split(content, []byte{'\n'}) {
				if pattern.Match(line) {
					pins = append(pins, SpecPathPin{Path: filepath.ToSlash(rel), Line: index + 1})
				}
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walk active Spec paths: %w", err)
		}
	}
	sort.Slice(pins, func(i, j int) bool {
		if pins[i].Path != pins[j].Path {
			return pins[i].Path < pins[j].Path
		}
		return pins[i].Line < pins[j].Line
	})
	return pins, nil
}
