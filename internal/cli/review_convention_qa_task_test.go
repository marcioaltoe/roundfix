// Suite: QA settlement convention in the immutable candidate prompt.
// Boundary IN: committed QA Task; OUT: locally built prompt, no provider.
package cli

import (
	"strings"
	"testing"

	"roundfix/internal/preflight"
)

func TestDeliveryConventionC2ExemptsTheQATask(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{}
	fixture := newReviewCommandFixture(t, "codex", runner)
	head := commitScopeFile(t, fixture, "docs/specs/example/task_05.md", "---\ntask: task_05\nspec: example\nstatus: completed\ntype: qa\ncomplexity: high\n---\n\n# Task 05: Final QA\n\n## Requirements\n\nAudit the committed candidate and write its QA Report.\n\n## Verification\n\nThe QA Report records the evidence.\n")
	diff, omitted, err := reviewScopedDiff(t.Context(), fixture.repository, fixture.baseCommit, head, []string{"docs/specs", "docs/history/specs"}, preflight.ExecGitRunner{})
	if err != nil {
		t.Fatal(err)
	}
	prompt := buildReviewPrompt(fixture.baseCommit, head, diff, omitted)
	const c2 = "The Daemon writes a Task file's status after the Task's Verification passes, its `## Carry-forward provenance` section for non-QA Tasks only, and `## Recorded paths` only when a Task changed a path it did not declare; a QA Task settles on its committed QA Report and carries none of these sections. The Agent writes a non-QA Task's `## Result` during its turn, while the Daemon holds the Task `in_progress`, so a Result that reports the status as `in_progress`, unchanged or Daemon-owned states what was true when it was written and does not contradict a `completed` status."
	if !strings.Contains(prompt, "Delivery Conventions (roundfix/delivery-conventions/v4)") || !strings.Contains(prompt, "C2. "+c2+"\n") || !strings.Contains(prompt, "+type: qa") {
		t.Fatalf("missing QA convention or candidate: %s", prompt)
	}
	if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
		t.Fatal("provider reached")
	}
}
