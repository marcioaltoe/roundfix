package baseline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func challengeModule(t *testing.T) document {
	t.Helper()
	data, err := os.ReadFile("assets/modules/spec-workflow.json")
	if err != nil {
		t.Fatal(err)
	}
	var module document
	if err := json.Unmarshal(data, &module); err != nil {
		t.Fatal(err)
	}
	return module
}
func challengeRule(t *testing.T, module document, id string) document {
	t.Helper()
	var found document
	count := 0
	for _, rule := range objectsOrEmpty(module["rules"]) {
		if rule["id"] == id {
			found = rule
			count++
		}
	}
	if count != 1 {
		t.Fatalf("rule %s occurs %d times", id, count)
	}
	return found
}
func normChallenge(s string) string { return strings.Join(strings.Fields(s), " ") }
func TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses(t *testing.T) {
	t.Parallel()
	module := challengeModule(t)
	rule := challengeRule(t, module, "rule.spec.rule-challenge")
	clauses := objectsOrEmpty(rule["clauses"])
	if len(clauses) != 2 {
		t.Fatalf("clauses = %d", len(clauses))
	}
	t.Run("the_rule", func(t *testing.T) {
		if rule["version"] != float64(1) || !reflect.DeepEqual(rule["coverage"], []any{"coverage.spec-workflow"}) {
			t.Errorf("rule metadata: %v", rule)
		}
		for i, want := range []struct{ id, force string }{{"clause.spec.rule-challenge-01-ask", "stop-and-ask"}, {"clause.spec.rule-challenge-02-answer", "mandatory"}} {
			c := clauses[i]
			if c["id"] != want.id || c["enforcement"] != want.force {
				t.Errorf("clause = %v", c)
			}
			if _, ok := c["replaces"]; ok {
				t.Error("unexpected replaces")
			}
		}
		var guide document
		for _, g := range objectsOrEmpty(module["supportingGuides"]) {
			if g["id"] == "guide.spec-routing" {
				guide = g
			}
		}
		if guide == nil {
			t.Fatal("missing guide")
		}
		rules, ok := guide["rules"].([]any)
		if !ok {
			t.Fatal("missing guide rules")
		}
		found := false
		for i := 1; i+2 < len(rules); i++ {
			if rules[i] == "rule.spec.rule-challenge" && rules[i-1] == "rule.spec.routing" && rules[i+1] == "rule.spec.open-authoring-question" && rules[i+2] == "rule.spec.project-constraints" {
				found = true
			}
		}
		if !found {
			t.Errorf("guide rules = %v", rules)
		}
		paths, err := filepath.Glob("assets/profiles/*.json")
		if err != nil || len(paths) == 0 {
			t.Fatalf("profiles: %v", err)
		}
		for _, path := range paths {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var p struct {
				Required []string `json:"requiredRules"`
			}
			if err := json.Unmarshal(data, &p); err != nil {
				t.Fatal(err)
			}
			for _, r := range p.Required {
				if r == "rule.spec.rule-challenge" {
					t.Errorf("%s requires rule", path)
				}
			}
		}
	})
	t.Run("names_no_record", func(t *testing.T) {
		re := regexp.MustCompile(`ADR-[0-9]|Spec [0-9]{4}`)
		for _, c := range clauses {
			if re.MatchString(c["guidance"].(string)) {
				t.Error("guidance names a repository record")
			}
		}
	})
	for _, tc := range []struct {
		name    string
		phrases []string
	}{
		{"ask/both_silences", []string{"never drop or narrow a solution in silence because a rule forbids it, and never contradict a rule in silence: put a Rule Challenge to the user."}},
		{"ask/what_a_rule_is", []string{"A rule is a decision the repository recorded: an accepted ADR, a Normative Clause of a guide under"}},
		{"ask/the_trigger", []string{"A rule blocks the goal when the solution you would recommend, or a part the goal needs, cannot exist without contradicting it, or when it stops you from stating the problem"}},
		{"ask/not_a_block", []string{"extra work, taste and friction are not a block."}},
		{"ask/the_question", []string{"name the rule and where it lives, what it imposes, up to three things it prevents for this goal, and your reading of why it exists and whether it looks out of date or badly written"}},
		{"ask/three_answers", []string{"then offer keep, change and retire, each with its consequence and the recommended one first."}},
		{"ask/binds_until_the_answer", []string{"The rule binds until the user answers: asking changes nothing"}},
		{"ask/never_an_answer", []string{"a recommendation, silence, a dismissal or a timeout is never an answer."}},
		{"ask/a_decided_change", []string{"Ask no Rule Challenge about a change the user already decided, because the request asks for it or it follows from an answer of the user that the artifact records"}},
		{"ask/a_kept_line_needs_a_new_fact", []string{"Ask none about a rule that carries a Kept Line unless the question names a fact the earlier one did not have."}},
		{"ask/at_most_three", []string{"Ask at most three in one Authoring Phase or brainstorming session"}},
		{"ask/the_others_keep_binding", []string{"list the others, which keep binding, in the artifact's open questions."}},
		{"ask/the_techspec_asks", []string{"The TechSpec author asks the same question before a design takes what an accepted ADR rejected, or otherwise contradicts its decision, unless the user already decided it."}},
		{"ask/implementation_guard_rail", []string{"In a Run, in a Task and in the QA gate a rule is a guard-rail: flag a conflict there and ask no Rule Challenge."}},
		{"answer/record_fields", []string{"one fielded record per question:", "plus `new fact` when the rule already carried a Kept Line."}},
		{"answer/none_when_no_rule_blocked", []string{"An Idea and a feature PRD always carry the section, which reads `None.` when no rule blocked the goal"}},
		{"answer/the_prd_carries_the_records", []string{"the PRD carries the Idea's records, and a question the TechSpec author asks is recorded in the PRD."}},
		{"answer/change_or_retire", []string{"On change or retire, the PRD states the change of the rule as a requirement, and a Task of the same Spec changes the rule together with the work that depended on it"}},
		{"answer/adr_lifecycle", []string{"an ADR through the lifecycle of `docs/agents/docs-layout.md`, any other rule in its own file."}},
		{"answer/kept_line_form", []string{"On keep, a Task of the Spec adds a Kept Line after the rule's text, `Kept (YYYY-MM-DD): <what was asked>; <why it stays>.`, which names no Spec"}},
		{"answer/newest_last", []string{"a rule keeps every Kept Line, the newest last"}},
		{"answer/no_spec_follows", []string{"when no Spec follows, the session that asked adds the line."}},
		{"answer/a_baseline_rule", []string{"A rule inside `setup-context-driven` markers changes only in the Baseline, and its Kept Line goes in `docs/agents/specific-repository.md`, under a heading of its own, naming the rule"}},
		{"answer/a_proposal_to_the_owner", []string{"a repository that does not own the Baseline sends the user's decision as one message to the Baseline's owner, obeys the rule until the Baseline changes, and records the message and the scope it left out."}},
		{"answer/nobody_present", []string{"When no person is present to answer, write no artifact for that phase and never answer for the user"}},
		{"answer/only_that_spec_waits", []string{"only that Spec waits, as a Backlog Entry that holds the whole question, names the Spec folder when one exists and carries `waiting_on` and `since` in its frontmatter, while other work goes on."}},
		{"answer/the_author_subagent", []string{"An Author Subagent returns the question to the session that started it."}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			i := 0
			if strings.HasPrefix(tc.name, "answer/") {
				i = 1
			}
			text := normChallenge(clauses[i]["guidance"].(string))
			for _, phrase := range tc.phrases {
				if n := strings.Count(text, normChallenge(phrase)); n != 1 {
					t.Errorf("phrase %q occurs %d times", phrase, n)
				}
			}
		})
	}
}
func TestTheSpecRoutingGuideRendersTheRuleChallengeClausesOnce(t *testing.T) {
	t.Parallel()
	module := challengeModule(t)
	clauses := objectsOrEmpty(challengeRule(t, module, "rule.spec.rule-challenge")["clauses"])
	route := objectsOrEmpty(challengeRule(t, module, "rule.spec.routing")["clauses"])
	constraints := objectsOrEmpty(challengeRule(t, module, "rule.spec.project-constraints")["clauses"])
	wanted := []string{route[len(route)-1]["guidance"].(string), clauses[0]["guidance"].(string), clauses[1]["guidance"].(string), constraints[0]["guidance"].(string)}
	check := func(t *testing.T, data []byte) {
		t.Helper()
		text := normChallenge(string(data))
		previous := -1
		for _, w := range wanted {
			w = normChallenge(w)
			if n := strings.Count(text, w); n != 1 {
				t.Fatalf("guidance occurs %d times", n)
			}
			at := strings.Index(text, w)
			if at <= previous {
				t.Fatal("guidance out of order")
			}
			previous = at
		}
	}
	t.Run("plan_postimage", func(t *testing.T) {
		plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), standardTypeScriptDecisions("make verify"))
		check(t, planPostimage(t, plan, "docs/agents/spec-routing.md").Content)
	})
	t.Run("formatter_golden", func(t *testing.T) {
		data, err := os.ReadFile("assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md")
		if err != nil {
			t.Fatal(err)
		}
		check(t, data)
	})
	t.Run("repository_guide", func(t *testing.T) {
		data, err := os.ReadFile("../../docs/agents/spec-routing.md")
		if err != nil {
			t.Fatal(err)
		}
		check(t, data)
	})
}
