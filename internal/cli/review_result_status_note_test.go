// Suite: Result Status Note dismissal through the review command.
// Boundary IN: committed Task trees and fake reviewer findings.
// Boundary OUT: sealed validator requests, review records, exit codes and stderr.
package cli

import (
	"fmt"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestReviewDismissesAResultStatusNoteOfACompletedTaskWithoutTheValidator(t *testing.T) {
	t.Parallel()
	tests := []struct{ name, body, anchor string }{
		{"0297 shape", "---\ntask: task_01\nspec: 0001-example\nstatus: completed\ntype: backend\n---\n\n## Requirements\nDo the work.\n\n## Result\nImplemented the backend slice.\n\nThe declared Task Verification and repository Verification were not run;\nthey remain Daemon-owned. Task status remains `in_progress`, and no commit,\npush or pull request was made.\n", "The declared Task Verification"},
		{"0290 shape", "---\ntask: task_01\nspec: 0001-example\nstatus: completed\ntype: docs\n---\n\n## Requirements\nDo the work.\n\n## Result\n\nImplemented the documentation slice while leaving the Daemon-owned status\n`in_progress` unchanged:\n", "## Result\n"},
		{"front matter anchor", validatorTaskBody, "status: completed"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			line := strings.Count(test.body[:strings.Index(test.body, test.anchor)], "\n") + 1
			if test.name == "0290 shape" {
				line++
			}
			finding := fmt.Sprintf("`%s:%d` — The Task is `completed` while its Result says `in_progress` or Daemon-owned. Failure: an invalid Task state transition.", validatorTaskPath, line)
			fixture, runner := validatorFixture(t, "- "+finding, nil)
			validatorWrite(t, fixture.repository, validatorTaskPath, test.body)
			gittest.Run(t, fixture.repository, "add", "docs")
			gittest.Run(t, fixture.repository, "commit", "--allow-empty", "-m", "test: Result Status Note")
			fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
			code, record, stderr := fixture.run(t)
			if code != exitOK || record.Outcome != reviewOutcomeFindingsDismissed || record.Validation.Validator != "not-needed" || len(runner.sealedRequests) != 0 || record.Validation.Conventions != deliveryConventionsVersion {
				t.Fatalf("exit=%d record=%+v calls=%d", code, record, len(runner.sealedRequests))
			}
			if len(record.FindingItems) != 1 {
				t.Fatalf("findings=%+v", record.FindingItems)
			}
			item := record.FindingItems[0]
			if item.ID != "F1" || item.Text != finding || !reviewFindingDismissedByValidation(item) || item.Validation.Rule != "convention:C2" || item.Validation.Reason != reviewResultStatusNoteReason {
				t.Fatalf("item=%+v", item)
			}
			if stderr != "roundfix: review finding F1 dismissed by validation (convention:C2): "+reviewResultStatusNoteReason+"\n" {
				t.Fatalf("stderr=%q", stderr)
			}
		})
	}
}

func TestReviewKeepsAContradictionThatIsNotAResultStatusNote(t *testing.T) {
	t.Parallel()
	const pair = "The Task is completed while its Result says in progress. Failure: an invalid Task state transition."
	tests := []struct {
		name, body, text string
		outside          bool
	}{
		{"failed Task", strings.Replace(validatorTaskBody, "status: completed", "status: failed", 1), pair, false},
		{"QA Task", strings.Replace(validatorTaskBody, "status: completed", "status: completed\ntype: qa", 1), pair, false},
		{"Result without a note", strings.Replace(validatorTaskBody, "Status is Daemon-owned.", "Implemented the slice.", 1), pair, false},
		{"note outside the Result", strings.Replace(strings.Replace(validatorTaskBody, "Status is Daemon-owned.", "Implemented the slice.", 1), "Do the work.", "Status is Daemon-owned.", 1), pair, false},
		{"finding names another failure", validatorTaskBody, "The Task is completed while its Result gives no evidence. Failure: an invalid Task state transition.", false},
		{"anchor outside the settlement region", validatorTaskBody, pair, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			line := 4
			if test.outside {
				line = 8
			}
			finding := fmt.Sprintf("`%s:%d` — %s", validatorTaskPath, line, test.text)
			fixture, runner := validatorFixture(t, "- "+finding, validatorAnswer(validatorVerdict{"F1", "stands", "", "a real state defect"}))
			validatorWrite(t, fixture.repository, validatorTaskPath, test.body)
			gittest.Run(t, fixture.repository, "add", "docs")
			gittest.Run(t, fixture.repository, "commit", "--allow-empty", "-m", "test: other contradiction")
			fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
			code, record, stderr := fixture.run(t)
			if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || len(record.FindingItems) != 1 || record.FindingItems[0].Validation.Status != "stands" || stderr != "" {
				t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
			}
			if test.outside {
				if len(runner.sealedRequests) != 0 || record.Validation.Validator != "not-needed" {
					t.Fatalf("record=%+v calls=%d", record, len(runner.sealedRequests))
				}
				return
			}
			questions := validatorQuestions(t, runner)
			if len(questions) != 1 || len(questions[0].EligibleRules) != 1 || questions[0].EligibleRules[0] != "convention:C2" || record.Validation.Validator != "ran" || record.FindingItems[0].Validation.Reason != "a real state defect" {
				t.Fatalf("record=%+v questions=%+v", record, questions)
			}
		})
	}
}
