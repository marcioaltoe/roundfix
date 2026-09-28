// Suite: pre-PR review provider scope.
// Invariant: roundfix review ignores legacy CodeRabbit request coherence unless CodeRabbit is selected.
// Boundary IN: review command configuration, candidate review, and persisted outcome.
// Boundary OUT: legacy fetch, watch, resolve, and their request publication.
package cli

import (
	"fmt"
	"path/filepath"
	"testing"

	"roundfix/internal/agent"
)

func TestReviewUnderCodexIgnoresCodeRabbitConfiguration(t *testing.T) {
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	writeIncoherentCodeRabbitReviewFixture(t, fixture)

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("Codex review exit=%d outcome=%q, want exit=%d outcome=%q; stderr=%q", code, record.Outcome, exitOK, reviewOutcomeReviewed, stderr)
	}
}

func TestReviewUnderNoneIgnoresCodeRabbitConfiguration(t *testing.T) {
	runner := &reviewCommandRunner{}
	fixture := newReviewCommandFixture(t, "none", runner)
	writeIncoherentCodeRabbitReviewFixture(t, fixture)

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeOmitted {
		t.Fatalf("disabled review exit=%d outcome=%q, want exit=%d outcome=%q; stderr=%q", code, record.Outcome, exitOK, reviewOutcomeOmitted, stderr)
	}
	if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 || runner.endCalls != 0 {
		t.Fatalf(
			"none policy used Agent runner: probes=%d prepares=%d prompts=%d ends=%d",
			runner.probeCalls,
			runner.prepareCalls,
			runner.preparedCalls,
			runner.endCalls,
		)
	}
}

func writeIncoherentCodeRabbitReviewFixture(t *testing.T, fixture reviewCommandFixture) {
	t.Helper()
	mustWrite(t, filepath.Join(fixture.repository, ".coderabbit.yaml"), `reviews:
  auto_review:
    enabled: false
`)
	mustWrite(t, filepath.Join(fixture.repository, ".roundfixrc.yml"), fmt.Sprintf(`defaults:
  artifact_dir: %q
pre_pr_review:
  provider: %s
review_source:
  request_review: false
`, fixture.artifactDir, fixture.provider))
}
