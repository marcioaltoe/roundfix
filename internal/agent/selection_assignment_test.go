package agent

import "testing"

func TestSelectionRuntimeDeferredStateMatchesBeforeEffortApplication(t *testing.T) {
	t.Parallel()

	assignment, state := runtimeDeferredSelectionState(t)
	if !selectionStateMatches(assignment, state) {
		t.Fatalf("effective state does not prove deferred assignment %#v", assignment)
	}
}

func TestSelectionRuntimeDeferredStateRejectsDifferentCurrentModel(t *testing.T) {
	t.Parallel()

	assignment, state := runtimeDeferredSelectionState(t)
	state.CurrentModel = "openrouter/deepseek/deepseek-v4-pro"
	if selectionStateMatches(assignment, state) {
		t.Fatalf("a different current model proved deferred assignment %#v", assignment)
	}
}

func runtimeDeferredSelectionState(t *testing.T) (SelectionAssignment, SelectionCapabilities) {
	t.Helper()

	state := parseMeasuredEffortCapabilities(t, measuredGrok45EffortCapabilities)
	assignment, err := PlanSelectionAssignment(
		RuntimeSpec{ID: "opencode", Model: "openrouter/x-ai/grok-4.5", ReasoningEffort: "high"},
		state,
	)
	if err != nil {
		t.Fatalf("plan deferred OpenCode Agent Selection: %v", err)
	}
	if state.ReasoningOption == nil || state.ReasoningOption.CurrentValue == assignment.ReasoningValue {
		t.Fatalf("fixture reasoning value must remain unapplied: state = %#v, assignment = %#v", state.ReasoningOption, assignment)
	}
	return assignment, state
}
