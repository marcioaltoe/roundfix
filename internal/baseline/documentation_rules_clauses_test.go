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

func documentationModule(t *testing.T, name string) document {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("assets", "modules", name+".json"))
	if err != nil {
		t.Fatal(err)
	}
	var d document
	if err := json.Unmarshal(data, &d); err != nil {
		t.Fatal(err)
	}
	return d
}
func documentationClause(t *testing.T, id string) document {
	t.Helper()
	module := "context-workflow"
	if strings.HasPrefix(id, "clause.core.") {
		module = "core"
	}
	if strings.HasPrefix(id, "clause.secondbrain.") {
		module = "secondbrain"
	}
	d := documentationModule(t, module)
	var result document
	count := 0
	for _, r := range objectsOrEmpty(d["rules"]) {
		for _, c := range objectsOrEmpty(r["clauses"]) {
			if c["id"] == id {
				result = c
				count++
			}
		}
	}
	if count != 1 {
		t.Fatalf("%s occurs %d times", id, count)
	}
	return result
}
func documentationNorm(s string) string { return strings.Join(strings.Fields(s), " ") }
func TestTheBaselineModulesCarryTheDocumentationRules(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct{ name, id, phrase string }{
		{"adr/amended_stays_accepted", "clause.context.adr-01-template", "An accepted ADR that a later accepted ADR changes in part, without replacing it, stays `accepted`."},
		{"adr/the_later_adr_writes_it", "clause.context.adr-01-template", "In the change that adds the later ADR, its header gains `amended_by: [ADR-NNNN]`"},
		{"adr/dated_notice", "clause.context.adr-01-template", "a dated notice that begins `Amended by ADR-NNNN (YYYY-MM-DD)` and says what changed."},
		{"adr/the_update_fills", "clause.context.adr-01-template", "`roundfix baseline update` fills the list from the notices an ADR already carries."},
		{"legacy/sentence_kept", "clause.context.adr-03-legacy-compatibility", "Do not rewrite existing ADRs solely to adopt lifecycle metadata."},
		{"legacy/the_exception", "clause.context.adr-03-legacy-compatibility", "The one exception is a header field that a guide of this repository defines for existing ADRs, such as `amended_by`: it is written into them and changes no other line"},
		{"legacy/stays_legacy", "clause.context.adr-03-legacy-compatibility", "a header that holds no lifecycle field leaves its ADR a legacy ADR."},
		{"layout/design", "clause.context.docs-one-job-per-directory", "`docs/design/` for design artifacts of the product, such as an interface, API or data design, and never a decision or a plan"},
		{"layout/lessons", "clause.context.docs-one-job-per-directory", "`docs/lessons/` for Lessons Learned"},
		{"layout/root_reference", "clause.context.docs-one-job-per-directory", "including one that would otherwise sit at the repository root, the `DESIGN.md` design contract excepted"},
		{"layout/variation", "clause.context.docs-one-job-per-directory", "A folder under `docs/` or a reference document at the root that this clause does not name is a variation: name its path and its job in `docs/agents/specific-repository.md`, or move its content to the directory whose job it is."},
		{"layout/doctor_reports", "clause.context.docs-one-job-per-directory", "`roundfix doctor` reports one that is neither."},
		{"layout/rest_kept", "clause.context.docs-one-job-per-directory", "Preserve repository-authored bytes outside setup markers, and add a repository rule only to `docs/agents/specific-repository.md`."},
		{"retirement/lesson_deleted", "clause.context.findings-09-archive", "and a Lesson Learned once an accepted ADR, a guide or a skill states what it teaches are deleted"},
		{"message/no_message", "clause.core.message-the-owning-project", "When a Finding closes, send no message about what it taught: the docs-layout guide says when that becomes a Lesson Learned in this repository."},
		{"lesson/record_path", "clause.context.lesson-01-record", "It is one file, `docs/lessons/YYYY-MM-DD-<kebab-slug>.md`, in this complete copyable form"},
		{"lesson/labels_as_written", "clause.context.lesson-01-record", "keep each label as written and write the rest in the language of the repository's documentation"},
		{"lesson/the_test", "clause.context.lesson-02-test", "Write a Lesson Learned only when every part of this test holds: something was tried, it went wrong, the cause is proven, something solved it, and no accepted ADR, reference document or glossary entry already states it."},
		{"lesson/otherwise_nothing", "clause.context.lesson-02-test", "Otherwise write nothing."},
		{"lesson/two_moments", "clause.context.lesson-02-test", "when a Finding closes, in the change that answers it or retires it, and when a Spec that needed a correction, a failed QA or a retry is archived."},
		{"lesson/in_that_change", "clause.context.lesson-02-test", "Write the file in that change, so the Pull Request that carries it shows the lesson to a person, who may delete it."},
		{"export/two_parts", "clause.secondbrain.export-01-standard-list", "This repository publishes to the Secondbrain only what `.secondbrain-export` lists, and that file has two parts."},
		{"export/markers", "clause.secondbrain.export-01-standard-list", "between the lines `# roundfix:standard-export:begin` and `# roundfix:standard-export:end`"},
		{"export/standard_items", "clause.secondbrain.export-01-standard-list", "the README, the glossary files, `docs/user-guide/`, `docs/references/`, `docs/design/`, `docs/lessons/` and each ADR whose header says `knowledge: true`, one path per line."},
		{"export/baseline_owns", "clause.secondbrain.export-01-standard-list", "The Baseline owns those lines: `roundfix baseline update` rewrites them, and the only line added by hand is the path of an ADR, in order, when the ADR is written with that mark."},
		{"export/never_process", "clause.secondbrain.export-01-standard-list", "Never publish process: `AGENTS.md`, `CLAUDE.md`, `docs/agents/`, the Spec folders, `docs/backlog/`, `docs/findings/`, `docs/handoffs/` and `docs/history/` stay out of both parts."},
		{"variation/own_lines", "clause.secondbrain.export-02-variation", "The second part is the repository's own lines, everything outside those two markers, which no Baseline operation changes."},
		{"variation/exclusion", "clause.secondbrain.export-02-variation", "or drop a standard item with an exclusion line: `!` followed by the item as the list writes it, or `!docs/adr/*` for every ADR."},
		{"variation/reason", "clause.secondbrain.export-02-variation", "Give each such line its reason in a comment line directly above it; one comment stands for the lines that follow it without a blank line."},
		{"variation/doctor", "clause.secondbrain.export-02-variation", "`roundfix doctor` reports a line without a reason, a process path, and a standard list that differs from what the update would write."},
		{"mark/add_it", "clause.secondbrain.export-03-knowledge-mark", "add `knowledge: true` or `knowledge: false` to its header, on the line after `status`."},
		{"mark/the_test", "clause.secondbrain.export-03-knowledge-mark", "It is `true` when the decision states a rule of the business or of the product that someone outside this repository's development can use, and `false` when it is about how this repository is built, tested, organized or delivered, or about a tool it uses."},
		{"mark/unmarked_unpublished", "clause.secondbrain.export-03-knowledge-mark", "An ADR without the mark is not published."},
		{"mark/existing_once", "clause.secondbrain.export-03-knowledge-mark", "The ADRs a repository already has are marked once, by the same test, by the Agent of the session that applies the Baseline update that brings this rule"},
		{"mark/one_reviewed_commit", "clause.secondbrain.export-03-knowledge-mark", "in one commit that a person reviews in a Pull Request, followed by a second update so the list names them."},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := documentationClause(t, tt.id)
			text := documentationNorm(c["guidance"].(string))
			if n := strings.Count(text, documentationNorm(tt.phrase)); n != 1 {
				t.Fatalf("phrase count = %d, want 1: %s", n, tt.phrase)
			}
			if tt.name == "adr/amended_stays_accepted" && strings.Index(text, tt.phrase) < strings.LastIndex(text, "```") {
				t.Fatal("amendment precedes code fence")
			}
			if tt.name == "retirement/lesson_deleted" && (!strings.HasPrefix(text, "A retirement deletes the file.") || len(c["replaces"].([]any)) == 0) {
				t.Fatal("retirement sentinel or replacements missing")
			}
		})
	}
	t.Run("lesson/eight_lines", func(t *testing.T) {
		s := documentationClause(t, "clause.context.lesson-01-record")["guidance"].(string)
		last := -1
		for _, p := range []string{"created: YYYY-MM-DD", "- Context:", "- Tried:", "- Went wrong:", "- Cause:", "- Solved by:", "- Applies when:", "- Does not apply when:", "- Links:"} {
			i := strings.Index(s, p)
			if i <= last || strings.Count(s, p) != 1 {
				t.Fatalf("%q missing, repeated or out of order", p)
			}
			last = i
		}
	})
	for _, tt := range []struct {
		name, module, rule, coverage, guide, previous string
		ids                                           []string
	}{
		{"the_lesson_rule", "context-workflow", "rule.context.lessons", "coverage.context-workflow", "guide.docs-layout", "rule.context.docs-layout", []string{"clause.context.lesson-01-record", "clause.context.lesson-02-test"}},
		{"the_export_rule", "secondbrain", "rule.secondbrain.knowledge-export", "coverage.secondbrain", "guide.secondbrain", "rule.secondbrain.hermes-escalation", []string{"clause.secondbrain.export-01-standard-list", "clause.secondbrain.export-02-variation", "clause.secondbrain.export-03-knowledge-mark"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := documentationModule(t, tt.module)
			count := 0
			for _, r := range objectsOrEmpty(d["rules"]) {
				if r["id"] != tt.rule {
					continue
				}
				count++
				if r["version"] != float64(1) || !reflect.DeepEqual(r["coverage"], []any{tt.coverage}) {
					t.Fatal("rule metadata differs")
				}
				cs := objectsOrEmpty(r["clauses"])
				if len(cs) != len(tt.ids) {
					t.Fatal("clause count differs")
				}
				for i, c := range cs {
					if c["id"] != tt.ids[i] || c["enforcement"] != "mandatory" {
						t.Fatal("clause order or force differs")
					}
					if _, ok := c["replaces"]; ok {
						t.Fatal("unexpected replaces")
					}
				}
			}
			if count != 1 {
				t.Fatalf("rule count %d", count)
			}
			for _, g := range objectsOrEmpty(d["supportingGuides"]) {
				if g["id"] == tt.guide {
					rules := g["rules"].([]any)
					for i, v := range rules {
						if v == tt.previous && i+1 < len(rules) && rules[i+1] == tt.rule && i+2 == len(rules) {
							return
						}
					}
					t.Fatal("guide rule ordering differs")
				}
			}
			t.Fatal("missing guide")
		})
	}
	t.Run("required_by_no_profile", func(t *testing.T) {
		files, err := filepath.Glob("assets/profiles/*.json")
		if err != nil || len(files) == 0 {
			t.Fatal("missing profiles", err)
		}
		for _, f := range files {
			b, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			var d document
			if err = json.Unmarshal(b, &d); err != nil {
				t.Fatal(err)
			}
			for _, r := range func() []any { a, _ := d["requiredRules"].([]any); return a }() {
				if r == "rule.context.lessons" || r == "rule.secondbrain.knowledge-export" {
					t.Errorf("%s requires %v", f, r)
				}
			}
		}
	})
	t.Run("names_no_record", func(t *testing.T) {
		re := regexp.MustCompile(`ADR-[0-9]|Spec [0-9]{4}`)
		for _, id := range documentationIDs {
			if re.MatchString(documentationClause(t, id)["guidance"].(string)) {
				t.Errorf("%s names a record", id)
			}
		}
	})
}

var documentationIDs = []string{"clause.context.adr-01-template", "clause.context.adr-03-legacy-compatibility", "clause.context.docs-one-job-per-directory", "clause.context.findings-09-archive", "clause.context.lesson-01-record", "clause.context.lesson-02-test", "clause.core.message-the-owning-project", "clause.secondbrain.export-01-standard-list", "clause.secondbrain.export-02-variation", "clause.secondbrain.export-03-knowledge-mark"}

func TestTheGuidesRenderTheDocumentationRulesOnce(t *testing.T) {
	t.Parallel()
	check := func(t *testing.T, read func(string) string) {
		t.Helper()
		for _, id := range documentationIDs {
			guide := "docs-layout"
			if strings.HasPrefix(id, "clause.core.") {
				guide = "agent-instructions"
			}
			if strings.HasPrefix(id, "clause.secondbrain.") {
				guide = "secondbrain"
			}
			text := documentationNorm(read(guide))
			want := documentationNorm(documentationClause(t, id)["guidance"].(string))
			if n := strings.Count(text, want); n != 1 {
				t.Errorf("%s occurs %d times in %s", id, n, guide)
			}
		}
		for _, pair := range [][3]string{{"docs-layout", "clause.context.inbox-03-extend-before-minting", "clause.context.lesson-01-record"}, {"secondbrain", "clause.secondbrain.baseline-owned-guidance", "clause.secondbrain.export-01-standard-list"}, {"docs-layout", "clause.context.lesson-01-record", "clause.context.lesson-02-test"}, {"secondbrain", "clause.secondbrain.export-01-standard-list", "clause.secondbrain.export-02-variation"}, {"secondbrain", "clause.secondbrain.export-02-variation", "clause.secondbrain.export-03-knowledge-mark"}} {
			text := documentationNorm(read(pair[0]))
			a := documentationNorm(documentationClause(t, pair[1])["guidance"].(string))
			b := documentationNorm(documentationClause(t, pair[2])["guidance"].(string))
			if strings.Index(text, a) >= strings.Index(text, b) {
				t.Errorf("%s order incorrect", pair[0])
			}
		}
	}
	t.Run("plan_postimage", func(t *testing.T) {
		decisions := standardTypeScriptDecisions("make verify")
		for i := range decisions {
			if decisions[i].ID == "secondbrain.enabled" {
				decisions[i].Value = true
			}
		}
		plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), decisions)
		check(t, func(g string) string { return string(planPostimage(t, plan, "docs/agents/"+g+".md").Content) })
	})
	t.Run("formatter_golden", func(t *testing.T) {
		check(t, func(g string) string {
			b, e := os.ReadFile("assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/" + g + ".md")
			if e != nil {
				t.Fatal(e)
			}
			return string(b)
		})
	})
	t.Run("repository_guide", func(t *testing.T) {
		check(t, func(g string) string {
			b, e := os.ReadFile("../../docs/agents/" + g + ".md")
			if e != nil {
				t.Fatal(e)
			}
			return string(b)
		})
	})
}
