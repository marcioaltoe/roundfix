// Suite: the checking stage waits for every member of the Head Check Set.
// Invariant: CLEAN never advances a pending or empty set, and failures name their source.
// Boundary IN: real queue persistence, fake-clock polling, scripted reports and a log buffer.
// Boundary OUT: GitHub and check recovery use injected fakes; no network or real-time waits.
package delivery

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"roundfix/internal/store"
)

func TestTheQueueWaitsForWorkflowRunsWhoseChecksHaveNotRegistered(t *testing.T) {
	t.Parallel()
	for _, passing := range []bool{false, true} {
		name := "backend fails"
		if passing {
			name = "every member passes"
		}
		t.Run(name, func(t *testing.T) {
			engine, runStore, item, clock, boundary := newCheckWaitFixture(t, time.Minute)
			var log bytes.Buffer
			engine.log = &log
			recovery := &fakeCheckRecovery{failures: []CheckFailure{{RunID: "3", Attempt: 1, Packages: []string{"roundfix/internal/changed"}}}}
			engine.checks = recovery
			polls := 0
			start := clock.Now()
			boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
				polls++
				if polls > 3 {
					t.Fatal("unexpected extra poll")
				}
				report := CheckReport{
					HeadSHA: "candidate", MergeState: "CLEAN",
					Checks: []PullRequestCheck{{Name: "Validate PR title", Workflow: "CI | Conventions", Bucket: "pass", Link: "https://github.com/example/repo/actions/runs/1/job/10"}},
					WorkflowRuns: []WorkflowRun{
						{ID: "1", Name: "CI | Conventions", Status: "completed", Conclusion: "success", Attempt: 1},
						{ID: "2", Name: "CI | Frontend", Status: "in_progress", Attempt: 1},
						{ID: "3", Name: "CI | Backend", Status: "in_progress", Attempt: 1},
					},
				}
				if polls >= 2 {
					for index, workflow := range []string{"CI | Frontend", "CI | Backend"} {
						report.Checks = append(report.Checks, PullRequestCheck{Name: "Lint, Type Check, Test & Build", Workflow: workflow, Bucket: "pending", Link: fmt.Sprintf("https://github.com/example/repo/actions/runs/%d/job/%d", index+2, (index+2)*10)})
					}
				}
				if polls == 3 {
					for index := 1; index < 3; index++ {
						report.Checks[index].Bucket = "pass"
						report.WorkflowRuns[index].Status = "completed"
						report.WorkflowRuns[index].Conclusion = "success"
					}
					if !passing {
						report.Checks[2].Bucket, report.Checks[2].State = "fail", "FAILURE"
						report.WorkflowRuns[2].Conclusion = "failure"
					}
				}
				return report, nil
			}
			if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
				t.Fatal(err)
			}
			stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
			if polls != 3 || clock.Now().Sub(start) != 30*time.Second {
				t.Fatalf("polls=%d elapsed=%s item=%+v", polls, clock.Now().Sub(start), stored)
			}
			wantLog := "roundfix: checks: Delivery Queue item example: merge state CLEAN\n" +
				"roundfix: checks: Delivery Queue item example: waiting for 2 check(s): CI | Frontend (in_progress); CI | Backend (in_progress)\n" +
				"roundfix: checks: Delivery Queue item example: waiting for 2 check(s): CI | Frontend / Lint, Type Check, Test & Build; CI | Backend / Lint, Type Check, Test & Build\n"
			if passing {
				wantLog += "roundfix: checks: Delivery Queue item example: 3 check(s) completed for candidate\n"
				if stored.Stage != store.DeliveryStageMerging || stored.Blocker != "" || recovery.inspections != 0 {
					t.Fatalf("item=%+v inspections=%d", stored, recovery.inspections)
				}
			} else {
				const want = "checks-failed: CI | Backend / Lint, Type Check, Test & Build (failure) https://github.com/example/repo/actions/runs/3/job/30"
				if stored.Stage != store.DeliveryStageParked || stored.Blocker != want || recovery.inspections != 1 || recovery.reruns != 0 || ClassifyPark(store.DeliveryQueue{}, stored).Class != ParkClassFinding {
					t.Fatalf("item=%+v inspections=%d reruns=%d", stored, recovery.inspections, recovery.reruns)
				}
			}
			if log.String() != wantLog {
				t.Fatalf("log=%q want %q", log.String(), wantLog)
			}
		})
	}
}

func TestAnExpectedWorkflowWithoutARunKeepsTheItemWaiting(t *testing.T) {
	t.Parallel()
	for _, arrives := range []bool{false, true} {
		name := "no run yet"
		if arrives {
			name = "the run arrives"
		}
		t.Run(name, func(t *testing.T) {
			engine, runStore, item, _, boundary := newCheckWaitFixture(t, time.Minute)
			var log bytes.Buffer
			engine.log = &log
			polls := 0
			boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
				polls++
				report := CheckReport{HeadSHA: "candidate", Checks: []PullRequestCheck{{Name: "coverage", Bucket: "pass"}}, ExpectedWorkflows: []string{".github/workflows/ci-verify.yml"}}
				if arrives && polls >= 3 {
					report.WorkflowRuns = []WorkflowRun{{ID: "7", Path: ".github/workflows/ci-verify.yml", Status: "completed", Conclusion: "success"}}
				}
				return report, nil
			}
			if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
				t.Fatal(err)
			}
			stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
			wantLog := "roundfix: checks: Delivery Queue item example: waiting for 1 check(s): .github/workflows/ci-verify.yml (no run yet)\n"
			if arrives {
				wantLog += "roundfix: checks: Delivery Queue item example: 2 check(s) completed for candidate\n"
				if stored.Stage != store.DeliveryStageMerging || polls != 3 {
					t.Fatalf("item=%+v polls=%d", stored, polls)
				}
			} else {
				const want = "checks-timeout: current-head checks not settled within delivery.check_timeout 1m0s; pending: .github/workflows/ci-verify.yml (no run yet)"
				if stored.Stage != store.DeliveryStageParked || stored.Blocker != want || polls != 5 {
					t.Fatalf("item=%+v polls=%d", stored, polls)
				}
			}
			if log.String() != wantLog {
				t.Fatalf("log=%q want %q", log.String(), wantLog)
			}
		})
	}
}

func TestAFailedWorkflowRunWithoutAFailingCheckParks(t *testing.T) {
	t.Parallel()
	for _, conclusion := range []string{"failure", "timed_out", "cancelled", "action_required", "startup_failure", "stale", "neutral", "skipped"} {
		t.Run(conclusion, func(t *testing.T) {
			engine, runStore, item, _, boundary := newCheckWaitFixture(t, time.Minute)
			var log bytes.Buffer
			engine.log = &log
			boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
				report := recoveryReport("pass")
				report.WorkflowRuns = []WorkflowRun{{ID: "7", Name: "CI | Verify", Status: "completed", Conclusion: conclusion, URL: "https://github.com/example/repo/actions/runs/7"}}
				return report, nil
			}
			if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
				t.Fatal(err)
			}
			stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
			if conclusion == "neutral" || conclusion == "skipped" {
				if stored.Stage != store.DeliveryStageMerging || log.String() != "roundfix: checks: Delivery Queue item example: 2 check(s) completed for candidate\n" {
					t.Fatalf("item=%+v log=%q", stored, log.String())
				}
			} else if want := "checks-failed: CI | Verify (" + conclusion + ") https://github.com/example/repo/actions/runs/7"; stored.Stage != store.DeliveryStageParked || stored.Blocker != want {
				t.Fatalf("item=%+v want %q", stored, want)
			}
		})
	}
}

func TestTheChecksTimeoutNamesWhatItWaitedFor(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		report CheckReport
		detail string
	}{
		{"no check reported", CheckReport{HeadSHA: "candidate"}, "; pending: no check reported"},
		{"merge state", CheckReport{HeadSHA: "candidate", MergeState: "BLOCKED", Checks: []PullRequestCheck{{Name: "verify", Bucket: "pass"}}}, "; merge state BLOCKED"},
		{"pending check", recoveryReport("pending"), "; pending: CI / verify"},
		{"unknown mergeability", CheckReport{HeadSHA: "candidate", MergeState: "CLEAN", Mergeable: "UNKNOWN", Checks: []PullRequestCheck{{Name: "verify", Bucket: "pass"}}}, "; merge state UNKNOWN"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			engine, runStore, item, _, boundary := newCheckWaitFixture(t, time.Minute)
			var log bytes.Buffer
			engine.log = &log
			boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) { return tc.report, nil }
			if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
				t.Fatal(err)
			}
			stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
			want := "checks-timeout: current-head checks not settled within delivery.check_timeout 1m0s" + tc.detail
			if stored.Stage != store.DeliveryStageParked || stored.Blocker != want {
				t.Fatalf("item=%+v want %q", stored, want)
			}
		})
	}
}

func TestAReadErrorIsLoggedOnce(t *testing.T) {
	t.Parallel()
	engine, runStore, item, _, boundary := newCheckWaitFixture(t, time.Minute)
	var log bytes.Buffer
	engine.log = &log
	polls := 0
	boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
		polls++
		if polls <= 2 {
			return CheckReport{}, errors.New("Actions unavailable")
		}
		return recoveryReport("pass"), nil
	}
	if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
	stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
	want := "roundfix: checks: Delivery Queue item example: read failed: Actions unavailable\nroundfix: checks: Delivery Queue item example: 1 check(s) completed for candidate\n"
	if stored.Stage != store.DeliveryStageMerging || polls != 3 || strings.Count(log.String(), "read failed:") != 1 || log.String() != want {
		t.Fatalf("item=%+v polls=%d log=%q", stored, polls, log.String())
	}
}

func TestEveryFailedHeadMemberIsNamed(t *testing.T) {
	t.Parallel()
	engine, runStore, item, _, boundary := newCheckWaitFixture(t, time.Minute)
	var log bytes.Buffer
	engine.log = &log
	boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
		return CheckReport{HeadSHA: "candidate", Checks: []PullRequestCheck{
			{Name: "build", Workflow: "CI", Bucket: "fail", State: "FAILURE", Link: "https://github.com/example/repo/actions/runs/1/job/10"},
			{Name: "lint", Workflow: "CI", Bucket: "cancelled", Link: "https://github.com/example/repo/actions/runs/1/job/11"},
			{Name: "coverage", Bucket: "cancel"},
			{Name: "test", Bucket: "pending", Link: "https://github.com/example/repo/actions/runs/2/job/20"},
		}, WorkflowRuns: []WorkflowRun{
			{ID: "1", Name: "CI", Status: "completed", Conclusion: "failure", URL: "https://github.com/example/repo/actions/runs/1"},
			{ID: "2", Path: ".github/workflows/test.yml", Status: "completed", Conclusion: "timed_out", URL: "https://github.com/example/repo/actions/runs/2"},
		}}, nil
	}
	if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
	stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
	const want = "checks-failed: CI / build (failure) https://github.com/example/repo/actions/runs/1/job/10; CI / lint (cancelled) https://github.com/example/repo/actions/runs/1/job/11; coverage (cancel); .github/workflows/test.yml (timed_out) https://github.com/example/repo/actions/runs/2"
	if stored.Stage != store.DeliveryStageParked || stored.Blocker != want {
		t.Fatalf("item=%+v want %q", stored, want)
	}
}

func TestAWorkflowRunWaitsAfterItsCheckPasses(t *testing.T) {
	t.Parallel()
	for _, conclusion := range []string{"success", "failure"} {
		t.Run(conclusion, func(t *testing.T) {
			engine, runStore, item, _, boundary := newCheckWaitFixture(t, time.Minute)
			var log bytes.Buffer
			engine.log = &log
			polls := 0
			boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
				polls++
				run := WorkflowRun{ID: "7", Path: ".github/workflows/verify.yml", Status: "in_progress", URL: "https://github.com/example/repo/actions/runs/7"}
				if polls >= 2 {
					run.Status, run.Conclusion = "completed", conclusion
				}
				return CheckReport{HeadSHA: "candidate", MergeState: "CLEAN", Checks: []PullRequestCheck{{Name: "test", Bucket: "pass", Link: "https://github.com/example/repo/actions/runs/7/job/70"}}, WorkflowRuns: []WorkflowRun{run}}, nil
			}
			if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
				t.Fatal(err)
			}
			stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
			if polls != 2 {
				t.Fatalf("item=%+v polls=%d", stored, polls)
			}
			if conclusion == "success" {
				if stored.Stage != store.DeliveryStageMerging || !strings.HasSuffix(log.String(), "1 check(s) completed for candidate\n") {
					t.Fatalf("item=%+v log=%q", stored, log.String())
				}
			} else if stored.Stage != store.DeliveryStageParked || stored.Blocker != "checks-failed: .github/workflows/verify.yml (failure) https://github.com/example/repo/actions/runs/7" {
				t.Fatalf("item=%+v", stored)
			}
		})
	}
}

func TestARequestedRerunWaitsForANewerWorkflowAttempt(t *testing.T) {
	t.Parallel()
	engine, runStore, item, _, boundary := newCheckWaitFixture(t, time.Minute)
	var log bytes.Buffer
	engine.log = &log
	recovery := &fakeCheckRecovery{failures: []CheckFailure{outsideFailure(1)}}
	engine.checks = recovery
	polls := 0
	boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
		polls++
		report := recoveryReport("pass")
		if polls == 1 {
			report = recoveryReport("fail")
		}
		run := WorkflowRun{ID: "42", Name: "CI", Status: "completed", Conclusion: "failure", Attempt: 1, URL: "https://github.com/example/repo/actions/runs/42"}
		if polls >= 3 {
			run.Attempt = 2
		}
		report.WorkflowRuns = []WorkflowRun{run}
		return report, nil
	}
	if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
	stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
	if polls != 3 || recovery.reruns != 1 || stored.Stage != store.DeliveryStageParked || stored.Blocker != "checks-failed: CI (failure) https://github.com/example/repo/actions/runs/42" || !strings.Contains(log.String(), "waiting for 1 check(s): CI (re-run requested)\n") {
		t.Fatalf("item=%+v polls=%d reruns=%d log=%q", stored, polls, recovery.reruns, log.String())
	}
}

func TestCheckReadErrorsKeepTheLastSuccessfulWaitAndResetTheirLog(t *testing.T) {
	t.Parallel()
	engine, runStore, item, _, boundary := newCheckWaitFixture(t, time.Minute)
	var log bytes.Buffer
	engine.log = &log
	polls := 0
	boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
		polls++
		if polls == 2 || polls >= 4 {
			return CheckReport{}, errors.New("Actions unavailable")
		}
		return recoveryReport("pending"), nil
	}
	if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
	stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
	const want = "checks-timeout: current-head checks not settled within delivery.check_timeout 1m0s; pending: CI / verify"
	if stored.Stage != store.DeliveryStageParked || stored.Blocker != want || polls != 5 || strings.Count(log.String(), "read failed: Actions unavailable\n") != 2 || strings.Count(log.String(), "waiting for 1 check(s): CI / verify\n") != 1 {
		t.Fatalf("item=%+v polls=%d log=%q", stored, polls, log.String())
	}
}
