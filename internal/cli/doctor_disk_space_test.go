// Boundary: temporary locations and injected filesystem capacity; no host disk measurement.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestTheDiskSpaceAncestorWalksPastAFile(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, path string }{
		{"existing directory", root},
		{"missing child", filepath.Join(root, "absent", "child")},
		{"file", file},
		{"below a file", filepath.Join(file, "x", "_pkg_.a")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ancestor, err := doctorDiskSpaceAncestor(tc.path)
			if err != nil || ancestor != root {
				t.Fatalf("ancestor of %q=%q err=%v want=%q", tc.path, ancestor, err, root)
			}
		})
	}
	t.Run("an unreadable directory still fails", func(t *testing.T) {
		if runtime.GOOS == "windows" || os.Geteuid() == 0 {
			t.Skip("directory modes do not deny lookups on Windows or for root")
		}
		locked := filepath.Join(root, "locked")
		if err := os.Mkdir(locked, 0o700); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := os.Chmod(locked, 0o700); err != nil {
				t.Error(err)
			}
		})
		if err := os.Chmod(locked, 0); err != nil {
			t.Fatal(err)
		}
		_, err := doctorDiskSpaceAncestor(filepath.Join(locked, "child"))
		if err == nil || errors.Is(err, os.ErrNotExist) {
			t.Fatalf("unreadable directory error=%v, want a stat error other than not-exist", err)
		}
	})
}

func TestDoctorWarnsWhenALocationIsLowOnDisk(t *testing.T) {
	t.Parallel()
	const mib int64 = 1024 * 1024
	for _, tc := range []struct {
		name                              string
		available                         int64
		unreadable, missing, lane, failed bool
	}{
		{name: "below the threshold", available: 10239 * mib},
		{name: "at the threshold", available: 10240 * mib},
		{name: "unreadable", unreadable: true},
		{name: "missing location", available: 10239 * mib, missing: true},
		{name: "healthy", available: 1 << 40},
		{name: "one byte below the threshold", available: 10240*mib - 1},
		{name: "before linux lane", available: 10239 * mib, lane: true},
		{name: "preserves a failed exit", available: 10239 * mib, failed: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			worktrees := filepath.Join(home, "worktrees")
			cache := filepath.Join(home, "gocache")
			if tc.missing {
				worktrees = filepath.Join(worktrees, "not-created")
				cache = filepath.Join(cache, "not-created")
			} else {
				for _, path := range []string{worktrees, cache} {
					if err := os.MkdirAll(path, 0o700); err != nil {
						t.Fatal(err)
					}
				}
			}
			loaded := roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: t.TempDir(), HomeDir: home}
			loaded.Config.Worktree.Location = "~/" + strings.TrimPrefix(worktrees, home+string(filepath.Separator))
			loaded.Config.Worktree.SharedCaches = []roundconfig.SharedCache{{Path: filepath.Join("{home}", strings.TrimPrefix(cache, home+string(filepath.Separator)))}}
			loaded.Config.Delivery.LinuxLane = tc.lane
			checker := newDoctorFakeHealthChecker(
				CheckResult{Name: HealthCheckNode, Status: CheckStatusOK},
				CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK},
				CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK},
			)
			if tc.failed {
				checker.node.Status = CheckStatusFailed
			}
			withDoctorFakeLoadedAndReadiness(t, checker, loaded, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.doctor.linuxLane = func(context.Context, roundconfig.Loaded) CheckResult {
					return CheckResult{Name: HealthCheckLinuxLane, Status: CheckStatusOK}
				}
			})
			var baseline, baselineErr bytes.Buffer
			baselineCode := runCLI(t, []string{"doctor"}, &baseline, &baselineErr)
			var reads []string
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.doctor.diskSpace = defaultDoctorDependencies().diskSpace
				deps.availableBytes = func(path string) (int64, error) {
					reads = append(reads, path)
					if tc.unreadable {
						return 0, errors.New("capacity denied")
					}
					return tc.available, nil
				}
			})
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			wantCode := exitOK
			if tc.failed {
				wantCode = exitRunFailed
			}
			if code != wantCode || code != baselineCode || stderr.String() != baselineErr.String() || stderr.Len() != 0 {
				t.Fatalf("code=%d baseline=%d want=%d stderr=%q baseline stderr=%q", code, baselineCode, wantCode, stderr.String(), baselineErr.String())
			}
			expectedReads := []string{worktrees, cache, os.TempDir()}
			if tc.missing {
				expectedReads = []string{home, home, os.TempDir()}
				for _, path := range []string{worktrees, cache} {
					if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
						t.Fatalf("missing location %q changed: %v", path, err)
					}
				}
			}
			if !reflect.DeepEqual(reads, expectedReads) {
				t.Fatalf("capacity reads=%q want=%q", reads, expectedReads)
			}
			line := ""
			if tc.unreadable {
				line = fmt.Sprintf("disk-space: warn (DR-DISK-UNREAD: could not read free space for worktree location %s: capacity denied; DR-DISK-UNREAD: could not read free space for Shared Build Cache %s: capacity denied; DR-DISK-UNREAD: could not read free space for temporary directory %s: capacity denied)\n", worktrees, cache, os.TempDir())
			} else if tc.available < 10240*mib {
				line = fmt.Sprintf("disk-space: warn (DR-DISK-LOW: worktree location %s has %d MiB free, below 10240 MiB; DR-DISK-LOW: Shared Build Cache %s has %d MiB free, below 10240 MiB; DR-DISK-LOW: temporary directory %s has %d MiB free, below 10240 MiB; next: roundfix gc)\n", worktrees, tc.available/mib, cache, tc.available/mib, os.TempDir(), tc.available/mib)
			}
			const storageLine = "storage: ok (no Run Database)\n"
			if !strings.Contains(baseline.String(), storageLine) {
				t.Fatalf("missing storage line in baseline: %q", baseline.String())
			}
			want := strings.Replace(baseline.String(), storageLine, storageLine+line, 1)
			if stdout.String() != want {
				t.Fatalf("stdout=%q want=%q", stdout.String(), want)
			}
		})
	}
}
