package verifyselect_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"roundfix/internal/gittest"
	"sort"
	"strings"
	"testing"
)

func newGitRepository(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	gittest.InitRepo(t, repo, "-q")
	runGit(t, repo, "config", "user.name", "Roundfix Test")
	runGit(t, repo, "config", "user.email", "roundfix@example.test")
	runGit(t, repo, "config", "commit.gpgsign", "false")
	return repo
}

func runGit(t *testing.T, repo string, args ...string) string {
	t.Helper()
	command := exec.CommandContext(t.Context(), "git", append([]string{"-C", repo, "-c", "core.fsmonitor=false"}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

func writeFile(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

func findRepositoryRoot(t *testing.T) string {
	t.Helper()
	directory, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd(): %v", err)
	}
	for {
		if _, statErr := os.Stat(filepath.Join(directory, "go.mod")); statErr == nil {
			return directory
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			t.Fatal("repository root not found")
		}
		directory = parent
	}
}

func declaredTests(t *testing.T, glob string) []string {
	t.Helper()
	files, err := filepath.Glob(glob)
	if err != nil {
		t.Fatalf("Glob(%q): %v", glob, err)
	}
	declaration := regexp.MustCompile(`(?m)^func (Test[A-Za-z0-9_]+)\s*\(`)
	seen := make(map[string]struct{})
	for _, path := range files {
		contents, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Fatalf("ReadFile(%q): %v", path, readErr)
		}
		for _, match := range declaration.FindAllSubmatch(contents, -1) {
			seen[string(match[1])] = struct{}{}
		}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
