package judge

// Suite: measured Spec planning.
// Invariant: every state contains only the measured cuts from accepted Sources.
// Boundary IN: temporary Specs, ADRs, and PlanSpec.
// Boundary OUT: network, credentials, logging, and CLI.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const englishContext = "The author reads the decision and the gate is the rule for the work.\n\n"
const measuredClaim = "ADR-0123 keeps the verification gate in the repository for every author (Spec 0123)."
const cleanClaim = "The cited decision keeps the verification gate in the repository for every author ."

func fixturePlan(t *testing.T, q Questions, prd, tech, adr string, stage Stage) Plan {
	t.Helper()
	root := t.TempDir()
	spec := filepath.Join(root, "spec")
	writeFixture(t, spec, "_prd.md", prd)
	if tech != "" {
		writeFixture(t, spec, "_techspec.md", tech)
	}
	writeFixture(t, root, "docs/adr/0123-decision.md", adr)
	plan, err := PlanSpec(q, root, spec, stage)
	if err != nil {
		t.Fatal(err)
	}
	return plan
}

func TestCitationClaimsFollowTheMeasuredExtraction(t *testing.T) {
	q := loadQuestions(t)
	adr := "---\nstatus: accepted\n---\n\n" + englishContext + "<The gate & the author>"
	for _, tc := range []struct {
		name, text string
		want       int
		line       int
	}{
		{"rule 1 fence", "```text\n" + measuredClaim + "\n```\n", 0, 0},
		{"rule 1 table", "  | " + measuredClaim + " |\n", 0, 0},
		{"rule 1 heading", "  # " + measuredClaim + "\n", 0, 0},
		{"rule 1 blank ends paragraph", measuredClaim + "\n\nADR-0123 mentions the unrelated rule.", 1, 3},
		{"rule 1 bullet ends paragraph", "ADR-0123 mentions the rule\n- " + measuredClaim, 1, 4},
		{"rule 1 numbered item", "ADR-0123 mentions the rule\n12. " + measuredClaim, 1, 4},
		{"rule 1 heading ends paragraph", "ADR-0123 mentions the rule\n## Claim\n" + measuredClaim, 1, 5},
		{"rule 1 table ends paragraph", "ADR-0123 mentions the rule\n| other |\n" + measuredClaim, 1, 5},
		{"rule 1 fence ends paragraph", "ADR-0123 mentions the rule\n```\nignored\n```\n" + measuredClaim, 1, 7},
		{"rule 2 joins lines", "ADR-0123 keeps the verification gate\nin the repository for every author (Spec 0123).", 1, 3},
		{"rule 2 capital sentence", "ADR-0123 mentions it. " + measuredClaim, 1, 3},
		{"rule 2 backtick sentence", "ADR-0123 mentions it! `" + measuredClaim + "`", 1, 3},
		{"rule 2 parenthesis sentence", "ADR-0123 mentions it? (" + measuredClaim + ")", 1, 3},
		{"rule 2 bracket sentence", "ADR-0123 mentions it. [" + measuredClaim + "]", 1, 3},
		{"rule 2 lowercase is not boundary", "ADR-0123 mentions it. authors say " + measuredClaim, 0, 0},
		{"rule 3 two tokens", measuredClaim + " with ADR-0124", 0, 0},
		{"rule 3 repeated token", measuredClaim + " with ADR-0123", 0, 0},
		{"rule 3 attribution verb", "ADR-0123 mentions the verification gate in the repository for every author.", 0, 0},
		{"rule 4 transformation", measuredClaim, 1, 3},
		{"rule 6 token line", "A paragraph begins on this line and continues\n" + measuredClaim, 1, 4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			plan := fixturePlan(t, q, englishContext+tc.text, "", adr, "prd")
			if len(plan.Pending) != tc.want {
				t.Fatalf("pending=%+v skips=%+v", plan.Pending, plan.Skipped)
			}
			if tc.want > 0 && plan.Pending[0].Line != tc.line {
				t.Fatalf("line=%d want=%d", plan.Pending[0].Line, tc.line)
			}
			if tc.name == "rule 2 joins lines" && plan.Pending[0].Text != cleanClaim {
				t.Fatalf("joined claim=%q want=%q", plan.Pending[0].Text, cleanClaim)
			}
		})
	}
	t.Run("rules 4 and 5 exact state and Unicode ADR cut", func(t *testing.T) {
		body := strings.TrimSpace(englishContext + strings.Repeat("the 雪 & <gate> ", q.Citation.DecisionRecordMaxChars))
		plan := fixturePlan(t, q, "---\nignored: ADR-9999 keeps text\n---\n"+englishContext+"- "+measuredClaim, "", "---\nstatus: accepted\n---\n\n"+body, "prd")
		if len(plan.Pending) != 1 {
			t.Fatalf("plan=%+v", plan)
		}
		var state map[string]string
		if err := json.Unmarshal(plan.Pending[0].state, &state); err != nil {
			t.Fatal(err)
		}
		if len(state) != 2 || state["claim"] != cleanClaim {
			t.Fatalf("claim=%q fields=%d", state["claim"], len(state))
		}
		if state["decision_record"] != string([]rune(body)[:q.Citation.DecisionRecordMaxChars]) {
			t.Fatalf("decision record differs: got %d characters", len([]rune(state["decision_record"])))
		}
		if plan.Pending[0].Line != 6 || strings.Contains(string(plan.Pending[0].state), `\u003c`) || strings.Contains(string(plan.Pending[0].state), `\u0026`) {
			t.Fatalf("state escaping or line: %+v %s", plan.Pending[0], plan.Pending[0].state)
		}
	})
	t.Run("rule 4 length boundaries are code points", func(t *testing.T) {
		for _, length := range []int{q.Citation.ClaimMinChars - 1, q.Citation.ClaimMinChars, q.Citation.ClaimMaxChars, q.Citation.ClaimMaxChars + 1} {
			prefix := "The cited decision keeps "
			claim := "ADR-0123 keeps " + strings.Repeat("é", length-len([]rune(prefix)))
			plan := fixturePlan(t, q, englishContext+claim, "", adr, "prd")
			want := 1
			if length < q.Citation.ClaimMinChars || length > q.Citation.ClaimMaxChars {
				want = 0
			}
			if len(plan.Pending) != want {
				t.Fatalf("length %d: %+v", length, plan)
			}
			if want == 0 && (len(plan.Skipped) != 1 || plan.Skipped[0].Reason != "claim length outside the measured 60-500 characters") {
				t.Fatalf("length skip=%+v", plan.Skipped)
			}
		}
	})
	t.Run("rule 3 proposed ADR and non-English ADR", func(t *testing.T) {
		for _, record := range []string{"---\nstatus: proposed\n---\n" + englishContext, "A decisão é uma regra para os autores e não está na sua documentação."} {
			plan := fixturePlan(t, q, englishContext+measuredClaim, "", record, "prd")
			if len(plan.Pending) != 0 || len(plan.Skipped) != 1 {
				t.Fatalf("plan=%+v", plan)
			}
			if !strings.HasPrefix(record, "---") && plan.Skipped[0].Reason != "decision record is not English" {
				t.Fatalf("skip=%+v", plan.Skipped)
			}
		}
	})
	t.Run("one pending per distinct state and stage", func(t *testing.T) {
		prd := englishContext + measuredClaim + "\n\n" + measuredClaim
		tech := englishContext + measuredClaim
		for _, stage := range []Stage{"", "prd", "techspec"} {
			plan := fixturePlan(t, q, prd, tech, adr, stage)
			if len(plan.Pending) != 1 {
				t.Fatalf("stage %q: %+v", stage, plan)
			}
			want := "_prd.md"
			if stage == "techspec" {
				want = "_techspec.md"
			}
			if filepath.Base(plan.Pending[0].Artifact) != want {
				t.Fatalf("stage %q: wrong artifact", stage)
			}
		}
	})
}

func TestGoalPairsFollowTheMeasuredSelection(t *testing.T) {
	q := loadQuestions(t)
	prd := englishContext + "## Goals\n- The author reads the\n  evidence for the work.\n\n2. The author keeps the rule & <gate>.\n## Other\n- This is not a goal."
	shortTitle := "Reader gate"
	longTitle := "**Reader gate with `evidence`**"
	longBody := strings.Repeat("The author reads the 雪 & <gate> for the work. ", q.Goal.SectionMaxChars)
	tech := englishContext + "## Coverage Map\n- Goals 1-2 → Reader gate with\n  evidence.\n- Goal 3 → Reader gate with evidence.\n- Goal 1 → Testing Approach.\n- Goal 1 → Short mechanism.\n- Goal 1 → Tiny.\n- Goal 1 → Missing mechanism.\n" +
		"## " + shortTitle + "\n" + strings.Repeat(englishContext, 5) + "\n### " + longTitle + "\n" + longBody + "\n## Testing Approach\n" + strings.Repeat(englishContext, 5) + "\n## Short mechanism\nThe author reads.\n## Tiny\n" + strings.Repeat(englishContext, 5)
	plan := fixturePlan(t, q, prd, tech, englishContext, "techspec")
	if len(plan.Pending) != 2 || len(plan.Skipped) != 5 {
		t.Fatalf("plan=%+v", plan)
	}
	for i, pending := range plan.Pending {
		var state map[string]string
		if err := json.Unmarshal(pending.state, &state); err != nil {
			t.Fatal(err)
		}
		wantGoal := "The author reads the evidence for the work."
		if i == 1 {
			wantGoal = "The author keeps the rule & <gate>."
		}
		if len(state) != 3 || state["goal"] != wantGoal || state["section_title"] != longTitle || state["section"] != string([]rune(strings.TrimSpace(longBody))[:q.Goal.SectionMaxChars]) {
			t.Fatalf("state=%v", state)
		}
		if pending.Line != 4 || strings.Contains(string(pending.state), `\u003c`) {
			t.Fatalf("pending=%+v", pending)
		}
	}
	if plan.Skipped[0].Reason != "Goal 3 is not in the PRD" {
		t.Fatalf("skip=%+v", plan.Skipped[0])
	}
	for _, skip := range plan.Skipped[1:] {
		if skip.Reason != "no named section qualifies" {
			t.Fatalf("skip=%+v", skip)
		}
	}
	t.Run("rule 3 section boundary and child body", func(t *testing.T) {
		for _, tc := range []struct {
			name, tech string
			want       int
		}{
			{"child belongs to parent", "## Reader mechanism\n### Child\n" + strings.Repeat(englishContext, 5), 1},
			{"next peer does not belong", "## Reader mechanism\nShort\n## Other mechanism\n" + strings.Repeat(englishContext, 5), 0},
			{"next higher does not belong", "### Reader mechanism\nShort\n## Other mechanism\n" + strings.Repeat(englishContext, 5), 0},
			{"top heading ends body", "### Reader mechanism\nShort\n# Other\n" + strings.Repeat(englishContext, 5), 0},
		} {
			t.Run(tc.name, func(t *testing.T) {
				plan := fixturePlan(t, q, prd, englishContext+"## Coverage Map\n- Goal 1 → Reader mechanism.\n"+tc.tech, englishContext, "techspec")
				if len(plan.Pending) != tc.want {
					t.Fatalf("plan=%+v", plan)
				}
			})
		}
	})
	t.Run("rule 2 table map and repeated states", func(t *testing.T) {
		plan := fixturePlan(t, q, prd, englishContext+"## Coverage Map\n| PRD Goals 1–2 | Reader mechanism |\n- Goals 1-2 -> Reader mechanism\n## Reader mechanism\n"+strings.Repeat(englishContext, 5), englishContext, "")
		if len(plan.Pending) != 2 {
			t.Fatalf("plan=%+v", plan)
		}
	})
	t.Run("rules 3 and 4 measured boundaries", func(t *testing.T) {
		for _, tc := range []struct {
			name                    string
			titleLen, bodyLen, want int
		}{
			{"short body", q.Goal.SectionTitleMinChars, q.Goal.SectionMinChars - 1, 0},
			{"minimum body", q.Goal.SectionTitleMinChars, q.Goal.SectionMinChars, 1},
			{"short title", q.Goal.SectionTitleMinChars - 1, q.Goal.SectionMinChars, 0},
		} {
			t.Run(tc.name, func(t *testing.T) {
				title := strings.Repeat("é", tc.titleLen)
				body := strings.Repeat("é", tc.bodyLen)
				plan := fixturePlan(t, q, prd, strings.Repeat(englishContext, 10)+"## Coverage Map\n- Goal 1 → "+title+"\n## "+title+"\n"+body, englishContext, "techspec")
				if len(plan.Pending) != tc.want {
					t.Fatalf("plan=%+v", plan)
				}
			})
		}
	})
	t.Run("prd stage omits goals", func(t *testing.T) {
		plan := fixturePlan(t, q, prd, tech, englishContext, "prd")
		if len(plan.Pending) != 0 {
			t.Fatalf("plan=%+v", plan)
		}
	})
}

func TestPlanSkipsANonEnglishArtifact(t *testing.T) {
	q := loadQuestions(t)
	portuguese := "A decisão é uma regra para os autores e não está na sua documentação.\n## Goals\n- Uma regra para os autores.\n"
	tech := englishContext + "## Coverage Map\n- Goal 1 → Reader mechanism.\n## Reader mechanism\n" + strings.Repeat(englishContext, 5)
	for _, stage := range []Stage{"", "prd", "techspec"} {
		plan := fixturePlan(t, q, portuguese, tech, englishContext, stage)
		if len(plan.Pending) != 0 || len(plan.SkippedArtifacts) != 1 || plan.SkippedArtifacts[0].Reason != "not English" {
			t.Fatalf("plan=%+v", plan)
		}
	}
	plan := fixturePlan(t, q, englishContext, portuguese, englishContext, "techspec")
	if len(plan.Pending) != 0 || len(plan.SkippedArtifacts) != 1 || plan.SkippedArtifacts[0].Artifact != "_techspec.md" || plan.SkippedArtifacts[0].Reason != "not English" {
		t.Fatalf("plan=%+v", plan)
	}
}

func TestPlanSkipsASymbolicLinkArtifact(t *testing.T) {
	q := loadQuestions(t)
	for _, name := range []string{"_prd.md", "_techspec.md"} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			spec := filepath.Join(root, "spec")
			writeFixture(t, spec, "_prd.md", englishContext)
			if name == "_prd.md" {
				if err := os.Remove(filepath.Join(spec, name)); err != nil {
					t.Fatal(err)
				}
			}
			target := writeFixture(t, root, "target.md", englishContext+measuredClaim)
			if err := os.Symlink(target, filepath.Join(spec, name)); err != nil {
				t.Fatal(err)
			}
			plan, err := PlanSpec(q, root, spec, "techspec")
			if (err != nil) != (name == "_prd.md") {
				t.Fatalf("error=%v", err)
			}
			if len(plan.Pending) != 0 || len(plan.SkippedArtifacts) == 0 || plan.SkippedArtifacts[0].Reason != "not a regular file in the Spec directory" {
				t.Fatalf("plan=%+v", plan)
			}
		})
	}
}
