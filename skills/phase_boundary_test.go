package skills

import (
	"io/fs"
	"strings"
	"testing"
)

func TestTheAuthoringSkillsStateThePhaseBoundary(t *testing.T) {
	t.Parallel()

	skills := []struct {
		name  string
		phase string
	}{
		{name: "write-idea", phase: "idea"},
		{name: "write-prd", phase: "prd"},
		{name: "write-techspec", phase: "techspec"},
		{name: "write-tasks", phase: "tasks"},
	}
	checks := []struct {
		name   string
		phrase func(string) string
		want   int
	}{
		{name: "records_its_own_phase", phrase: func(phase string) string {
			return "run `roundfix phase start " + phase + " --spec <slug>`"
		}, want: 1},
		{name: "an_unknown_command_continues", phrase: func(string) string {
			return "When the installed Roundfix does not know the command, say so in the report and continue."
		}, want: 1},
		{name: "the_author_subagent_returns", phrase: func(string) string {
			return "return your report and end no session"
		}, want: 1},
		{name: "the_session_writes_the_handoff", phrase: func(string) string {
			return "write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed"
		}, want: 1},
		{name: "the_user_is_told", phrase: func(string) string {
			return "tell the user what stays armed"
		}, want: 1},
		{name: "no_phase_ends_the_session", phrase: func(string) string {
			return "When the phase ends, the session goes on: no phase of a Spec ends the session, whether a person is present or not."
		}, want: 1},
		{name: "the_session_goes_on", phrase: func(string) string {
			return "go on in this session with the next phase or the next queued work, handing an Authoring Phase to its Author Subagent as its skill's `Declared model` section says, and write no Session Handoff for this phase"
		}, want: 1},
		{name: "the_subagent_session_goes_on", phrase: func(string) string {
			return "the session that started you goes on with the next phase"
		}, want: 1},
		{name: "the_session_ends_when_asked_or_stopped", phrase: func(string) string {
			return "End the session only when the user asks or the work stops because nothing remains or what remains waits on the user"
		}, want: 1},
		{name: "the_next_session_follows_up", phrase: func(string) string {
			return "The next session starts with the Roundfix skill's `follow-up` Skill Mode, which claims the handoff."
		}, want: 1},
		{name: "no_old_boundary", phrase: func(string) string {
			return "When the phase ends, a Phase Boundary follows"
		}, want: 0},
		{name: "no_stop_after_the_phase", phrase: func(string) string {
			return "start no further phase"
		}, want: 0},
	}

	for copy, source := range roundfixLayoutSources() {
		copy, source := copy, source
		for _, authoringSkill := range skills {
			authoringSkill := authoringSkill
			t.Run(copy+"/"+authoringSkill.name, func(t *testing.T) {
				t.Parallel()
				data, err := fs.ReadFile(source, authoringSkill.name+"/SKILL.md")
				if err != nil {
					t.Fatal(err)
				}
				text := string(data)
				if got := strings.Count(text, "## Phase boundary"); got != 1 {
					t.Fatalf("Phase boundary headings = %d, want 1", got)
				}
				declaredAt := strings.Index(text, "## Declared model")
				boundaryAt := strings.Index(text, "## Phase boundary")
				processAt := strings.Index(text, "## Process")
				if declaredAt < 0 || boundaryAt < 0 || processAt < 0 || declaredAt >= boundaryAt || boundaryAt >= processAt {
					t.Fatalf("headings out of order: declared=%d boundary=%d process=%d", declaredAt, boundaryAt, processAt)
				}

				section := roundfixWhitespace(roundfixSection(text, "Phase boundary"))
				for _, check := range checks {
					check := check
					t.Run(check.name, func(t *testing.T) {
						t.Parallel()
						phrase := roundfixWhitespace(check.phrase(authoringSkill.phase))
						if got := strings.Count(section, phrase); got != check.want {
							t.Errorf("phrase occurrence = %d, want %d", got, check.want)
						}
					})
				}
			})
		}
	}
}

func TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary(t *testing.T) {
	t.Parallel()

	checks := []struct {
		name   string
		place  func(string) string
		phrase string
		want   int
	}{
		{name: "follow_up_records_its_phase", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 1)
		}, phrase: "record the phase with `roundfix phase start follow-up`", want: 1},
		{name: "delivery_is_recorded", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 7)
		}, phrase: "run `roundfix phase start delivery --spec <slug>`", want: 1},
		{name: "release_is_recorded", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 7)
		}, phrase: "before a release run `roundfix phase start release`", want: 1},
		{name: "the_handoff_at_a_boundary", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed", want: 1},
		{name: "the_next_session_claims_the_handoff", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "the next session starts with this mode, which claims the handoff", want: 1},
		{name: "no_phase_ends_the_session", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "No phase of a Spec ends the session, whether a person is present or not", want: 1},
		{name: "the_session_follows_the_queue", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "after an Authoring Phase, a delivery or a release, go on with the next queued work", want: 1},
		{name: "the_session_ends_when_asked_or_stopped", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "End the session only when the user asks or the work stops because nothing remains or what remains waits on the user", want: 1},
		{name: "a_wake_up_returns_to_the_session", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "A scheduled wake-up returns to the same session with its whole context.", want: 1},
		{name: "a_fix_stays_in_the_session", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "A fix or a refactor is authored in the session", want: 1},
		{name: "the_turn_still_arms_or_says", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "arm a detached Run, the Delivery Queue or a scheduled wake-up, or say plainly that work stops and what remains", want: 1},
		{name: "unattended_asks_no_question", place: func(text string) string {
			return roundfixSection(text, "--unattended")
		}, phrase: "This mode asks no question about where a fix or a refactor is authored", want: 1},
		{name: "unattended_no_phase_ends_the_session", place: func(text string) string {
			return roundfixSection(text, "--unattended")
		}, phrase: "With nobody present no phase of a Spec ends the session either", want: 1},
		{name: "no_boundary_after_each_phase", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "after each Authoring Phase", want: 0},
		{name: "no_stop_after_the_phase", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "start no further phase", want: 0},
		{name: "unattended_no_boundary", place: func(text string) string {
			return roundfixSection(text, "--unattended")
		}, phrase: "ends at a Phase Boundary", want: 0},
		{name: "status_records_no_phase", place: func(text string) string {
			return roundfixSection(text, "--status")
		}, phrase: "roundfix phase start", want: 0},
	}

	for copy, source := range roundfixLayoutSources() {
		copy, source := copy, source
		t.Run(copy, func(t *testing.T) {
			t.Parallel()
			data, err := fs.ReadFile(source, "roundfix/references/modes.md")
			if err != nil {
				t.Fatal(err)
			}
			text := string(data)
			for _, check := range checks {
				check := check
				t.Run(check.name, func(t *testing.T) {
					t.Parallel()
					place := roundfixWhitespace(check.place(text))
					phrase := roundfixWhitespace(check.phrase)
					if got := strings.Count(place, phrase); got != check.want {
						t.Errorf("phrase occurrence = %d, want %d", got, check.want)
					}
				})
			}
		})
	}
}
