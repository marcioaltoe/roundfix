// Suite: reviewer selection in public review records.
// Invariant: reviewed rounds name their selection; ceiling records name none.
// Boundary IN: temporary Git repositories, review CLI output, and disposition ledger.
// Boundary OUT: real reviewers, user artifacts, and the network.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestARoundRecordNamesItsSelectionEvenWhenItIsThePreferred(t *testing.T) {
	t.Parallel()
	runner := newLineageRunner("No findings.", "No findings.")
	fixture := newReviewCommandFixture(t, "codex", runner)
	for round := 1; round <= 2; round++ {
		if round == 2 {
			commitLineageFile(t, fixture, "delta.txt", "round two\n")
		}
		var stdout, stderr bytes.Buffer
		code := runCLI(t, []string{"review", "--base", fixture.baseCommit}, &stdout, &stderr)
		if code != exitOK || stderr.Len() != 0 {
			t.Fatalf("round %d: exit=%d stdout=%s stderr=%s", round, code, &stdout, &stderr)
		}
		var record struct {
			Lineage map[string]json.RawMessage `json:"lineage"`
		}
		if err := json.Unmarshal(stdout.Bytes(), &record); err != nil {
			t.Fatal(err)
		}
		if got, ok := record.Lineage["selection"]; !ok || string(got) != "0" {
			t.Errorf("round %d: selection=%s present=%v, want 0 present=true", round, got, ok)
		}
		if got := string(record.Lineage["round"]); got != fmt.Sprint(round) {
			t.Fatalf("round=%s, want %d", got, round)
		}
	}
}

func TestACeilingRecordNamesNoSelection(t *testing.T) {
	t.Parallel()
	fixture, runner, second := roundTwoLineageFixture(t)
	fix := commitLineageFile(t, fixture, "fix.txt", "final correction\n")
	for _, outcome := range []reviewOutcome{reviewOutcomeBlocked, reviewOutcomeCeilingClosed} {
		t.Run(string(outcome), func(t *testing.T) {
			if outcome == reviewOutcomeCeilingClosed {
				disposeLineageFinding(t, "--fixed-by", fix)
			}
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"review", "--base", fixture.baseCommit}, &stdout, &stderr)
			wantCode := exitOK
			wantStderr := ""
			if outcome == reviewOutcomeBlocked {
				wantCode = exitPreflight
				wantStderr = fmt.Sprintf("roundfix: review blocked: review round ceiling reached: round 2 reviewed %s; dispose finding F1 with roundfix review dispose\n", second.HeadCommit)
			}
			if code != wantCode || stderr.String() != wantStderr {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
			}
			var record struct {
				Outcome reviewOutcome   `json:"outcome"`
				Lineage json.RawMessage `json:"lineage"`
			}
			if err := json.Unmarshal(stdout.Bytes(), &record); err != nil {
				t.Fatal(err)
			}
			if record.Outcome != outcome {
				t.Fatalf("outcome=%s, want %s", record.Outcome, outcome)
			}
			wantLineage := fmt.Sprintf(`{"round":2,"reviewedHead":%q,"sessionOpen":false,"continued":false}`, second.HeadCommit)
			if string(record.Lineage) != wantLineage {
				t.Errorf("lineage=%s, want %s", record.Lineage, wantLineage)
			}
			if runner.preparedCalls != 2 {
				t.Errorf("reviewer calls=%d, want only the two reviewed rounds", runner.preparedCalls)
			}
			t.Logf("stdout: %s\nstderr: %s\nexit: %d", &stdout, &stderr, code)
		})
	}
}
