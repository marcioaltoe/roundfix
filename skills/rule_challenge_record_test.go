package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTheTemplatesCarryTheRuleChallengeRecord(t *testing.T) {
	t.Parallel()
	for _, c := range []struct{ name, skill, path, next string }{{"idea", "write-idea", "references/idea-template.md", "## Decisions"}, {"prd", "write-prd", "references/prd-template.md", "## Glossary"}} {
		t.Run(c.name+"/heading_and_position", func(t *testing.T) {
			assertIdeaStageSkillPhrases(t, c.skill, c.path, "## Rule challenges")
			for _, root := range []string{filepath.Join("..", ".agents", "skills"), "."} {
				b, err := os.ReadFile(filepath.Join(root, c.skill, c.path))
				if err != nil {
					t.Fatal(err)
				}
				s := string(b)
				pos := strings.Index(s, "## Rule challenges")
				if pos < 0 {
					t.Error("missing section")
					continue
				}
				rest := s[pos+len("## Rule challenges"):]
				i := strings.Index(rest, "\n## ")
				if i < 0 || !strings.HasPrefix(rest[i+1:], c.next+"\n") {
					t.Errorf("wrong next heading in %s", root)
				}
			}
		})
		t.Run(c.name+"/fields", func(t *testing.T) {
			assertIdeaStageSkillPhrases(t, c.skill, c.path, "- rule: ADR-NNNN", "- imposes: <what the rule requires>", "- prevents: <what it prevents for this goal>", "- asked: YYYY-MM-DD, "+c.name, "- answer: change", "- effect: <the requirement", "new fact: only when the rule already carried a Kept Line")
		})
		t.Run(c.name+"/none", func(t *testing.T) {
			assertIdeaStageSkillPhrases(t, c.skill, c.path, "the single entry `None.` when no rule blocked the goal")
		})
		t.Run(c.name+"/phase", func(t *testing.T) { assertIdeaStageSkillPhrases(t, c.skill, c.path, "asked: YYYY-MM-DD, "+c.name) })
	}
}

func TestTheSkillModesNeverAnswerARuleChallenge(t *testing.T) {
	t.Parallel()
	t.Run("unattended", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "roundfix", "references/modes.md", "A Rule Challenge is the user's too: this mode asks it and never answers it.", "A Spec whose Rule Challenge gets no answer writes no artifact for that phase")
	})
	t.Run("follow_up_step_five", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "roundfix", "references/modes.md", "its Waiting Declaration naming the Idea stage. A Spec whose author returns a Rule Challenge goes on only after the user answered it; without an answer it waits as `--unattended` says.")
		for _, root := range []string{filepath.Join("..", ".agents", "skills"), "."} {
			b, err := os.ReadFile(filepath.Join(root, "roundfix/references/modes.md"))
			if err != nil {
				t.Fatal(err)
			}
			step := roundfixNumberedItem(roundfixSection(string(b), "--follow-up"), 5)
			if !strings.Contains(step, "A Spec whose author returns a Rule Challenge") {
				t.Errorf("%s: missing from step 5", root)
			}
		}
	})
}
