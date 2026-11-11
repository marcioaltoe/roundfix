// Suite: Machine Verification Capacity entries and leases.
// Invariant: ticket order and proven owner death alone control admission.
// Boundary: temporary homes and read-only process identity checks; no signals.
package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/testwait"
)

func TestMachineVerificationCapacityServesWaitersInTicketOrder(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	a := acquireMachineTestLease(t, home, "A", 1)
	b := startMachineTestAcquisition(t, home, "B", 1)
	bw := testwait.Until(t, "B waiting", b.waits, b.done)
	if len(bw.Holding) != 1 || len(bw.Ahead) != 0 {
		t.Fatalf("B wait = %+v", bw)
	}
	c := startMachineTestAcquisition(t, home, "C", 1)
	cw := testwait.Until(t, "C waiting", c.waits, c.done)
	if len(cw.Holding) != 1 || len(cw.Ahead) != 1 || cw.Ahead[0].Holder.RunID != "B" {
		t.Fatalf("C wait = %+v", cw)
	}
	releaseMachineTestLease(t, a)
	bl := b.granted(t)
	if bl.Waited() <= 0 {
		t.Error("B reports no wait")
	}
	view := readMachineTestView(t, home)
	if len(view.Entries) != 2 || view.Entries[0].Holder.RunID != "B" || view.Entries[0].State != "holding" || view.Entries[1].Holder.RunID != "C" || view.Entries[1].State != "waiting" {
		t.Fatalf("after B grant: %+v", view)
	}
	releaseMachineTestLease(t, bl)
	cl := c.granted(t)
	releaseMachineTestLease(t, cl)
	if view := readMachineTestView(t, home); len(view.Entries) != 0 {
		t.Fatalf("leases remain: %+v", view)
	}
}

func TestMachineVerificationCapacityKeepsALaterWaiterBehindALiveWaiterAhead(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	ahead := writeMachineTestEntry(t, home, 1, "waiting", os.Getpid(), "")
	later := startMachineTestAcquisition(t, home, "later", 1)
	wait := testwait.Until(t, "later waiter behind live ticket", later.waits, later.done)
	if len(wait.Holding) != 0 || len(wait.Ahead) != 1 || wait.Ahead[0].Path != ahead {
		t.Fatalf("wait = %+v", wait)
	}
	later.cancel()
	testwait.Until[struct{}, struct{}](t, "cancelled later waiter", later.done, nil)
	if later.lease != nil || !errors.Is(later.err, context.Canceled) {
		t.Fatalf("acquisition granted or not cancelled: %v, %v", later.lease, later.err)
	}
	if _, err := os.Stat(ahead); err != nil {
		t.Fatalf("live waiter removed: %v", err)
	}
}

func TestMachineVerificationCapacityNeverHoldsMoreThanItsUnits(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	a := acquireMachineTestLease(t, home, "A", 2)
	b := acquireMachineTestLease(t, home, "B", 2)
	c := startMachineTestAcquisition(t, home, "C", 2)
	wait := testwait.Until(t, "third acquisition waiting", c.waits, c.done)
	if wait.Capacity != 2 || len(wait.Holding) != 2 || len(wait.Ahead) != 0 {
		t.Fatalf("third wait = %+v", wait)
	}
	releaseMachineTestLease(t, a)
	cl := c.granted(t)
	view := readMachineTestView(t, home)
	if len(view.Entries) != 2 {
		t.Fatalf("expected two holders: %+v", view)
	}
	for _, entry := range view.Entries {
		if entry.State != "holding" {
			t.Fatalf("not holding: %+v", entry)
		}
	}
	releaseMachineTestLease(t, b)
	releaseMachineTestLease(t, cl)
}

func TestMachineVerificationCapacityOffWritesNothing(t *testing.T) {
	t.Parallel()
	for _, capacity := range []int{0, -1} {
		t.Run(fmt.Sprint(capacity), func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			lease, err := AcquireMachineVerification(t.Context(), home, machineTestHolder("off"), capacity, func(MachineVerificationWait) { t.Error("off acquisition waited") })
			if err != nil {
				t.Fatal(err)
			}
			if lease.Waited() != 0 {
				t.Error("off lease reports a wait")
			}
			releaseMachineTestLease(t, lease)
			releaseMachineTestLease(t, lease)
			files, err := os.ReadDir(home)
			if err != nil || len(files) != 0 {
				t.Fatalf("off capacity wrote files: %v, %v", files, err)
			}
		})
	}
}

func TestACancelledMachineVerificationWaitLeavesNoEntry(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	holder, err := AcquireMachineVerification(ctx, home, machineTestHolder("holder"), 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { releaseMachineTestLease(t, holder) })
	waiter := startMachineTestAcquisition(t, home, "waiter", 1)
	testwait.Until(t, "waiter queued", waiter.waits, waiter.done)
	waiter.cancel()
	testwait.Until[struct{}, struct{}](t, "waiter cancellation cleanup", waiter.done, nil)
	if !errors.Is(waiter.err, context.Canceled) || !strings.HasPrefix(waiter.err.Error(), "wait for the Machine Verification Capacity") {
		t.Fatalf("cancel error = %v", waiter.err)
	}
	view := readMachineTestView(t, home)
	if len(view.Entries) != 1 || view.Entries[0].Holder.RunID != "holder" {
		t.Fatalf("cancelled waiter remains: %+v", view)
	}
	cancel() // Release must work even after the holder's acquisition context ends.
	releaseMachineTestLease(t, holder)
	releaseMachineTestLease(t, holder)
	files, err := filepath.Glob(filepath.Join(machineVerificationDirectory(home), "*.json"))
	if err != nil || len(files) != 0 {
		t.Fatalf("entry remains after release: %v, %v", files, err)
	}
}

func TestReadMachineVerificationLeavesOutAnOwnerProvenGoneAndWritesNothing(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	// This positive PID is outside the host's allocated process range; check the
	// same public liveness boundary before using it as the absent-owner fixture.
	const absentPID = 1 << 30
	if gone, _ := OwnerGone(t.Context(), absentPID, ""); !gone {
		t.Fatal("absent PID fixture is not proven absent")
	}
	dead := writeMachineTestEntry(t, home, 1, "holding", absentPID, "")
	live := writeMachineTestEntry(t, home, 2, "waiting", os.Getpid(), "")
	dir := machineVerificationDirectory(home)
	broken := filepath.Join(dir, "broken.json")
	if err := os.WriteFile(broken, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	before := machineTestDirectorySnapshot(t, dir)
	view := readMachineTestView(t, home)
	after := machineTestDirectorySnapshot(t, dir)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("reader wrote files: before=%v after=%v", before, after)
	}
	if len(view.Entries) != 1 || view.Entries[0].Path != live || len(view.Unreadable) != 1 || view.Unreadable[0] != broken {
		t.Fatalf("reader view = %+v; dead = %s", view, dead)
	}
	missing := t.TempDir()
	if view := readMachineTestView(t, missing); len(view.Entries) != 0 || len(view.Unreadable) != 0 {
		t.Fatalf("missing directory view = %+v", view)
	}
	if files, err := os.ReadDir(missing); err != nil || len(files) != 0 {
		t.Fatalf("reader created directory: %v, %v", files, err)
	}
}

func TestMachineVerificationSkipsUnreadableEntriesWithoutRemovingOrCountingThem(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	dir := machineVerificationDirectory(home)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{"broken.json": "{", "foreign.json": `{"schema":"foreign"}`} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	before := machineTestDirectorySnapshot(t, dir)
	lease := acquireMachineTestLease(t, home, "holder", 1)
	view := readMachineTestView(t, home)
	if len(view.Entries) != 1 || len(view.Unreadable) != 2 {
		t.Fatalf("view = %+v", view)
	}
	releaseMachineTestLease(t, lease)
	after := machineTestDirectorySnapshot(t, dir)
	for name, stamp := range before {
		if after[name] != stamp {
			t.Errorf("unreadable entry %s changed", name)
		}
	}
}

func TestDescribeMachineVerificationHolder(t *testing.T) {
	t.Parallel()
	holder := machineTestHolder("run_1")
	if got, want := DescribeMachineVerificationHolder(holder), "Run run_1 of test-spec in test-repository"; got != want {
		t.Fatalf("description = %q, want %q", got, want)
	}
	holder.Kind = MachineVerificationHolderDelivery
	holder.RunID = ""
	if got, want := DescribeMachineVerificationHolder(holder), "Delivery Queue item test-spec in test-repository"; got != want {
		t.Fatalf("description = %q, want %q", got, want)
	}
}

func machineTestHolder(id string) MachineVerificationHolder {
	return MachineVerificationHolder{Kind: MachineVerificationHolderRun, RunID: id, Repository: "test-repository", SpecSlug: "test-spec", Step: "test step"}
}

func acquireMachineTestLease(t *testing.T, home, id string, capacity int) *MachineVerificationLease {
	t.Helper()
	lease, err := AcquireMachineVerification(t.Context(), home, machineTestHolder(id), capacity, func(MachineVerificationWait) { t.Errorf("immediate acquisition %s waited", id) })
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { releaseMachineTestLease(t, lease) })
	if lease.Waited() != 0 {
		t.Errorf("immediate acquisition %s reports a wait", id)
	}
	return lease
}

func releaseMachineTestLease(t *testing.T, lease *MachineVerificationLease) {
	t.Helper()
	if err := lease.Release(); err != nil {
		t.Fatal(err)
	}
}

type machineTestAcquisition struct {
	waits  chan MachineVerificationWait
	done   chan struct{}
	cancel context.CancelFunc
	lease  *MachineVerificationLease
	err    error
}

func startMachineTestAcquisition(t *testing.T, home, id string, capacity int) *machineTestAcquisition {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	result := &machineTestAcquisition{waits: make(chan MachineVerificationWait, 1), done: make(chan struct{}), cancel: cancel}
	go func() {
		defer close(result.done)
		result.lease, result.err = AcquireMachineVerification(ctx, home, machineTestHolder(id), capacity, func(wait MachineVerificationWait) { result.waits <- wait })
	}()
	t.Cleanup(func() {
		cancel()
		testwait.Until[struct{}, struct{}](t, "acquisition exit", result.done, nil)
		if result.lease != nil {
			releaseMachineTestLease(t, result.lease)
		}
	})
	return result
}

func (result *machineTestAcquisition) granted(t *testing.T) *MachineVerificationLease {
	t.Helper()
	testwait.Until[struct{}, struct{}](t, "capacity grant", result.done, nil)
	if result.err != nil || result.lease == nil {
		t.Fatalf("grant: lease=%v err=%v", result.lease, result.err)
	}
	return result.lease
}

func writeMachineTestEntry(t *testing.T, home string, ticket int64, state string, pid int, identity string) string {
	t.Helper()
	dir := machineVerificationDirectory(home)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, fmt.Sprintf("%020d.json", ticket))
	// Write the public wire format independently of the production serializer.
	acquired := ""
	if state == "holding" {
		acquired = `,"acquired_at":"2026-11-10T22:41:30Z"`
	}
	data := fmt.Sprintf(`{"schema":%q,"ticket":%d,"state":%q,"kind":"run","run_id":"fixture","repository":"fixture-repository","spec":"fixture-spec","step":"fixture-step","owner_pid":%d,"owner_identity":%q,"since":"2026-11-10T22:41:10Z"%s}`, machineVerificationSchema, ticket, state, pid, identity, acquired)
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sequence"), []byte(fmt.Sprint(ticket)), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func readMachineTestView(t *testing.T, home string) MachineVerificationView {
	t.Helper()
	view, err := ReadMachineVerification(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	return view
}

type machineTestFileStamp struct {
	Size     int64
	Modified time.Time
}

func machineTestDirectorySnapshot(t *testing.T, dir string) map[string]machineTestFileStamp {
	t.Helper()
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	result := make(map[string]machineTestFileStamp)
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			t.Fatal(err)
		}
		result[file.Name()] = machineTestFileStamp{info.Size(), info.ModTime()}
	}
	return result
}

func TestMachineVerificationCapacityBoundsEightConcurrentAcquisitions(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	type outcome struct {
		lease *MachineVerificationLease
		err   error
	}
	grants := make(chan outcome, 8)
	var workers sync.WaitGroup
	for i := range 8 {
		workers.Go(func() {
			lease, err := AcquireMachineVerification(ctx, home, machineTestHolder(fmt.Sprint(i)), 2, nil)
			grants <- outcome{lease, err}
		})
	}
	done := make(chan struct{})
	go func() { workers.Wait(); close(done) }()
	t.Cleanup(func() {
		cancel()
		testwait.Until[struct{}, struct{}](t, "all concurrent acquisitions exited", done, nil)
		for {
			select {
			case grant := <-grants:
				if grant.lease != nil {
					releaseMachineTestLease(t, grant.lease)
				}
			default:
				return
			}
		}
	})
	for range 8 {
		grant := testwait.Until[outcome, struct{}](t, "concurrent capacity grant", grants, nil)
		if grant.err != nil || grant.lease == nil {
			t.Fatalf("concurrent grant: %v, %v", grant.lease, grant.err)
		}
		// Register cleanup before assertions so even a broken capacity oracle leaves
		// no live leases or acquisition goroutines behind.
		t.Cleanup(func() { releaseMachineTestLease(t, grant.lease) })
		holding := 0
		for _, entry := range readMachineTestView(t, home).Entries {
			if entry.State == "holding" {
				holding++
			}
		}
		if holding > 2 {
			t.Fatalf("capacity 2 has %d holders", holding)
		}
		releaseMachineTestLease(t, grant.lease)
	}
	testwait.Until[struct{}, struct{}](t, "eight acquisitions ended", done, nil)
	if view := readMachineTestView(t, home); len(view.Entries) != 0 {
		t.Fatalf("entries remain after eight grants: %+v", view)
	}
}
