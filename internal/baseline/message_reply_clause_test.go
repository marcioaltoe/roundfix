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

func TestTheReceivingClauseSendsOneCompletionNoticeOnlyWhenAMessageAsksForIt(t *testing.T) {
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
		"A message may also ask to be told when the work it asks for is done; ask for that only when the sending project will act on the notice, and do not wait or watch for it either: the receiving project sends it once that work is delivered.",
		"send no acknowledgment, adoption notice or completion notice to a message that declared none, unless it asks to be told when its work is done.",
		"A message that asks to be told when its work is done receives one completion notice and no acknowledgment or adoption notice: keep that request in the record of its content,",
		"once the work recorded from it is delivered, when the release that carries that work is published or, when that work releases nothing, when it merges, send the notice with `memory_message_send`, naming what changed and the version or merge that carries it;",
		"when that work is closed without being delivered, the notice says so instead.",
		"The session that follows that delivery sends the notice; a Daemon-assigned turn never sends one.",
	} {
		if count := strings.Count(instructions, phrase); count != 1 {
			t.Errorf("agent instructions contain %q %d times, want once", phrase, count)
		}
	}
	if forbidden := "completion notice to a message that declared none. An empty inbox is rest"; strings.Contains(instructions, forbidden) {
		t.Errorf("agent instructions still contain %q", forbidden)
	}
}

func TestTheReceivingClauseAnswersQuestionsAndBoundsTheDeclaredPrecedence(t *testing.T) {
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
		"the receiving project owns what follows, the answer to a direct question the message asks included.",
		"so it takes precedence for one message per sender at a time: triage first, from each sending project, the oldest message that declares its work blocked until this project answers",
		"then the rest oldest first, a later blocking message from the same sender included.",
		"The precedence orders the triage and never replaces the maintainer's judgment of what this project handles first.",
		"Answer each direct question a message asks with `memory_message_send` as soon as the answer exists, whether or not the message declared blocked work",
	} {
		if count := strings.Count(instructions, phrase); count != 1 {
			t.Errorf("agent instructions contain %q %d times, want once", phrase, count)
		}
	}
	for _, forbidden := range []string{
		"Triage first each message whose sender declares its work blocked",
		"as soon as the answer exists when its message declared blocked work",
	} {
		if strings.Contains(instructions, forbidden) {
			t.Errorf("agent instructions still contain %q", forbidden)
		}
	}
}
