// Suite: selective verification
// Invariant: every declared change class selects exactly the test sets it can affect.
// Boundary IN: path classification, local Git change discovery, and set definitions.
// Boundary OUT: Makefile wiring and the whole-suite partition contract.
package verifyselect_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/verifyselect"
)

func TestClassifyPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want verifyselect.Set
	}{
		{name: "Baseline package", path: "internal/baseline/catalog.go", want: verifyselect.BaselineSet},
		{name: "Baseline ACP package", path: "internal/baselineacp/session.go", want: verifyselect.BaselineSet},
		{name: "embedded skill", path: "skills/roundfix/SKILL.md", want: verifyselect.BaselineSet},
		{name: "canonical skill", path: ".agents/skills/roundfix/SKILL.md", want: verifyselect.BaselineSet},
		{name: "Baseline CLI source", path: "internal/cli/baseline_plan.go", want: verifyselect.BaselineSet},
		{name: "Baseline CLI test", path: "internal/cli/baseline_plan_test.go", want: verifyselect.BaselineSet},
		{name: "module manifest", path: "go.mod", want: verifyselect.BothSets},
		{name: "module sums", path: "go.sum", want: verifyselect.BothSets},
		{name: "Makefile", path: "Makefile", want: verifyselect.BothSets},
		{name: "core source", path: "internal/app/version.go", want: verifyselect.CoreSet},
		{name: "core test", path: "internal/app/version_test.go", want: verifyselect.CoreSet},
		{name: "documentation", path: "docs/user-guide/commands.md", want: verifyselect.NoSet},
		{name: "unrelated data", path: "testdata/input.json", want: verifyselect.NoSet},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := verifyselect.ClassifyPath(test.path); got != test.want {
				t.Fatalf("ClassifyPath(%q) = %v, want %v", test.path, got.Names(), test.want.Names())
			}
		})
	}
}

func TestSelectClassifiesFixtureChangeSets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		paths []string
		want  verifyselect.Set
	}{
		{name: "core only", paths: []string{"internal/app/version.go"}, want: verifyselect.CoreSet},
		{name: "Baseline only", paths: []string{"skills/roundfix/SKILL.md"}, want: verifyselect.BaselineSet},
		{name: "both path classes", paths: []string{"internal/app/version.go", "internal/baseline/catalog.go"}, want: verifyselect.BothSets},
		{name: "module file", paths: []string{"go.mod"}, want: verifyselect.BothSets},
		{name: "documentation only", paths: []string{"docs/user-guide/commands.md"}, want: verifyselect.NoSet},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			selection := verifyselect.SelectPaths(test.paths)
			if selection.Sets != test.want {
				t.Fatalf("SelectPaths(%v) = %v, want %v", test.paths, selection.Sets.Names(), test.want.Names())
			}
		})
	}
}

func TestSelectListsCommittedUnstagedAndUntrackedPaths(t *testing.T) {
	repo := newGitRepository(t)
	writeFile(t, repo, "base.txt", "base\n")
	runGit(t, repo, "add", "base.txt")
	runGit(t, repo, "commit", "-q", "-m", "base")
	base := strings.TrimSpace(runGit(t, repo, "rev-parse", "HEAD"))

	writeFile(t, repo, "internal/app/committed.go", "package app\n")
	runGit(t, repo, "add", "internal/app/committed.go")
	runGit(t, repo, "commit", "-q", "-m", "core change")
	writeFile(t, repo, "base.txt", "unstaged\n")
	writeFile(t, repo, "skills/example/SKILL.md", "untracked\n")

	selection, err := verifyselect.Select(t.Context(), repo, base)
	if err != nil {
		t.Fatalf("Select() error = %v", err)
	}
	wantPaths := []string{"base.txt", "internal/app/committed.go", "skills/example/SKILL.md"}
	if !reflect.DeepEqual(selection.Paths, wantPaths) {
		t.Fatalf("Select() paths = %v, want %v", selection.Paths, wantPaths)
	}
	if selection.Sets != verifyselect.BothSets {
		t.Fatalf("Select() sets = %v, want %v", selection.Sets.Names(), verifyselect.BothSets.Names())
	}
}

func TestSelectFailsSafeToBothSets(t *testing.T) {
	t.Run("unresolvable base", func(t *testing.T) {
		repo := newGitRepository(t)
		writeFile(t, repo, "base.txt", "base\n")
		runGit(t, repo, "add", "base.txt")
		runGit(t, repo, "commit", "-q", "-m", "base")

		selection, err := verifyselect.Select(t.Context(), repo, "refs/heads/does-not-exist")
		if err == nil {
			t.Fatal("Select() error = nil, want unresolved-base diagnostic")
		}
		if selection.Sets != verifyselect.BothSets {
			t.Fatalf("Select() sets = %v, want %v", selection.Sets.Names(), verifyselect.BothSets.Names())
		}
	})

	t.Run("failing Git invocation", func(t *testing.T) {
		selection, err := verifyselect.Select(t.Context(), t.TempDir(), "main")
		if err == nil {
			t.Fatal("Select() error = nil, want Git diagnostic")
		}
		if selection.Sets != verifyselect.BothSets {
			t.Fatalf("Select() sets = %v, want %v", selection.Sets.Names(), verifyselect.BothSets.Names())
		}
	})
}

func TestPackagesExposeEachSet(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, repo, "go.mod", "module example.test/fixture\n\ngo 1.26\n")
	writeFile(t, repo, "core/core.go", "package core\n")
	writeFile(t, repo, "internal/baseline/base.go", "package baseline\n")
	writeFile(t, repo, "internal/baselineacp/acp.go", "package baselineacp\n")
	writeFile(t, repo, "skills/check/check.go", "package check\n")

	core, err := verifyselect.Packages(t.Context(), repo, verifyselect.CoreSet)
	if err != nil {
		t.Fatalf("Packages(core) error = %v", err)
	}
	baseline, err := verifyselect.Packages(t.Context(), repo, verifyselect.BaselineSet)
	if err != nil {
		t.Fatalf("Packages(baseline) error = %v", err)
	}
	if want := []string{"./core"}; !reflect.DeepEqual(core, want) {
		t.Fatalf("Packages(core) = %v, want %v", core, want)
	}
	wantBaseline := []string{"./internal/baseline", "./internal/baselineacp", "./skills/check"}
	if !reflect.DeepEqual(baseline, wantBaseline) {
		t.Fatalf("Packages(baseline) = %v, want %v", baseline, wantBaseline)
	}
}

func TestBaselineCLITestPatternIsDerived(t *testing.T) {
	repoRoot := findRepositoryRoot(t)
	names, err := verifyselect.BaselineCLITestNames(repoRoot)
	if err != nil {
		t.Fatalf("BaselineCLITestNames() error = %v", err)
	}
	patternText, err := verifyselect.BaselineCLITestPattern(repoRoot)
	if err != nil {
		t.Fatalf("BaselineCLITestPattern() error = %v", err)
	}
	pattern, err := regexp.Compile(patternText)
	if err != nil {
		t.Fatalf("BaselineCLITestPattern() returned invalid regexp %q: %v", patternText, err)
	}

	wantNames := declaredTests(t, filepath.Join(repoRoot, "internal", "cli", "baseline_*_test.go"))
	if !reflect.DeepEqual(names, wantNames) {
		t.Fatalf("BaselineCLITestNames() = %v, want declarations %v", names, wantNames)
	}
	allCLITests := declaredTests(t, filepath.Join(repoRoot, "internal", "cli", "*_test.go"))
	wantSet := make(map[string]bool, len(wantNames))
	for _, name := range wantNames {
		wantSet[name] = true
	}
	for _, name := range allCLITests {
		if got, want := pattern.MatchString(name), wantSet[name]; got != want {
			t.Errorf("pattern.MatchString(%q) = %t, want %t", name, got, want)
		}
	}
	if pattern.MatchString("TestNotDeclared") {
		t.Fatal("Baseline CLI pattern matches an undeclared test")
	}
}

func TestRunPrintsRequestedDefinitions(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, repo, "go.mod", "module example.test/fixture\n\ngo 1.26\n")
	writeFile(t, repo, "core/core.go", "package core\n")
	writeFile(t, repo, "internal/baseline/base.go", "package baseline\n")
	writeFile(t, repo, "internal/cli/baseline_fixture_test.go", "package cli\n\nimport \"testing\"\n\nfunc TestBaselineFixture(t *testing.T) {}\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := verifyselect.Run(t.Context(), []string{"-repo", repo, "-packages", "baseline"}, &stdout, &stderr); code != 0 {
		t.Fatalf("Run(packages) code = %d, stderr = %q", code, stderr.String())
	}
	if got, want := stdout.String(), "./internal/baseline\n"; got != want {
		t.Fatalf("Run(packages) stdout = %q, want %q", got, want)
	}

	stdout.Reset()
	stderr.Reset()
	if code := verifyselect.Run(t.Context(), []string{"-repo", repo, "-baseline-cli-pattern"}, &stdout, &stderr); code != 0 {
		t.Fatalf("Run(pattern) code = %d, stderr = %q", code, stderr.String())
	}
	if got, want := strings.TrimSpace(stdout.String()), "^(TestBaselineFixture)$"; got != want {
		t.Fatalf("Run(pattern) stdout = %q, want %q", got, want)
	}
}

func newGitRepository(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runGit(t, repo, "init", "-q")
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
