// Suite: provider prompt admission through disposable Git repositories.
// Boundary IN: committed candidates and configured review providers.
// Boundary OUT: stderr, persisted review records and prepared provider prompts.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
)

func TestReviewCommandRefusesAClaudePromptOverHalfItsWindow(t *testing.T) {
	for _, round := range []int{1, 2} {
		t.Run(fmt.Sprint(round), func(t *testing.T) {
			runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: sessionReviewFinding}}}}
			fixture := newReviewCommandFixture(t, "claude", runner)
			writeReviewCommandProfileConfig(t, fixture.repository, "claude", fixture.artifactDir, "claude", "claude")
			if round == 2 {
				code, _, stderr := fixture.run(t)
				if code != exitRunFailed {
					t.Fatalf("seed review exit=%d stderr=%q", code, stderr)
				}
			}
			probes, prepares, prompts, ends := runner.probeCalls, runner.prepareCalls, runner.preparedCalls, runner.endCalls
			commitScopeFile(t, fixture, "large.txt", strings.Repeat("x", claudeReviewContextTokens))
			code, record, stderr := fixture.run(t)
			want := fmt.Sprintf("review prompt too large: %d estimated tokens after omitting %d path(s) exceeds the claude review budget of %d tokens, half of its %d-token context window", record.EstimatedPromptTokens, len(record.OmittedPaths), claudeReviewPromptBudget, claudeReviewContextTokens)
			assertBlockedReviewCommand(t, code, record, stderr, want)
			if record.Reason != want || stderr != "roundfix: review blocked: "+want+"\n" {
				t.Fatalf("reason=%q stderr=%q", record.Reason, stderr)
			}
			if record.EstimatedPromptTokens <= claudeReviewPromptBudget || record.AnswerPath != "" {
				t.Fatalf("over-budget admission record=%+v", record)
			}
			if runner.probeCalls != probes || runner.prepareCalls != prepares || runner.preparedCalls != prompts || runner.endCalls != ends {
				t.Fatalf("refusal reached runner: probes=%d prepares=%d prompts=%d ends=%d", runner.probeCalls-probes, runner.prepareCalls-prepares, runner.preparedCalls-prompts, runner.endCalls-ends)
			}
			fields := reviewPromptBoundRecordFields(t, fixture)
			if _, ok := fields["estimatedPromptTokens"]; !ok {
				t.Fatal("persisted record omitted Claude estimate")
			}
		})
	}
}

func TestReviewCommandSendsAClaudePromptBeyondTheCodexByteBound(t *testing.T) {
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings."}}}}
	fixture := newReviewCommandFixture(t, "claude", runner)
	writeReviewCommandProfileConfig(t, fixture.repository, "claude", fixture.artifactDir, "claude", "claude")
	commitScopeFile(t, fixture, "large.txt", strings.Repeat("x", reviewDiffBound))
	commitScopeFile(t, fixture, "docs/specs/example/_prd.md", "## Decisions\n\nKeep the candidate bounded.\n")
	commitScopeFile(t, fixture, "docs/specs/example/_techspec.md", "Technical context carried in the prompt.\n")
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if record.DiffBytes <= reviewDiffBound || record.EstimatedPromptTokens > claudeReviewPromptBudget || record.EstimatedPromptTokens != estimateReviewPromptTokens(runner.request.Prompt) {
		t.Fatalf("diffBytes=%d estimate=%d prompt estimate=%d", record.DiffBytes, record.EstimatedPromptTokens, estimateReviewPromptTokens(runner.request.Prompt))
	}
	if !strings.Contains(runner.request.Prompt, "--- BEGIN SPEC CONTEXT: example ---") {
		t.Fatal("estimate must include Spec context")
	}
	if runner.prepareCalls != 1 || runner.preparedCalls != 1 {
		t.Fatalf("prepares=%d prompts=%d, want one each", runner.prepareCalls, runner.preparedCalls)
	}
}

func TestReviewCommandKeepsTheCodexDiffByteBound(t *testing.T) {
	runner := &reviewCommandRunner{}
	fixture := newReviewCommandFixture(t, "codex", runner)
	commitScopeFile(t, fixture, "large.txt", strings.Repeat("x", reviewDiffBound))
	commitScopeFile(t, fixture, "docs/specs/example/_prd.md", "## Decisions\n\nKeep the candidate bounded.\n")
	commitScopeFile(t, fixture, "docs/specs/example/_techspec.md", "Technical context carried in the prompt.\n")
	code, record, stderr := fixture.run(t)
	want := fmt.Sprintf("review diff too large: %d bytes after omitting %d path(s) exceeds the review bound of %d bytes", record.DiffBytes, len(record.OmittedPaths), reviewDiffBound)
	assertBlockedReviewCommand(t, code, record, stderr, want)
	if record.Reason != want || record.EstimatedPromptTokens != 0 {
		t.Fatalf("Codex record=%+v", record)
	}
	if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
		t.Fatal("Codex byte refusal reached runner")
	}
	if _, ok := reviewPromptBoundRecordFields(t, fixture)["estimatedPromptTokens"]; ok {
		t.Fatal("Codex record contains a prompt estimate")
	}
}

func TestEstimateReviewPromptTokens(t *testing.T) {
	for _, test := range []struct {
		prompt string
		want   int
	}{
		{"", 0},
		{"a", 1},
		{"ab", 1},
		{"abc", 2},
		{"界", 2},
	} {
		if got := estimateReviewPromptTokens(test.prompt); got != test.want {
			t.Errorf("estimate(%q)=%d, want %d", test.prompt, got, test.want)
		}
	}
}

func reviewPromptBoundRecordFields(t *testing.T, fixture reviewCommandFixture) map[string]json.RawMessage {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(content, &fields); err != nil {
		t.Fatal(err)
	}
	return fields
}
