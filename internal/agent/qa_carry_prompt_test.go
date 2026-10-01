package agent

import (
	"strings"
	"testing"
)

func TestQAContractKeepsCarriedRowsAndDeclaresInputs(t *testing.T) {
	prompt, err := BuildQAPrompt(sampleQAPromptRequest())
	if err != nil {
		t.Fatalf("BuildQAPrompt returned error: %v", err)
	}

	for _, line := range []string{
		"- A seeded row whose status is carried (established by: …; head: …) is not executed again: keep its identifier, status and provenance, and count it as passed. Execute every other row; the seeded Row carry-forward section names why each prior row re-runs.",
		"- Declare inputs on every row you execute: repository_path for repository content, commit_range for a row that reads Task commits, their authorization or changed-file scope. Never write evidence_snapshots; the Daemon records it.",
	} {
		if got := strings.Count(prompt, line); got != 1 {
			t.Fatalf("expected prompt to contain %q exactly once, got %d:\n%s", line, got, prompt)
		}
	}
}
