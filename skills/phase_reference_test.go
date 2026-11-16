package skills

import (
	"io/fs"
	"strings"
	"testing"
)

func TestThePhaseReferenceDocumentsTheCommandAndTheRecord(t *testing.T) {
	t.Parallel()
	for sourceName, source := range roundfixLayoutSources() {
		t.Run(sourceName, func(t *testing.T) {
			data, err := fs.ReadFile(source, "roundfix/references/phase.md")
			if err != nil {
				t.Fatal(err)
			}
			text := strings.Join(strings.Fields(string(data)), " ")
			for _, tc := range []struct{ name, phrase string }{
				{"the_command", "roundfix phase start <phase> [--spec <slug>]"},
				{"the_seven_phases", "`follow-up`, `idea`, `prd`, `techspec`, `tasks`, `delivery` or `release`"},
				{"the_spec_option", "required for every phase but `follow-up` and `release`"},
				{"no_configuration", "It loads no Project Config and no User Config."},
				{"an_unknown_command_continues", "When the installed Roundfix does not know the command, say so and continue."},
				{"the_record", "the directory `phase-record` under Roundfix Home, one JSON object per line with the schema `roundfix/phase-record/v1`"},
				{"the_notice_refuses_nothing", "The Boundary Notice refuses nothing: standard output and the exit code stay as they are."},
			} {
				t.Run(tc.name, func(t *testing.T) {
					phrase := strings.Join(strings.Fields(tc.phrase), " ")
					if count := strings.Count(text, phrase); count != 1 {
						t.Fatalf("%q occurs %d times, want exactly one", phrase, count)
					}
					if tc.name == "the_command" && !strings.Contains(string(data), "```bash\n"+tc.phrase+"\n```") {
						t.Fatal("usage must be in a bash block")
					}
				})
			}
		})
	}
}
