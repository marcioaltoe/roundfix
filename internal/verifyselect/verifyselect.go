// Package verifyselect maps repository changes to the test sets they can affect.
package verifyselect

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Set is a bit set of test sets selected by a path or change list.
type Set uint8

// NoSet selects no tests.
const NoSet Set = 0

const (
	// CoreSet selects the core tests.
	CoreSet Set = 1 << iota
	// BaselineSet selects the Baseline and skill tests.
	BaselineSet
	// BothSets selects every test set.
	BothSets = CoreSet | BaselineSet
)

// Selection is the changed-path evidence and the test sets it selects.
type Selection struct {
	Paths []string
	Sets  Set
}

// Names returns the selected set names in execution order.
func (sets Set) Names() []string {
	names := make([]string, 0, 2)
	if sets&CoreSet != 0 {
		names = append(names, "core")
	}
	if sets&BaselineSet != 0 {
		names = append(names, "baseline")
	}
	return names
}

// ClassifyPath returns the test sets a repository-relative path can affect.
func ClassifyPath(path string) Set {
	path = strings.TrimPrefix(filepath.ToSlash(filepath.Clean(path)), "./")
	switch {
	case path == "go.mod", path == "go.sum", path == "Makefile":
		return BothSets
	case inDirectory(path, "internal/baseline"),
		inDirectory(path, "internal/baselineacp"),
		inDirectory(path, "skills"),
		inDirectory(path, ".agents/skills"),
		strings.HasPrefix(path, "internal/cli/baseline_"):
		return BaselineSet
	case inDirectory(path, "cmd"), inDirectory(path, "internal"), strings.HasSuffix(path, ".go"):
		return CoreSet
	default:
		return BothSets
	}
}

// SelectPaths classifies and combines a fixture or previously listed change set.
func SelectPaths(paths []string) Selection {
	unique := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		path = strings.TrimPrefix(filepath.ToSlash(filepath.Clean(path)), "./")
		if path == "." || path == "" {
			continue
		}
		unique[path] = struct{}{}
	}

	selected := Selection{Paths: make([]string, 0, len(unique))}
	for path := range unique {
		selected.Paths = append(selected.Paths, path)
		selected.Sets |= ClassifyPath(path)
	}
	sort.Strings(selected.Paths)
	return selected
}

// ChangedPaths lists committed changes since the merge base with baseRef,
// staged changes, unstaged changes, and untracked files. Returned paths are
// repository-relative, unique, and sorted.
func ChangedPaths(ctx context.Context, repoRoot, baseRef string) ([]string, error) {
	repoRoot = strings.TrimSpace(repoRoot)
	baseRef = strings.TrimSpace(baseRef)
	if repoRoot == "" {
		return nil, errors.New("list changed paths: repository root is required")
	}
	if baseRef == "" {
		return nil, errors.New("list changed paths: base ref is required")
	}

	mergeBaseOutput, err := runGit(ctx, repoRoot, "merge-base", "HEAD", baseRef)
	if err != nil {
		return nil, fmt.Errorf("list changed paths: resolve merge base: %w", err)
	}
	mergeBase := strings.TrimSpace(string(mergeBaseOutput))
	if mergeBase == "" {
		return nil, errors.New("list changed paths: resolve merge base: git returned an empty revision")
	}

	commands := [][]string{
		{"diff", "--name-only", "-z", "--no-renames", mergeBase, "HEAD"},
		{"diff", "--cached", "--name-only", "-z", "--no-renames"},
		{"diff", "--name-only", "-z", "--no-renames"},
		{"ls-files", "--others", "--exclude-standard", "-z"},
	}
	paths := make([]string, 0)
	for _, args := range commands {
		output, runErr := runGit(ctx, repoRoot, args...)
		if runErr != nil {
			return nil, fmt.Errorf("list changed paths: %w", runErr)
		}
		paths = append(paths, splitNUL(output)...)
	}
	return SelectPaths(paths).Paths, nil
}

// Select lists and classifies changes. If Git cannot produce the change list,
// the returned selection contains both sets together with the diagnostic error.
func Select(ctx context.Context, repoRoot, baseRef string) (Selection, error) {
	paths, err := ChangedPaths(ctx, repoRoot, baseRef)
	if err != nil {
		return Selection{Sets: BothSets}, err
	}
	selection := SelectPaths(paths)
	if selection.Sets != BaselineSet {
		return selection, nil
	}

	listed, root, err := listGoPackages(ctx, repoRoot)
	if err != nil {
		selection.Sets = BothSets
		return selection, fmt.Errorf("select dependent test sets: %w", err)
	}
	importsBaseline, err := coreImportsBaseline(listed, root)
	if err != nil {
		selection.Sets = BothSets
		return selection, fmt.Errorf("select dependent test sets: %w", err)
	}
	if importsBaseline {
		selection.Sets |= CoreSet
	}
	return selection, nil
}

// Packages returns the concrete go test package arguments owned by one set.
func Packages(ctx context.Context, repoRoot string, set Set) ([]string, error) {
	if set != CoreSet && set != BaselineSet {
		return nil, fmt.Errorf("list packages: set must be core or baseline")
	}
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return nil, errors.New("list packages: repository root is required")
	}

	listed, root, err := listGoPackages(ctx, repoRoot)
	if err != nil {
		return nil, fmt.Errorf("list packages: %w", err)
	}

	packages := make([]string, 0)
	for _, listedPackage := range listed {
		relative, relErr := filepath.Rel(root, listedPackage.Dir)
		if relErr != nil {
			return nil, fmt.Errorf("list packages: make %q relative to repository: %w", listedPackage.Dir, relErr)
		}
		if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return nil, fmt.Errorf("list packages: package directory %q is outside repository", listedPackage.Dir)
		}
		packageSet := packageSetForDirectory(filepath.ToSlash(relative))
		if packageSet != set {
			continue
		}
		if relative == "." {
			packages = append(packages, ".")
			continue
		}
		packages = append(packages, "./"+filepath.ToSlash(relative))
	}
	sort.Strings(packages)
	return packages, nil
}

// BaselineCLITestNames returns the tests declared in
// internal/cli/baseline_*_test.go.
func BaselineCLITestNames(repoRoot string) ([]string, error) {
	pattern := filepath.Join(repoRoot, "internal", "cli", "baseline_*_test.go")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("derive Baseline CLI tests: glob files: %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("derive Baseline CLI tests: no files match %s", pattern)
	}

	names := make(map[string]struct{})
	fileSet := token.NewFileSet()
	for _, path := range files {
		parsed, parseErr := parser.ParseFile(fileSet, path, nil, 0)
		if parseErr != nil {
			return nil, fmt.Errorf("derive Baseline CLI tests: parse %s: %w", path, parseErr)
		}
		for _, declaration := range parsed.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Recv != nil || !strings.HasPrefix(function.Name.Name, "Test") {
				continue
			}
			names[function.Name.Name] = struct{}{}
		}
	}

	result := make([]string, 0, len(names))
	for name := range names {
		result = append(result, name)
	}
	sort.Strings(result)
	return result, nil
}

// BaselineCLITestPattern returns an anchored go test -run pattern containing
// exactly the Baseline CLI test declarations.
func BaselineCLITestPattern(repoRoot string) (string, error) {
	names, err := BaselineCLITestNames(repoRoot)
	if err != nil {
		return "", err
	}
	quoted := make([]string, len(names))
	for index, name := range names {
		quoted[index] = regexp.QuoteMeta(name)
	}
	return "^(" + strings.Join(quoted, "|") + ")$", nil
}

// Run executes the verify-select command and returns its process exit code.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("verify-select", flag.ContinueOnError)
	flags.SetOutput(stderr)
	baseRef := flags.String("base", "main", "base ref used to select changed tests")
	repoRoot := flags.String("repo", ".", "repository root")
	packageSet := flags.String("packages", "", "print packages in core or baseline")
	baselinePattern := flags.Bool("baseline-cli-pattern", false, "print the Baseline internal/cli test pattern")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "verify-select: positional arguments are not supported")
		return 2
	}
	if *packageSet != "" && *baselinePattern {
		fmt.Fprintln(stderr, "verify-select: -packages and -baseline-cli-pattern are mutually exclusive")
		return 2
	}

	if *baselinePattern {
		pattern, err := BaselineCLITestPattern(*repoRoot)
		if err != nil {
			fmt.Fprintf(stderr, "verify-select: %v\n", err)
			return 1
		}
		fmt.Fprintln(stdout, pattern)
		return 0
	}

	if *packageSet != "" {
		set, ok := parseSet(*packageSet)
		if !ok {
			fmt.Fprintf(stderr, "verify-select: unknown package set %q; want core or baseline\n", *packageSet)
			return 2
		}
		packages, err := Packages(ctx, *repoRoot, set)
		if err != nil {
			fmt.Fprintf(stderr, "verify-select: %v\n", err)
			return 1
		}
		for _, packagePath := range packages {
			fmt.Fprintln(stdout, packagePath)
		}
		return 0
	}

	selection, err := Select(ctx, *repoRoot, *baseRef)
	if err != nil {
		fmt.Fprintf(stderr, "verify-select: %v; selecting core and baseline\n", err)
	}
	for _, name := range selection.Sets.Names() {
		fmt.Fprintln(stdout, name)
	}
	return 0
}

func inDirectory(path, directory string) bool {
	return path == directory || strings.HasPrefix(path, directory+"/")
}

type goListPackage struct {
	Dir          string
	ImportPath   string
	Imports      []string
	TestImports  []string
	XTestImports []string
}

func decodePackages(output []byte) ([]goListPackage, error) {
	decoder := json.NewDecoder(bytes.NewReader(output))
	packages := make([]goListPackage, 0)
	for {
		var listedPackage goListPackage
		if err := decoder.Decode(&listedPackage); err != nil {
			if errors.Is(err, io.EOF) {
				return packages, nil
			}
			return nil, fmt.Errorf("decode go list output: %w", err)
		}
		packages = append(packages, listedPackage)
	}
}

func listGoPackages(ctx context.Context, repoRoot string) ([]goListPackage, string, error) {
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, "", fmt.Errorf("resolve repository root: %w", err)
	}
	output, err := runCommand(ctx, root, "go", "list", "-json", "./...")
	if err != nil {
		return nil, "", err
	}
	listed, err := decodePackages(output)
	if err != nil {
		return nil, "", err
	}
	return listed, root, nil
}

func coreImportsBaseline(packages []goListPackage, root string) (bool, error) {
	baselinePackages := make(map[string]struct{})
	for _, listedPackage := range packages {
		relative, err := filepath.Rel(root, listedPackage.Dir)
		if err != nil {
			return false, fmt.Errorf("make %q relative to repository: %w", listedPackage.Dir, err)
		}
		if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return false, fmt.Errorf("package directory %q is outside repository", listedPackage.Dir)
		}
		if packageSetForDirectory(filepath.ToSlash(relative)) == BaselineSet {
			baselinePackages[listedPackage.ImportPath] = struct{}{}
		}
	}

	for _, listedPackage := range packages {
		relative, err := filepath.Rel(root, listedPackage.Dir)
		if err != nil {
			return false, fmt.Errorf("make %q relative to repository: %w", listedPackage.Dir, err)
		}
		if packageSetForDirectory(filepath.ToSlash(relative)) != CoreSet {
			continue
		}
		for _, imports := range [][]string{listedPackage.Imports, listedPackage.TestImports, listedPackage.XTestImports} {
			for _, imported := range imports {
				if _, ok := baselinePackages[imported]; ok {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func packageSetForDirectory(directory string) Set {
	if inDirectory(directory, "internal/baseline") ||
		inDirectory(directory, "internal/baselineacp") ||
		inDirectory(directory, "skills") {
		return BaselineSet
	}
	return CoreSet
}

func parseSet(name string) (Set, bool) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "core":
		return CoreSet, true
	case "baseline":
		return BaselineSet, true
	default:
		return NoSet, false
	}
}

func splitNUL(output []byte) []string {
	parts := bytes.Split(output, []byte{0})
	paths := make([]string, 0, len(parts))
	for _, part := range parts {
		if len(part) != 0 {
			paths = append(paths, string(part))
		}
	}
	return paths
}

func runGit(ctx context.Context, repoRoot string, args ...string) ([]byte, error) {
	gitArgs := []string{"--no-optional-locks", "-C", repoRoot, "-c", "core.fsmonitor=false"}
	gitArgs = append(gitArgs, args...)
	return runCommand(ctx, "", "git", gitArgs...)
}

func runCommand(ctx context.Context, workDir, name string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = workDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return nil, fmt.Errorf("%s %s: %s: %w", name, strings.Join(args, " "), detail, err)
	}
	return stdout.Bytes(), nil
}
