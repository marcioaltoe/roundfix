package cli

// Suite: archive-root identity through the public review command.
// Boundary IN: disposable Git candidate and fake review provider.
// Boundary OUT: network and daemon settlement.

import (
	"strings"
	"testing"

	"roundfix/internal/agent"
)

func TestReviewIgnoresAProseReadmeInTheArchiveRoot(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings."}}}}
	f := newReviewCommandFixture(t, "codex", runner)
	commitScopeFile(t, f, "docs/history/specs/README.md", "# Archived Specs\n\nA prose index.\n")
	code, record, diag := runLineageReview(t, f)
	if code != exitOK || runner.preparedCalls != 1 || len(record.ArchivedSpecs) != 0 || len(record.SkippedSpecs) != 0 || !strings.Contains(runner.request.Prompt, "A prose index.") {
		t.Fatalf("exit=%d provider calls=%d record=%+v stderr=%s", code, runner.preparedCalls, record, diag)
	}
}

func TestReviewStillRefusesAMalformedArchiveRecord(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings."}}}}
	f := newReviewCommandFixture(t, "codex", runner)
	commitScopeFile(t, f, "docs/history/specs/0256-malformed.md", "---\nschema: roundfix/archive-record/v1\n---\n\n# Missing record fields\n")
	code, _, diag := runLineageReview(t, f)
	if code == exitOK || runner.preparedCalls != 0 || !strings.Contains(diag, "read review archive record:") {
		t.Fatalf("exit=%d provider calls=%d stderr=%s", code, runner.preparedCalls, diag)
	}
}
