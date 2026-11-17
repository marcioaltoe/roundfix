package baseline

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestTheOutliveTheTurnClauseStatesWhenASessionEnds(t *testing.T) {
	t.Parallel()

	phrases := []struct {
		name   string
		phrase string
	}{
		{name: "must_arm_something_that_outlives_its_turn", phrase: "must arm something that outlives its own turn before that turn ends: a detached Run, a scheduled wake-up, or a daemon"},
		{name: "never_end_a_turn_stating_work_will_continue", phrase: "Never end a turn stating that work will continue"},
		{name: "when_no_mechanism_is_armed_say_what_remains", phrase: "When no such mechanism is armed, say plainly that work is stopping and what remains."},
		{name: "session_handoff_names_boundary_state", phrase: "write the session handoff with the Spec, the phase that ended, the next phase, each open decision and what stays armed"},
		{name: "user_is_told_what_stays_armed", phrase: "tell the user what stays armed"},
		{name: "no_phase_ends_the_session", phrase: "No phase of a Spec ends the session, whether a person is present or not"},
		{name: "the_session_goes_on", phrase: "after an authoring phase, a delivery or a release, the session goes on with the next queued work"},
		{name: "authoring_skill_says_where_it_runs", phrase: "each authoring skill says when a phase runs in a subagent with its own context"},
		{name: "the_session_ends_when_asked_or_stopped", phrase: "The session ends only when the user asks or the work stops because nothing remains or what remains waits on the user"},
	}
	oldBoundaryPhrases := []struct {
		name   string
		phrase string
	}{
		{name: "after_each_authoring_phase", phrase: "after each authoring phase"},
		{name: "start_no_further_phase_in_that_session", phrase: "start no further phase in that session"},
		{name: "the_way_across_a_boundary", phrase: "the way across a boundary"},
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
	t.Run("no_old_boundary", func(t *testing.T) {
		t.Parallel()
		for _, want := range oldBoundaryPhrases {
			want := want
			if got := strings.Count(guidance, want.phrase); got != 0 {
				t.Fatalf("clause contains old phrase %d times, want zero", got)
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
					if got := strings.Count(text, literal); got != 0 {
						t.Fatalf("guide contains old literal %d times, want zero", got)
					}
				})
			}
		}
	})
}
