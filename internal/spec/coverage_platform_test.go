// Suite: platform-neutral repository test-function coverage
// Invariant: release platforms share one deterministic record, anchored to Go on the host.
// Boundary IN: release matrix, go list, Go test files, and host go test -list.
// Boundary OUT: foreign test binaries and production behavior.

package spec

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestCoverageCollectionListsEachPlatformOnlyTest(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	cmd := exec.Command("go", "mod", "init", "roundfix")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("initialize fixture module: %v\n%s", err, output)
	}
	files := map[string]string{
		"dist/npm/platforms.json": `[{"goos":"windows"},{"goos":"linux"},{"goos":"darwin"},{"goos":"linux"}]`,
		"fixture.go":              "package fixture\n",
		"common_test.go": `package fixture
import check "testing"
func TestCommon(t *check.T) {}
func TestUnnamed(*check.T) {}
func TestMain(m *check.M) {}
func Testlower(t *check.T) {}
func TestWrong(t *check.B) {}
func TestResult(t *check.T) int { return 0 }
func TestMany(a, b *check.T) {}
func TestGeneric[T any](t *check.T) {}
type helper struct{}
func (helper) TestMethod(t *check.T) {}
`,
		"external_test.go":             "package fixture_test\nimport verify \"testing\"\nfunc TestExternal(t *verify.T) {}\n",
		"platform_darwin_test.go":      "package fixture\nimport \"testing\"\nfunc TestDarwin(t *testing.T) {}\nfunc TestAcrossFiles(t *testing.T) {}\n",
		"platform_linux_test.go":       "package fixture\nimport \"testing\"\nfunc TestLinux(t *testing.T) {}\nfunc TestAcrossFiles(t *testing.T) {}\n",
		"platform_windows_test.go":     "//go:build windows\n\npackage fixture\nimport \"testing\"\nfunc TestWindows(t *testing.T) {}\nfunc TestAcrossFiles(t *testing.T) {}\n",
		"unix_test.go":                 "//go:build unix\n\npackage fixture\nimport \"testing\"\nfunc TestUnix(t *testing.T) {}\n",
		"docscontract_test.go":         "//go:build docscontract\n\npackage fixture\nimport \"testing\"\nfunc TestDocsContract(t *testing.T) {}\n",
		"empty/empty.go":               "package empty\n",
		"onlylinux/only_linux.go":      "package onlylinux\n",
		"onlylinux/only_linux_test.go": "package onlylinux\nimport \"testing\"\nfunc TestOnlyPackage(t *testing.T) {}\n",
		"docs/ignored/ignored_test.go": "package docs\nimport \"testing\"\nfunc TestIgnored(t *testing.T) {}\n",
	}
	for name, data := range files {
		writeCoverageFixture(t, root, name, data)
	}
	actual, err := collectCoverageRecord(root)
	if err != nil {
		t.Fatal(err)
	}
	want := CoverageRecord{
		Platforms: []string{"darwin", "linux", "windows"},
		Packages:  map[string][]string{"roundfix": {"TestAcrossFiles", "TestCommon", "TestDocsContract", "TestExternal", "TestUnnamed"}, "roundfix/empty": {}, "roundfix/onlylinux": {}},
		PlatformTests: map[string]map[string][]string{
			"roundfix":           {"TestDarwin": {"darwin"}, "TestLinux": {"linux"}, "TestUnix": {"darwin", "linux"}, "TestWindows": {"windows"}},
			"roundfix/onlylinux": {"TestOnlyPackage": {"linux"}},
		},
	}
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("record = %#v, want %#v", actual, want)
	}
	first, err := marshalCoverageRecord(actual)
	if err != nil {
		t.Fatal(err)
	}
	second, err := collectCoverageRecord(root)
	if err != nil {
		t.Fatal(err)
	}
	data, err := marshalCoverageRecord(second)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, data) {
		t.Fatal("consecutive static collections produced different bytes")
	}
	t.Logf("all release platforms collected without foreign execution on %s", runtime.GOOS)
}

func TestCoverageCollectionMatchesTheToolchainOnThisPlatform(t *testing.T) {
	t.Parallel()
	root := filepath.Join("..", "..")
	actual, err := collectCoverageRecordFor(root, []string{runtime.GOOS})
	if err != nil {
		t.Fatal(err)
	}
	toolchain, err := collectToolchainCoverage(root)
	if err != nil {
		t.Fatal(err)
	}
	// The unchanged toolchain parser describes this host's one platform.
	toolchain.Platforms = []string{runtime.GOOS}
	comparison := compareCoverageRecords(toolchain, actual)
	for _, difference := range append(comparison.Regressions, comparison.Additions...) {
		t.Error(difference)
	}
	t.Logf("static collection and go test -list compared on %s", runtime.GOOS)
}

func TestCompareCoverageRecordsReportsAPlatformRegression(t *testing.T) {
	t.Parallel()
	platforms := []string{"darwin", "linux", "windows"}
	record := func(common []string, limited map[string][]string) CoverageRecord {
		return CoverageRecord{Platforms: platforms, Packages: map[string][]string{"roundfix/example": common}, PlatformTests: map[string]map[string][]string{"roundfix/example": limited}}
	}
	for _, tc := range []struct {
		name                   string
		recorded, actual       CoverageRecord
		regressions, additions []string
	}{
		{"lost one platform", record([]string{"TestKept"}, nil), record(nil, map[string][]string{"TestKept": {"darwin", "windows"}}), []string{`coverage regression: package "roundfix/example" no longer executes "TestKept" on linux`}, nil},
		{"lost limited test", record(nil, map[string][]string{"TestUnix": {"darwin", "linux"}}), record(nil, nil), []string{`coverage regression: package "roundfix/example" no longer executes "TestUnix" on darwin, linux`}, nil},
		{"lost all platforms", record([]string{"TestGone"}, nil), record(nil, nil), []string{`coverage regression: package "roundfix/example" no longer executes "TestGone"`}, nil},
		{"gained limited test", record(nil, nil), record(nil, map[string][]string{"TestNew": {"windows"}}), nil, []string{`coverage addition: package "roundfix/example" now executes "TestNew" on windows`}},
		{"gained all platforms", record(nil, nil), record([]string{"TestNew"}, nil), nil, []string{`coverage addition: package "roundfix/example" now executes "TestNew"`}},
		{"lost release platform", CoverageRecord{Platforms: platforms, Packages: map[string][]string{}}, CoverageRecord{Platforms: []string{"darwin", "linux"}, Packages: map[string][]string{}}, []string{`coverage regression: platform "windows" is no longer listed`}, nil},
		{"gained release platform", CoverageRecord{Platforms: []string{"darwin", "linux"}, Packages: map[string][]string{}}, CoverageRecord{Platforms: platforms, Packages: map[string][]string{}}, nil, []string{`coverage addition: platform "windows" is now listed`}},
		{"lost package", CoverageRecord{Platforms: platforms, Packages: map[string][]string{"roundfix/empty": {}}}, CoverageRecord{Platforms: platforms, Packages: map[string][]string{}}, []string{`coverage regression: package "roundfix/empty" is no longer listed`}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := compareCoverageRecords(tc.recorded, tc.actual)
			if !equalStrings(got.Regressions, tc.regressions) || !equalStrings(got.Additions, tc.additions) {
				t.Fatalf("comparison = %#v, want regressions %q additions %q", got, tc.regressions, tc.additions)
			}
		})
	}
}

func TestCoverageRecordRefusesAMalformedPlatformEntry(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                  string
		platforms, builtOn    []string
		packagePath, testName string
		common                []string
	}{
		{name: "unsorted platforms", platforms: []string{"linux", "darwin"}},
		{name: "duplicate platforms", platforms: []string{"linux", "linux"}},
		{name: "uppercase platform", platforms: []string{"Linux"}},
		{name: "empty token", platforms: []string{""}},
		{name: "punctuation token", platforms: []string{"linux-amd64"}},
		{name: "missing package", platforms: []string{"darwin", "linux", "windows"}, packagePath: "absent", testName: "TestLimited", builtOn: []string{"linux"}},
		{name: "empty subset", platforms: []string{"darwin", "linux", "windows"}, testName: "TestLimited"},
		{name: "whole set", platforms: []string{"darwin", "linux", "windows"}, testName: "TestLimited", builtOn: []string{"darwin", "linux", "windows"}},
		{name: "unsorted subset", platforms: []string{"darwin", "linux", "windows"}, testName: "TestLimited", builtOn: []string{"linux", "darwin"}},
		{name: "duplicate subset", platforms: []string{"darwin", "linux", "windows"}, testName: "TestLimited", builtOn: []string{"linux", "linux"}},
		{name: "foreign platform", platforms: []string{"darwin", "linux", "windows"}, testName: "TestLimited", builtOn: []string{"freebsd"}},
		{name: "no platforms", testName: "TestLimited", builtOn: []string{"linux"}},
		{name: "test in both places", platforms: []string{"darwin", "linux", "windows"}, testName: "TestLimited", builtOn: []string{"linux"}, common: []string{"TestLimited"}},
		{name: "invalid test name", platforms: []string{"darwin", "linux", "windows"}, testName: "Testlower", builtOn: []string{"linux"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			record := CoverageRecord{Platforms: tc.platforms, Packages: map[string][]string{"roundfix": tc.common}}
			if tc.testName != "" {
				packagePath := tc.packagePath
				if packagePath == "" {
					packagePath = "roundfix"
				}
				record.PlatformTests = map[string]map[string][]string{packagePath: {tc.testName: tc.builtOn}}
			}
			if err := validateCoverageRecord(record); err == nil {
				t.Fatal("malformed platform entry accepted")
			}
			// Exercise reader validation without normalization hiding disorder.
			path := filepath.Join(t.TempDir(), "record.json")
			encoded, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, encoded, 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := readCoverageRecord(path); err == nil {
				t.Fatal("reader accepted malformed platform entry")
			}
		})
	}
	t.Run("marshal sorts without mutating", func(t *testing.T) {
		record := CoverageRecord{Platforms: []string{"windows", "linux", "darwin"}, Packages: map[string][]string{"roundfix": {"TestZulu", "TestAlpha"}}, PlatformTests: map[string]map[string][]string{"roundfix": {"TestUnix": {"linux", "darwin"}}}}
		data, err := marshalCoverageRecord(record)
		if err != nil {
			t.Fatal(err)
		}
		again, err := marshalCoverageRecord(record)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, again) {
			t.Fatal("marshals differ")
		}
		if record.Platforms[0] != "windows" || record.PlatformTests["roundfix"]["TestUnix"][0] != "linux" || record.Packages["roundfix"][0] != "TestZulu" {
			t.Fatal("marshal mutated input")
		}
		var normalized CoverageRecord
		if err := json.Unmarshal(data, &normalized); err != nil {
			t.Fatal(err)
		}
		if err := validateCoverageRecord(normalized); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCoveragePlatformsComeFromTheReleaseMatrix(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeCoverageFixture(t, root, "dist/npm/platforms.json", `[{"goos":"windows"},{"goos":"freebsd"},{"goos":"linux"},{"goos":"freebsd"},{"goos":"plan9"}]`)
	got, err := coveragePlatforms(root)
	if err != nil {
		t.Fatal(err)
	}
	if !equalStrings(got, []string{"freebsd", "linux", "plan9", "windows"}) {
		t.Fatalf("platforms = %q", got)
	}
	for _, tc := range []struct{ name, matrix string }{
		{"empty", `[]`}, {"null", `null`}, {"invalid JSON", `[`}, {"missing goos", `[{}]`}, {"invalid goos", `[{"goos":"Linux"}]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeCoverageFixture(t, root, "dist/npm/platforms.json", tc.matrix)
			if _, err := coveragePlatforms(root); err == nil {
				t.Fatal("malformed matrix accepted")
			}
		})
	}
	t.Run("unreadable", func(t *testing.T) {
		if _, err := coveragePlatforms(t.TempDir()); err == nil || !strings.Contains(err.Error(), "read release matrix") {
			t.Fatalf("error = %v, want unreadable matrix", err)
		}
	})
}

func writeCoverageFixture(t *testing.T, root, name, data string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}
