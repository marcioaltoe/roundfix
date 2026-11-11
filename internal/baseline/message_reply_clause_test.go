package baseline

import (
	"strings"
	"testing"
)

func TestTheMessageClausesAskForAReplyOnlyForBlockedWork(t *testing.T) {
	t.Parallel()
	decisions := standardTypeScriptDecisions("make verify")
	for i := range decisions {
		if decisions[i].ID == "secondbrain.enabled" {
			decisions[i].Value = true
		}
	}
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), decisions)
	instructions := strings.Join(strings.Fields(string(planPostimage(t, plan, "docs/agents/agent-instructions.md").Content)), " ")
	for _, phrase := range []string{
		"A message expects no reply unless the sending project's own work depends on the answer.",
		"say so in the message's subject and first lines, in the message's own words and language",
		"after the work already running, which is never interrupted for it",
		"send no acknowledgment, adoption notice or completion notice to a message that declared none",
	} {
		if count := strings.Count(instructions, phrase); count != 1 {
			t.Errorf("agent instructions contain %q %d times, want once", phrase, count)
		}
	}
	for _, forbidden := range []string{"triage them oldest first", "when the message asks for a reply"} {
		if strings.Contains(instructions, forbidden) {
			t.Errorf("agent instructions still contain %q", forbidden)
		}
	}
}
