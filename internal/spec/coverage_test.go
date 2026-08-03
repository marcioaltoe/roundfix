// Suite: repository test-function coverage equivalence
// Invariant: every repository package retains every recorded top-level test function.
// Boundary IN: go list, go test -list, and the Spec-owned coverage record.
// Boundary OUT: test bodies, production behavior, exported APIs, and external systems.

package spec

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

// The record travels with its Spec: it was recorded under the active tree
// and moved to the archive when Spec 0071 closed. Resolution tries the
// archived home first and falls back to the active path, so the harness
// survives the archive without weakening — a missing record at both homes
// still fails.
const coverageRecordArchivedPath = "docs/specs/_archived/0071-verification-cost/coverage-record.json"
const coverageRecordActivePath = "docs/specs/0071-verification-cost/coverage-record.json"

// CoverageRecord is the deterministic set of top-level test functions the Go
// suite discovers, grouped by every package in the repository package list.
type CoverageRecord struct {
	Packages map[string][]string `json:"packages"`
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

	recordPath := filepath.Join(repoRoot, coverageRecordArchivedPath)
	if _, statErr := os.Stat(recordPath); statErr != nil {
		recordPath = filepath.Join(repoRoot, coverageRecordActivePath)
	}
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

// collectCoverageRecord asks the toolchain for every test name in one pass.
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
func collectCoverageRecord(repoRoot string) (CoverageRecord, error) {
	listOutput, err := runGo(repoRoot, "test", "-buildvcs=false", "-list", "^Test", "./...")
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
			record.Packages[packagePath] = pending
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
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go %s: %w\n%s", strings.Join(args, " "), err, output)
	}
	return string(output), nil
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
	normalized := CoverageRecord{Packages: make(map[string][]string, len(record.Packages))}
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

		recordedSet := stringSet(recordedTests)
		actualSet := stringSet(actualTests)
		for _, testName := range recordedTests {
			if _, ok := actualSet[testName]; !ok {
				comparison.Regressions = append(
					comparison.Regressions,
					fmt.Sprintf(
						"coverage regression: package %q no longer executes %q",
						packagePath,
						testName,
					),
				)
			}
		}
		for _, testName := range actualTests {
			if _, ok := recordedSet[testName]; !ok {
				comparison.Additions = append(
					comparison.Additions,
					fmt.Sprintf(
						"coverage addition: package %q now executes %q",
						packagePath,
						testName,
					),
				)
			}
		}
	}
	return comparison
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
