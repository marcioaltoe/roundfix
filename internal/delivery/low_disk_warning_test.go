// Boundary: real queue persistence and stage transitions; injected capacity, workflow and forge.
package delivery

import (
	"bytes"
	"context"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

type freeSpaceCheckerFunc func(context.Context) []string

func (check freeSpaceCheckerFunc) LowFreeSpace(ctx context.Context) []string { return check(ctx) }

type freeSpaceTestWorkflow struct {
	*fakeDeliveryWorkflow
	beforeRun, beforeGate func()
}

func (workflow freeSpaceTestWorkflow) RunSpec(ctx context.Context, workDir, slug string) (RunResult, error) {
	workflow.beforeRun()
	return workflow.fakeDeliveryWorkflow.RunSpec(ctx, workDir, slug)
}

func (workflow freeSpaceTestWorkflow) Gate(ctx context.Context, workDir, slug, head string) (GateResult, error) {
	workflow.beforeGate()
	return workflow.fakeDeliveryWorkflow.Gate(ctx, workDir, slug, head)
}

func TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate(t *testing.T) {
	t.Parallel()
	const premise = "premise-changed: internal/example.go"
	const runLocation = "worktree location /worktrees has 10239 MiB free, below 10240 MiB"
	const gateLocation = "temporary directory /tmp has 7 MiB free, below 10240 MiB"
	const runWarning = "low-disk: " + runLocation
	const gateWarning = "low-disk: " + gateLocation + ", " + runLocation
	const laneReason = "image not present; with detail"
	const lane = "linux-lane-skipped: " + laneReason
	for _, tc := range []struct {
		name, before, atRun, atGate, laneSkip string
		locations                             [][]string
		noChecker                             bool
	}{
		{
			name: "low before the Run and the gate", before: premise + "; low-disk: old location",
			locations: [][]string{{runLocation}, {gateLocation, runLocation}},
			atRun:     premise + "; " + runWarning, atGate: premise + "; " + gateWarning,
		},
		{
			name: "healthy at the gate", locations: [][]string{{runLocation}, nil},
			atRun: runWarning,
		},
		{
			name: "low with a lane-skip warning present", before: premise + "; low-disk: old location; " + lane,
			locations: [][]string{{runLocation}, {gateLocation, runLocation}}, laneSkip: laneReason,
			atRun: premise + "; " + runWarning + "; " + lane, atGate: premise + "; " + gateWarning + "; " + lane,
		},
		{
			name: "healthy removes only low-disk", before: premise + "; low-disk: old location; " + lane,
			locations: [][]string{nil, nil}, laneSkip: laneReason,
			atRun: premise + "; " + lane, atGate: premise + "; " + lane,
		},
		{
			name: "unchanged low-disk", before: runWarning,
			locations: [][]string{{runLocation}, {runLocation}}, atRun: runWarning, atGate: runWarning,
		},
		{
			name: "no checker", before: premise + "; low-disk: old location; " + lane, noChecker: true, laneSkip: laneReason,
			atRun: premise + "; low-disk: old location; " + lane, atGate: premise + "; low-disk: old location; " + lane,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			runStore := openDeliveryEngineStore(t, ctx)
			root := filepath.Join(t.TempDir(), "repo")
			queue, err := runStore.CreateDeliveryQueue(ctx, root, []string{"spec-one"})
			if err != nil {
				t.Fatal(err)
			}
			item := queue.Items[0]
			item.Stage = store.DeliveryStageRunning
			item.Branch = "roundfix/deliver-spec-one"
			item.Worktree = filepath.Join(t.TempDir(), "item")
			item.Warning = tc.before
			recordDeliveryItemWorkspace(t, ctx, runStore, root, &item)
			if err := runStore.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
				t.Fatal(err)
			}
			var log bytes.Buffer
			checks, runs, gates, publications := 0, 0, 0, 0
			assertStored := func(want string, stage store.DeliveryStage) {
				t.Helper()
				got := readDeliveryQueue(t, ctx, runStore, root).Items[0]
				if got.Warning != want || got.Stage != stage {
					t.Fatalf("stored warning=%q stage=%q want warning=%q stage=%q", got.Warning, got.Stage, want, stage)
				}
			}
			fake := newFakeDeliveryWorkflow()
			fake.gates[item.SpecSlug] = GateResult{Passed: true, LaneSkip: tc.laneSkip}
			workflow := freeSpaceTestWorkflow{
				fakeDeliveryWorkflow: fake,
				beforeRun: func() {
					runs++
					assertStored(tc.atRun, store.DeliveryStageRunning)
					if !tc.noChecker && checks != 1 {
						t.Fatalf("checker calls before Run=%d want=1", checks)
					}
				},
				beforeGate: func() {
					gates++
					assertStored(tc.atGate, store.DeliveryStageGating)
					if !tc.noChecker && checks != 2 {
						t.Fatalf("checker calls before gate=%d want=2", checks)
					}
				},
			}
			var checker FreeSpaceChecker
			if !tc.noChecker {
				checker = freeSpaceCheckerFunc(func(callCtx context.Context) []string {
					if callCtx != ctx {
						t.Fatal("checker lost the queue context")
					}
					if checks >= len(tc.locations) {
						t.Fatalf("unexpected checker call %d", checks+1)
					}
					result := tc.locations[checks]
					checks++
					return result
				})
			}
			body := "deliver spec-one"
			if tc.laneSkip != "" {
				body += "\n\nLinux Lane: not run. The Delivery Queue gate skipped it: image not present; with detail. This candidate's Go packages were not tested on Linux before this pull request."
			}
			boundary := newFakeDeliveryBoundary()
			boundary.findPullRequest = func(req PullRequestRequest) PullRequestResult {
				publications++
				if req.Body != body {
					t.Fatalf("PR body=%q want=%q", req.Body, body)
				}
				assertStored(tc.atGate, store.DeliveryStagePublishing)
				pr := PullRequest{Number: "1", State: "OPEN", HeadBranch: req.HeadBranch, HeadSHA: boundary.remoteHeads[req.HeadBranch]}
				boundary.pullRequests[req.HeadBranch] = pr
				return PullRequestResult{PullRequest: pr, Created: true}
			}
			engine := NewEngine(runStore, EngineDependencies{
				Workspace: workflow, Runner: workflow, Reviewer: workflow, Archiver: workflow,
				Gate: workflow, Authorizer: workflow, Publication: workflow,
				PullRequests: boundary, Revalidator: cleanTestRevalidator{}, FreeSpace: checker, Log: &log,
			})
			if _, err := engine.Run(ctx, root); err != nil {
				t.Fatal(err)
			}
			assertStored(tc.atGate, store.DeliveryStageMerged)
			if runs != 1 || gates != 1 || publications != 1 {
				t.Fatalf("Run calls=%d gate calls=%d PR calls=%d want=1 each", runs, gates, publications)
			}
			var gotLogs, wantLogs []string
			for _, line := range strings.Split(log.String(), "\n") {
				if strings.Contains(line, ": low-disk:") {
					gotLogs = append(gotLogs, line)
				}
			}
			for i, locations := range tc.locations {
				if len(locations) > 0 {
					phase := []string{"Implement Run", "repository gate"}[i]
					wantLogs = append(wantLogs, "roundfix: warning: Delivery Queue item spec-one: low-disk: "+strings.Join(locations, ", ")+" before its "+phase+"; the queue continues")
				}
			}
			if !reflect.DeepEqual(gotLogs, wantLogs) {
				t.Fatalf("warning log=%q want=%q", gotLogs, wantLogs)
			}
		})
	}
}
