// Suite: Implement usage summary.
// Invariant: every terminal outcome prints stored usage immediately after its outcome.
// Boundary IN: public CLI, temporary repository/Home, fake Agent result.
// Boundary OUT: parsing usage from an ACP stream and price computation.
package cli

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	runworktree "roundfix/internal/worktree"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

type usageImplementRunner struct {
	*implementFakeRunner
	usage agent.TurnUsage
}

func (runner *usageImplementRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	result, err := runner.implementFakeRunner.Run(ctx, req, sink)
	result.Usage = runner.usage
	return result, err
}

func TestImplementSummaryEndsWithTheTokensLine(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, state string
		code        int
		err         error
	}{
		{name: "clean", state: store.StateClean, code: 0},
		{name: "unresolved", state: store.StateUnresolved, code: 1, err: errors.New("agent failed")},
		{name: "stopped", state: store.StateStopped, code: 0, err: agent.StopError{Err: context.Canceled}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}, errByTask: map[string]error{"task_01": tc.err}}
			withImplementCollaborators(t, runner)
			var out, diag bytes.Buffer
			code := runCLI(t, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &out, &diag)
			lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
			if code != tc.code || len(lines) < 2 || !strings.HasPrefix(lines[len(lines)-2], tc.state+":") || lines[len(lines)-1] != "Tokens: none reported by 1 prompt(s); cost not reported" {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diag)
			}
		})
	}
}

func TestImplementSummaryPrintsRecordedUsage(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}, {id: "task_02", needs: []string{"task_01"}}})
	runner := &usageImplementRunner{implementFakeRunner: &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted, "task_02": spec.StatusCompleted}}, usage: agent.TurnUsage{Basis: agent.UsageBasisRequestSum, TotalTokens: 5639755, Readings: 46}}
	withImplementCollaborators(t, runner)
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &out, &diag)
	want := "Clean: all 2 Task(s) completed.\nTokens: 11279510 from 2 of 2 prompt(s); cost not reported\n"
	if code != 0 || !strings.HasSuffix(out.String(), want) {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diag)
	}
	runID := implementRunIDFromStderr(t, diag.String())
	var shown, showDiag bytes.Buffer
	if code := runCLI(t, []string{"runs", "show", runID}, &shown, &showDiag); code != 0 || !strings.HasSuffix(shown.String(), strings.SplitN(want, "\n", 2)[1]) {
		t.Fatalf("show exit=%d stdout=%s stderr=%s home=%s", code, &shown, &showDiag, home)
	}
}

func TestUnreportedUsageIsNeverPrintedAsZero(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}, {id: "task_02", needs: []string{"task_01"}}})
	withImplementCollaborators(t, &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted, "task_02": spec.StatusCompleted}})
	var implemented, diag bytes.Buffer
	if code := runCLI(t, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &implemented, &diag); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, &diag)
	}
	runID := implementRunIDFromStderr(t, diag.String())
	writer, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	queue, err := writer.CreateDeliveryQueue(t.Context(), repo, []string{implementTestSlug})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.RunID = runID
	if err := writer.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"runs", "show", runID}, {"deliver", "status"}} {
		var out, diag bytes.Buffer
		if code := runCLI(t, args, &out, &diag); code != 0 {
			t.Fatalf("%v exit=%d stderr=%s", args, code, &diag)
		}
		assertUnreportedOutput(t, out.String())
	}
	assertUnreportedOutput(t, implemented.String())
}

func assertUnreportedOutput(t *testing.T, output string) {
	t.Helper()
	if !strings.Contains(output, "none reported by 2 prompt(s)") || strings.Contains(output, "0 tokens") {
		t.Fatalf("unreported usage lost: %s", output)
	}
}

func TestImplementTokenSummaryCoversSetupIntegrationAndPush(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"setup failure", "integration pending", "integration budget", "push failure", "push success"} {
		t.Run(name, func(t *testing.T) {
			home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}}
			_, _, pusher, _ := withImplementCollaborators(t, runner)
			state, code, tokenLine := store.StateFailed, exitRunFailed, "Tokens: none reported by 1 prompt(s); cost not reported"
			switch name {
			case "setup failure":
				tokenLine = "Tokens: no prompts recorded"
				updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
					deps.createRunWorktree = func(context.Context, runworktree.CreateOptions) (runworktree.Ref, error) {
						return runworktree.Ref{}, errors.New("setup failed")
					}
				})
			case "integration pending":
				state = store.StateIntegrationPending
				updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
					deps.integrateRunWorktree = func(context.Context, runworktree.Ref, string, string) (runworktree.IntegrationResult, error) {
						return runworktree.IntegrationResult{Mode: runworktree.ModePending}, nil
					}
				})
			case "integration budget":
				state = store.StateBudgetExceeded
				mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), "budget:\n  max_run_duration: 1h\n")
				gitImplement(t, repo, "add", ".roundfixrc.yml")
				gitImplement(t, repo, "commit", "-m", "configure usage fixture budget")
				now := time.Now()
				updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
					deps.implementBudgetNow = func() time.Time { return now }
					deps.integrateRunWorktree = func(context.Context, runworktree.Ref, string, string) (runworktree.IntegrationResult, error) {
						now = now.Add(2 * time.Hour)
						return runworktree.IntegrationResult{}, context.DeadlineExceeded
					}
				})
			case "push failure", "push success":
				configureImplementAutoPush(t, repo, true)
				configureImplementUpstream(t, repo, "origin", "ma/widget-flow")
				if name == "push failure" {
					pusher.err = errors.New("push failed")
				} else {
					state = store.StateClean
					code = exitOK
				}
			}
			var out, diag bytes.Buffer
			gotCode := runCLI(t, []string{"implement", "--spec", implementTestSlug, "--no-input"}, &out, &diag)
			lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
			index := -1
			for i, line := range lines {
				if strings.HasPrefix(line, state+":") {
					index = i
					break
				}
			}
			if gotCode != code || index < 0 || index+1 >= len(lines) || lines[index+1] != tokenLine || strings.Count(out.String(), "Tokens:") != 1 {
				t.Fatalf("case=%s exit=%d stdout=%s stderr=%s", name, gotCode, &out, &diag)
			}
			if name == "push failure" && lines[index] != "Failed: 1 completed, 0 failed, 0 skipped, 0 pending." {
				t.Fatalf("failure discarded settled task count: %s", &out)
			}
			if name == "push success" && (index+2 >= len(lines) || lines[index+2] != "pushed origin/ma/widget-flow") {
				t.Fatalf("tokens must precede push: %s", &out)
			}
			runID := onlyImplementRunID(t, home)
			if run := implementRunFromStore(t, home, runID); run.State != state {
				t.Fatalf("state=%s want=%s", run.State, state)
			}
		})
	}
}
