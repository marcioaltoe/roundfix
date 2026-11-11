// Boundary: disposable Homes, real capacity entries, CLI readers and shell Verification.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/store"
)

func TestTheRepositoryGateHoldsTheMachineCapacityWhileItRuns(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		exit int
	}{{"passing gate", 0}, {"failing gate", 7}} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := newDeliveryItemFixture(t)
			f.workflow.loaded.Config.Runs.MaxVerifying = 1
			f.workflow.loaded.Config.Delivery.LinuxLane = false
			// The lane reads the fetched default, not the owner's effective config.
			mustWrite(t, filepath.Join(f.repo, ".roundfixrc.yml"), "delivery:\n  linux_lane: false\n")
			gittest.Run(t, f.repo, "add", ".roundfixrc.yml")
			gittest.Run(t, f.repo, "commit", "-m", "disable fixture lane")
			f.head = strings.TrimSpace(gittest.Run(t, f.repo, "rev-parse", "HEAD"))
			gittest.Run(t, f.repo, "update-ref", "refs/remotes/origin/main", f.head)
			dir := filepath.Join(f.home, ".roundfix", "verification-capacity")
			record := filepath.Join(t.TempDir(), "capacity-during-gate")
			quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'" }
			f.workflow.loaded.Config.Defaults.Verification = fmt.Sprintf(`for entry in %s/*.json; do test -f "$entry" || continue; sed -n 's/.*"state":"\([^"]*\)".*/state: \1/p' "$entry"; sed -n 's/.*"kind":"\([^"]*\)".*/kind: \1/p' "$entry"; done > %s; exit %d`, quote(dir), quote(record), tc.exit)
			result, err := f.workflow.Gate(commandContextForTest(t, t.Context()), f.repo, "widget", f.head)
			if err != nil || result.Passed != (tc.exit == 0) {
				t.Fatalf("gate=%+v err=%v", result, err)
			}
			if got := mustRead(t, record); got != "state: holding\nkind: delivery\n" {
				t.Fatalf("capacity during gate=%q", got)
			}
			entries, err := filepath.Glob(filepath.Join(dir, "*.json"))
			if err != nil || len(entries) != 0 {
				t.Fatalf("remaining entries=%v err=%v", entries, err)
			}
		})
	}
}

// Write the wire format directly, so these readers never acquire or repair a lease.
func seedMachineCapacityStatus(t *testing.T, home string) (string, string) {
	t.Helper()
	dir := filepath.Join(home, ".roundfix", "verification-capacity")
	mustMkdir(t, dir)
	repoOne, repoTwo := filepath.Join(t.TempDir(), "holder"), filepath.Join(t.TempDir(), "waiter")
	identity, err := store.OwnerProcessIdentity(t.Context(), os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range []struct {
		ticket                                              int
		state, kind, run, repo, spec, step, since, acquired string
	}{
		{8, "waiting", "run", "run_20261110T221208Z_946b774b8b696f09", repoTwo, "0002-waiter", "Task task_03 Verification attempt 2", "2026-11-10T22:50:02Z", ""},
		{7, "holding", "delivery", "", repoOne, "0001-holder", "repository gate", "2026-11-10T22:41:10Z", "2026-11-10T22:41:30Z"},
	} {
		record := map[string]any{"schema": "roundfix/machine-verification-entry/v1", "ticket": row.ticket, "state": row.state, "kind": row.kind, "run_id": row.run, "repository": row.repo, "spec": row.spec, "step": row.step, "since": row.since, "owner_pid": os.Getpid(), "owner_identity": identity}
		if row.acquired != "" {
			record["acquired_at"] = row.acquired
		}
		data, err := json.Marshal(record)
		if err != nil {
			t.Fatal(err)
		}
		mustWrite(t, filepath.Join(dir, fmt.Sprintf("%020d.json", row.ticket)), string(data)+"\n")
	}
	return repoOne, repoTwo
}

func machineCapacityStatusBlock(repoOne, repoTwo string) string {
	return fmt.Sprintf("Machine Verification Capacity: 1 of 1 held, 1 waiting\n  holding: Delivery Queue item 0001-holder in %s, repository gate, since 2026-11-10T22:41:30Z\n  waiting 1: Run run_20261110T221208Z_946b774b8b696f09 of 0002-waiter in %s, Task task_03 Verification attempt 2, since 2026-11-10T22:50:02Z\n", repoOne, repoTwo)
}

// Reuse the existing read-only command oracle for SQLite's live-reader
// coordination sidecars, and compare the capacity directory strictly as well.
type machineCapacityHomeSnapshot struct {
	home     map[string]string
	capacity map[string]doctorPathSnapshot
}

func snapshotMachineCapacityHome(t *testing.T, home string) machineCapacityHomeSnapshot {
	t.Helper()
	dir := filepath.Join(home, ".roundfix", "verification-capacity")
	var capacity map[string]doctorPathSnapshot
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		capacity = snapshotDoctorPath(t, dir)
	}
	return machineCapacityHomeSnapshot{home: snapshotUsageHome(t, home), capacity: capacity}
}

func assertMachineCapacityHomeUnchanged(t *testing.T, home string, before machineCapacityHomeSnapshot) {
	t.Helper()
	if after := snapshotMachineCapacityHome(t, home); !reflect.DeepEqual(before, after) {
		t.Fatal("reader changed Home or capacity entries")
	}
}

func TestRunsListMachineCapacityTranscript(t *testing.T) {
	t.Parallel()
	home, _ := withCLIWorkspace(t)
	one, two := seedMachineCapacityStatus(t, home)
	before := snapshotMachineCapacityHome(t, home)
	assertRunsListTranscript(t, []string{"runs", "list", "--all", "--format", "plain"}, "No Runs found.\n"+machineCapacityStatusBlock(one, two), "", exitOK)
	assertMachineCapacityHomeUnchanged(t, home, before)
}

func TestRunsListNamesTheMachineCapacityHolderAndWaiters(t *testing.T) {
	t.Parallel()
	for _, database := range []bool{false, true} {
		for _, format := range []string{"plain", "table"} {
			t.Run(fmt.Sprintf("database=%t/%s", database, format), func(t *testing.T) {
				t.Parallel()
				home, _ := withCLIWorkspace(t)
				if database {
					db, err := store.Open(t.Context(), home)
					if err != nil {
						t.Fatal(err)
					}
					t.Cleanup(func() {
						if err := db.Close(); err != nil {
							t.Error(err)
						}
					})
				}
				one, two := seedMachineCapacityStatus(t, home)
				before := snapshotMachineCapacityHome(t, home)
				assertRunsListTranscript(t, []string{"runs", "list", "--format", format}, "No Runs found.\n"+machineCapacityStatusBlock(one, two), "", exitOK)
				assertMachineCapacityHomeUnchanged(t, home, before)
			})
		}
	}
}

func TestDeliverStatusNamesTheMachineCapacityHolderAndWaiters(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	db, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.CreateDeliveryQueue(t.Context(), repo, []string{"0003-queued"}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	one, two := seedMachineCapacityStatus(t, home)
	before := snapshotMachineCapacityHome(t, home)
	assertRunsListTranscript(t, []string{"deliver", "status"}, "0003-queued\tqueued\t-\t-\n"+machineCapacityStatusBlock(one, two)+"Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none\nUsage: no Runs recorded\n", "", exitOK)
	assertMachineCapacityHomeUnchanged(t, home, before)
}

func TestDoctorReportsTheMachineVerificationCapacity(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"empty", "holder and waiter", "off", "unreadable entry", "unreadable directory"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			home, repo := withCLIWorkspace(t)
			loaded := roundconfig.Loaded{HomeDir: home, GitRoot: repo, Config: roundconfig.Builtin()}
			loaded.Config.Delivery.LinuxLane = true
			checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
			withDoctorFakeLoadedAndReadiness(t, checker, loaded, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			var baseline, baselineErr bytes.Buffer
			if code := runCLI(t, []string{"doctor"}, &baseline, &baselineErr); code != exitOK || baselineErr.Len() != 0 {
				t.Fatalf("baseline exit=%d stderr=%q", code, baselineErr.String())
			}
			detail := "runs.max_verifying 1; 0 held, 0 waiting"
			status := "ok"
			switch name {
			case "holder and waiter", "off":
				one, two := seedMachineCapacityStatus(t, home)
				setting := "runs.max_verifying 1"
				if name == "off" {
					loaded.Config.Runs.MaxVerifying = 0
					setting = "off"
				}
				detail = fmt.Sprintf("%s; 1 held, 1 waiting; holding: Delivery Queue item 0001-holder in %s, repository gate; waiting 1: Run run_20261110T221208Z_946b774b8b696f09 of 0002-waiter in %s, Task task_03 Verification attempt 2", setting, one, two)
			case "unreadable entry":
				path := filepath.Join(home, ".roundfix", "verification-capacity", "broken.json")
				mustMkdir(t, filepath.Dir(path))
				mustWrite(t, path, "{")
				status = "warn"
				detail += "; unreadable entry " + path + "; next: inspect " + path + " and remove it if no Roundfix process wrote it"
			case "unreadable directory":
				path := filepath.Join(home, ".roundfix", "verification-capacity")
				mustMkdir(t, filepath.Dir(path))
				mustWrite(t, path, "not a directory")
				_, err := store.ReadMachineVerification(t.Context(), home)
				if err == nil {
					t.Fatal("wanted directory read error")
				}
				status = "warn"
				detail = "read Machine Verification Capacity: " + err.Error()
			}
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.doctor.loadConfig = func(roundconfig.LoadOptions) (roundconfig.Loaded, error) { return loaded, nil }
				deps.doctor.verificationCapacity = defaultDoctorDependencies().verificationCapacity
				deps.doctor.diskSpace = func(context.Context, roundconfig.Loaded) []CheckResult {
					return []CheckResult{{Name: HealthCheckDiskSpace, Status: CheckStatusOK}}
				}
				deps.doctor.linuxLane = func(context.Context, roundconfig.Loaded) CheckResult {
					return CheckResult{Name: HealthCheckLinuxLane, Status: CheckStatusOK}
				}
			})
			before := snapshotMachineCapacityHome(t, home)
			line := "disk-space: ok\nverification-capacity: " + status + " (" + detail + ")\nlinux-lane: ok\n"
			want := strings.Replace(baseline.String(), "codex: ok\n", line+"codex: ok\n", 1)
			assertRunsListTranscript(t, []string{"doctor"}, want, "", exitOK)
			assertMachineCapacityHomeUnchanged(t, home, before)
		})
	}
}

func TestMachineCapacityStatusReadFailuresAndOff(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"empty", "off", "unreadable entry", "unreadable directory"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			home, _ := withCLIWorkspace(t)
			want, diag := "No Runs found.\n", ""
			switch name {
			case "off":
				one, two := seedMachineCapacityStatus(t, home)
				mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "runs:\n  max_verifying: 0\n")
				want += strings.Replace(machineCapacityStatusBlock(one, two), "1 of 1 held", "off, 1 held", 1)
			case "unreadable entry":
				path := filepath.Join(home, ".roundfix", "verification-capacity", "broken.json")
				mustMkdir(t, filepath.Dir(path))
				mustWrite(t, path, "{")
				want += "Machine Verification Capacity: 0 of 1 held, 0 waiting\n  unreadable: " + path + "\n"
			case "unreadable directory":
				path := filepath.Join(home, ".roundfix", "verification-capacity")
				mustMkdir(t, filepath.Dir(path))
				mustWrite(t, path, "not a directory")
				_, err := store.ReadMachineVerification(t.Context(), home)
				if err == nil {
					t.Fatal("wanted directory read error")
				}
				diag = "roundfix: warning: read Machine Verification Capacity: " + err.Error() + "\n"
			}
			before := snapshotMachineCapacityHome(t, home)
			assertRunsListTranscript(t, []string{"runs", "list"}, want, diag, exitOK)
			assertMachineCapacityHomeUnchanged(t, home, before)
		})
	}
}

// Cancel synchronously when acquisition reports its first wait: no timer or
// background process determines when this test stops waiting.
type machineCapacityCancelLog struct {
	bytes.Buffer
	cancel context.CancelFunc
}

func (log *machineCapacityCancelLog) Write(data []byte) (int, error) {
	n, err := log.Buffer.Write(data)
	log.cancel()
	return n, err
}

func TestRepositoryGateNamesCapacityHoldersBeforeACancelledWait(t *testing.T) {
	t.Parallel()
	f := newDeliveryItemFixture(t)
	one, _ := seedMachineCapacityStatus(t, f.home)
	mustWrite(t, filepath.Join(f.home, ".roundfix", "verification-capacity", "sequence"), "8\n")
	f.workflow.loaded.Config.Runs.MaxVerifying = 1
	f.workflow.loaded.Config.Defaults.Verification = "exit 99"
	ctx, cancel := context.WithCancel(commandContextForTest(t, t.Context()))
	defer cancel()
	log := &machineCapacityCancelLog{cancel: cancel}
	f.workflow.log = log
	_, err := f.workflow.Gate(ctx, f.repo, "widget", f.head)
	if !errors.Is(err, context.Canceled) || !strings.HasPrefix(err.Error(), "wait for the Machine Verification Capacity: ") {
		t.Fatalf("gate error=%v", err)
	}
	want := fmt.Sprintf("roundfix: Delivery Queue item widget: waiting for the Machine Verification Capacity; 1 of 1 held by Delivery Queue item 0001-holder in %s, repository gate\n", one)
	if log.String() != want {
		t.Fatalf("log=%q want=%q", log.String(), want)
	}
	view, err := store.ReadMachineVerification(t.Context(), f.home)
	if err != nil || len(view.Entries) != 2 {
		t.Fatalf("entries=%+v err=%v", view, err)
	}
	output := filepath.Join(f.workflow.loaded.Config.Defaults.ArtifactDir, "delivery", "widget", "repository-gate.log")
	if _, err := os.Stat(output); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Verification started during wait: %v", err)
	}
}
