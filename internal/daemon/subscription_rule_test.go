package daemon

import (
	"encoding/json"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/config"
	"roundfix/internal/runevent"
)

func TestSubscriptionOnlyRefusalActivatesTheFallback(t *testing.T) {
	t.Parallel()
	for _, phase := range []string{"prepare", "prompt"} {
		t.Run(phase, func(t *testing.T) {
			sink := &captureEventSink{}
			refusal := &agent.SelectionFailureError{Runtime: "opencode", Reason: config.SubscriptionOnlyReason + ": refused before Agent work"}
			runner := &fallbackBoundaryRunner{emitOutputByModel: map[string]bool{"fallback-model": true}}
			if phase == "prepare" {
				runner.prepareErrByModel = map[string]error{"preferred-model": refusal}
			} else {
				runner.runErrByModel = map[string]error{"preferred-model": refusal}
			}
			owner := fallbackBoundaryOwner(t, sink, runner)
			result, err := owner.Run(t.Context(), agent.ExecuteRequest{RunID: "run-fallback-boundary", Prompt: "do the work"})
			if err != nil || result.Output != "fallback output" {
				t.Fatalf("fallback result = %+v, error = %v", result, err)
			}
			want := "fallback-model"
			if phase == "prompt" {
				want = "preferred-model,fallback-model"
			}
			if got := strings.Join(runner.runModels(), ","); got != want {
				t.Fatalf("prompt models = %q, want %q", got, want)
			}
			events := eventsOfKind(sink, runevent.KindDaemonAgentSelectionFallback)
			if len(events) != 1 {
				t.Fatalf("fallback events = %+v, want one", events)
			}
			var payload struct {
				ReasonCode string `json:"reason_code"`
			}
			if err := json.Unmarshal(events[0].Payload, &payload); err != nil {
				t.Fatal(err)
			}
			if payload.ReasonCode != config.SubscriptionOnlyReason {
				t.Fatalf("reason_code = %q", payload.ReasonCode)
			}
			if got := countAgentStatusEvents(sink, agent.AgentWorkStartedStatus); got != 1 {
				t.Fatalf("work-started events = %d, want only fallback work", got)
			}
		})
	}
}
