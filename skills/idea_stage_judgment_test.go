package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// assertIdeaStageSkillPhrases checks both the canonical and shipped copy. Task
// authors can reuse it for related authoring references and templates.
func assertIdeaStageSkillPhrases(t *testing.T, skill, relative string, phrases ...string) {
	t.Helper()
	for _, root := range []string{filepath.Join("..", ".agents", "skills"), "."} {
		path := filepath.Join(root, skill, relative)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		text := strings.Join(strings.Fields(string(data)), " ")
		for _, phrase := range phrases {
			phrase = strings.Join(strings.Fields(phrase), " ")
			if count := strings.Count(text, phrase); count != 1 {
				t.Errorf("%s: phrase %q occurs %d times; want once", path, phrase, count)
			}
		}
	}
}

func TestTheSpecReferenceDocumentsTheIdeaStage(t *testing.T) {
	t.Parallel()
	assertIdeaStageSkillPhrases(t, "roundfix", "references/spec.md",
		"roundfix spec judge --stage idea (--item <reference> [--messages <path>] | --request <path>) [--format <text|json>]",
		"The judge receives the project's name and the one item's kind, type, origin, title and the first 4,000 characters of its text, after the secret scan, in whatever language it is written.",
		"Of a request file it receives only the first line, as the title, and the first 4,000 characters; nothing else is sent.",
		"`shape_open` (the right shape of the solution is still open)",
		"`unverified_assumption` (it rests on an assumption about its users, the business or the market that it gives no evidence for)",
		"`market_question` (deciding to build it depends on a market or competitive question)",
		"An answer is yes at 0.70 or more, no at 0.30 or less and unsure between.",
		"The verdict is `indicated` when any criterion is yes, `not indicated` when every criterion is no and `undecided` otherwise.",
		"An item typed `fix`, `perf`, `refactor` or `finding` and a Waiting Item are not judged and nothing is sent",
		"roundfix/idea-stage-judgment/v1",
	)
}
