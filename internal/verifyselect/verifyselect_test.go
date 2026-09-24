// Suite: selective verification
// Invariant: every declared change class selects exactly the test sets it can affect.
// Boundary IN: path classification, local Git change discovery, and set definitions.
// Boundary OUT: Makefile wiring and the whole-suite partition contract.
package verifyselect_test

import (
	"bytes"
	"fmt"
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
		{name: "Markdown in core package", path: "internal/app/README.md", want: verifyselect.CoreSet},
		{name: "documentation", path: "docs/user-guide/commands.md", want: verifyselect.NoSet},
		{name: "root Markdown", path: "README.md", want: verifyselect.NoSet},
		{name: "unrelated data", path: "testdata/input.json", want: verifyselect.BothSets},
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

func TestFixtureChangesSelectTheOwningSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want verifyselect.Set
	}{
		{name: "core package fixture", path: "internal/app/testdata/version.golden", want: verifyselect.CoreSet},
		{name: "Baseline package fixture", path: "internal/baseline/testdata/profile.golden", want: verifyselect.BaselineSet},
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

func TestUnknownNonDocumentationPathsFailSafe(t *testing.T) {
	t.Parallel()

	selection := verifyselect.SelectPaths([]string{".github/workflows/unknown.yml"})
	if selection.Sets != verifyselect.BothSets {
		t.Fatalf("SelectPaths() sets = %v, want %v", selection.Sets.Names(), verifyselect.BothSets.Names())
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

func TestSelectListsStagedPaths(t *testing.T) {
	repo := newGitRepository(t)
	writeFile(t, repo, "base.txt", "base\n")
	runGit(t, repo, "add", "base.txt")
	runGit(t, repo, "commit", "-q", "-m", "base")
	base := strings.TrimSpace(runGit(t, repo, "rev-parse", "HEAD"))

	writeFile(t, repo, "internal/app/staged.go", "package app\n")
	runGit(t, repo, "add", "internal/app/staged.go")

	selection, err := verifyselect.Select(t.Context(), repo, base)
	if err != nil {
		t.Fatalf("Select() error = %v", err)
	}
	wantPaths := []string{"internal/app/staged.go"}
	if !reflect.DeepEqual(selection.Paths, wantPaths) {
		t.Fatalf("Select() paths = %v, want %v", selection.Paths, wantPaths)
	}
	if selection.Sets != verifyselect.CoreSet {
		t.Fatalf("Select() sets = %v, want %v", selection.Sets.Names(), verifyselect.CoreSet.Names())
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

func TestBaselineChangesAlsoRunCoreImporters(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, repo, "go.mod", "module example.test/fixture\n\ngo 1.26\n")
	writeFile(t, repo, "internal/baseline/base.go", "package baseline\n")
	writeFile(t, repo, "core/direct/direct.go", "package direct\n\nimport _ \"example.test/fixture/internal/baseline\"\n")
	writeFile(t, repo, "core/transitive/transitive.go", "package transitive\n\nimport _ \"example.test/fixture/core/direct\"\n")
	writeFile(t, repo, "core/testonly/testonly.go", "package testonly\n")
	writeFile(t, repo, "core/testonly/testonly_test.go", "package testonly_test\n\nimport _ \"example.test/fixture/internal/baseline\"\n")
	writeFile(t, repo, "core/unrelated/unrelated.go", "package unrelated\n")

	packages, err := verifyselect.Packages(t.Context(), repo, verifyselect.BaselineSet)
	if err != nil {
		t.Fatalf("Packages(baseline) error = %v", err)
	}
	want := []string{"./core/direct", "./core/testonly", "./core/transitive", "./internal/baseline"}
	if !reflect.DeepEqual(packages, want) {
		t.Fatalf("Packages(baseline) = %v, want %v", packages, want)
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

func TestPartitionCoversEveryTestExactlyOnce(t *testing.T) {
	repoRoot := findRepositoryRoot(t)

	packages := listedPackages(t, repoRoot)
	corePackages, err := verifyselect.Packages(t.Context(), repoRoot, verifyselect.CoreSet)
	if err != nil {
		t.Fatalf("Packages(core) error = %v", err)
	}
	baselinePackages, err := verifyselect.Packages(t.Context(), repoRoot, verifyselect.BaselineSet)
	if err != nil {
		t.Fatalf("Packages(baseline) error = %v", err)
	}
	corePackageSet := makeNameSet(corePackages)
	baselinePackageSet := makeNameSet(baselinePackages)

	cliTests := listedTests(t, repoRoot, "./internal/cli")
	baselinePatternText, err := verifyselect.BaselineCLITestPattern(repoRoot)
	if err != nil {
		t.Fatalf("BaselineCLITestPattern() error = %v", err)
	}
	baselinePattern, err := regexp.Compile(baselinePatternText)
	if err != nil {
		t.Fatalf("BaselineCLITestPattern() returned invalid regexp %q: %v", baselinePatternText, err)
	}
	coreCLITests := make(nameSet)
	baselineCLITests := make(nameSet)
	for _, name := range cliTests {
		if baselinePattern.MatchString(name) {
			baselineCLITests[name] = struct{}{}
			continue
		}
		coreCLITests[name] = struct{}{}
	}

	t.Run("tree satisfies the partition", func(t *testing.T) {
		if err := partitionError(packages, corePackageSet, baselinePackageSet); err != nil {
			t.Errorf("package partition: %v", err)
		}
		if err := partitionError(cliTests, coreCLITests, baselineCLITests); err != nil {
			t.Errorf("internal/cli test partition: %v", err)
		}
	})

	t.Run("package omitted from both sets is named", func(t *testing.T) {
		missing := packages[0]
		coreWithout := cloneNameSet(corePackageSet)
		baselineWithout := cloneNameSet(baselinePackageSet)
		delete(coreWithout, missing)
		delete(baselineWithout, missing)

		err := partitionError(packages, coreWithout, baselineWithout)
		if err == nil {
			t.Fatalf("partitionError() error = nil after removing %q from both package sets", missing)
		}
		want := fmt.Sprintf("%q is selected by 0 sets, want exactly 1", missing)
		if err.Error() != want {
			t.Fatalf("partitionError() error = %q, want %q", err, want)
		}
	})

	t.Run("CLI test selected by both invocations is named", func(t *testing.T) {
		overlap := cliTests[0]
		coreWithOverlap := cloneNameSet(coreCLITests)
		baselineWithOverlap := cloneNameSet(baselineCLITests)
		coreWithOverlap[overlap] = struct{}{}
		baselineWithOverlap[overlap] = struct{}{}

		err := partitionError(cliTests, coreWithOverlap, baselineWithOverlap)
		if err == nil {
			t.Fatalf("partitionError() error = nil after selecting %q in both CLI invocations", overlap)
		}
		want := fmt.Sprintf("%q is selected by 2 sets, want exactly 1", overlap)
		if err.Error() != want {
			t.Fatalf("partitionError() error = %q, want %q", err, want)
		}
	})
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

type nameSet map[string]struct{}

func listedPackages(t *testing.T, repoRoot string) []string {
	t.Helper()
	command := exec.CommandContext(t.Context(), "go", "list", "-f", "{{.Dir}}", "./...")
	command.Dir = repoRoot
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("go list ./...: %v\n%s", err, output)
	}

	packages := make([]string, 0)
	for _, directory := range strings.Fields(string(output)) {
		relative, relErr := filepath.Rel(repoRoot, directory)
		if relErr != nil {
			t.Fatalf("make package directory %q relative to repository: %v", directory, relErr)
		}
		if relative == "." {
			packages = append(packages, ".")
			continue
		}
		packages = append(packages, "./"+filepath.ToSlash(relative))
	}
	if len(packages) == 0 {
		t.Fatal("go list ./... returned no packages")
	}
	sort.Strings(packages)
	return packages
}

func listedTests(t *testing.T, repoRoot, packagePath string) []string {
	t.Helper()
	command := exec.CommandContext(t.Context(), "go", "test", "-list", "^Test", packagePath)
	command.Dir = repoRoot
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("go test -list ^Test %s: %v\n%s", packagePath, err, output)
	}

	tests := make([]string, 0)
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Test") {
			tests = append(tests, line)
		}
	}
	if len(tests) == 0 {
		t.Fatalf("go test -list ^Test %s returned no tests", packagePath)
	}
	sort.Strings(tests)
	return tests
}

func makeNameSet(names []string) nameSet {
	set := make(nameSet, len(names))
	for _, name := range names {
		set[name] = struct{}{}
	}
	return set
}

func cloneNameSet(source nameSet) nameSet {
	clone := make(nameSet, len(source))
	for name := range source {
		clone[name] = struct{}{}
	}
	return clone
}

func partitionError(items []string, sets ...nameSet) error {
	for _, item := range items {
		memberships := 0
		for _, set := range sets {
			if _, ok := set[item]; ok {
				memberships++
			}
		}
		if memberships != 1 {
			return fmt.Errorf("%q is selected by %d sets, want exactly 1", item, memberships)
		}
	}
	return nil
}
