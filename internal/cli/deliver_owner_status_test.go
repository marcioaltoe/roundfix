// Suite: Delivery Queue owner status.
// Invariant: status reports proven owner absence without changing the queue.
// Boundary IN: public CLI dispatch and the real queue store in temporary homes.
// Boundary OUT: process identity proof, injected through the owner controller.
package cli

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/store"
)

func TestDeliverStatusReproducesSurfaceTranscriptOwnerNotRunning(t *testing.T) {
	t.Parallel()
	const pid = 1 << 30
	newDeliveryOwnerStatusQueue(t, pid, "recorded-owner-identity", store.DeliveryStageChecking)
	// Even an unavailable identity controller must not hide a proven-absent PID.
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.ownerProcesses = nil
	})
	assertDeliveryOwnerStatus(t,
		fmt.Sprintf("Owner: Delivery Queue owner PID %d is not running; run roundfix deliver resume\n", pid), "")
}

func TestDeliverStatusNamesAnOwnerWithADifferentIdentity(t *testing.T) {
	t.Parallel()
	pid := os.Getpid()
	const identity = "recorded-owner-identity"
	newDeliveryOwnerStatusQueue(t, pid, identity, store.DeliveryStageChecking)
	processes := &fakeDeliveryOwnerProcesses{proveErr: fmt.Errorf("identity mismatch: %w", store.ErrOwnerProcessIdentityUnproven)}
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.ownerProcesses = processes
	})
	assertDeliveryOwnerStatus(t,
		fmt.Sprintf("Owner: Delivery Queue owner PID %d has a different process identity; run roundfix deliver resume\n", pid), "")
	if processes.provedPID != pid || processes.identity != identity || processes.terminatedPID != 0 {
		t.Fatalf("owner proof or termination = %+v", processes)
	}
}

func TestDeliverStatusPrintsNoOwnerLineForAProvenOwner(t *testing.T) {
	t.Parallel()
	pid := os.Getpid()
	const identity = "recorded-owner-identity"
	newDeliveryOwnerStatusQueue(t, pid, identity, store.DeliveryStageChecking)
	processes := &fakeDeliveryOwnerProcesses{}
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.ownerProcesses = processes
	})
	assertDeliveryOwnerStatus(t, "", "")
	if processes.provedPID != pid || processes.identity != identity || processes.terminatedPID != 0 {
		t.Fatalf("owner proof or termination = %+v", processes)
	}
}

func TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner(t *testing.T) {
	t.Parallel()
	newDeliveryOwnerStatusQueue(t, 0, "", store.DeliveryStageChecking)
	assertDeliveryOwnerStatus(t, "", "")
}

func TestDeliverStatusPrintsNoOwnerLineWhenNoItemCanAdvance(t *testing.T) {
	t.Parallel()
	newDeliveryOwnerStatusQueue(t, 1<<30, "recorded-owner-identity", store.DeliveryStageMerged, store.DeliveryStageParked)
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 || strings.Contains(stdout.String(), "Owner:") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestDeliverStatusReportsAnUncheckableOwnerOnStderr(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name       string
		controller OwnerProcessController
		diagnostic string
	}{
		{"unreadable identity", &fakeDeliveryOwnerProcesses{proveErr: store.ErrOwnerIdentityUnreadable}, store.ErrOwnerIdentityUnreadable.Error()},
		{"other error", &fakeDeliveryOwnerProcesses{proveErr: fmt.Errorf("process table unavailable")}, "process table unavailable"},
		{"missing controller", nil, "Delivery Queue owner process controller is required"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			pid := os.Getpid()
			newDeliveryOwnerStatusQueue(t, pid, "recorded-owner-identity", store.DeliveryStageChecking)
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
				dependencies.ownerProcesses = test.controller
			})
			assertDeliveryOwnerStatus(t, "", fmt.Sprintf("roundfix: Delivery Queue owner PID %d could not be checked: %s\n", pid, test.diagnostic))
		})
	}
}

func TestDeliverStatusLeavesTheOwnerRecord(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name     string
		pid      int
		proofErr error
	}{
		{"absent", 1 << 30, nil},
		{"different identity", os.Getpid(), store.ErrOwnerProcessIdentityUnproven},
		{"proven", os.Getpid(), nil},
		{"unreadable", os.Getpid(), store.ErrOwnerIdentityUnreadable},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			home, repo := newDeliveryOwnerStatusQueue(t, test.pid, "unchanged-owner-identity", store.DeliveryStageChecking)
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
				dependencies.ownerProcesses = &fakeDeliveryOwnerProcesses{proveErr: test.proofErr}
			})
			before := openDeliveryQueueForCLI(t, home, repo)
			var stdout, stderr bytes.Buffer
			if code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr); code != exitOK {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			after := openDeliveryQueueForCLI(t, home, repo)
			if !reflect.DeepEqual(before, after) {
				t.Fatalf("status changed queue: before=%+v after=%+v", before, after)
			}
		})
	}
}

func TestDeliverStatusPrintsOwnerAfterWarningsAndBeforeParks(t *testing.T) {
	t.Parallel()
	const pid = 1 << 30
	home, repo := newDeliveryOwnerStatusQueue(t, pid, "recorded-owner-identity", store.DeliveryStageChecking, store.DeliveryStageParked)
	writer, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	queue, found, err := writer.DeliveryQueue(t.Context(), repo)
	if err != nil || !found {
		t.Fatalf("read queue: found=%v err=%v", found, err)
	}
	item := queue.Items[0]
	item.Warning = "premise-changed: internal/example.go"
	if err := writer.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	want := "0300-example\tchecking\t-\t-\n0301-example\tparked\tchecks-timeout\t-\n" +
		"Warning: 0300-example premise-changed: internal/example.go\n" +
		fmt.Sprintf("Owner: Delivery Queue owner PID %d is not running; run roundfix deliver resume\n", pid) +
		"Park: 0301-example environment: resolve the blocker, then run roundfix deliver retry 0301-example\n" +
		"Limits: deadline none, retries per item none, concurrency 1, tokens none\n" +
		"Usage: no Runs recorded\n" +
		"Pending question: 0301-example parked checks-timeout\n" +
		"Answer: resolve the blocker, then run roundfix deliver retry 0301-example\n"
	if code != exitOK || stderr.Len() != 0 || stdout.String() != want {
		t.Fatalf("exit=%d stdout=%q stderr=%q want=%q", code, stdout.String(), stderr.String(), want)
	}
}

func newDeliveryOwnerStatusQueue(t *testing.T, pid int, identity string, stages ...store.DeliveryStage) (string, string) {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	writer, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	var slugs []string
	for index := range stages {
		slugs = append(slugs, fmt.Sprintf("%04d-example", 300+index))
	}
	queue, err := writer.CreateDeliveryQueue(t.Context(), repo, slugs)
	if err != nil {
		t.Fatal(err)
	}
	for index, item := range queue.Items {
		item.Stage = stages[index]
		if item.Stage == store.DeliveryStageParked {
			item.Blocker = delivery.BlockerChecksTimeout
		}
		if err := writer.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
			t.Fatal(err)
		}
	}
	if pid > 0 {
		if err := writer.ClaimDeliveryQueueOwner(t.Context(), repo, pid, identity); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return home, repo
}

func assertDeliveryOwnerStatus(t *testing.T, ownerLine, diagnostic string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	want := "0300-example\tchecking\t-\t-\n" + ownerLine +
		"Limits: deadline none, retries per item none, concurrency 1, tokens none\n" +
		"Usage: no Runs recorded\n"
	if code != exitOK || stdout.String() != want || stderr.String() != diagnostic {
		t.Fatalf("exit=%d stdout=%q stderr=%q; want stdout=%q stderr=%q", code, stdout.String(), stderr.String(), want, diagnostic)
	}
}
