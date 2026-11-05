// Suite: unresolved Run parks preserve failed and skipped Tasks.
// Boundary IN: delivery engine transitions and the real queue store.
// Boundary OUT: Implement and workspace operations use the existing fake workflow.
package delivery

import (
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestARunUnresolvedParkNamesItsFailedAndSkippedTasks(t *testing.T) {
	t.Parallel()
	const (
		slug          = "0077-example"
		secretReason  = "Task commit carries a secret: fixture.txt:2: credential-assignment"
		genericAnswer = "resolve the blocker, then run roundfix deliver retry " + slug
		secretAnswer  = "if the failed Task's text prescribes the value the blocker locates, amend it in /worktrees/0077-example to one no Secret Pattern matches (letters only, or fewer than 20 characters); then run roundfix deliver retry 0077-example, whose rerun returns any Secret Match to the Agent as Verification Feedback"
	)
	for _, tt := range []struct {
		name    string
		result  RunResult
		blocker string
		answer  string
		class   ParkClass
	}{
		{
			name:    "secret refusal",
			result:  RunResult{FailedTasks: []TaskFailure{{Task: "task_01", Reason: secretReason}}, SkippedTasks: []string{"task_02", "task_03"}},
			blocker: "run-unresolved: task_01 failed: " + secretReason + "; skipped: task_02, task_03",
			answer:  secretAnswer, class: ParkClassFinding,
		},
		{
			name:    "two failures without skips",
			result:  RunResult{FailedTasks: []TaskFailure{{Task: "task_01", Reason: "first failure"}, {Task: "task_02", Reason: "second failure"}}},
			blocker: "run-unresolved: task_01 failed: first failure; task_02 failed: second failure",
			answer:  genericAnswer, class: ParkClassFinding,
		},
		{
			name:    "long reason",
			result:  RunResult{FailedTasks: []TaskFailure{{Task: "task_01", Reason: strings.Repeat("x", 400)}}},
			blocker: "run-unresolved: task_01 failed: " + strings.Repeat("x", 300) + "…",
			answer:  genericAnswer, class: ParkClassFinding,
		},
		{
			name:    "nothing failed or skipped",
			blocker: "run-unresolved", answer: genericAnswer, class: ParkClassFinding,
		},
		{
			name:    "normalized unicode reason",
			result:  RunResult{FailedTasks: []TaskFailure{{Task: "task_01", Reason: " \n\t" + strings.Repeat("é", 299) + "\t\n tail "}}},
			blocker: "run-unresolved: task_01 failed: " + strings.Repeat("é", 299) + " …",
			answer:  genericAnswer, class: ParkClassFinding,
		},
		{
			name:    "runtime infrastructure takes precedence",
			result:  RunResult{FailedTasks: []TaskFailure{{Task: "task_01", Reason: secretReason}}, RuntimeInfrastructure: "task_01 lost its rollout"},
			blocker: "runtime-infrastructure: task_01 lost its rollout",
			answer:  "run roundfix deliver retry 0077-example; a retry from runtime-infrastructure is not counted", class: ParkClassEnvironment,
		},
		{
			name:    "QA environment partial takes precedence",
			result:  RunResult{FailedTasks: []TaskFailure{{Task: "task_01", Reason: secretReason}}, QAEnvironmentPartial: true},
			blocker: "qa-environment-partial",
			answer:  "(cd /worktrees/0077-example && roundfix reconcile run-unresolved-detail --carry-forward), satisfy the environment-blocked QA rows, run roundfix archive 0077-example --qa-override --approval <source> --reason <text>, then run roundfix deliver retry 0077-example", class: ParkClassEnvironment,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			runStore := openDeliveryEngineStore(t, ctx)
			const root = "/repo-unresolved-detail"
			if _, err := runStore.CreateDeliveryQueue(ctx, root, []string{slug}); err != nil {
				t.Fatal(err)
			}
			workflow := newFakeDeliveryWorkflow()
			tt.result.RunID = "run-unresolved-detail"
			tt.result.Outcome = RunOutcomeUnresolved
			workflow.runs[slug] = tt.result
			workflow.recordWorkspace = func(branch, worktree string) error {
				_, _, _, err := runStore.RecordDeliveryQueueItemWorktree(ctx, root, slug, branch, worktree)
				return err
			}
			engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())
			if _, err := engine.Run(ctx, root); err != nil {
				t.Fatal(err)
			}
			queue := readDeliveryQueue(t, ctx, runStore, root)
			item := queue.Items[0]
			if item.Stage != store.DeliveryStageParked || item.Blocker != tt.blocker || item.RunID != tt.result.RunID {
				t.Fatalf("parked item = %+v, want blocker %q and Run %q", item, tt.blocker, tt.result.RunID)
			}
			if got := ClassifyPark(queue, item); got.Class != tt.class || got.Next != tt.answer {
				t.Fatalf("classification = %+v, want class %q and next %q", got, tt.class, tt.answer)
			}
			question, found := PendingQuestionFor(queue)
			if !found || question.SpecSlug != slug || question.Blocker != tt.blocker || question.Answer != tt.answer || question.Waiting != 0 {
				t.Fatalf("Pending Question = %+v, found=%v", question, found)
			}
		})
	}
}
