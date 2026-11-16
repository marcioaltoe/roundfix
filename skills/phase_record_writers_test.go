package skills

import (
	"io/fs"
	"strings"
	"testing"
)

func TestTheRoundfixSkillNamesThePromptEntriesAndTheirRetention(t *testing.T) {
	t.Parallel()
	for sourceName, source := range roundfixLayoutSources() {
		t.Run(sourceName, func(t *testing.T) {
			for _, tc := range []struct{ name, path, phrase string }{
				{"the_writers", "phase", "The Daemon appends a prompt entry when a prompt of an Agent Session ends, and `roundfix review` after each prompt of a pre-PR review"},
				{"a_prompt_that_never_returns", "phase", "A prompt that never returns leaves no entry."},
				{"a_failed_write_warns", "phase", "one warning is printed and the Run or the review ends as it would have"},
				{"retention_removes_day_files", "storage", "Run Retention also removes each Phase Record day file dated before its cutoff"},
				{"the_report_names_the_count", "storage", "`Phase Record day files removed: <n>`, or `Phase Record day files removable: <n>` in a dry run"},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := fs.ReadFile(source, "roundfix/references/"+tc.path+".md")
					if err != nil {
						t.Fatal(err)
					}
					text, phrase := strings.Join(strings.Fields(string(data)), " "), strings.Join(strings.Fields(tc.phrase), " ")
					if count := strings.Count(text, phrase); count != 1 {
						t.Fatalf("%q occurs %d times, want exactly once", phrase, count)
					}
				})
			}
		})
	}
}
