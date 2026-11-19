// Boundary IN: PRD bytes and disposable Git repositories.
// Boundary OUT: real repositories, network, and authoring sessions.
package speccheck

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

const delegatedEntry = `- delegated: 2026-11-20, prd
  - question: Which behavior?
  - options: Retry; Fail fast
  - adopted: Retry — preserves behavior
  - under: Project Config
  - answer: none yet
`

func delegatedPRD(t *testing.T, text string) []byte {
	t.Helper()
	root := clarificationRepo(t, text)
	content, err := os.ReadFile(filepath.Join(root, clarificationFixturePRD))
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func TestDelegatedDecisionsParsesTheEntryForm(t *testing.T) {
	t.Parallel()
	entry := DelegatedDecision{Line: 3, Date: "2026-11-20", Phase: "prd", Question: "Which behavior?", Options: "Retry; Fail fast", Adopted: "Retry — preserves behavior", Under: "Project Config", Answer: "none yet"}
	answered := entry
	answered.Answer = `2026-11-22, changed — "Fail fast"`
	second := entry
	second.Line = 10
	for _, tc := range []struct {
		name, text string
		want       []DelegatedDecision
	}{
		{"a_complete_entry", "## Decisions\n\n" + delegatedEntry, []DelegatedDecision{entry}},
		{"an_answered_entry", "## Decisions\n\n" + strings.Replace(delegatedEntry, "none yet", answered.Answer, 1), []DelegatedDecision{answered}},
		{"two_entries_in_one_section", "## Decisions\n\n" + delegatedEntry + "\n" + delegatedEntry, []DelegatedDecision{entry, second}},
		{"a_fenced_entry_is_not_read", "## Decisions\n\n~~~markdown\n" + delegatedEntry + "~~~\n", nil},
		{"an_entry_outside_the_decisions_section_is_not_read", delegatedEntry + "## Decisions\n## Other\n" + delegatedEntry, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := DelegatedDecisions(delegatedPRD(t, tc.text))
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("entries=%#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestADelegatedDecisionThatCannotBeReadNamesItsFirstProblem(t *testing.T) {
	t.Parallel()
	headProblem := "the head needs a YYYY-MM-DD date and the phase prd, techspec or tasks"
	cases := []struct{ name, from, to, problem string }{
		{"invalid_date", "2026-11-20", "2026-02-30", headProblem},
		{"unknown_phase", ", prd", ", idea", headProblem},
		{"repeated_question", "  - question: Which behavior?", "  - question: Which behavior?\n  - question: Twice?", "question appears twice"},
		{"adopted_without_em_dash", "Retry — preserves", "Retry - preserves", "adopted needs the option, an em dash and the author's reason"},
		{"unknown_grant", "under: Project Config", "under: the maintainer", "under must be Project Config or the user's answer of YYYY-MM-DD"},
		{"unknown_answer_form", "answer: none yet", `answer: 2026-11-22, maybe — "x"`, "answer must be none yet, or YYYY-MM-DD, kept or changed, an em dash and the user's words"},
	}
	for _, field := range []string{"question", "options", "adopted", "under", "answer"} {
		for _, line := range strings.Split(delegatedEntry, "\n") {
			if strings.HasPrefix(line, "  - "+field+":") {
				cases = append(cases, struct{ name, from, to, problem string }{"missing_" + field, line + "\n", "", field + " is missing"})
			}
		}
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := DelegatedDecisions(delegatedPRD(t, "## Decisions\n"+strings.Replace(delegatedEntry, tc.from, tc.to, 1)))
			if len(got) != 1 || got[0].Problem != tc.problem {
				t.Fatalf("entries=%+v, want problem %q", got, tc.problem)
			}
		})
	}
}

func TestTheClarificationRecordCheckReportsAnUnreadableDelegatedDecision(t *testing.T) {
	t.Parallel()
	const route = "Autonomous route: no person present; decisions taken from the request.\n"
	for _, name := range []string{"an_unreadable_entry_names_its_line", "a_readable_entry_beside_the_autonomous_route", "only_a_delegated_decision_still_lacks_the_record", "before_the_horizon"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			entry := delegatedEntry
			if name == "an_unreadable_entry_names_its_line" || name == "before_the_horizon" {
				entry = strings.Replace(entry, "2026-11-20", "2026-02-30", 1)
			}
			text := "## Decisions\n\n" + entry
			if name != "only_a_delegated_decision_still_lacks_the_record" {
				text += "\n" + route
			}
			var root string
			if name == "before_the_horizon" {
				root = t.TempDir()
				gittest.InitRepo(t, root, "-b", "main")
				writeReceiptFixture(t, root, clarificationFixturePRD, text)
				receiptCommit(t, root)
				writeReceiptFixture(t, root, ClarificationGuidePath, "guide\n")
			} else {
				root = clarificationRepo(t, text)
			}
			receiptCommit(t, root)
			for _, stage := range []Stage{StagePRD, StageTechSpec, StageTasks, StageAll} {
				result := clarificationCheck(t, root, stage)
				got := receiptCodeFindings(result, CodeDecisionsUnrecorded)
				var want []Finding
				switch name {
				case "an_unreadable_entry_names_its_line":
					want = []Finding{{Code: CodeDecisionsUnrecorded, Severity: SeverityGap,
						Summary: clarificationFixturePRD + " has a Delegated Decision that cannot be read: the head needs a YYYY-MM-DD date and the phase prd, techspec or tasks",
						Where:   []Location{{Path: clarificationFixturePRD, Line: 3}},
						Fix:     `Write the entry as - delegated: YYYY-MM-DD, <prd|techspec|tasks> followed by the lines - question:, - options:, - adopted: <option> — <reason>, - under: and - answer:, each indented two spaces, following the write-prd rule "Clarification Record".`,
					}}
				case "only_a_delegated_decision_still_lacks_the_record":
					want = []Finding{{Code: CodeDecisionsUnrecorded, Severity: SeverityGap,
						Summary: clarificationFixturePRD + " has no Clarification Record", Where: []Location{{Path: clarificationFixturePRD, Line: 1}},
						Fix: `Under ## Decisions, write one - YYYY-MM-DD: <question> — <answer> line for each question a person answered, or the line Autonomous route: no person present; decisions taken from <source>., following the write-prd rule "Clarification Record".`,
					}}
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("stage=%s findings=%#v, want %#v", stage, got, want)
				}
				PromoteGaps(&result)
				for _, finding := range receiptCodeFindings(result, CodeDecisionsUnrecorded) {
					if finding.Severity != SeverityError {
						t.Fatalf("strict finding=%+v", finding)
					}
				}
			}
		})
	}
}

func TestDelegatedDecisionsPreserveTheSectionAndValidationBoundaries(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, text, problem string
		count, line         int
	}{
		{"a_valid_calendar_leap_day", strings.Replace(delegatedEntry, "2026-11-20", "2024-02-29", 1), "", 1, 2},
		{"a_dated_grant", strings.Replace(delegatedEntry, "under: Project Config", "under: the user's answer of 2024-02-29", 1), "", 1, 2},
		{"an_invalid_grant_date", strings.Replace(delegatedEntry, "under: Project Config", "under: the user's answer of 2026-02-30", 1), "under must be Project Config or the user's answer of YYYY-MM-DD", 1, 2},
		{"a_kept_answer", strings.Replace(delegatedEntry, "answer: none yet", "answer: 2026-11-22, kept — Retry", 1), "", 1, 2},
		{"an_invalid_answer_date", strings.Replace(delegatedEntry, "answer: none yet", "answer: 2026-02-30, kept — Retry", 1), "answer must be none yet, or YYYY-MM-DD, kept or changed, an em dash and the user's words", 1, 2},
		{"an_empty_reason", strings.Replace(delegatedEntry, "Retry — preserves behavior", "Retry — ", 1), "adopted needs the option, an em dash and the author's reason", 1, 2},
		{"the_first_missing_field", "- delegated: 2026-11-20, prd\n", "question is missing", 1, 2},
		{"missing_precedes_repeated", strings.Replace(delegatedEntry, "  - question: Which behavior?\n", "", 1) + "  - options: Duplicate\n", "question is missing", 1, 2},
		{"a_non_field_ends_the_entry", strings.Replace(delegatedEntry, "  - options:", "a paragraph\n  - options:", 1), "options is missing", 1, 2},
		{"a_malformed_head_begins_an_entry", "- delegated:nothing\n", "the head needs a YYYY-MM-DD date and the phase prd, techspec or tasks", 1, 2},
		{"a_heading_in_a_fence_does_not_end_the_section", "```markdown\n## Other\n```\n" + delegatedEntry, "", 1, 5},
		{"a_shorter_fence_does_not_close", "````markdown\n```\n" + delegatedEntry, "", 0, 0},
		{"a_real_heading_ends_the_section", "## Other\n" + delegatedEntry, "", 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := DelegatedDecisions(delegatedPRD(t, "## Decisions\n"+tc.text))
			if len(got) != tc.count {
				t.Fatalf("entries=%+v, want %d", got, tc.count)
			}
			if len(got) > 0 && (got[0].Problem != tc.problem || got[0].Line != tc.line) {
				t.Fatalf("entry=%+v, want problem %q line %d", got[0], tc.problem, tc.line)
			}
		})
	}
}
