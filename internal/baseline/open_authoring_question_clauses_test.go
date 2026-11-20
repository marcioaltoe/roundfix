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

func TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses(t *testing.T) {
	t.Parallel()
	module := challengeModule(t)
	rule := challengeRule(t, module, "rule.spec.open-authoring-question")
	clauses := objectsOrEmpty(rule["clauses"])
	if len(clauses) != 2 {
		t.Fatalf("clauses = %d", len(clauses))
	}
	t.Run("the_rule", func(t *testing.T) {
		if rule["version"] != float64(1) || !reflect.DeepEqual(rule["coverage"], []any{"coverage.spec-workflow"}) {
			t.Errorf("rule metadata: %v", rule)
		}
		for i, want := range []struct{ id, force string }{{"clause.spec.open-authoring-question-01-wait-or-follow", "stop-and-ask"}, {"clause.spec.open-authoring-question-02-record", "mandatory"}} {
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
		for i := 0; i+1 < len(rules); i++ {
			if rules[i] == "rule.spec.open-authoring-question" && rules[i+1] == "rule.spec.project-constraints" {
				found = true
			}
		}
		if !found {
			t.Errorf("guide rules = %v", rules)
		}
		all := objectsOrEmpty(module["rules"])
		found = false
		for i := 0; i+1 < len(all); i++ {
			if all[i]["id"] == "rule.spec.open-authoring-question" && all[i+1]["id"] == "rule.spec.artifacts" {
				found = true
			}
		}
		if !found {
			t.Error("rule must immediately precede artifacts")
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
				if r == "rule.spec.open-authoring-question" {
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
	for _, tc := range []struct{ name, phrase string }{
		{"ask/what_the_question_is", "An Open Authoring Question is a question the author of a PRD, a TechSpec or a Task Graph needs answered to finish its phase and that no source settles"},
		{"ask/the_question_contract", "with two or three options, the recommended one first and its reason"},
		{"ask/the_reserved_statement", "and say whether it is a Reserved Question and why."},
		{"ask/the_present_user", "When the user is present, ask: a question the user can see waits for the user"},
		{"ask/adopts_nothing", "a recommendation, silence, a dismissal or a timeout adopts nothing."},
		{"ask/three_ways_told", "Treat the user as absent only when the brief that started the session says so, the user said so, or the session runs as a continuation that an unattended flow armed."},
		{"ask/the_wait", "With the user absent the Spec waits, unless following is in force"},
		{"ask/the_preference_grant", "reports `follow` as the project's Authoring Question Preference"},
		{"ask/the_answer_grant", "quotes the user's answer, given before the flow, that the session follows the author."},
		{"ask/never_writes_the_preference", "Never write that preference into Project Config yourself."},
		{"ask/reserved_idea_stage", "the Idea stage and the choice an Idea-stage judgment leaves to the user;"},
		{"ask/reserved_rule_challenge", "a Rule Challenge;"},
		{"ask/reserved_authorization", "every authorization, which is a commit, a push, a Pull Request, a merge, a release, a Governed Path, a destructive Git operation, a database statement"},
		{"ask/reserved_scope", "removes one of their non-goals or contradicts a recorded answer of the user"},
		{"ask/reserved_spend", "that needs a paid service, a paid plan or a quota the repository does not already use"},
		{"ask/reserved_data", "that sends repository or client content to a service or a person the repository does not already send it to;"},
		{"ask/reserved_shape", "a question without two or three options and a recommended one;"},
		{"ask/reserved_unclassified", "and a question you cannot classify."},
		{"ask/two_readers", "An author that a session started returns the question and never adopts its own recommendation; the session checks the author's statement against this list, and when either reads the question as reserved or cannot tell, the Spec waits."},
		{"record/adopts_the_recommended_option", "When following is in force and an Open Authoring Question is not a Reserved Question, adopt the option the author recommended and go on with the phase."},
		{"record/the_entry_head", "the line `- delegated: YYYY-MM-DD, <phase>`, where the phase is `prd`, `techspec` or `tasks`"},
		{"record/the_five_fields", "followed by five lines indented two spaces, `- question:`, `- options:`, `- adopted:` with the option, an em dash and the author's reason"},
		{"record/the_grant_and_answer_fields", "`- under:` with `Project Config` or `the user's answer of YYYY-MM-DD`, and `- answer:` with `none yet`."},
		{"record/never_a_dated_answer", "Never write it as a dated answer of a person or under the autonomous-route line."},
		{"record/later_phases_in_the_prd", "A question of the TechSpec or the Tasks phase is recorded in the PRD, and that phase's artifact names the entry."},
		{"record/report_and_status", "List every Delegated Decision of the flow, and their count, in the session's report;"},
		{"record/inside_the_authorization", "The Spec goes on through authoring and delivery inside the authorization it already holds"},
		{"record/the_later_answer", "write it into the entry's `answer` line as the date, `kept` or `changed`, an em dash and the user's words."},
		{"record/after_delivery_starts", "leave the plan on the default branch as it is and make the change a correction of the Spec in delivery or a Backlog Entry"},
		{"record/the_waiting_entry", "only that Spec waits, as a Backlog Entry that names the Spec folder when one exists"},
		{"record/the_resume", "When the user answers, record the answer as a dated line of the PRD's `## Decisions`, let the Spec adopt the entry and resume the phase."},
		{"record/the_other_waits", "The waits of the Idea stage and of a Rule Challenge keep their own rules."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			i := 0
			if strings.HasPrefix(tc.name, "record/") {
				i = 1
			}
			if n := strings.Count(normChallenge(clauses[i]["guidance"].(string)), normChallenge(tc.phrase)); n != 1 {
				t.Errorf("phrase %q occurs %d times", tc.phrase, n)
			}
		})
	}
}
func TestTheSpecRoutingGuideRendersTheOpenAuthoringQuestionClausesOnce(t *testing.T) {
	t.Parallel()
	module := challengeModule(t)
	clauses := objectsOrEmpty(challengeRule(t, module, "rule.spec.open-authoring-question")["clauses"])
	route := objectsOrEmpty(challengeRule(t, module, "rule.spec.rule-challenge")["clauses"])
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
