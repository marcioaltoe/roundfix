package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTheAuthoringSkillsPutARuleChallengeToTheUser(t *testing.T) {
	t.Parallel()
	skills := []string{"brainstorming", "write-idea", "write-prd", "write-techspec"}
	t.Run("heading_and_position", func(t *testing.T) {
		for _, skill := range skills {
			next := "## Declared model"
			if skill == "brainstorming" {
				next = "## Idea-stage judgment"
			}
			if skill == "write-prd" || skill == "write-techspec" {
				next = "## Open authoring questions"
			}
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "## Rule Challenge A rule guides the solution and never bounds the goal in silence.")
			for _, root := range []string{filepath.Join("..", ".agents", "skills"), "."} {
				data, err := os.ReadFile(filepath.Join(root, skill, "SKILL.md"))
				if err != nil {
					t.Fatal(err)
				}
				text := string(data)
				if strings.Count(text, "\n## Rule Challenge\n") != 1 {
					t.Errorf("%s/%s: heading must occur once", root, skill)
				}
				start := strings.Index(text, "## Rule Challenge\n")
				if start < 0 {
					t.Errorf("%s/%s: heading missing", root, skill)
					continue
				}
				rest := text[start+len("## Rule Challenge\n"):]
				pos := strings.Index(rest, "\n## ")
				if pos < 0 || !strings.HasPrefix(rest[pos+1:], next+"\n") {
					t.Errorf("%s/%s: wrong next heading", root, skill)
				}
			}
		}
	})
	t.Run("what_a_rule_is", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "A rule guides the solution and never bounds the goal in silence.")
		}
	})
	t.Run("the_trigger", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "A rule blocks the goal when the solution you would recommend, or a part the goal needs, cannot exist without contradicting it, or when it stops you from stating the problem.")
		}
	})
	t.Run("not_a_block", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "Extra work, taste and friction are not a block.")
		}
	})
	t.Run("the_question", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "Name the rule and where it lives, what it imposes, up to three things it prevents for this goal, and your reading of why the rule exists and whether it looks out of date or badly written.")
		}
	})
	t.Run("three_answers", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "Offer keep, change and retire, each with its consequence, the recommended one first")
		}
	})
	t.Run("binds_until_the_answer", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "The rule binds until the user answers: asking changes nothing, and a recommendation, silence, a dismissal or a timeout is never an answer.")
		}
	})
	t.Run("a_decided_change", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "Ask no Rule Challenge about a change the user already decided")
		}
	})
	t.Run("a_kept_line_needs_a_new_fact", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "unless the question names a fact the earlier one did not have. Record that fact as")
		}
	})
	t.Run("at_most_three", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "Ask at most three in one Authoring Phase or brainstorming session.")
		}
	})
	t.Run("never_answers_for_the_user", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "Never answer for the user. An Author Subagent returns the question, with its options and the one it recommends, to the session that started it.")
		}
	})
	t.Run("nobody_present", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "When no person is present, write no artifact for this phase")
		}
	})
	t.Run("change_or_retire", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "Change or retire: the PRD states the change of the rule as a requirement, and a Task of the same Spec changes the rule together with the work that depended on it.")
		}
	})
	t.Run("keep", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "Keep: the solution stays inside the rule, and the PRD requires a Task to add the Kept Line after the rule's text. The line names no Spec.")
		}
	})
	t.Run("a_baseline_rule", func(t *testing.T) {
		for _, skill := range skills {
			assertIdeaStageSkillPhrases(t, skill, "SKILL.md", "sends the user's decision as one ai-memory message to the Baseline's owner, obeys the rule until the Baseline changes, and records the message and the scope it left out.")
		}
	})
	t.Run("the_core_text_is_shared", func(t *testing.T) {
		var core string
		for _, skill := range skills {
			data, err := os.ReadFile(filepath.Join("..", ".agents", "skills", skill, "SKILL.md"))
			if err != nil {
				t.Fatal(err)
			}
			text := string(data)
			start := strings.Index(text, "## Rule Challenge\n")
			endText := "records the message and the scope it left out."
			end := strings.Index(text, endText)
			if start < 0 || end < start {
				t.Fatalf("%s: core missing", skill)
			}
			got := text[start : end+len(endText)]
			if core == "" {
				core = got
			} else if got != core {
				t.Errorf("%s: core differs", skill)
			}
		}
	})
	t.Run("stage/brainstorming", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "brainstorming", "SKILL.md", "In this conversation, when the direction the user wants meets a rule that blocks it, say so and ask the Rule Challenge instead of steering the user away from that direction. Hand the question and the user's answer to the skill the outcome routes to, which records them in the `## Rule challenges` section of its artifact; when no Spec follows a keep answer, add the Kept Line yourself.")
	})
	t.Run("stage/write-idea", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "write-idea", "SKILL.md", "In this stage, list the rules that touch the idea in the codebase sweep of step 4, and never drop an alternative of the Opportunity Scan only because a rule forbids it: keep it, name the rule beside it, and ask the Rule Challenge when it is the direction you would recommend. Record every Rule Challenge in the `## Rule challenges` section of `_idea.md`, in the fielded form the template shows; the section reads `None.` when no rule blocked the goal.")
	})
	t.Run("stage/write-prd", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "write-prd", "SKILL.md", "In this stage, read every ADR and rule the Project Constraints snapshot will list as an obligation and as a candidate: one that blocks a goal becomes a Rule Challenge before the snapshot is written. Copy the `## Rule challenges` records of `_idea.md` into the PRD unchanged and add the ones this stage asks, in the fielded form the template shows; the section reads `None.` when no rule blocked the goal. A change or retire answer becomes a Core Feature that states the change of the rule, and a keep answer a Core Feature that requires the Kept Line.")
	})
	t.Run("stage/write-techspec", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "write-techspec", "SKILL.md", "In this stage the question has a narrower trigger: ask it before the design takes what an accepted ADR rejected, or otherwise contradicts its decision, when neither the request, the PRD's `## Rule challenges` section nor another recorded answer of the user decides it. Record it in the PRD's `## Rule challenges` section with the phase `techspec`, adding the section to a minimal PRD, which otherwise omits it, and end that ADR's line under `## ADR Alternatives` with `Rule Challenge of YYYY-MM-DD: <answer>`. An amendment or a successor the user already decided needs no question; the line names where the user decided.")
	})
	t.Run("pointers/write-prd_inputs", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "write-prd", "SKILL.md", "Copy its `## Idea-stage judgment` section into the PRD unchanged, and its `## Rule challenges` records as `## Rule Challenge` below says.")
	})
	t.Run("pointers/brainstorming_explore", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "brainstorming", "SKILL.md", "Check project state first — files, docs, `CONTEXT.md`, the accepted ADRs under `docs/adr/`, the guides under `docs/agents/`, active Specs")
	})
	t.Run("pointers/write-techspec_read_first", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "write-techspec", "SKILL.md", "(respect prior decisions or challenge them explicitly — never contradict one silently, and ask the user first as `## Rule Challenge` below says)")
	})
	t.Run("pointers/write-techspec_rejected_alternatives", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "write-techspec", "SKILL.md", "say so and supersede or amend the ADR, after the Rule Challenge when the user has not decided it.")
	})
}
