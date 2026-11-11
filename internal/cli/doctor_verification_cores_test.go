// Boundary: temporary Run Databases and injected readiness checks; no ACP adapter starts.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

func TestDoctorReportsTheVerificationCoreBudget(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, gomaxprocs        string
		maxVerifying, maxActive int
		active                  int
		unreadable, missing     bool
		want                    string
	}{
		{
			name: "default", maxVerifying: 1, active: 2,
			want: "verification-cores: ok (6 per Verification: half of 12 logical CPUs divided by runs.max_verifying 1; 2 active Implement Runs)\n",
		},
		{
			name: "capacity off", maxActive: 3,
			want: "verification-cores: ok (2 per Verification: half of 12 logical CPUs divided by runs.max_active 3 while runs.max_verifying is 0; 0 active Implement Runs)\n",
		},
		{
			name: "undivided with three Runs", active: 3,
			want: "verification-cores: warn (6 per Verification: half of 12 logical CPUs, undivided while runs.max_verifying and runs.max_active are 0; 3 active Implement Runs; 3 at once may take 18 of 12 logical CPUs; next: lower GOMAXPROCS in the environment or set runs.max_verifying in User Config)\n",
		},
		{
			name: "explicit within the machine", gomaxprocs: "12", maxVerifying: 1, active: 1,
			want: "verification-cores: ok (12 per Verification: GOMAXPROCS 12 in the environment; 1 active Implement Runs)\n",
		},
		{
			name: "explicit beyond the machine", gomaxprocs: "8", maxVerifying: 2, active: 1,
			want: "verification-cores: warn (8 per Verification: GOMAXPROCS 8 in the environment; 1 active Implement Runs; 2 at once may take 16 of 12 logical CPUs; next: lower GOMAXPROCS in the environment or set runs.max_verifying in User Config)\n",
		},
		{
			name: "unreadable Run Database", maxVerifying: 1, unreadable: true,
			want: "verification-cores: warn (6 per Verification: half of 12 logical CPUs divided by runs.max_verifying 1; active Implement Runs unknown: open Run Database: denied; next: inspect the Run Database with roundfix runs list)\n",
		},
		{
			name: "default dependency without a Run Database", maxVerifying: 1, missing: true,
			want: "verification-cores: ok (6 per Verification: half of 12 logical CPUs divided by runs.max_verifying 1; 0 active Implement Runs)\n",
		},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			home, repo := withCLIWorkspace(t)
			withVerificationCoreBudgetInputs(t, row.gomaxprocs)
			loaded := roundconfig.Loaded{HomeDir: home, GitRoot: repo, Config: roundconfig.Builtin()}
			loaded.Config.Runs.MaxVerifying = row.maxVerifying
			loaded.Config.Runs.MaxActive = row.maxActive
			loaded.Config.Delivery.LinuxLane = true
			if !row.missing {
				seedDoctorCoreBudgetRuns(t, home, row.active)
			}
			checker := newDoctorFakeHealthChecker(
				CheckResult{Name: HealthCheckNode, Status: CheckStatusOK},
				CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK},
				CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK},
			)
			withDoctorFakeLoadedAndReadiness(t, checker, loaded, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.doctor.diskSpace = func(context.Context, roundconfig.Loaded) []CheckResult {
					return []CheckResult{{Name: HealthCheckDiskSpace, Status: CheckStatusOK}}
				}
				if !row.missing {
					deps.doctor.verificationCapacity = func(context.Context, roundconfig.Loaded) CheckResult {
						return CheckResult{Name: HealthCheckVerificationCapacity, Status: CheckStatusOK}
					}
				}
				deps.doctor.linuxLane = func(context.Context, roundconfig.Loaded) CheckResult {
					return CheckResult{Name: HealthCheckLinuxLane, Status: CheckStatusOK}
				}
				if row.unreadable {
					deps.gc.openStorageReader = func(context.Context, string) (*store.Store, error) {
						return nil, errors.New("open Run Database: denied")
					}
				}
			})
			before := snapshotDoctorPath(t, home)
			var baseline, baselineErr bytes.Buffer
			if code := runCLI(t, []string{"doctor"}, &baseline, &baselineErr); code != exitOK || baselineErr.Len() != 0 {
				t.Fatalf("baseline exit = %d; stderr = %s", code, baselineErr.String())
			}
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.doctor.verificationCores = defaultDoctorDependencies().verificationCores
			})
			var stdout, stderr bytes.Buffer
			if code := runCLI(t, []string{"doctor"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
				t.Fatalf("doctor exit = %d; stdout = %s; stderr = %s", code, stdout.String(), stderr.String())
			}
			// Comparing all stdout proves the complete line, its order and nil opt-out.
			want := strings.Replace(baseline.String(), "linux-lane: ok\n", row.want+"linux-lane: ok\n", 1)
			if stdout.String() != want {
				t.Fatalf("stdout = %q, want %q", stdout.String(), want)
			}
			checker.node.Status = CheckStatusFailed
			stdout.Reset()
			stderr.Reset()
			if code := runCLI(t, []string{"doctor"}, &stdout, &stderr); code != exitRunFailed || stderr.Len() != 0 {
				t.Fatalf("failed readiness exit = %d; stderr = %s", code, stderr.String())
			}
			if got := stdout.String(); got != strings.Replace(want, "node: ok\n", "node: failed\n", 1) {
				t.Fatalf("cores check changed failed readiness output: %q", got)
			}
			if after := snapshotDoctorPath(t, home); !reflect.DeepEqual(after, before) {
				t.Fatalf("doctor mutated temporary home: before = %#v; after = %#v", before, after)
			}
		})
	}
}

func seedDoctorCoreBudgetRuns(t *testing.T, home string, active int) {
	t.Helper()
	db, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if db != nil {
			if err := db.Close(); err != nil {
				t.Error(err)
			}
		}
	})
	for i := range active + 1 {
		run, err := db.CreateRun(t.Context(), store.CreateRunRequest{
			Kind: store.KindImplement, GitRoot: filepath.Join(home, fmt.Sprintf("repo-%d", i)),
			LocalBranch: "feat/cores", SpecSlug: fmt.Sprintf("%04d-cores", i+1),
		})
		if err != nil {
			t.Fatal(err)
		}
		if i == active {
			if _, err := db.CompleteRun(t.Context(), run.ID, store.StateClean); err != nil {
				t.Fatal(err)
			}
		}
	}
	if _, err := db.CreateRun(t.Context(), store.CreateRunRequest{
		Kind: store.KindFetch, GitRoot: filepath.Join(home, "review"), LocalBranch: "feat/review",
		HeadRepository: "fixture/review", HeadBranch: "feat/review", PRNumber: "1",
		HeadSHA: "review-head", ArtifactDir: filepath.Join(home, "artifacts"),
	}); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db = nil
}
