// Boundary: fake build and Verification commands, disposable repositories and injected capacity.
package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/delivery"
)

func TestTheItemBuildOnAFullDiskIsADiskFullFailure(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, output, path string }{
		{"path in the log", "write {root}/x/_pkg_.a: no space left on device", "{root}/x/_pkg_.a"},
		{"no path in the log", "compile: writing output: write $WORK/b001/_pkg_.a: no space left on device", ""},
		{"other failure", "build failed", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := filepath.ToSlash(t.TempDir())
			tc.output = strings.ReplaceAll(tc.output, "{root}", root)
			tc.path = strings.ReplaceAll(tc.path, "{root}", root)
			f := newDeliveryItemFixture(t)
			f.workflow.loaded.Config.Worktree.Location = f.home
			f.workflow.loaded.Config.Delivery.ItemBinary.Build = fmt.Sprintf("printf '%%s\\n' '%s' >&2; exit 7", tc.output)
			deps := commandDependenciesForContext(commandContextForTest(t, t.Context()))
			deps.availableBytes = func(path string) (int64, error) {
				if path == f.home {
					return 3*1024*1024 + 99, nil
				}
				return 7 * 1024 * 1024, nil
			}
			ctx := contextWithCommandDependencies(t.Context(), deps)
			_, err := f.workflow.stepExecutable(ctx, f.repo, "widget", "implement")
			logPath := filepath.Join(f.workflow.loaded.Config.Defaults.ArtifactDir, "delivery", "widget", "item-binary-build.log")
			wantText := fmt.Sprintf("build item binary: verification command %q failed; diagnostics: %s: exit status 7", f.workflow.loaded.Config.Delivery.ItemBinary.Build, logPath)
			if err == nil || err.Error() != wantText {
				t.Fatalf("error=%v want=%q", err, wantText)
			}
			var full *delivery.DiskFullError
			if tc.name == "other failure" {
				if errors.As(err, &full) {
					t.Fatalf("other failure became disk-full: %+v", full)
				}
			} else {
				wantPath, wantBytes := tc.path, int64(7*1024*1024)
				if wantPath == "" {
					wantPath, wantBytes = f.home, 3*1024*1024+99
				}
				if !errors.As(err, &full) || full.Path != wantPath || full.AvailableBytes != wantBytes {
					t.Fatalf("error=%v full=%+v want path=%s bytes=%d", err, full, wantPath, wantBytes)
				}
				var commandErr *daemon.VerificationCommandError
				if !errors.As(err, &commandErr) || full.Error() != commandErr.Error() {
					t.Fatalf("lost underlying error: %v", err)
				}
			}
			data, readErr := os.ReadFile(logPath)
			if readErr != nil || string(data) != tc.output+"\n" {
				t.Fatalf("log=%q err=%v", data, readErr)
			}
			f.assertStarts(t)
		})
	}
}

func TestTheRepositoryGateOnAFullDiskParksDiskFull(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, output, blocker string }{
		{"full disk", "write {root}/x/_pkg_.a: no space left on device", "disk-full: {root}/x/_pkg_.a: 2 MiB free"},
		{"other failure", "gate failed", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := filepath.ToSlash(t.TempDir())
			tc.output = strings.ReplaceAll(tc.output, "{root}", root)
			tc.blocker = strings.ReplaceAll(tc.blocker, "{root}", root)
			f := newDeliveryItemFixture(t)
			f.workflow.loaded.Config.Defaults.Verification = fmt.Sprintf("printf '%%s\\n' '%s' >&2; exit 9", tc.output)
			deps := commandDependenciesForContext(commandContextForTest(t, t.Context()))
			deps.availableBytes = func(string) (int64, error) { return 2*1024*1024 + 1, nil }
			result, err := f.workflow.Gate(contextWithCommandDependencies(t.Context(), deps), f.repo, "widget", f.head)
			logPath := filepath.Join(f.workflow.loaded.Config.Defaults.ArtifactDir, "delivery", "widget", "repository-gate.log")
			wantReason := fmt.Sprintf("verification command %q failed; diagnostics: %s: exit status 9", f.workflow.loaded.Config.Defaults.Verification, logPath)
			if err != nil || result.Passed || result.Blocker != tc.blocker || result.Reason != wantReason {
				t.Fatalf("result=%+v err=%v want reason=%q blocker=%q", result, err, wantReason, tc.blocker)
			}
		})
	}
}

func TestDiskFullFailureUsesOnlyTheLogTailAndPreservesTheError(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, log, errorText, wantPath string
		enospc, missing                bool
	}{
		{name: "last path wins", log: "write {dir}/first: no space left on device\nwrite {dir}/last: NO SPACE LEFT ON DEVICE", wantPath: "{dir}/last"},
		{name: "error path fallback", log: "no space left on device", errorText: "write {dir}/error: no space left on device", wantPath: "{dir}/error"},
		{name: "wrapped ENOSPC", errorText: "write {dir}/error: no space left on device", enospc: true, missing: true, wantPath: "{dir}/error"},
		{name: "outside tail", log: "write /first: no space left on device\n" + strings.Repeat("x", 64*1024)},
		{name: "unreadable log", missing: true},
		{name: "text alone is not ENOSPC", errorText: "no space left on device", missing: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			tc.log = strings.ReplaceAll(tc.log, "{dir}", filepath.ToSlash(dir))
			tc.errorText = strings.ReplaceAll(tc.errorText, "{dir}", filepath.ToSlash(dir))
			tc.wantPath = strings.ReplaceAll(tc.wantPath, "{dir}", filepath.ToSlash(dir))
			logPath := filepath.Join(dir, "failure.log")
			if !tc.missing {
				if err := os.WriteFile(logPath, []byte(tc.log), 0600); err != nil {
					t.Fatal(err)
				}
			}
			cause := errors.New(tc.errorText)
			if tc.enospc {
				cause = fmt.Errorf("write %s: %w", filepath.ToSlash(filepath.Join(dir, "error")), syscall.ENOSPC)
			}
			var reads []string
			ctx := contextWithCommandDependencies(t.Context(), commandDependencies{
				availableBytes: func(path string) (int64, error) { reads = append(reads, path); return 1024*1024 + 123, nil },
			})
			full := diskFullFailure(ctx, roundconfig.Loaded{}, cause, logPath)
			if tc.wantPath == "" {
				if full != nil || len(reads) != 0 {
					t.Fatalf("full=%+v reads=%v", full, reads)
				}
				return
			}
			if full == nil || full.Path != tc.wantPath || full.AvailableBytes != 1024*1024+123 || full.Error() != cause.Error() || !errors.Is(full, cause) {
				t.Fatalf("full=%+v cause=%v", full, cause)
			}
			ancestor, err := doctorDiskSpaceAncestor(tc.wantPath)
			if err != nil || len(reads) != 1 || reads[0] != ancestor {
				t.Fatalf("reads=%v ancestor=%s err=%v", reads, ancestor, err)
			}
		})
	}
}

func TestADiskFullPathMayHoldSpaces(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, text, want string }{
		{"spaces", "write /Users/me/Google Drive/repo/_pkg_.a: no space left on device", "/Users/me/Google Drive/repo/_pkg_.a"},
		{"open with spaces", "open /Volumes/My Disk/go-build/x: no space left on device", "/Volumes/My Disk/go-build/x"},
		{"no spaces", "write /build/x/_pkg_.a: no space left on device", "/build/x/_pkg_.a"},
		{"last wins", "write /first one: no space left on device\nwrite /build/last: no space left on device", "/build/last"},
		{"no path", "compile: writing output: write $WORK/b001/_pkg_.a: no space left on device", ""},
	} {
		if got := lastDiskFullPath(tc.text); got != tc.want {
			t.Errorf("%s: lastDiskFullPath(%q) = %q, want %q", tc.name, tc.text, got, tc.want)
		}
	}
}

func TestADiskFullPathBelowAFileIsReadAtItsNearestDirectory(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "x", "_pkg_.a")
	logPath := filepath.Join(root, "failure.log")
	if err := os.WriteFile(logPath, []byte(fmt.Sprintf("write %s: no space left on device\n", path)), 0o600); err != nil {
		t.Fatal(err)
	}
	cause := errors.New("build failed")
	const available int64 = 5*1024*1024 + 7
	var reads []string
	ctx := contextWithCommandDependencies(t.Context(), commandDependencies{
		availableBytes: func(path string) (int64, error) {
			reads = append(reads, path)
			return available, nil
		},
	})
	full := diskFullFailure(ctx, roundconfig.Loaded{}, cause, logPath)
	if full == nil || full.Path != path || full.AvailableBytes != available || !errors.Is(full, cause) {
		t.Fatalf("full=%+v cause=%v want path=%s bytes=%d", full, cause, path, available)
	}
	if !reflect.DeepEqual(reads, []string{root}) {
		t.Fatalf("capacity reads=%q want=%q", reads, []string{root})
	}
}
