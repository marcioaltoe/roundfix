package cli

// Boundary: real CLI dispatch, local configuration and release installation;
// release lookup, downloads and the installed child are fixture dependencies.
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
)

func upgradeFixtureNotice() string {
	return fmt.Sprintf("recommendations: snapshot %s; %d current, 0 differ, 0 pinned\n", roundconfig.ModelRecommendationSnapshotVersion, len(roundconfig.RequiredWorkCategories()))
}

func prepareNoticeUpgrade(t *testing.T, outcome string) *upgradeFake {
	t.Helper()
	withCLIWorkspace(t)
	fake := newUpgradeFake(t)
	fake.releaseTag = "v1.0.0"
	switch outcome {
	case "no release":
		fake.releaseErr = app.ErrNoReleases
	case "available", "installed", "failed":
		fake.releaseTag = "v1.1.0"
		fake.assets = []app.ReleaseAsset{fake.platformAsset("roundfix_darwin_arm64", []byte("new binary"))}
		if outcome == "failed" {
			fake.contentByURL = map[string][]byte{}
		}
	}
	withUpgradeFakeDeps(t, fake)
	return fake
}

func TestUpgradeNoticeNamesPendingHistory(t *testing.T) {
	t.Parallel()
	for _, outcome := range []string{"current", "installed"} {
		t.Run(outcome, func(t *testing.T) {
			prepareNoticeUpgrade(t, outcome)
			repo := commandEnvironmentForTest(t).workDir
			mustMkdir(t, filepath.Join(repo, "docs/specs"))
			for _, slug := range []string{"0001-first", "0002-second", "0003-third"} {
				dir := filepath.Join(repo, "docs/history/specs", slug)
				mustMkdir(t, dir)
				mustWrite(t, filepath.Join(dir, "_prd.md"), "# Legacy Spec\n")
			}
			dir := filepath.Join(repo, "docs/history/findings")
			mustMkdir(t, dir)
			mustWrite(t, filepath.Join(dir, "retired.md"), "---\nstatus: closed\nreason: resolved\n---\n\n# Retired Finding\n\nFull detail.\n")
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"upgrade"}, &stdout, &stderr)
			wantOut := "already current 1.0.0\n"
			if outcome == "installed" {
				wantOut = "upgraded 1.0.0 → 1.1.0\n"
			}
			wantErr := upgradeFixtureNotice() + "roundfix: history: 4 unit(s) pending sanitize (3 Legacy Archive Folder(s), findings); run roundfix baseline update to plan them\n"
			if code != exitOK || stdout.String() != wantOut || stderr.String() != wantErr {
				t.Fatalf("exit=%d stdout=%q stderr=%q; want stdout=%q stderr=%q", code, &stdout, &stderr, wantOut, wantErr)
			}
		})
	}
}

func TestUpgradeNoticeOutsideARepositoryNamesBaselineUpdate(t *testing.T) {
	t.Parallel()
	prepareNoticeUpgrade(t, "current")
	env := commandEnvironmentForTest(t)
	setCommandEnvironmentForTest(t, env.homeDir, t.TempDir())
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"upgrade"}, &stdout, &stderr)
	wantErr := upgradeFixtureNotice() + "roundfix: history: outside a repository; run roundfix baseline update in each adopted repository to plan its pending history\n"
	if code != exitOK || stdout.String() != "already current 1.0.0\n" || stderr.String() != wantErr {
		t.Fatalf("exit=%d stdout=%q stderr=%q; want stderr=%q", code, &stdout, &stderr, wantErr)
	}
}

func TestUpgradeNoticeIsSilentWithoutPendingHistory(t *testing.T) {
	t.Parallel()
	prepareNoticeUpgrade(t, "current")
	repo := commandEnvironmentForTest(t).workDir
	mustMkdir(t, filepath.Join(repo, "docs/specs"))
	for _, kind := range []string{"findings", "backlog"} {
		dir := filepath.Join(repo, "docs/history", kind)
		mustMkdir(t, dir)
		mustWrite(t, filepath.Join(dir, "reduced.md"), "---\nstatus: closed\n---\n\n# Reduced entry\n\nFull text in Git at `"+strings.Repeat("a", 40)+"`: `docs/history/"+kind+"/reduced.md`.\n")
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"upgrade"}, &stdout, &stderr)
	if code != exitOK || stdout.String() != "already current 1.0.0\n" || stderr.String() != upgradeFixtureNotice() {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
	}
}

func TestUpgradeNoticeHistoryInventoryFailures(t *testing.T) {
	t.Parallel()
	for _, failure := range []string{"archive", "kind"} {
		t.Run(failure, func(t *testing.T) {
			prepareNoticeUpgrade(t, "current")
			repo := commandEnvironmentForTest(t).workDir
			mustMkdir(t, filepath.Join(repo, "docs/specs"))
			mustMkdir(t, filepath.Join(repo, "docs/history"))
			if failure == "archive" {
				mustWrite(t, filepath.Join(repo, "docs/history/specs"), "not a directory\n")
			} else {
				dir := filepath.Join(repo, "docs/history/findings")
				mustMkdir(t, dir)
				if err := os.Symlink("missing", filepath.Join(dir, "broken\nlink.md")); err != nil {
					t.Fatal(err)
				}
			}
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"upgrade"}, &stdout, &stderr)
			notice := strings.TrimPrefix(stderr.String(), upgradeFixtureNotice())
			if code != exitOK || stdout.String() != "already current 1.0.0\n" || !strings.HasPrefix(stderr.String(), upgradeFixtureNotice()) || !strings.HasPrefix(notice, "roundfix: history not checked: ") || strings.Count(notice, "\n") != 1 {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
			}
		})
	}
}

func TestUpgradeNoticeIsSilentWithExternalSpecRoot(t *testing.T) {
	t.Parallel()
	prepareNoticeUpgrade(t, "current")
	repo := commandEnvironmentForTest(t).workDir
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), fmt.Sprintf("specs:\n  root: %q\n", t.TempDir()))
	dir := filepath.Join(repo, "docs/history/findings")
	mustMkdir(t, dir)
	mustWrite(t, filepath.Join(dir, "retired.md"), "# Pending Finding\n")
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"upgrade"}, &stdout, &stderr)
	if code != exitOK || stdout.String() != "already current 1.0.0\n" || stderr.String() != upgradeFixtureNotice() {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
	}
}

func TestUpgradeWritesTheNoticeOnEveryReleaseOutcome(t *testing.T) {
	t.Parallel()
	for _, outcome := range []string{"no release", "current", "available", "installed"} {
		t.Run(outcome, func(t *testing.T) {
			prepareNoticeUpgrade(t, outcome)
			args := []string{"upgrade"}
			if outcome == "available" {
				args = append(args, "--check")
			}
			var stdout, stderr bytes.Buffer
			if code := runCLI(t, args, &stdout, &stderr); code != exitOK || stderr.String() != upgradeFixtureNotice() {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
			}
		})
	}
	t.Run("in process differences use the shared renderer", func(t *testing.T) {
		prepareNoticeUpgrade(t, "current")
		writeProfilesCheckFixture(t)
		var checkOut, checkErr, stdout, stderr bytes.Buffer
		if code := runCLI(t, []string{"profiles", "check"}, &checkOut, &checkErr); code != exitOK {
			t.Fatal(checkErr.String())
		}
		if code := runCLI(t, []string{"upgrade"}, &stdout, &stderr); code != exitOK || stderr.String() != checkOut.String() {
			t.Fatalf("exit=%d notice=%q want=%q", code, &stderr, &checkOut)
		}
	})
}

func TestUpgradeAsksTheInstalledExecutableAfterAnInstall(t *testing.T) {
	t.Parallel()
	for _, outcome := range []string{"no release", "current", "available", "installed"} {
		t.Run(outcome, func(t *testing.T) {
			fake := prepareNoticeUpgrade(t, outcome)
			environment := commandEnvironmentForTest(t)
			calls := 0
			childOutput := "recommendations: snapshot installed-version; 5 current, 0 differ, 0 pinned\n"
			var stdout, stderr bytes.Buffer
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.upgrade.installedProfilesCheck = func(ctx context.Context, path, dir string) ([]byte, error) {
					calls++
					if path != fake.executablePath || dir != environment.workDir {
						t.Fatalf("child path=%q dir=%q", path, dir)
					}
					deadline, ok := ctx.Deadline()
					if remaining := time.Until(deadline); !ok || remaining <= 0 || remaining > 10*time.Second {
						t.Fatalf("child deadline=%v present=%v", deadline, ok)
					}
					if stdout.String() != "upgraded 1.0.0 → 1.1.0\n" {
						t.Fatalf("outcome must precede child, stdout=%q", &stdout)
					}
					return []byte(childOutput), nil
				}
			})
			args := []string{"upgrade"}
			if outcome == "available" {
				args = append(args, "--check")
			}
			if code := runCLI(t, args, &stdout, &stderr); code != exitOK {
				t.Fatalf("exit=%d stderr=%q", code, &stderr)
			}
			wantCalls := 0
			if outcome == "installed" {
				wantCalls = 1
				if stderr.String() != childOutput {
					t.Fatalf("notice=%q want=%q", &stderr, childOutput)
				}
			}
			if calls != wantCalls {
				t.Fatalf("child calls=%d want=%d", calls, wantCalls)
			}
		})
	}
}

func TestUpgradeKeepsItsExitCodeWhenTheCheckFails(t *testing.T) {
	t.Parallel()
	for _, failure := range []string{"child", "timeout", "configuration", "work directory"} {
		t.Run(failure, func(t *testing.T) {
			outcome := "installed"
			if failure == "configuration" {
				outcome = "current"
			}
			prepareNoticeUpgrade(t, outcome)
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.upgrade.installedProfilesCheck = func(context.Context, string, string) ([]byte, error) {
					if failure == "timeout" {
						return []byte("discard child output\n"), context.DeadlineExceeded
					}
					return []byte("discard child output\n"), errors.New("child failed\nwith detail")
				}
			})
			env := commandEnvironmentForTest(t)
			wantReason := "child failed with detail"
			switch failure {
			case "timeout":
				wantReason = "context deadline exceeded"
			case "configuration":
				env.homeDirErr = errors.New("home unavailable")
				wantReason = "resolve User Config home: home unavailable"
			case "work directory":
				env.workDirErr = errors.New("directory unavailable")
				wantReason = "resolve process working directory: directory unavailable"
			}
			var stdout, stderr bytes.Buffer
			code := runWithContext(context.Background(), []string{"upgrade"}, &stdout, &stderr, env)
			wantOut := "upgraded 1.0.0 → 1.1.0\n"
			if outcome == "current" {
				wantOut = "already current 1.0.0\n"
			}
			wantErr := "roundfix: recommendations not checked: " + wantReason + "\n"
			if code != exitOK || stdout.String() != wantOut || stderr.String() != wantErr {
				t.Fatalf("exit=%d stdout=%q stderr=%q; want exit=0 stdout=%q stderr=%q", code, &stdout, &stderr, wantOut, wantErr)
			}
		})
	}
	t.Run("malformed local config", func(t *testing.T) {
		prepareNoticeUpgrade(t, "current")
		env := commandEnvironmentForTest(t)
		mustMkdir(t, filepath.Join(env.homeDir, ".roundfix"))
		mustWrite(t, filepath.Join(env.homeDir, ".roundfix", "config.yml"), "profiles: [\n")
		var stdout, stderr bytes.Buffer
		if code := runCLI(t, []string{"upgrade"}, &stdout, &stderr); code != exitOK || stdout.String() != "already current 1.0.0\n" || !strings.HasPrefix(stderr.String(), "roundfix: recommendations not checked: ") || strings.Count(stderr.String(), "\n") != 1 {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
		}
	})
}

func TestUpgradeWritesNoNoticeOnHelpUsageErrorOrFailure(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		args []string
		code int
	}{
		{"help", []string{"upgrade", "--help"}, exitOK},
		{"usage", []string{"upgrade", "--unknown"}, exitPreflight},
		{"failure", []string{"upgrade"}, exitRunFailed},
	} {
		t.Run(test.name, func(t *testing.T) {
			prepareNoticeUpgrade(t, "failed")
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.upgrade.installedProfilesCheck = func(context.Context, string, string) ([]byte, error) { t.Fatal("unexpected child"); return nil, nil }
			})
			var stdout, stderr bytes.Buffer
			if code := runCLI(t, test.args, &stdout, &stderr); code != test.code {
				t.Fatalf("exit=%d want=%d", code, test.code)
			}
			if strings.Contains(stderr.String(), "recommendations") {
				t.Fatalf("unexpected notice: %q", &stderr)
			}
		})
	}
}
