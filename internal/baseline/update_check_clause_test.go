// Suite: Cross-project continuity
// Invariant: a resolved Baseline profile tells each session to run the Update Check before its first task.
// Boundary IN: Baseline decision-plan postimages.
// Boundary OUT: applying a Plan and refreshing a live repository.

package baseline

import (
	"strings"
	"testing"
)

func TestTheSessionStartClauseRunsTheUpdateCheck(t *testing.T) {
	t.Parallel()
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), standardTypeScriptDecisions("make verify"))
	instructions := strings.Join(strings.Fields(string(planPostimage(t, plan, "docs/agents/agent-instructions.md").Content)), " ")
	for _, phrase := range []string{
		"run `roundfix update check` once and relay its report",
		"With `update.mode` set to `auto`, run that mode before the first task",
		"In a Daemon-assigned turn, skip the check.",
	} {
		if got := strings.Count(instructions, phrase); got != 1 {
			t.Errorf("agent instructions contain %q %d times, want exactly once", phrase, got)
		}
	}
}
