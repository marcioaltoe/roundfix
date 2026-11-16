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
	}{
		{name: "records_its_own_phase", phrase: func(phase string) string {
			return "run `roundfix phase start " + phase + " --spec <slug>`"
		}},
		{name: "an_unknown_command_continues", phrase: func(string) string {
			return "When the installed Roundfix does not know the command, say so in the report and continue."
		}},
		{name: "the_author_subagent_returns", phrase: func(string) string {
			return "return your report and end no session"
		}},
		{name: "the_session_writes_the_handoff", phrase: func(string) string {
			return "write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed"
		}},
		{name: "the_user_is_told", phrase: func(string) string {
			return "tell the user what stays armed"
		}},
		{name: "no_further_phase", phrase: func(string) string {
			return "start no further phase in this session"
		}},
		{name: "the_next_session_follows_up", phrase: func(string) string {
			return "The next session starts with the Roundfix skill's `follow-up` Skill Mode, which claims the handoff."
		}},
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
						if got := strings.Count(section, phrase); got != 1 {
							t.Errorf("phrase occurrence = %d, want 1", got)
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
	}{
		{name: "follow_up_records_its_phase", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 1)
		}, phrase: "record the phase with `roundfix phase start follow-up`"},
		{name: "delivery_is_recorded", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 7)
		}, phrase: "run `roundfix phase start delivery --spec <slug>`"},
		{name: "release_is_recorded", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 7)
		}, phrase: "before a release run `roundfix phase start release`"},
		{name: "the_handoff_at_a_boundary", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed"},
		{name: "no_further_phase", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "start no further phase in this session"},
		{name: "the_next_session_claims_the_handoff", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "the next session starts with this mode, which claims the handoff"},
		{name: "a_wake_up_never_crosses", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "A scheduled wake-up returns to the same session with its whole context, so it keeps a phase going and never crosses a Phase Boundary."},
		{name: "a_fix_stays_in_the_session", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "A fix or a refactor is authored in the session"},
		{name: "the_turn_still_arms_or_says", place: func(text string) string {
			return roundfixNumberedItem(roundfixSection(text, "--follow-up"), 8)
		}, phrase: "arm a detached Run, the Delivery Queue or a scheduled wake-up, or say plainly that work stops and what remains"},
		{name: "unattended_asks_no_question", place: func(text string) string {
			return roundfixSection(text, "--unattended")
		}, phrase: "This mode asks no question about where a fix or a refactor is authored"},
		{name: "unattended_ends_at_the_boundary", place: func(text string) string {
			return roundfixSection(text, "--unattended")
		}, phrase: "With nobody present the session still ends at a Phase Boundary"},
		{name: "unattended_wake_up", place: func(text string) string {
			return roundfixSection(text, "--unattended")
		}, phrase: "A scheduled wake-up keeps a phase going and never crosses a Phase Boundary, because it returns to the same session."},
		{name: "status_records_no_phase", place: func(text string) string {
			return roundfixSection(text, "--status")
		}, phrase: "roundfix phase start"},
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
					want := 1
					if check.name == "status_records_no_phase" {
						want = 0
					}
					if got := strings.Count(place, phrase); got != want {
						t.Errorf("phrase occurrence = %d, want %d", got, want)
					}
				})
			}
		})
	}
}
