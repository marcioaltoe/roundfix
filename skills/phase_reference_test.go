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
				{"the_record_is_never_read_back", "The command never reads the record back: whatever the session recorded before, it writes `boundary` as `false` and prints nothing on standard error."},
				{"no_boundary_notice", "Boundary Notice"},
			} {
				t.Run(tc.name, func(t *testing.T) {
					phrase := strings.Join(strings.Fields(tc.phrase), " ")
					want := 1
					if tc.name == "no_boundary_notice" {
						want = 0
					}
					if count := strings.Count(text, phrase); count != want {
						t.Fatalf("%q occurs %d times, want %d", phrase, count, want)
					}
					if tc.name == "the_command" && !strings.Contains(string(data), "```bash\n"+tc.phrase+"\n```") {
						t.Fatal("usage must be in a bash block")
					}
				})
			}
		})
	}
}
