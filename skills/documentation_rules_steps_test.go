package skills

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestTheSkillsCarryTheDocumentationSteps(t *testing.T) {
	t.Parallel()
	normalize := func(s string) string { return strings.Join(strings.Fields(s), " ") }
	read := func(t *testing.T, root, path string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	roots := []string{"../.agents/skills", "."}
	cases := []struct {
		name, path string
		phrases    []string
	}{
		{"archive/five_signs", "archive-spec/SKILL.md", []string{"adopts a Finding;", "holds a QA Report whose verdict is not", "the Task Graph holds a Corrective Task;", "lists a declaration other than", "a Task of the Run failed its Verification and ran again, or the Delivery Queue item was retried."}},
		{"archive/no_sign_nothing", "archive-spec/SKILL.md", []string{"With no sign, write nothing."}},
		{"archive/commit_with_archive", "archive-spec/SKILL.md", []string{"commit it with the archive, so the same Pull Request shows it to a person, who may delete it."}},
		{"archive/report_the_failed_part", "archive-spec/SKILL.md", []string{"say in the report that the test was applied and which part failed."}},
		{"archive/never_fill", "archive-spec/SKILL.md", []string{"Never write a lesson to fill the folder."}},
		{"archive/queue_defers", "archive-spec/SKILL.md", []string{"When the Delivery Queue writes the archive, no Agent runs this step there"}},
		{"archive/skip_without_rule", "archive-spec/SKILL.md", []string{"A repository whose docs-layout guide has no Lesson Learned rule skips this section."}},
		{"corrections/skip_without_rule", "roundfix/references/corrections.md", []string{"A repository whose docs-layout guide has no Lesson Learned rule skips this section."}},
		{"corrections/cheapest_moment", "roundfix/references/corrections.md", []string{"A Correction is where a delivery's lesson is cheapest to state, because its declarer has just named the cause."}},
		{"corrections/in_the_correcting_commit", "roundfix/references/corrections.md", []string{"in the commit that makes the correction, on the Item Branch, so the delivery Pull Request shows it."}},
		{"corrections/after_merge", "roundfix/references/corrections.md", []string{"When the item has already merged, or the queue made the correction with no session, write the file in the session's next change to the repository."}},
		{"modes/step3_documentation", "roundfix/references/modes.md", []string{"member too: the lines the Standard Export List gains and loses"}},
		{"modes/step4_stops_exporting", "roundfix/references/modes.md", []string{"Ask that same question, naming the paths, when"}},
		{"modes/marking", "roundfix/references/modes.md", []string{"mark each ADR that lacks the mark by that rule's test, changing no other line."}},
		{"modes/review_in_pr", "roundfix/references/modes.md", []string{"say that the marks are a judgment for a person to review in the Pull Request, in a commit that holds only the marks and the list."}},
		{"modes/auto_marks_nothing", "roundfix/references/modes.md", []string{"mark nothing and report the count."}},
		{"baseline/documentation_part", "roundfix/references/baseline.md", []string{"part of its plan that the Plan Digest covers and the same confirmation applies."}},
		{"baseline/marks_none", "roundfix/references/baseline.md", []string{"It reports the count of ADRs without a Knowledge Mark and marks none."}},
		{"baseline/marker_refused", "roundfix/references/baseline.md", []string{"A missing or repeated marker is refused: nothing is written to that file, and the rest of the plan applies."}},
		{"baseline/no_history", "roundfix/references/baseline.md", []string{"does not skip this part."}},
		{"setup/two_checks", "roundfix/references/setup.md", []string{"reports two documentation checks as warnings that never change its exit."}},
		{"setup/repairs", "roundfix/references/setup.md", []string{"repairs the list, an Amendment List and a retired ADR, and the rest is an edit by hand."}},
		{"release/sentence", "roundfix/references/release.md", []string{"are planning evidence too, and the export list"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, root := range roots {
				text := normalize(read(t, root, tc.path))
				for _, phrase := range tc.phrases {
					if strings.Count(text, normalize(phrase)) != 1 {
						t.Errorf("%s: expected phrase once: %s", root, phrase)
					}
				}
			}
		})
	}
	t.Run("archive/before_steps", func(t *testing.T) {
		for _, root := range roots {
			text := read(t, root, "archive-spec/SKILL.md")
			a, b := strings.Index(text, "## Lesson Learned"), strings.Index(text, "## Steps")
			if a < 0 || b < 0 || a >= b {
				t.Errorf("%s: lesson must precede steps", root)
			}
		}
	})
	t.Run("modes/five_steps", func(t *testing.T) {
		for _, root := range roots {
			text := read(t, root, "roundfix/references/modes.md")
			start := strings.Index(text, "## --update")
			if start < 0 {
				t.Errorf("%s: missing update", root)
				continue
			}
			text = text[start+len("## --update"):]
			if end := strings.Index(text, "\n## "); end >= 0 {
				text = text[:end]
			}
			matches := regexp.MustCompile(`(?m)^([0-9]+)\. `).FindAllStringSubmatch(text, -1)
			if len(matches) != 5 {
				t.Errorf("%s: got %d steps", root, len(matches))
				continue
			}
			for i, m := range matches {
				if m[1] != string(rune('1'+i)) {
					t.Errorf("%s: unexpected step %s", root, m[1])
				}
			}
		}
	})
}
