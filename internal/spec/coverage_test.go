// Suite: repository test-function coverage equivalence
// Invariant: every repository package retains every recorded top-level test function.
// Boundary IN: go list, go test -list, and the Spec-owned coverage record.
// Boundary OUT: test bodies, production behavior, exported APIs, and external systems.

package spec

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

var updateCoverageRecord = flag.Bool(
	"update-coverage-record",
	false,
	"re-record the repository test function coverage record",
)

const coverageRecordPath = "docs/references/coverage-record.json"

// CoverageRecord lists top-level tests across the release platforms. Packages
// holds common tests; PlatformTests holds each limited test's platform set.
type CoverageRecord struct {
	Platforms     []string                       `json:"platforms,omitempty"`
	Packages      map[string][]string            `json:"packages"`
	PlatformTests map[string]map[string][]string `json:"platformTests,omitempty"`
}

type coverageComparison struct {
	Regressions []string
	Additions   []string
}

func TestCoverageEquivalence(t *testing.T) {
	t.Parallel()
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	actual, err := collectCoverageRecord(repoRoot)
	if err != nil {
		t.Fatalf("collect coverage record: %v", err)
	}

	recordPath := filepath.Join(repoRoot, coverageRecordPath)
	if *updateCoverageRecord {
		if err := writeCoverageRecord(recordPath, actual); err != nil {
			t.Fatalf("re-record coverage at %s: %v", recordPath, err)
		}
		t.Logf("re-recorded %s", recordPath)
	}

	recorded, err := readCoverageRecord(recordPath)
	if err != nil {
		t.Fatalf(
			"read coverage record: %v; regenerate deliberately with -update-coverage-record",
			err,
		)
	}
	comparison := compareCoverageRecords(recorded, actual)
	for _, addition := range comparison.Additions {
		t.Log(addition)
	}
	for _, regression := range comparison.Regressions {
		t.Error(regression)
	}
}

func TestCompareCoverageRecordsReportsMissingTest(t *testing.T) {
	t.Parallel()

	recorded := CoverageRecord{Packages: map[string][]string{
		"roundfix/internal/spec": {"TestKept", "TestRemoved"},
	}}
	actual := CoverageRecord{Packages: map[string][]string{
		"roundfix/internal/spec": {"TestKept"},
	}}

	comparison := compareCoverageRecords(recorded, actual)
	want := []string{
		`coverage regression: package "roundfix/internal/spec" no longer executes "TestRemoved"`,
	}
	if !equalStrings(comparison.Regressions, want) {
		t.Fatalf("regressions = %q, want %q", comparison.Regressions, want)
	}
	if len(comparison.Additions) != 0 {
		t.Fatalf("additions = %q, want none", comparison.Additions)
	}
}

func TestCompareCoverageRecordsReportsAddedTestWithoutRegression(t *testing.T) {
	t.Parallel()

	recorded := CoverageRecord{Packages: map[string][]string{
		"roundfix/internal/spec": {"TestKept"},
	}}
	actual := CoverageRecord{Packages: map[string][]string{
		"roundfix/internal/spec": {"TestAdded", "TestKept"},
	}}

	comparison := compareCoverageRecords(recorded, actual)
	if len(comparison.Regressions) != 0 {
		t.Fatalf("regressions = %q, want none", comparison.Regressions)
	}
	want := []string{
		`coverage addition: package "roundfix/internal/spec" now executes "TestAdded"`,
	}
	if !equalStrings(comparison.Additions, want) {
		t.Fatalf("additions = %q, want %q", comparison.Additions, want)
	}
}

func TestMarshalCoverageRecordIsDeterministic(t *testing.T) {
	t.Parallel()

	record := CoverageRecord{Packages: map[string][]string{
		"roundfix/internal/spec": {"TestZulu", "TestAlpha"},
		"roundfix/internal/app":  {},
	}}
	first, err := marshalCoverageRecord(record)
	if err != nil {
		t.Fatalf("first marshal: %v", err)
	}
	second, err := marshalCoverageRecord(record)
	if err != nil {
		t.Fatalf("second marshal: %v", err)
	}
	if string(first) != string(second) {
		t.Fatalf("consecutive marshals differ:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	if strings.Index(string(first), "TestAlpha") > strings.Index(string(first), "TestZulu") {
		t.Fatalf("test names are not sorted:\n%s", first)
	}
}

// coveragePlatforms uses the release matrix as the sole platform authority.
func coveragePlatforms(repoRoot string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(repoRoot, "dist/npm/platforms.json"))
	if err != nil {
		return nil, fmt.Errorf("read release matrix: %w", err)
	}
	var matrix []struct {
		GOOS string `json:"goos"`
	}
	if err := json.Unmarshal(data, &matrix); err != nil {
		return nil, fmt.Errorf("decode release matrix: %w", err)
	}
	set := map[string]struct{}{}
	for _, entry := range matrix {
		if !isCoveragePlatform(entry.GOOS) {
			return nil, fmt.Errorf("invalid release platform %q", entry.GOOS)
		}
		set[entry.GOOS] = struct{}{}
	}
	platforms := sortedCoverageKeys(set)
	if len(platforms) == 0 {
		return nil, fmt.Errorf("release matrix is empty")
	}
	return platforms, nil
}

func collectCoverageRecord(repoRoot string) (CoverageRecord, error) {
	platforms, err := coveragePlatforms(repoRoot)
	if err != nil {
		return CoverageRecord{}, err
	}
	return collectCoverageRecordFor(repoRoot, platforms)
}

func collectCoverageRecordFor(repoRoot string, platforms []string) (CoverageRecord, error) {
	record := CoverageRecord{Platforms: append([]string{}, platforms...), Packages: map[string][]string{}, PlatformTests: map[string]map[string][]string{}}
	if len(platforms) == 0 {
		return CoverageRecord{}, fmt.Errorf("collection platforms are empty")
	}
	if err := validateCoverageRecord(record); err != nil {
		return CoverageRecord{}, err
	}
	// Cache by filename across platforms; go list owns all build constraints.
	files := map[string][]string{}
	tests := map[string]map[string][]string{}
	for _, platform := range platforms {
		cmd := exec.Command("go", "list", "-tags", "docscontract", "-f", "{{.ImportPath}}\t{{.Dir}}\t{{join .TestGoFiles \",\"}}\t{{join .XTestGoFiles \",\"}}", "./...")
		cmd.Dir = repoRoot
		cmd.Env = append(os.Environ(), "GOOS="+platform, "CGO_ENABLED=0", "GOWORK=off")
		// stdout only: go writes "go: downloading" progress to stderr when a
		// platform needs a module the host has not fetched.
		output, err := cmd.Output()
		if err != nil {
			return CoverageRecord{}, fmt.Errorf("list coverage packages on %s: %w\n%s", platform, err, commandStderr(err))
		}
		for _, line := range strings.Split(strings.TrimSuffix(string(output), "\n"), "\n") {
			fields := strings.Split(line, "\t")
			if len(fields) != 4 {
				return CoverageRecord{}, fmt.Errorf("malformed package listing on %s: %q", platform, line)
			}
			packagePath := fields[0]
			if strings.HasPrefix(packagePath, "roundfix/docs/") {
				continue
			}
			if tests[packagePath] == nil {
				tests[packagePath] = map[string][]string{}
			}
			names := map[string]struct{}{}
			for _, fileList := range fields[2:] {
				if fileList == "" {
					continue
				}
				for _, filename := range strings.Split(fileList, ",") {
					path := filepath.Join(fields[1], filename)
					parsed, ok := files[path]
					if !ok {
						parsed, err = coverageFileTests(path)
						if err != nil {
							return CoverageRecord{}, err
						}
						files[path] = parsed
					}
					for _, name := range parsed {
						names[name] = struct{}{}
					}
				}
			}
			for name := range names {
				tests[packagePath][name] = append(tests[packagePath][name], platform)
			}
		}
	}
	for packagePath, names := range tests {
		record.Packages[packagePath] = []string{}
		for name, builtOn := range names {
			if len(builtOn) == len(platforms) {
				record.Packages[packagePath] = append(record.Packages[packagePath], name)
			} else {
				if record.PlatformTests[packagePath] == nil {
					record.PlatformTests[packagePath] = map[string][]string{}
				}
				record.PlatformTests[packagePath][name] = builtOn
			}
		}
		sort.Strings(record.Packages[packagePath])
	}
	if err := validateCoverageRecord(record); err != nil {
		return CoverageRecord{}, fmt.Errorf("validate static coverage: %w", err)
	}
	return record, nil
}

func coverageFileTests(path string) ([]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse coverage file %s: %w", path, err)
	}
	testingName := ""
	for _, imp := range file.Imports {
		importPath, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			return nil, fmt.Errorf("decode import in %s: %w", path, err)
		}
		if importPath == "testing" {
			testingName = "testing"
			if imp.Name != nil {
				testingName = imp.Name.Name
			}
		}
	}
	var names []string
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || !isTestFunctionName(fn.Name.Name) || fn.Type.TypeParams.NumFields() != 0 || fn.Type.Results.NumFields() != 0 || fn.Type.Params.NumFields() != 1 {
			continue
		}
		pointer, ok := fn.Type.Params.List[0].Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		selector, ok := pointer.X.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "T" {
			continue
		}
		qualifier, ok := selector.X.(*ast.Ident)
		if !ok || testingName == "" || testingName == "_" || qualifier.Name != testingName {
			continue
		}
		names = append(names, fn.Name.Name)
	}
	sort.Strings(names)
	return names, nil
}

func sortedCoverageKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func isCoveragePlatform(platform string) bool {
	if platform == "" {
		return false
	}
	for index, ch := range platform {
		if (ch < 'a' || ch > 'z') && (index == 0 || ch < '0' || ch > '9') {
			return false
		}
	}
	return true
}

// collectToolchainCoverage asks the toolchain for every test name in one pass.
//
// `go test -list` accepts a package pattern, so `./...` answers for the whole
// repository at once. Asking package by package instead meant one `go test`
// invocation per package — each compiling that package's test binary before
// printing names it already knew — which cost this package roughly 30s of the
// suite. One invocation costs under a second.
//
// The output interleaves: a package's test names print before its own
// terminating `ok <pkg>` or `? <pkg>` line, so names accumulate until a
// terminator names the package they belong to.
func collectToolchainCoverage(repoRoot string) (CoverageRecord, error) {
	// -tags docscontract keeps the pull-request-boundary domain enumerated:
	// without it the moved tests would vanish from the record silently.
	listOutput, err := runGo(repoRoot, "test", "-buildvcs=false", "-tags", "docscontract", "-list", "^Test", "./...")
	if err != nil {
		return CoverageRecord{}, fmt.Errorf("list repository tests: %w", err)
	}

	record := CoverageRecord{Packages: map[string][]string{}}
	var pending []string
	for _, line := range strings.Split(listOutput, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if packagePath, terminated := coveragePackageTerminator(line); terminated {
			sort.Strings(pending)
			if !strings.HasPrefix(packagePath, "roundfix/docs/") {
				record.Packages[packagePath] = pending
			}
			pending = nil
			continue
		}
		if strings.HasPrefix(line, "Test") {
			pending = append(pending, line)
		}
	}
	if len(pending) > 0 {
		return CoverageRecord{}, fmt.Errorf("listed tests with no terminating package line: %v", pending)
	}
	if err := validateCoverageRecord(record); err != nil {
		return CoverageRecord{}, fmt.Errorf("validate collected coverage: %w", err)
	}
	return record, nil
}

// coveragePackageTerminator recognises the lines `go test -list` prints to
// close out one package: `ok  <pkg> <elapsed>` when it has tests, and
// `?  <pkg> [no test files]` when it has none.
func coveragePackageTerminator(line string) (string, bool) {
	for _, prefix := range []string{"ok ", "? ", "FAIL "} {
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return "", false
		}
		return fields[1], true
	}
	return "", false
}

func runGo(repoRoot string, args ...string) (string, error) {
	cmd := exec.Command("go", args...)
	cmd.Dir = repoRoot
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("go %s: %w\n%s", strings.Join(args, " "), err, commandStderr(err))
	}
	return string(output), nil
}

func commandStderr(err error) []byte {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.Stderr
	}
	return nil
}

func listedTestNames(output string) []string {
	var names []string
	for _, line := range strings.Split(output, "\n") {
		name := strings.TrimSpace(line)
		if token.IsIdentifier(name) && isTestFunctionName(name) {
			names = append(names, name)
		}
	}
	return names
}

func isTestFunctionName(name string) bool {
	const prefix = "Test"
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) {
		return true
	}
	next, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(next)
}

func readCoverageRecord(path string) (CoverageRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return CoverageRecord{}, fmt.Errorf("read %s: %w", path, err)
	}
	var record CoverageRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return CoverageRecord{}, fmt.Errorf("decode %s: %w", path, err)
	}
	if err := validateCoverageRecord(record); err != nil {
		return CoverageRecord{}, fmt.Errorf("validate %s: %w", path, err)
	}
	return record, nil
}

func writeCoverageRecord(path string, record CoverageRecord) error {
	data, err := marshalCoverageRecord(record)
	if err != nil {
		return fmt.Errorf("encode coverage record: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func marshalCoverageRecord(record CoverageRecord) ([]byte, error) {
	normalized := CoverageRecord{Platforms: append([]string{}, record.Platforms...), Packages: make(map[string][]string, len(record.Packages)), PlatformTests: map[string]map[string][]string{}}
	sort.Strings(normalized.Platforms)
	for packagePath, tests := range record.PlatformTests {
		normalized.PlatformTests[packagePath] = map[string][]string{}
		for name, platforms := range tests {
			normalized.PlatformTests[packagePath][name] = append([]string{}, platforms...)
			sort.Strings(normalized.PlatformTests[packagePath][name])
		}
	}
	for packagePath, tests := range record.Packages {
		normalized.Packages[packagePath] = append([]string{}, tests...)
		sort.Strings(normalized.Packages[packagePath])
	}
	if err := validateCoverageRecord(normalized); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(normalized, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal coverage record: %w", err)
	}
	return append(data, '\n'), nil
}

func validateCoverageRecord(record CoverageRecord) error {
	for index, platform := range record.Platforms {
		if !isCoveragePlatform(platform) || (index > 0 && record.Platforms[index-1] >= platform) {
			return fmt.Errorf("platforms are not strictly sorted lower-case tokens")
		}
	}
	platformSet := stringSet(record.Platforms)
	for packagePath, tests := range record.PlatformTests {
		common, exists := record.Packages[packagePath]
		if !exists {
			return fmt.Errorf("platform tests package %q is missing from packages", packagePath)
		}
		commonSet := stringSet(common)
		for name, platforms := range tests {
			if !token.IsIdentifier(name) || !isTestFunctionName(name) {
				return fmt.Errorf("package %s has invalid platform test name %q", packagePath, name)
			}
			if _, exists := commonSet[name]; exists {
				return fmt.Errorf("package %s test %s is listed in both places", packagePath, name)
			}
			if len(platforms) == 0 || len(platforms) >= len(record.Platforms) {
				return fmt.Errorf("package %s test %s platforms must be a nonempty proper subset", packagePath, name)
			}
			for index, platform := range platforms {
				if _, exists := platformSet[platform]; !exists {
					return fmt.Errorf("package %s test %s has unknown platform %q", packagePath, name, platform)
				}
				if index > 0 && platforms[index-1] >= platform {
					return fmt.Errorf("package %s test %s platforms are not strictly sorted", packagePath, name)
				}
			}
		}
	}
	if record.Packages == nil {
		return fmt.Errorf("packages map is missing")
	}
	for packagePath, tests := range record.Packages {
		if packagePath == "" {
			return fmt.Errorf("package path is empty")
		}
		for index, testName := range tests {
			if !token.IsIdentifier(testName) || !isTestFunctionName(testName) {
				return fmt.Errorf("package %s has invalid test name %q", packagePath, testName)
			}
			if index > 0 && tests[index-1] >= testName {
				return fmt.Errorf("package %s test names are not strictly sorted", packagePath)
			}
		}
	}
	return nil
}

func compareCoverageRecords(recorded, actual CoverageRecord) coverageComparison {
	packageSet := make(map[string]struct{}, len(recorded.Packages)+len(actual.Packages))
	for packagePath := range recorded.Packages {
		packageSet[packagePath] = struct{}{}
	}
	for packagePath := range actual.Packages {
		packageSet[packagePath] = struct{}{}
	}
	packages := make([]string, 0, len(packageSet))
	for packagePath := range packageSet {
		packages = append(packages, packagePath)
	}
	sort.Strings(packages)

	var comparison coverageComparison
	for _, platform := range coveragePlatformDifference(coverageComparisonPlatforms(recorded), coverageComparisonPlatforms(actual)) {
		comparison.Regressions = append(comparison.Regressions, fmt.Sprintf("coverage regression: platform %q is no longer listed", platform))
	}
	for _, platform := range coveragePlatformDifference(coverageComparisonPlatforms(actual), coverageComparisonPlatforms(recorded)) {
		comparison.Additions = append(comparison.Additions, fmt.Sprintf("coverage addition: platform %q is now listed", platform))
	}
	for _, packagePath := range packages {
		recordedTests, wasRecorded := recorded.Packages[packagePath]
		actualTests, isPresent := actual.Packages[packagePath]
		if !isPresent {
			comparison.Regressions = append(
				comparison.Regressions,
				fmt.Sprintf("coverage regression: package %q is no longer listed", packagePath),
			)
		}
		if !wasRecorded {
			comparison.Additions = append(
				comparison.Additions,
				fmt.Sprintf("coverage addition: package %q is now listed", packagePath),
			)
		}

		recordedByTest := coverageTestPlatforms(recorded, packagePath, recordedTests)
		actualByTest := coverageTestPlatforms(actual, packagePath, actualTests)
		names := map[string]struct{}{}
		for name := range recordedByTest {
			names[name] = struct{}{}
		}
		for name := range actualByTest {
			names[name] = struct{}{}
		}
		for _, name := range sortedCoverageKeys(names) {
			lost := coveragePlatformDifference(recordedByTest[name], actualByTest[name])
			gained := coveragePlatformDifference(actualByTest[name], recordedByTest[name])
			if len(lost) > 0 {
				message := fmt.Sprintf("coverage regression: package %q no longer executes %q", packagePath, name)
				if !equalStrings(lost, coverageComparisonPlatforms(recorded)) {
					message += " on " + strings.Join(lost, ", ")
				}
				comparison.Regressions = append(comparison.Regressions, message)
			}
			if len(gained) > 0 {
				message := fmt.Sprintf("coverage addition: package %q now executes %q", packagePath, name)
				if !equalStrings(gained, coverageComparisonPlatforms(actual)) {
					message += " on " + strings.Join(gained, ", ")
				}
				comparison.Additions = append(comparison.Additions, message)
			}
		}
	}
	return comparison
}

func coverageComparisonPlatforms(record CoverageRecord) []string {
	if len(record.Platforms) == 0 {
		return []string{""}
	}
	return record.Platforms
}

func coverageTestPlatforms(record CoverageRecord, packagePath string, common []string) map[string][]string {
	tests := map[string][]string{}
	for _, name := range common {
		tests[name] = coverageComparisonPlatforms(record)
	}
	for name, platforms := range record.PlatformTests[packagePath] {
		tests[name] = platforms
	}
	return tests
}

func coveragePlatformDifference(left, right []string) []string {
	rightSet := stringSet(right)
	var difference []string
	for _, platform := range left {
		if _, exists := rightSet[platform]; !exists {
			difference = append(difference, platform)
		}
	}
	sort.Strings(difference)
	return difference
}

func stringSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func TestCoverageRecordCountsNoPackageUnderDocs(t *testing.T) {
	record, err := readCoverageRecord(filepath.Join("..", "..", coverageRecordPath))
	if err != nil {
		t.Fatal(err)
	}
	for packagePath := range record.Packages {
		if strings.HasPrefix(packagePath, "roundfix/docs/") {
			t.Errorf("documentation package counted: %s", packagePath)
		}
	}
	actual, err := collectCoverageRecord(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	for packagePath := range actual.Packages {
		if strings.HasPrefix(packagePath, "roundfix/docs/") {
			t.Errorf("documentation package collected: %s", packagePath)
		}
	}
}
