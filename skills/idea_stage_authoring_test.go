package skills

import "testing"

func TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage(t *testing.T) {
	t.Parallel()

	t.Run("three criteria", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "brainstorming", "SKILL.md",
			"the right shape of the solution is still open",
			"the request rests on an assumption about its users, the business or the market that it gives no evidence for",
			"deciding to build it depends on a market or competitive question",
		)
		assertIdeaStageSkillPhrases(t, "write-idea", "SKILL.md",
			"market/competitive questions",
			"the right _shape_ of the solution is still open",
			"unexamined business assumptions",
		)
		assertIdeaStageSkillPhrases(t, "write-prd", "SKILL.md",
			"the right shape of the solution is still open",
			"the request rests on an assumption about its users, the business or the market that it gives no evidence for",
			"deciding to build it depends on a market or competitive question",
		)
	})

	t.Run("three verdict rules", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "brainstorming", "SKILL.md",
			"Any criterion yes: the Idea stage is indicated.",
			"Every criterion no: go to `write-prd`.",
			"Anything else, or any doubt: ask the user whether the Idea stage runs, showing the answers.",
		)
		assertIdeaStageSkillPhrases(t, "write-prd", "SKILL.md",
			"Indicated, any criterion yes: stop this skill and start `write-idea`, showing the verdict and the answers.",
			"Undecided, no criterion yes and not every one no, or any doubt: ask the user whether the Idea stage runs, showing the same answers.",
			"Not indicated, every criterion no: continue.",
		)
		assertIdeaStageSkillPhrases(t, "write-idea", "SKILL.md",
			"reads indicated",
			"user chose it for an undecided one",
			"go straight to `grilling`/`grill-with-docs` → `write-prd`",
		)
	})

	t.Run("request file with what reaches the judge", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "brainstorming", "SKILL.md",
			"only the project's name, the request's first line and its first 4,000 characters reach the judge, after the secret scan",
		)
		assertIdeaStageSkillPhrases(t, "write-prd", "SKILL.md",
			"only the project's name, the request's first line and its first 4,000 characters reach the judge, after the secret scan",
		)
	})

	t.Run("never present your own answers as the judge's", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "brainstorming", "SKILL.md",
			"Never present your own answers as the judge's",
		)
		assertIdeaStageSkillPhrases(t, "write-prd", "SKILL.md",
			"Never present your own answers as the judge's",
		)
	})

	t.Run("dismissal option", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "brainstorming", "SKILL.md",
			"the user may dismiss the stage at its first question",
		)
		assertIdeaStageSkillPhrases(t, "write-idea", "SKILL.md",
			"Give that first question one last option that dismisses the stage",
		)
		assertIdeaStageSkillPhrases(t, "write-prd", "SKILL.md",
			"The user may dismiss the stage at its first question; then this skill resumes.",
		)
	})

	t.Run("waiting rule", func(t *testing.T) {
		assertIdeaStageSkillPhrases(t, "write-idea", "SKILL.md",
			"the feature stays a Backlog Entry that records the judgment and waits, with `waiting_on` and `since` in its frontmatter",
		)
		assertIdeaStageSkillPhrases(t, "write-prd", "SKILL.md",
			"write no PRD: keep the feature as a Backlog Entry that records the judgment and waits, with `waiting_on` and `since` in its frontmatter",
		)
	})

	t.Run("every key of the record in both templates", func(t *testing.T) {
		keys := []string{
			"- item:",
			"- judged-by:",
			"- shape_open:",
			"- unverified_assumption:",
			"- market_question:",
			"- verdict:",
			"- stage:",
			"- shape:",
		}
		assertIdeaStageSkillPhrases(t, "write-idea", "references/idea-template.md", keys...)
		assertIdeaStageSkillPhrases(t, "write-prd", "references/prd-template.md", keys...)
	})
}

func TestTheFollowUpModeRunsTheIdeaStageJudgment(t *testing.T) {
	t.Parallel()

	assertIdeaStageSkillPhrases(t, "roundfix", "references/modes.md",
		"For each planned feature that has no Idea, run the Idea-stage judgment as `docs/agents/spec-routing.md` says, with `roundfix spec judge --stage idea --item <reference>` and, for a message entry, the same messages file, or with `--request <path>` for a feature the user asked for in the session, and show its verdict and answers in the plan: an indicated feature starts at `write-idea`, an undecided one waits for the user's choice, and a not indicated one goes to its PRD.",
		"The questions of an indicated Idea stage, and the choice for an undecided one, are asked here, while the user is present and before any Spec is authored; the user may dismiss an indicated stage at its first question.",
		"A feature whose Idea stage is indicated or undecided is authored only after the user answered; without an answer it stays a Backlog Entry that records the judgment and waits, its Waiting Declaration naming the Idea stage.",
		"The Idea stage is the user's: this mode asks its questions and never answers them, and it never decides an undecided Idea-stage judgment. A feature whose judgment is indicated or undecided and gets no answer becomes or stays a Backlog Entry that records the judgment, with a Waiting Declaration whose `waiting_on` names the user's answers to the Idea stage, and the report names it.",
	)
}
