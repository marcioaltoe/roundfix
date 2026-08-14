package spec

// Suite: Spec archive eligibility and historical replay
// Invariant: only pass or fully declared rows archive, unproven actions remain visible, and the prior corpus is unchanged.
// Boundary IN: archived Task status, QA metadata, and the archive QA precondition.
// Boundary OUT: command I/O and filesystem movement in internal/cli/archive_test.go.

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestArchiveDirAnswersEveryRetiredKind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		kind ArchiveKind
		want string
	}{
		{name: "Specs", kind: ArchiveKindSpec, want: "docs/history/specs"},
		{name: "findings", kind: ArchiveKindFinding, want: "docs/history/findings"},
		{name: "ADRs", kind: ArchiveKindADR, want: "docs/history/adr"},
		{name: "backlog entries", kind: ArchiveKindBacklog, want: "docs/history/backlog"},
		{name: "Review Artifacts", kind: ArchiveKindReview, want: "docs/history/reviews"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := ArchiveDir(test.kind); got != test.want {
				t.Fatalf("ArchiveDir(%q) = %q, want %q", test.kind, got, test.want)
			}
		})
	}
}

func TestArchiveDirRejectsAnUnknownKind(t *testing.T) {
	t.Parallel()

	const unknown ArchiveKind = "invented"
	if got := ArchiveDir(unknown); got != "" {
		t.Fatalf("ArchiveDir(%q) = %q, want empty rejection", unknown, got)
	}
}

func TestArchiveSpecRootDefaultLayout(t *testing.T) {
	t.Parallel()

	// The built-in <repo>/docs/specs root resolves to the repository's default
	// docs/history/specs directory, whether the path is repository-relative or an
	// absolute path belonging to the repository.
	tests := []struct {
		name      string
		specsRoot string
		want      string
	}{
		{
			name:      "repository-relative default root",
			specsRoot: "docs/specs",
			want:      filepath.ToSlash(filepath.Join("docs", "history", "specs")),
		},
		{
			name:      "absolute repository default root",
			specsRoot: filepath.Join("/repo", "docs", "specs"),
			want:      filepath.ToSlash(filepath.Join("/repo", "docs", "history", "specs")),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := ArchiveSpecRoot(test.specsRoot, true); got != test.want {
				t.Fatalf("ArchiveSpecRoot(%q, true) = %q, want %q", test.specsRoot, got, test.want)
			}
		})
	}
}

func TestArchiveSpecRootNonDefaultKeepsArchiveBesideActiveRoot(t *testing.T) {
	t.Parallel()

	// A configured non-default Spec Root keeps its archive beside the active
	// root rather than under the repository's default docs/history/specs. Even a
	// non-default root whose path ends in docs/specs must NOT be classified as
	// the repository default.
	tests := []struct {
		name      string
		specsRoot string
		want      string
	}{
		{
			name:      "configured non-default root ending in docs specs",
			specsRoot: filepath.FromSlash("nested/docs/specs"),
			want:      filepath.ToSlash(filepath.Join("nested/docs/specs", "_archived")),
		},
		{
			name:      "configured non-default root with a plain name",
			specsRoot: filepath.FromSlash("configured-specs"),
			want:      filepath.ToSlash(filepath.Join("configured-specs", "_archived")),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := ArchiveSpecRoot(test.specsRoot, false); got != test.want {
				t.Fatalf("ArchiveSpecRoot(%q, false) = %q, want %q", test.specsRoot, got, test.want)
			}
		})
	}
}

func TestArchiveSpecRootExternalKeepsArchiveBesideActiveRoot(t *testing.T) {
	t.Parallel()

	// An external Spec Root keeps its archive beside the active root rather
	// than under the referring repository's default docs/history/specs. Issue 009:
	// an external root whose path happens to end in docs/specs must NOT be
	// misclassified as the repository default.
	tests := []struct {
		name      string
		specsRoot string
		want      string
	}{
		{
			name:      "external root ending in docs specs",
			specsRoot: filepath.FromSlash("/tmp/other/docs/specs"),
			want:      filepath.ToSlash(filepath.Join("/tmp/other/docs/specs", "_archived")),
		},
		{
			name:      "external root with a plain name",
			specsRoot: filepath.FromSlash("/tmp/other/external-specs"),
			want:      filepath.ToSlash(filepath.Join("/tmp/other/external-specs", "_archived")),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := ArchiveSpecRoot(test.specsRoot, false); got != test.want {
				t.Fatalf("ArchiveSpecRoot(%q, false) = %q, want %q", test.specsRoot, got, test.want)
			}
		})
	}
}

const (
	spec0058ReplaySlug    = "0058-npm-trusted-publishing-and-release-preflight"
	spec0058ReleaseAction = "a maintainer publishes a tagged release and records the run"
)

var spec0058SourceReport = archiveTestPath(
	ArchiveKindSpec,
	spec0058ReplaySlug,
	"qa",
	"qa-report-2026-08-01-04.md",
)

var spec0058SourceReportProvenance = filepath.ToSlash(filepath.Join(
	"_archived",
	"specs",
	spec0058ReplaySlug,
	"qa",
	"qa-report-2026-08-01-04.md",
))

func TestSpec0058ReplayArchivesDeclaredUnreachableRelease(t *testing.T) {
	t.Parallel()
	repositoryRoot := archiveTestRepositoryRoot(t)
	specsRoot, specDir := prepareSpec0058Replay(t, repositoryRoot, "accepted", true)

	prdBefore := archiveTestReadFile(t, filepath.Join(specDir, "_prd.md"))
	if strings.Contains(prdBefore, "qa_override:") {
		t.Fatalf("replay PRD retained qa_override:\n%s", prdBefore)
	}
	report, err := ReadQAReport(specDir)
	if err != nil {
		t.Fatalf("ReadQAReport(replay): %v", err)
	}
	if report.Verdict != VerdictPartial || report.RowsBlockedDeclared != 1 || report.RowsBlockedEnvironment != 0 || report.RowsBlockedFinding != 0 {
		t.Fatalf("replay QA Report = %+v, want partial with one declared-only blocked row", report)
	}

	result, err := Archive(ArchiveRequest{
		SpecsRoot:  specsRoot,
		Slug:       spec0058ReplaySlug,
		ArchivedAt: time.Date(2026, time.August, 4, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Archive(Spec 0058 replay): %v", err)
	}
	if _, err := os.Stat(specDir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("active replay still exists after archive: %v", err)
	}
	archivedPRD := archiveTestReadFile(t, filepath.Join(result.ArchivedDir, "_prd.md"))
	var frontmatter struct {
		Status     string   `yaml:"status"`
		QAOverride bool     `yaml:"qa_override"`
		Unproven   []string `yaml:"unproven"`
	}
	frontmatterBytes, _, err := splitFrontmatter([]byte(archivedPRD))
	if err != nil {
		t.Fatalf("parse archived replay PRD: %v", err)
	}
	if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
		t.Fatalf("decode archived replay PRD: %v", err)
	}
	if frontmatter.Status != "archived" || frontmatter.QAOverride {
		t.Fatalf("archived replay frontmatter = %+v, want archived without qa_override", frontmatter)
	}
	if !reflect.DeepEqual(frontmatter.Unproven, []string{spec0058ReleaseAction}) {
		t.Fatalf("archived replay unproven = %q, want release action %q", frontmatter.Unproven, spec0058ReleaseAction)
	}
}

func TestSpec0058ReplayRecordsFixtureProvenance(t *testing.T) {
	t.Parallel()
	repositoryRoot := archiveTestRepositoryRoot(t)
	provenance := archiveTestReadFile(t, filepath.Join(repositoryRoot, "internal", "spec", "testdata", "archive-replay-0058", "PROVENANCE.md"))
	for _, want := range []string{
		spec0058SourceReportProvenance,
		"The original Spec 0058 PRD has no `## Unreachable Acceptance` section",
		"The declaration overlay is added by Spec 0070",
	} {
		if !strings.Contains(provenance, want) {
			t.Errorf("replay provenance does not contain %q", want)
		}
	}
}

func TestSpec0058ReplayReportsWronglyDeclaredReachableRow(t *testing.T) {
	t.Parallel()
	repositoryRoot := archiveTestRepositoryRoot(t)
	specsRoot, specDir := prepareSpec0058Replay(t, repositoryRoot, "wrongly-declared", true)
	reportPath, err := NewestQAReport(specDir)
	if err != nil {
		t.Fatalf("NewestQAReport(wrongly declared replay): %v", err)
	}
	reportContent := archiveTestReadFile(t, reportPath)
	if !strings.Contains(reportContent, "wrongly-declared-row finding") {
		t.Fatalf("wrongly declared replay did not report the finding:\n%s", reportContent)
	}

	_, err = Archive(ArchiveRequest{SpecsRoot: specsRoot, Slug: spec0058ReplaySlug})
	if err == nil {
		t.Fatal("Archive accepted a replay whose declared row was reachable")
	}
	if !strings.Contains(err.Error(), `newest QA Report verdict is "fail"; expected "pass"`) {
		t.Fatalf("Archive wrongly declared replay error = %q", err)
	}
	archiveTestAssertReplayRemainsActive(t, specsRoot)
}

func TestSpec0058ReplayRefusesUnmatchedBlockedRow(t *testing.T) {
	t.Parallel()
	repositoryRoot := archiveTestRepositoryRoot(t)
	specsRoot, _ := prepareSpec0058Replay(t, repositoryRoot, "unmatched", false)

	_, err := Archive(ArchiveRequest{SpecsRoot: specsRoot, Slug: spec0058ReplaySlug})
	if err == nil {
		t.Fatal("Archive accepted a replay with a blocked row and no declaration")
	}
	for _, want := range []string{
		"rows_blocked_declared is 1",
		"Spec declares 0 unreachable acceptances",
		"shortfall is 1",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Archive unmatched replay error %q does not contain %q", err, want)
		}
	}
	archiveTestAssertReplayRemainsActive(t, specsRoot)
}

func TestArchivedPassCorpusRemainsArchiveEligible(t *testing.T) {
	t.Parallel()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	pattern := archiveTestRepositoryPath(filepath.Join(filepath.Dir(testFile), "..", ".."), ArchiveKindSpec, "*", "qa", "qa-report-*.md")
	reportPaths, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("find archived QA Reports: %v", err)
	}
	if len(reportPaths) == 0 {
		t.Fatal("archived QA Report corpus is empty")
	}

	seen := make(map[string]struct{})
	passSpecs := 0
	for _, reportPath := range reportPaths {
		specDir := filepath.Dir(filepath.Dir(reportPath))
		if _, ok := seen[specDir]; ok {
			continue
		}
		seen[specDir] = struct{}{}

		report, err := ReadQAReport(specDir)
		if err != nil {
			t.Fatalf("ReadQAReport(%q): %v", specDir, err)
		}
		if report.Verdict != VerdictPass {
			continue
		}
		passSpecs++

		taskPaths, err := filepath.Glob(filepath.Join(specDir, "task_*.md"))
		if err != nil {
			t.Fatalf("find archived Tasks for %q: %v", specDir, err)
		}
		if len(taskPaths) == 0 {
			t.Fatalf("archived pass Spec %q has no Tasks", specDir)
		}
		for _, taskPath := range taskPaths {
			content, err := os.ReadFile(taskPath)
			if err != nil {
				t.Fatalf("read archived Task %q: %v", taskPath, err)
			}
			frontmatterBytes, _, err := splitFrontmatter(content)
			if err != nil {
				t.Fatalf("parse archived Task %q: %v", taskPath, err)
			}
			var frontmatter struct {
				Status string `yaml:"status"`
			}
			if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
				t.Fatalf("parse archived Task %q frontmatter: %v", taskPath, err)
			}
			if status := NormalizeStatus(frontmatter.Status); status != string(StatusCompleted) {
				t.Fatalf("archived pass Task %q status is %q; expected %q", taskPath, status, StatusCompleted)
			}
		}

		unproven, err := archiveUnprovenActions(specDir, report)
		if err != nil {
			t.Fatalf("archive precondition rejected archived pass Spec %q: %v", specDir, err)
		}
		if len(unproven) != 0 {
			t.Fatalf("archive precondition added unproven actions to pass Spec %q: %q", specDir, unproven)
		}
	}
	if passSpecs == 0 {
		t.Fatal("archived corpus has no pass Specs")
	}
	t.Logf("checked archive eligibility for %d archived pass Specs", passSpecs)
}

func TestArchivedQAOverrideCorpusIncludesFailedSpec(t *testing.T) {
	t.Parallel()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	pattern := archiveTestRepositoryPath(filepath.Join(filepath.Dir(testFile), "..", ".."), ArchiveKindSpec, "*", "_prd.md")
	prdPaths, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("find archived Spec PRDs: %v", err)
	}
	for _, prdPath := range prdPaths {
		content, err := os.ReadFile(prdPath)
		if err != nil {
			t.Fatalf("read archived Spec PRD %q: %v", prdPath, err)
		}
		frontmatterBytes, _, err := splitFrontmatter(content)
		if err != nil {
			t.Fatalf("parse archived Spec PRD %q: %v", prdPath, err)
		}
		var frontmatter struct {
			Status     string `yaml:"status"`
			QAOverride bool   `yaml:"qa_override"`
		}
		if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
			t.Fatalf("parse archived Spec PRD %q frontmatter: %v", prdPath, err)
		}
		if !frontmatter.QAOverride || frontmatter.Status != "archived" {
			continue
		}
		report, err := ReadQAReport(filepath.Dir(prdPath))
		if errors.Is(err, ErrNoQAReport) {
			continue
		}
		if err != nil {
			t.Fatalf("read QA override Spec report for %q: %v", prdPath, err)
		}
		if report.Verdict == VerdictFail {
			t.Logf("found archived failed QA override Spec %s", filepath.Base(filepath.Dir(prdPath)))
			return
		}
	}
	t.Fatal("archived corpus has no qa_override Spec whose newest QA verdict is fail")
}

func prepareSpec0058Replay(t *testing.T, repositoryRoot string, reportFixture string, addDeclaration bool) (string, string) {
	t.Helper()
	fixtureRoot := filepath.Join(repositoryRoot, "internal", "spec", "testdata", "archive-replay-0058")
	sourceDir := filepath.Join(repositoryRoot, filepath.FromSlash(filepath.Dir(filepath.Dir(spec0058SourceReport))))
	specsRoot := filepath.Join(t.TempDir(), "docs", "specs")
	specDir := filepath.Join(specsRoot, spec0058ReplaySlug)
	archiveTestCopyTree(t, sourceDir, specDir)

	prdPath := filepath.Join(specDir, "_prd.md")
	prd := archiveTestReadFile(t, prdPath)
	if strings.Contains(prd, "## Unreachable Acceptance") {
		t.Fatal("original Spec 0058 PRD unexpectedly contains an unreachable declaration")
	}
	for _, original := range []string{"status: archived", "archived: 2026-08-01", "qa_override: true"} {
		if !strings.Contains(prd, original) {
			t.Fatalf("original Spec 0058 PRD does not contain %q", original)
		}
	}
	prd = strings.Replace(prd, "status: archived", "status: active", 1)
	prd = strings.Replace(prd, "archived: 2026-08-01\n", "", 1)
	prd = strings.Replace(prd, "qa_override: true\n", "", 1)
	if addDeclaration {
		prd += "\n" + archiveTestReadFile(t, filepath.Join(fixtureRoot, "unreachable-acceptance.md"))
	}
	if err := os.WriteFile(prdPath, []byte(prd), 0o644); err != nil {
		t.Fatalf("write replay PRD: %v", err)
	}

	reportSource := filepath.Join(fixtureRoot, reportFixture, "qa-report-2026-08-04.md")
	reportTarget := filepath.Join(specDir, "qa", filepath.Base(reportSource))
	if err := os.WriteFile(reportTarget, []byte(archiveTestReadFile(t, reportSource)), 0o644); err != nil {
		t.Fatalf("write replay QA Report: %v", err)
	}
	return specsRoot, specDir
}

func archiveTestRepositoryRoot(t *testing.T) string {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(testFile), "..", ".."))
}

func archiveTestReadFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	return string(content)
}

func archiveTestCopyTree(t *testing.T, source string, destination string) {
	t.Helper()
	err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, 0o644)
	})
	if err != nil {
		t.Fatalf("copy replay artifacts from %q to %q: %v", source, destination, err)
	}
}

func archiveTestAssertReplayRemainsActive(t *testing.T, specsRoot string) {
	t.Helper()
	activeDir := filepath.Join(specsRoot, spec0058ReplaySlug)
	archivedDir := filepath.Join(ArchiveSpecRoot(specsRoot, false), spec0058ReplaySlug)
	if _, err := os.Stat(activeDir); err != nil {
		t.Fatalf("refused replay active directory: %v", err)
	}
	if _, err := os.Stat(archivedDir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("refused replay archived directory exists: %v", err)
	}
}

func archiveTestPath(kind ArchiveKind, elements ...string) string {
	parts := append([]string{filepath.FromSlash(ArchiveDir(kind))}, elements...)
	return filepath.ToSlash(filepath.Join(parts...))
}

func archiveTestRepositoryPath(repoRoot string, kind ArchiveKind, elements ...string) string {
	return filepath.Join(repoRoot, filepath.FromSlash(archiveTestPath(kind, elements...)))
}
