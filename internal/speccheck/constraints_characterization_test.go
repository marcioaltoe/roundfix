// Suite: Spec Consistency Check characterization corpus
// Invariant: report-authored regressions and repository-wide findings remain observable, while active Specs carry no errors.
// Boundary IN: public speccheck API, replay fixtures, and every active and archived repository Spec
// Boundary OUT: CLI rendering and exit-code policy
package speccheck_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

const (
	replay0058QA001  = "replay-0058-qa-001"
	replay0058QA004  = "replay-0058-qa-004"
	replay0056F001   = "replay-0056-f-001"
	replay0056F002   = "replay-0056-f-002"
	replay0060Task03 = "replay-0060-task-03"
	corpusBudget     = time.Second
)

func TestCheckReplay0060Task03RefusesWorkIndependentVerification(t *testing.T) {
	t.Parallel()

	const findingPath = "docs/findings/_archived/2026-07-31-a-rehearsal-task-can-settle-completed-without-rehearsing.md"
	result := checkFixture(t, replay0060Task03)
	finding := requireReplayFinding(t, findingPath, result, "SC-VERIFY-WORK-INDEPENDENT", "cannot distinguish Task work from no work")
	assertReplayLocations(t, findingPath, finding,
		speccheck.Location{Path: "docs/specs/" + replay0060Task03 + "/task_03.md", Line: 44},
	)
	report := speccheck.RenderText(result)
	for _, want := range []string{
		"SC-VERIFY-WORK-INDEPENDENT",
		"docs/specs/" + replay0060Task03 + "/task_03.md:44",
		"fix: Add a declared Verification command that asserts this Task's own effect.",
	} {
		if !strings.Contains(report, want) {
			t.Errorf("replay of %s: text finding does not contain %q:\n%s", findingPath, want, report)
		}
	}
	provenance := readReplayFile(t, replay0060Task03, "README.md")
	for _, want := range []string{findingPath, "exact Verification commands", "status is `pending`"} {
		if !strings.Contains(provenance, want) {
			t.Errorf("replay provenance does not contain %q:\n%s", want, provenance)
		}
	}
}

func TestCheckReplay0060Task03RefusesContradictoryRequirementsAndUndeclaredRehearsal(t *testing.T) {
	t.Parallel()

	const findingPath = "docs/findings/2026-07-31-a-rehearsal-task-can-settle-completed-without-rehearsing.md"
	result := checkFixture(t, replay0060Task03)

	contradiction := requireReplayFinding(t, findingPath, result, speccheck.CodeRequirementContradictory, "commit")
	assertReplayLocations(t, findingPath, contradiction,
		speccheck.Location{Path: "docs/specs/" + replay0060Task03 + "/task_03.md", Line: 13},
		speccheck.Location{Path: "docs/specs/" + replay0060Task03 + "/task_03.md", Line: 15},
	)

	rehearsal := requireReplayFinding(t, findingPath, result, speccheck.CodeRehearsalUndeclared, "Rehearsal Cases")
	assertReplayLocations(t, findingPath, rehearsal,
		speccheck.Location{Path: "docs/specs/" + replay0060Task03 + "/task_03.md", Line: 9},
	)
}

func TestCheckReplay0058QA001FromReport(t *testing.T) {
	t.Parallel()

	const reportPath = "docs/specs/_archived/0058-npm-trusted-publishing-and-release-preflight/qa/qa-report-2026-07-31.md"
	result := checkFixture(t, replay0058QA001)
	finding := requireReplayFinding(t, reportPath, result, speccheck.CodeCoverageUnmapped, "Core Feature 2")
	assertReplayLocations(t, reportPath, finding,
		speccheck.Location{Path: "docs/specs/" + replay0058QA001 + "/_prd.md", Line: 17},
		speccheck.Location{Path: "docs/specs/" + replay0058QA001 + "/_techspec.md", Line: 12},
	)
}

func TestCheckReplay0058QA004FromReport(t *testing.T) {
	t.Parallel()

	const reportPath = "docs/specs/_archived/0058-npm-trusted-publishing-and-release-preflight/qa/qa-report-2026-08-01.md"
	result := checkFixture(t, replay0058QA004)
	finding := requireReplayFinding(t, reportPath, result, speccheck.CodeVocabularyUndocumented, "publish:")
	assertReplayLocations(t, reportPath, finding,
		speccheck.Location{Path: "docs/specs/" + replay0058QA004 + "/workflow.sh", Line: 4},
		speccheck.Location{Path: "docs/specs/" + replay0058QA004 + "/runbook.md", Line: 1},
	)
	emitted := readReplayFile(t, replay0058QA004, "workflow.sh")
	documented := readReplayFile(t, replay0058QA004, "runbook.md")
	for _, token := range []string{"identity:", "publish:", "registry:", "runtime:", "undetermined:"} {
		if !strings.Contains(emitted, token) {
			t.Errorf("replay of %s: workflow does not emit %q", reportPath, token)
		}
	}
	if strings.Contains(documented, "publish:") {
		t.Errorf("replay of %s: runbook unexpectedly documents publish:", reportPath)
	}
	for _, token := range []string{"identity:", "registry:", "runtime:", "undetermined:"} {
		if !strings.Contains(documented, token) {
			t.Errorf("replay of %s: runbook does not document %q", reportPath, token)
		}
	}
}

func TestCheckReplay0056F001FromReport(t *testing.T) {
	t.Parallel()

	const reportPath = "docs/specs/_archived/0056-profiles-configure-merge-semantics/qa/qa-report-2026-08-01.md"
	result := checkFixture(t, replay0056F001)

	unlisted := requireReplayFinding(t, reportPath, result, speccheck.CodeADRUnlisted, "ADR-0086")
	assertReplayLocations(t, reportPath, unlisted,
		speccheck.Location{Path: "docs/specs/" + replay0056F001 + "/_techspec.md", Line: 12},
		speccheck.Location{Path: "docs/specs/" + replay0056F001 + "/_prd.md", Line: 11},
	)

	related := requireReplayFinding(t, reportPath, result, speccheck.CodeADRRelated, "ADR-0055")
	assertReplayLocations(t, reportPath, related,
		speccheck.Location{Path: "docs/adr/0055-exact-capability-proof.md", Line: 7},
		speccheck.Location{Path: "docs/specs/" + replay0056F001 + "/_prd.md", Line: 11},
	)
	relatedADR := readFixtureRepositoryFile(t, "docs", "adr", "0055-exact-capability-proof.md")
	for _, citation := range []string{"ADR-0039", "ADR-0049"} {
		if !strings.Contains(relatedADR, citation) {
			t.Errorf("replay of %s: related ADR does not cite %s", reportPath, citation)
		}
	}
}

func TestCheckReplay0056F002FromReport(t *testing.T) {
	t.Parallel()

	const reportPath = "docs/specs/_archived/0056-profiles-configure-merge-semantics/qa/qa-report-2026-08-01.md"
	result := checkFixture(t, replay0056F002)
	finding := requireReplayFinding(t, reportPath, result, speccheck.CodeCoverageUnmapped, "Core Feature 6")
	assertReplayLocations(t, reportPath, finding,
		speccheck.Location{Path: "docs/specs/" + replay0056F002 + "/_prd.md", Line: 16},
		speccheck.Location{Path: "docs/specs/" + replay0056F002 + "/_techspec.md", Line: 14},
	)
	if techSpec := readReplayFile(t, replay0056F002, "_techspec.md"); !strings.Contains(techSpec, "only added or replaced categories") {
		t.Errorf("replay of %s: TechSpec does not record the narrowed proof scope", reportPath)
	}
}

func TestCheckReplayReadmeProvenance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		slug       string
		reportPath string
	}{
		{
			name:       "0058 QA-001 report",
			slug:       replay0058QA001,
			reportPath: "docs/specs/_archived/0058-npm-trusted-publishing-and-release-preflight/qa/qa-report-2026-07-31.md",
		},
		{
			name:       "0058 QA-004 report",
			slug:       replay0058QA004,
			reportPath: "docs/specs/_archived/0058-npm-trusted-publishing-and-release-preflight/qa/qa-report-2026-08-01.md",
		},
		{
			name:       "0056 F-001 report",
			slug:       replay0056F001,
			reportPath: "docs/specs/_archived/0056-profiles-configure-merge-semantics/qa/qa-report-2026-08-01.md",
		},
		{
			name:       "0056 F-002 report",
			slug:       replay0056F002,
			reportPath: "docs/specs/_archived/0056-profiles-configure-merge-semantics/qa/qa-report-2026-08-01.md",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(fixtureSpecRoot, tt.slug, "README.md")
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read replay provenance %q: %v", path, err)
			}
			for _, required := range []string{tt.reportPath, "authored from the report", "not recovered from Git"} {
				if !strings.Contains(string(content), required) {
					t.Errorf("%s does not record %q", path, required)
				}
			}
		})
	}
}

func TestCheckCorpusGolden(t *testing.T) {
	repoRoot := characterizationRepositoryRoot(t)
	activeRoot := filepath.Join(repoRoot, "docs", "specs")
	archivedRoot := materializeArchivedCorpus(t, filepath.Join(activeRoot, "_archived"))
	want := readCorpusGolden(t)

	got := corpusGolden{
		Schema:   want.Schema,
		Update:   want.Update,
		Active:   sweepCorpus(t, activeRoot, repoRoot),
		Archived: sweepCorpus(t, archivedRoot, repoRoot),
	}
	if !reflect.DeepEqual(got, want) {
		actual, err := json.MarshalIndent(got, "", "  ")
		if err != nil {
			t.Fatalf("render actual corpus counts: %v", err)
		}
		t.Errorf("Spec corpus finding counts changed; inspect the detector change, then deliberately update testdata/corpus-golden.json:\n%s", actual)
	}
}

func TestCheckCorpusBudget(t *testing.T) {
	if !dedicatedCorpusBudgetRun() {
		t.Skip("wall-clock budget requires a dedicated run: go test ./internal/speccheck -run '^TestCheckCorpusBudget$'")
	}

	repoRoot := characterizationRepositoryRoot(t)
	activeRoot := filepath.Join(repoRoot, "docs", "specs")
	archivedRoot := materializeArchivedCorpus(t, filepath.Join(activeRoot, "_archived"))

	started := time.Now()
	sweepCorpus(t, activeRoot, repoRoot)
	sweepCorpus(t, archivedRoot, repoRoot)
	elapsed := time.Since(started)
	t.Logf("full Spec corpus sweep completed in %s (budget %s)", elapsed, corpusBudget)

	if elapsed >= corpusBudget {
		t.Errorf("full Spec corpus sweep took %s, want under %s", elapsed, corpusBudget)
	}
}

func dedicatedCorpusBudgetRun() bool {
	// An ordinary package sweep shares the machine with other packages and
	// cannot make a meaningful wall-clock assertion. Accept only the dedicated
	// selectors used by this Task's focused check and the serial gate step.
	run := flag.Lookup("test.run")
	if run == nil {
		return false
	}
	switch run.Value.String() {
	case "Budget", "^TestCheckCorpusBudget$":
		return true
	default:
		return false
	}
}

func TestCheckActiveCorpusHasNoErrors(t *testing.T) {
	repoRoot := characterizationRepositoryRoot(t)
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	active, err := spec.ListActive(specsRoot)
	if err != nil {
		t.Fatalf("list active Specs: %v", err)
	}

	for _, activeSpec := range active {
		result, err := speccheck.Check(specsRoot, repoRoot, activeSpec.Slug)
		if err != nil {
			t.Errorf("%s: check active Spec: %v", activeSpec.Slug, err)
			continue
		}
		for _, finding := range result.Findings {
			switch finding.Severity {
			case speccheck.SeverityError:
				t.Errorf("%s: %s: %s", activeSpec.Slug, finding.Code, finding.Summary)
			case speccheck.SeverityGap:
				t.Logf("%s: %s gap remains visible: %s", activeSpec.Slug, finding.Code, finding.Summary)
			}
		}
	}
}

type corpusGolden struct {
	Schema   string         `json:"schema"`
	Update   string         `json:"update"`
	Active   map[string]int `json:"active"`
	Archived map[string]int `json:"archived"`
}

func requireReplayFinding(t *testing.T, reportPath string, result speccheck.Result, code, summaryFragment string) speccheck.Finding {
	t.Helper()

	findings := findingsWithCode(result, code)
	if len(findings) != 1 {
		t.Fatalf("replay of %s: %s findings = %#v, want exactly one", reportPath, code, findings)
	}
	if !strings.Contains(findings[0].Summary, summaryFragment) {
		t.Fatalf("replay of %s: summary = %q, want %q", reportPath, findings[0].Summary, summaryFragment)
	}
	return findings[0]
}

func assertReplayLocations(t *testing.T, reportPath string, finding speccheck.Finding, want ...speccheck.Location) {
	t.Helper()

	for _, location := range want {
		if !hasExactLocation(finding, location.Path, location.Line) {
			t.Errorf("replay of %s: %s locations = %#v, want %#v", reportPath, finding.Code, finding.Where, location)
		}
	}
}

func readReplayFile(t *testing.T, slug, name string) string {
	t.Helper()

	return readFixtureRepositoryFile(t, "docs", "specs", slug, name)
}

func readFixtureRepositoryFile(t *testing.T, pathElements ...string) string {
	t.Helper()

	path := filepath.Join(append([]string{"testdata", "repo"}, pathElements...)...)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read replay fixture %q: %v", path, err)
	}
	return string(content)
}

func characterizationRepositoryRoot(t *testing.T) string {
	t.Helper()

	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(testFile), "..", ".."))
}

func readCorpusGolden(t *testing.T) corpusGolden {
	t.Helper()

	path := filepath.Join("testdata", "corpus-golden.json")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read corpus golden %q: %v", path, err)
	}
	var golden corpusGolden
	if err := json.Unmarshal(content, &golden); err != nil {
		t.Fatalf("parse corpus golden %q: %v", path, err)
	}
	if golden.Schema != "roundfix-speccheck-corpus/v1" || strings.TrimSpace(golden.Update) == "" {
		t.Fatalf("corpus golden %q must declare its schema and update path", path)
	}
	return golden
}

func sweepCorpus(t *testing.T, specsRoot, repoRoot string) map[string]int {
	t.Helper()

	counts := make(map[string]int, len(corpusFindingCodes))
	for _, code := range corpusFindingCodes {
		counts[code] = 0
	}
	entries, err := os.ReadDir(specsRoot)
	if err != nil {
		t.Fatalf("read Spec corpus %q: %v", specsRoot, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		if _, err := os.Stat(filepath.Join(specsRoot, entry.Name(), "_prd.md")); errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			t.Fatalf("inspect Spec %q: %v", entry.Name(), err)
		}
		result, err := speccheck.Check(specsRoot, repoRoot, entry.Name())
		if err != nil {
			t.Fatalf("Check(%q) in corpus %q: %v", entry.Name(), specsRoot, err)
		}
		for _, finding := range result.Findings {
			if _, known := counts[finding.Code]; !known {
				t.Fatalf("Check(%q) returned uncharacterized code %q", entry.Name(), finding.Code)
			}
			counts[finding.Code]++
		}
	}
	return counts
}

var corpusFindingCodes = []string{
	speccheck.CodeConstraintMissing,
	speccheck.CodeConstraintUnreasoned,
	speccheck.CodeConstraintSource,
	speccheck.CodeToolingUnauthorized,
	speccheck.CodeToolingUnbounded,
	speccheck.CodeADRUnlisted,
	speccheck.CodeADRRelated,
	speccheck.CodeCoverageUnmapped,
	speccheck.CodeCoverageUntasked,
	speccheck.CodeLoopOrderDivergent,
	speccheck.CodeFindingLifecycle,
	speccheck.CodeRollupMember,
	speccheck.CodeArchiveLicense,
	speccheck.CodeReferenceUnresolved,
	speccheck.CodeVocabularyUndocumented,
	speccheck.CodeVerifyWorkIndependent,
	speccheck.CodeRequirementContradictory,
	speccheck.CodeRehearsalUndeclared,
}

func materializeArchivedCorpus(t *testing.T, sourceRoot string) string {
	t.Helper()

	targetRoot := filepath.Join(t.TempDir(), "docs", "specs")
	err := filepath.WalkDir(sourceRoot, func(sourcePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(sourceRoot, sourcePath)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetRoot, relative)
		if entry.IsDir() {
			return os.MkdirAll(targetPath, 0o755)
		}
		if entry.Name() == "_prd.md" {
			content, err := os.ReadFile(sourcePath)
			if err != nil {
				return err
			}
			content, err = activeCorpusPRD(content)
			if err != nil {
				return fmt.Errorf("rehost archived PRD %q: %w", sourcePath, err)
			}
			return os.WriteFile(targetPath, content, 0o644)
		}
		if err := os.Link(sourcePath, targetPath); err == nil {
			return nil
		}
		return copyCorpusFile(sourcePath, targetPath)
	})
	if err != nil {
		t.Fatalf("materialize archived Spec corpus: %v", err)
	}
	return targetRoot
}

func activeCorpusPRD(content []byte) ([]byte, error) {
	const statusPrefix = "\nstatus:"
	statusStart := bytes.Index(content, []byte(statusPrefix))
	if statusStart < 0 {
		return nil, errors.New("frontmatter has no status field")
	}
	statusStart++
	statusEnd := bytes.IndexByte(content[statusStart:], '\n')
	if statusEnd < 0 {
		return nil, errors.New("frontmatter status field has no line ending")
	}
	statusEnd += statusStart

	active := make([]byte, 0, len(content))
	active = append(active, content[:statusStart]...)
	active = append(active, "status: active"...)
	active = append(active, content[statusEnd:]...)
	return active, nil
}

func copyCorpusFile(sourcePath, targetPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()
	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(target, source); err != nil {
		target.Close()
		return err
	}
	return target.Close()
}
