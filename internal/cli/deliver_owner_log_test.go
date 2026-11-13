// Suite: Queue Owner logging and owner-record lifetime.
// Boundary: production workflow wiring, a real capacity lease and queue store
// in temporary homes; the owner engine is injected and no provider is called.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/store"
)

// Release on the wait message itself so the gate wakes on the lease's release,
// rather than a timer or a background test goroutine.
type machineCapacityReleaseLog struct {
	bytes.Buffer
	lease *store.MachineVerificationLease
}

func (log *machineCapacityReleaseLog) Write(data []byte) (int, error) {
	n, err := log.Buffer.Write(data)
	if log.lease != nil {
		err = errors.Join(err, log.lease.Release())
		log.lease = nil
	}
	return n, err
}

func TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate(t *testing.T) {
	t.Parallel()
	f := newDeliveryItemFixture(t)
	setCommandEnvironmentForTest(t, f.home, f.repo)
	mustWrite(t, filepath.Join(f.repo, ".roundfixrc.yml"), "delivery:\n  linux_lane: false\n")
	gittest.Run(t, f.repo, "add", ".roundfixrc.yml")
	gittest.Run(t, f.repo, "commit", "-m", "disable lane")
	f.head = strings.TrimSpace(gittest.Run(t, f.repo, "rev-parse", "HEAD"))
	gittest.Run(t, f.repo, "update-ref", "refs/remotes/origin/main", f.head)
	f.workflow.loaded.UserConfigPath = filepath.Join(f.home, ".roundfix", "config.yml")
	f.workflow.loaded.Config.Runs.MaxVerifying = 1
	f.workflow.loaded.Config.Defaults.Verification = "exit 0"
	ctx := commandContextForTest(t, t.Context())
	holder := store.MachineVerificationHolder{
		Kind: store.MachineVerificationHolderRun, Repository: f.repo,
		RunID: "run_20261112T101500Z_0123456789abcdef", SpecSlug: "0096-other",
		Step: "Task task_qa Verification attempt 1",
	}
	lease, err := store.AcquireMachineVerification(ctx, f.home, holder, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := lease.Release(); err != nil {
			t.Error(err)
		}
	})
	log := &machineCapacityReleaseLog{lease: lease}
	dependencies := commandDeliveryEngineDependencies(f.workflow.store, f.workflow.loaded, log)
	workflow, ok := dependencies.Gate.(*commandDeliveryWorkflow)
	if !ok {
		t.Fatalf("gate=%T, want *commandDeliveryWorkflow", dependencies.Gate)
	}
	if workflow.log != log || dependencies.Log != log {
		t.Fatalf("workflow type=%T engine log=%T; expected a command workflow with the owner's log", dependencies.Gate, dependencies.Log)
	}
	result, err := dependencies.Gate.Gate(ctx, f.repo, "widget", f.head)
	if err != nil || !result.Passed || result.LaneSkip != "" {
		t.Fatalf("gate=%+v err=%v", result, err)
	}
	want := fmt.Sprintf("roundfix: Delivery Queue item widget: waiting for the Machine Verification Capacity; 1 of 1 held by Run %s of 0096-other in %s, Task task_qa Verification attempt 1\n", holder.RunID, f.repo)
	if log.String() != want {
		t.Fatalf("log=%q want=%q", log.String(), want)
	}
	view, err := store.ReadMachineVerification(ctx, f.home)
	if err != nil || len(view.Entries) != 0 || len(view.Unreadable) != 0 {
		t.Fatalf("remaining capacity=%+v err=%v", view, err)
	}
}

func TestADeliveryWorkflowWithoutALogWritesToStandardError(t *testing.T) {
	t.Parallel()
	workflow := &commandDeliveryWorkflow{}
	if got := workflow.console(); got != os.Stderr {
		t.Fatalf("console=%v, want os.Stderr", got)
	}
	var log bytes.Buffer
	workflow.log = &log
	if got := workflow.console(); got != &log {
		t.Fatalf("console=%v, want the supplied log", got)
	}
}

func TestAQueueOwnerThatStopsWithAnItemLeftKeepsItsOwnerRecord(t *testing.T) {
	t.Parallel()
	for _, outcome := range []string{"panic", "error", "idle"} {
		t.Run(outcome, func(t *testing.T) {
			t.Parallel()
			home, repo := newDeliveryOwnerStatusQueue(t, 0, "", store.DeliveryStageGating)
			engineErr := errors.New("owner engine stopped")
			var claimedIdentity string
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
				dependencies.newDeliveryEngine = func(db *store.Store, _ roundconfig.Loaded) deliveryEngine {
					return &retryCommandDeliveryEngine{run: func(ctx context.Context, gitRoot string) (delivery.EngineResult, error) {
						queue, found, err := db.DeliveryQueue(ctx, gitRoot)
						if err != nil || !found || queue.OwnerPID != os.Getpid() || queue.OwnerIdentity == "" {
							t.Fatalf("claimed queue=%+v found=%v err=%v", queue, found, err)
						}
						claimedIdentity = queue.OwnerIdentity
						switch outcome {
						case "panic":
							fmt.Fprintln(io.Writer(nil), "owner workflow message")
						case "error":
							return delivery.EngineResult{}, engineErr
						case "idle":
							item := queue.Items[0]
							item.Stage = store.DeliveryStageMerged
							if err := db.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
								t.Fatal(err)
							}
						}
						return delivery.EngineResult{}, nil
					}}
				}
			})
			var stderr bytes.Buffer
			code := -1
			var caught any
			func() {
				defer func() { caught = recover() }()
				code = runDeliveryOwner(commandContextForTest(t, t.Context()), nil, io.Discard, &stderr, &detachChild{}, commandEnvironmentForTest(t))
			}()
			if claimedIdentity == "" {
				t.Fatalf("engine was not called: exit=%d stderr=%q panic=%v", code, stderr.String(), caught)
			}
			if outcome == "panic" {
				if caught == nil {
					t.Fatal("expected the engine's nil-writer panic")
				}
			} else if caught != nil {
				t.Fatalf("unexpected panic: %v", caught)
			}
			queue := openDeliveryQueueForCLI(t, home, repo)
			if outcome == "idle" {
				if code != exitOK || stderr.Len() != 0 || queue.OwnerPID != 0 || queue.OwnerIdentity != "" || queue.Items[0].Stage != store.DeliveryStageMerged {
					t.Fatalf("exit=%d stderr=%q queue=%+v", code, stderr.String(), queue)
				}
			} else {
				if outcome == "error" && (code != exitPreflight || !strings.Contains(stderr.String(), engineErr.Error())) {
					t.Fatalf("engine error: exit=%d stderr=%q", code, stderr.String())
				}
				if queue.OwnerPID != os.Getpid() || queue.OwnerIdentity != claimedIdentity || queue.Items[0].Stage != store.DeliveryStageGating {
					t.Fatalf("owner stopped with an item left: PID=%d identity=%q stage=%s; want PID=%d identity=%q stage=gating", queue.OwnerPID, queue.OwnerIdentity, queue.Items[0].Stage, os.Getpid(), claimedIdentity)
				}
			}
		})
	}
}
