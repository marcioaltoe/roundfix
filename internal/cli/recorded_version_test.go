// Suite: Recorded Roundfix Version command boundaries
// Invariant: commands refuse before Project Config, while recovery commands report and continue.
// Boundary IN: public command streams, exit codes, and repository config bytes.
// Boundary OUT: live readiness and release transport, replaced by existing hermetic seams.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

const futureRecordedRoundfixVersion = "999999.0.0"

func TestACommandRefusesARepositoryRecordedByANewerRoundfix(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{{"spec", "check", "x"}, {"follow-up"}, {"init", "--scope", "project", "--force"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			t.Parallel()
			repo, refusal := newRecordedVersionRepository(t)
			path := filepath.Join(repo, ".roundfixrc.yml")
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			code := runCLIContext(t, t.Context(), args, &stdout, &stderr)
			output := stdout.String() + stderr.String()
			if code != 2 || !strings.Contains(output, refusal) || strings.Contains(output, "future_section") {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
			}
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, after) {
				t.Fatalf("Project Config changed: %q", after)
			}
		})
	}
}

func TestReadOnlyCommandsWarnOfARepositoryRecordedByANewerRoundfix(t *testing.T) {
	t.Parallel()
	t.Run("runs list", func(t *testing.T) {
		t.Parallel()
		_, refusal := newRecordedVersionRepository(t)
		var stdout, stderr bytes.Buffer
		code := runCLIContext(t, t.Context(), []string{"runs", "list"}, &stdout, &stderr)
		wantWarning := "roundfix: warning: " + refusal + "; Project Config not read\n"
		if code != 0 || stdout.String() != "No Runs found.\n" || stderr.String() != wantWarning {
			t.Fatalf("exit=%d stdout=%q stderr=%q; want whole Surface Transcript 1, warning=%q", code, &stdout, &stderr, wantWarning)
		}
	})
	t.Run("update check", func(t *testing.T) {
		t.Parallel()
		_, refusal := newRecordedVersionRepository(t)
		configureUpdateCheckTest(t, "0.0.0-dev", func(context.Context) ([]app.Release, error) {
			t.Fatal("development build reached release transport")
			return nil, nil
		})
		var stdout, stderr bytes.Buffer
		code := runCLIContext(t, t.Context(), []string{"update", "check"}, &stdout, &stderr)
		wantWarning := "roundfix: warning: " + refusal + "; Project Config not read\n"
		releaseLine := "Roundfix 0.0.0-dev is a development build; no release comparison (update.mode notify).\n"
		if code != 0 || stderr.String() != wantWarning || !strings.HasPrefix(stdout.String(), releaseLine) {
			t.Fatalf("exit=%d stdout=%q stderr=%q; want release line=%q, warning=%q", code, &stdout, &stderr, releaseLine, wantWarning)
		}
	})
	t.Run("stop", func(t *testing.T) {
		t.Parallel()
		repo, refusal := newRecordedVersionRepository(t)
		home := commandEnvironmentForTest(t).homeDir
		runStore, err := store.Open(t.Context(), home)
		if err != nil {
			t.Fatal(err)
		}
		active, createErr := runStore.CreateRun(t.Context(), store.CreateRunRequest{
			Kind:           store.KindResolve,
			HeadRepository: "owner/project",
			HeadBranch:     "feature/review",
			BaseRepository: "owner/project",
			PRNumber:       "123",
			GitRoot:        repo,
			LocalBranch:    "feature/review",
			HeadSHA:        "abc123",
			ArtifactDir:    filepath.Join(repo, ".roundfix"),
		})
		closeErr := runStore.Close()
		if createErr != nil {
			t.Fatal(createErr)
		}
		if closeErr != nil {
			t.Fatal(closeErr)
		}
		var stdout, stderr bytes.Buffer
		code := runCLIContext(t, t.Context(), []string{"stop", active.ID}, &stdout, &stderr)
		wantWarning := "roundfix: warning: " + refusal + "; Project Config not read\n"
		if code != 0 || stderr.String() != wantWarning || !strings.Contains(stdout.String(), "Stop Request recorded;") {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
		}
		assertStopRequested(t, home, active.ID)
	})
	t.Run("doctor", func(t *testing.T) {
		t.Parallel()
		_, refusal := newRecordedVersionRepository(t)
		checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
		withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
			return profileProofResult{}
		})
		updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
			deps.doctor.loadConfig = func(options roundconfig.LoadOptions) (roundconfig.Loaded, error) {
				if !options.AllowNewerRepository {
					t.Fatal("Doctor must allow newer repository")
				}
				return roundconfig.Load(options)
			}
		})
		var stdout, stderr bytes.Buffer
		code := runCLIContext(t, t.Context(), []string{"doctor"}, &stdout, &stderr)
		firstLine := "baseline-version: failed (" + refusal + "; Project Config not read; next: roundfix upgrade)\n"
		if code != 1 || !strings.HasPrefix(stdout.String(), firstLine+"node: ok\n") || stderr.Len() != 0 {
			t.Fatalf("exit=%d stdout=%q stderr=%q; want first line=%q then node: ok", code, &stdout, &stderr, firstLine)
		}
		if checker.nodeCalls != 1 || checker.acpxCalls != 1 || checker.adapterCalls == 0 || checker.codexCalls != 1 {
			t.Fatalf("independent checks did not all run: %#v", checker)
		}
	})
}

func newRecordedVersionRepository(t *testing.T) (repo, refusal string) {
	t.Helper()
	_, repo = withCLIWorkspace(t)
	path := filepath.Join(repo, roundconfig.SetupManifestPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"roundfixVersion":"`+futureRecordedRoundfixVersion+`"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".roundfixrc.yml"), []byte("future_section:\n  enabled: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	refusal = fmt.Sprintf("Setup Manifest %q records Roundfix %s, newer than this binary's %s; a newer roundfix wrote this repository's Baseline, so use that binary or run 'roundfix upgrade'", path, futureRecordedRoundfixVersion, app.NormalizeVersion(app.Version))
	return repo, refusal
}
