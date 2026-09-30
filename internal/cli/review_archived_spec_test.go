// Suite: archived Spec-aware pre-PR review.
// Invariant: a review preserves every archived Spec named by its candidate and identifies the required correction only for findings.
// Boundary IN: candidate Spec discovery, review record persistence and review stderr.
// Boundary OUT: real ACP adapters, Delivery Queue transitions and corrective Spec authoring.
package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
)

func TestReviewRecordsTheSpecsACandidateArchives(t *testing.T) {
	const (
		carriedSlug = "0177-carried-archive"
		skippedSlug = "0176-skipped-archive"
	)
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "Findings:\n- internal/example.go:1: correction needed", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	writeArchivedReviewSpec(t, fixture.repository, carriedSlug, true)
	writeArchivedReviewSpec(t, fixture.repository, skippedSlug, false)
	gittest.Run(t, fixture.repository, "add", "docs/history/specs")
	gittest.Run(t, fixture.repository, "commit", "-m", "archive Specs in candidate")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))

	code, fresh, stderr := fixture.run(t)

	if code != exitRunFailed || fresh.Outcome != reviewOutcomeFindings {
		t.Fatalf("fresh archived-Spec review exit=%d record=%+v stderr=%q", code, fresh, stderr)
	}
	wantArchived := []string{skippedSlug, carriedSlug}
	if !reflect.DeepEqual(fresh.ArchivedSpecs, wantArchived) {
		t.Fatalf("fresh archived Specs = %v, want %v", fresh.ArchivedSpecs, wantArchived)
	}
	if !reflect.DeepEqual(fresh.Specs, []string{carriedSlug}) || !reflect.DeepEqual(fresh.SkippedSpecs, []string{skippedSlug}) {
		t.Fatalf("fresh Spec context = carried:%v skipped:%v", fresh.Specs, fresh.SkippedSpecs)
	}

	code, reused, stderr := fixture.run(t)

	if code != exitRunFailed || !reused.Reused {
		t.Fatalf("reused archived-Spec review exit=%d record=%+v stderr=%q", code, reused, stderr)
	}
	if !reflect.DeepEqual(reused.ArchivedSpecs, wantArchived) {
		t.Fatalf("reused archived Specs = %v, want %v", reused.ArchivedSpecs, wantArchived)
	}
	if runner.preparedCalls != 1 {
		t.Fatalf("reviewer prompt calls = %d, want 1 across fresh and reused reviews", runner.preparedCalls)
	}
}

func TestReviewRecordsNoArchivedSpecForAnActiveSpec(t *testing.T) {
	const slug = "0178-active-spec"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	specDir := filepath.Join(fixture.repository, "docs", "specs", slug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create active Spec directory: %v", err)
	}
	mustWrite(t, filepath.Join(specDir, "_prd.md"), "# Active Spec\n\n## Decisions\n\n- Keep active context visible.\n")
	mustWrite(t, filepath.Join(specDir, "_techspec.md"), "# Active technical design\n")
	gittest.Run(t, fixture.repository, "add", filepath.ToSlash(filepath.Join("docs", "specs", slug)))
	gittest.Run(t, fixture.repository, "commit", "-m", "change active Spec")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("active-Spec review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if record.ArchivedSpecs == nil || len(record.ArchivedSpecs) != 0 {
		t.Fatalf("active-Spec archived Specs = %#v, want []", record.ArchivedSpecs)
	}
	recordBytes, err := os.ReadFile(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
	if err != nil {
		t.Fatalf("read persisted review record: %v", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(recordBytes, &fields); err != nil {
		t.Fatalf("decode persisted review record fields: %v", err)
	}
	if string(fields["archivedSpecs"]) != "[]" {
		t.Fatalf("review archivedSpecs JSON = %s, want []", fields["archivedSpecs"])
	}
}

func TestReviewNamesTheCorrectiveSpecForFindingsAfterArchive(t *testing.T) {
	const slug = "0176-corrective-source"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "Findings:\n- internal/example.go:1: correction needed", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	writeArchivedReviewSpec(t, fixture.repository, slug, true)
	gittest.Run(t, fixture.repository, "add", filepath.ToSlash(filepath.Join("docs", "history", "specs", slug)))
	gittest.Run(t, fixture.repository, "commit", "-m", "archive reviewed Spec")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))

	code, record, stderr := fixture.run(t)

	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings {
		t.Fatalf("review findings after archive exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	for _, want := range []string{slug, "an archived Spec is never corrected in place", "author a corrective Spec with its own authorization and QA gate"} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("review stderr = %q, want %q", stderr, want)
		}
	}
}

func TestReviewAfterArchiveWithoutFindingsNamesNoCorrectiveSpec(t *testing.T) {
	const slug = "0176-no-correction-needed"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	writeArchivedReviewSpec(t, fixture.repository, slug, true)
	gittest.Run(t, fixture.repository, "add", filepath.ToSlash(filepath.Join("docs", "history", "specs", slug)))
	gittest.Run(t, fixture.repository, "commit", "-m", "archive reviewed Spec")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("review without findings after archive exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if strings.Contains(stderr, slug) || strings.Contains(stderr, "corrective Spec") {
		t.Fatalf("no-findings review stderr names a correction: %q", stderr)
	}
}

func TestDeliveryReviewResultCarriesArchivedSpecs(t *testing.T) {
	record := newReviewRecord(
		"/repo",
		"base",
		"head",
		roundconfig.PrePRReview{Provider: "codex", Source: "project"},
		reviewOutcomeFindings,
	)
	record.Findings = "internal/example.go:1: correction needed"
	record.ArchivedSpecs = []string{"0176-one", "0177-two"}

	result, err := deliveryReviewResult(record, "head")

	if err != nil {
		t.Fatalf("map delivery review result: %v", err)
	}
	if result.Outcome != delivery.ReviewOutcomeFindings || !reflect.DeepEqual(result.ArchivedSpecs, record.ArchivedSpecs) {
		t.Fatalf("delivery review result = %+v, want findings with archived Specs %v", result, record.ArchivedSpecs)
	}
}

func writeArchivedReviewSpec(t *testing.T, repository string, slug string, complete bool) {
	t.Helper()
	specDir := filepath.Join(repository, "docs", "history", "specs", slug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create archived Spec directory: %v", err)
	}
	mustWrite(t, filepath.Join(specDir, "_prd.md"), "# Archived Spec\n\n## Decisions\n\n- Keep archived context visible.\n")
	if complete {
		mustWrite(t, filepath.Join(specDir, "_techspec.md"), "# Archived technical design\n")
	}
}
