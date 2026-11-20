package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTheClarificationRecordGuideCarriesTheDelegatedDecisionForm(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, file, phrase string }{
		{"guide/heading", "clarification-record.md", "## Delegated Decisions"},
		{"guide/third_entry", "clarification-record.md", "It is a third entry of the record and never one of the two shapes above"},
		{"guide/head_line", "clarification-record.md", "- delegated: YYYY-MM-DD, <prd|techspec|tasks>"},
		{"guide/question_line", "clarification-record.md", "- question: <the question, in one line>"},
		{"guide/options_line", "clarification-record.md", "- options: <each option, separated by semicolons>"},
		{"guide/adopted_line", "clarification-record.md", "- adopted: <the adopted option> — <the author's reason>"},
		{"guide/under_line", "clarification-record.md", "- under: <Project Config | the user's answer of YYYY-MM-DD>"},
		{"guide/answer_line", "clarification-record.md", "- answer: none yet"},
		{"guide/the_grant", "clarification-record.md", "`Project Config` when the project's Authoring Question Preference is `follow`"},
		{"guide/kept_form", "clarification-record.md", "- answer: YYYY-MM-DD, kept — <the user's words>"},
		{"guide/changed_form", "clarification-record.md", "- answer: YYYY-MM-DD, changed — <the user's words>"},
		{"guide/only_before_delivery", "clarification-record.md", "Write the answer into the entry only while no Run and no delivery of the Spec has started."},
		{"guide/the_check", "clarification-record.md", "`SC-DECISIONS-UNRECORDED` also reports, as a gap, an entry that begins"},
		{"guide/never_counts", "clarification-record.md", "A Delegated Decision never satisfies the rule above: the PRD still needs a dated answer or the autonomous-route line."},
		{"template/third_entry", "prd-template.md", "A Delegated Decision, an option a session adopted with nobody present, is a third entry and never one of those two shapes:"},
		{"template/head_line", "prd-template.md", "- delegated: YYYY-MM-DD, <prd|techspec|tasks>"},
		{"template/adopted_line", "prd-template.md", "- adopted: <the adopted option> — <the author's reason>"},
		{"template/under_line", "prd-template.md", "- under: <Project Config | the user's answer of YYYY-MM-DD>"},
		{"template/answer_line", "prd-template.md", "- answer: none yet"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assertIdeaStageSkillPhrases(t, "write-prd", filepath.Join("references", tc.file), tc.phrase)
			for _, root := range []string{filepath.Join("..", ".agents", "skills"), "."} {
				path := filepath.Join(root, "write-prd", "references", tc.file)
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				text := string(data)
				if tc.name == "guide/heading" {
					var last string
					for _, line := range strings.Split(text, "\n") {
						if strings.HasPrefix(line, "## ") {
							last = line
						}
					}
					if last != "## Delegated Decisions" {
						t.Errorf("%s: last heading=%q", path, last)
					}
				}
				if tc.name == "template/third_entry" {
					const block = `     A Delegated Decision, an option a session adopted with nobody present,
     is a third entry and never one of those two shapes:
     - delegated: YYYY-MM-DD, <prd|techspec|tasks>
       - question: <the question>
       - options: <each option, separated by semicolons>
       - adopted: <the adopted option> — <the author's reason>
       - under: <Project Config | the user's answer of YYYY-MM-DD>
       - answer: none yet
     See references/clarification-record.md.`
					if !strings.Contains(text, block) {
						t.Errorf("%s: entry must immediately precede guide reference", path)
					}
				}
			}
		})
	}
}
