package baseline

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestTheOutliveTheTurnClauseStatesThePhaseBoundary(t *testing.T) {
	t.Parallel()

	phrases := []struct {
		name   string
		phrase string
	}{
		{name: "must_arm_something_that_outlives_its_turn", phrase: "must arm something that outlives its own turn before that turn ends: a detached Run, a scheduled wake-up, or a daemon"},
		{name: "never_end_a_turn_stating_work_will_continue", phrase: "Never end a turn stating that work will continue"},
		{name: "when_no_mechanism_is_armed_say_what_remains", phrase: "When no such mechanism is armed, say plainly that work is stopping and what remains."},
		{name: "phase_boundary_is_recorded", phrase: "At a phase boundary of a Spec, after each authoring phase and after a Spec is delivered and released"},
		{name: "session_handoff_names_boundary_state", phrase: "write the session handoff with the Spec, the phase that ended, the next phase, each open decision and what stays armed"},
		{name: "user_is_told_what_stays_armed", phrase: "tell the user what stays armed"},
		{name: "no_further_phase_starts", phrase: "start no further phase in that session"},
		{name: "wake_up_stays_in_the_session", phrase: "A scheduled wake-up returns to the same session with its whole context, so it keeps a phase going and is never the way across a boundary."},
	}

	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Module("autonomous-work")
	if !ok {
		t.Fatal("missing autonomous-work module")
	}
	var module document
	if err := json.Unmarshal(asset.Data, &module); err != nil {
		t.Fatal(err)
	}
	var clause document
	for _, rule := range objectsOrEmpty(module["rules"]) {
		for _, candidate := range objectsOrEmpty(rule["clauses"]) {
			if candidate["id"] == "clause.autonomous.loop-07-outlive-the-turn" {
				clause = candidate
			}
		}
	}
	if len(clause) == 0 {
		t.Fatal("missing outlive-the-turn clause")
	}
	if clause["enforcement"] != "mandatory" {
		t.Fatalf("clause enforcement = %v, want mandatory", clause["enforcement"])
	}
	if _, exists := clause["replaces"]; exists {
		t.Fatal("clause has replaces")
	}
	guidance, ok := clause["guidance"].(string)
	if !ok {
		t.Fatal("clause guidance is not a string")
	}

	golden, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md")
	if !ok {
		t.Fatal("missing Standard TypeScript autonomous-work golden")
	}
	local, err := os.ReadFile("../../docs/agents/autonomous-work.md")
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range phrases {
		want := want
		t.Run(want.name, func(t *testing.T) {
			t.Parallel()
			if got := strings.Count(guidance, want.phrase); got != 1 {
				t.Fatalf("clause contains phrase %d times, want exactly once", got)
			}
			literal := strings.Join(strings.Fields(want.phrase), " ")
			for name, data := range map[string][]byte{
				"Standard TypeScript golden": golden.Data,
				"repository guide":           local,
			} {
				name, data := name, data
				t.Run(name, func(t *testing.T) {
					t.Parallel()
					text := strings.Join(strings.Fields(string(data)), " ")
					if got := strings.Count(text, literal); got != 1 {
						t.Fatalf("guide contains literal %d times, want exactly once", got)
					}
				})
			}
		})
	}
}
