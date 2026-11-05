package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
)

// seedOneRunForBrowser gives the browser and the cockpit one Run to read.
func seedOneRunForBrowser(t *testing.T, homeDir, repoDir string) store.Run {
	t.Helper()
	return seedRunsForList(t, homeDir, []runListSeed{{
		kind:      store.KindResolve,
		state:     store.StateActive,
		gitRoot:   repoDir,
		branch:    "feature/terminal",
		prNumber:  "353",
		createdAt: time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC),
	}})[0]
}

func assertNamesTheBufferNotTheWrapper(t *testing.T, text string) {
	t.Helper()
	if !strings.Contains(text, "output *bytes.Buffer is not a terminal file") {
		t.Fatalf("expected the refusal to name *bytes.Buffer, got %q", text)
	}
	if strings.Contains(text, "commandEnvironmentWriter") {
		t.Fatalf("the terminal program received the environment wrapper: %q", text)
	}
}

func TestTerminalProgramsReceiveTheTerminalFile(t *testing.T) {
	t.Parallel()
	t.Run("run browser", func(t *testing.T) {
		t.Parallel()
		withCLIWorkspace(t)
		output := commandWriter(&bytes.Buffer{}, commandEnvironmentForTest(t))

		_, err := defaultRunBrowserSession(context.Background(), output, nil, nil)

		if err == nil {
			t.Fatal("expected the Run Browser to refuse a buffer")
		}
		assertNamesTheBufferNotTheWrapper(t, err.Error())
	})
	t.Run("attach cockpit", func(t *testing.T) {
		t.Parallel()
		homeDir, repoDir := withCLIWorkspace(t)
		run := seedOneRunForBrowser(t, homeDir, repoDir)
		environment := commandEnvironmentForTest(t)
		var stderr bytes.Buffer
		loaded, err := loadRunReaderCommandConfig(environment, &stderr)
		if err != nil {
			t.Fatalf("load config: %v", err)
		}
		reader, err := store.OpenReader(context.Background(), homeDir)
		if err != nil {
			t.Fatalf("open reader: %v", err)
		}
		defer func() { _ = reader.Close() }()
		stdout := commandWriter(&bytes.Buffer{}, environment)

		code := runAttachCockpit(context.Background(), loaded, reader, run, configuredAttachCapacities(loaded), stdout, &stderr)

		if code != exitRunFailed {
			t.Fatalf("attach cockpit exit=%d stderr=%q; want %d", code, stderr.String(), exitRunFailed)
		}
		assertNamesTheBufferNotTheWrapper(t, stderr.String())
	})
	t.Run("owning cockpit", func(t *testing.T) {
		t.Parallel()
		homeDir, repoDir := withCLIWorkspace(t)
		run := seedOneRunForBrowser(t, homeDir, repoDir)
		setCommandEnvForTest(t, "ROUNDFIX_TUI", "always")
		stderr := commandWriter(&bytes.Buffer{}, commandEnvironmentForTest(t))
		if !liveTUIEnabled(stderr) {
			t.Fatal("expected ROUNDFIX_TUI=always to select the owning cockpit")
		}
		runStore, err := store.Open(context.Background(), homeDir)
		if err != nil {
			t.Fatalf("open run store: %v", err)
		}
		defer func() { _ = runStore.Close() }()

		ui, err := startRunUI(context.Background(), roundtui.LiveRunView{}, run.ID, homeDir, runStore, stderr, false)
		if err != nil {
			t.Fatalf("start run UI: %v", err)
		}
		defer ui.Close(context.Background())
		var cockpitErr error
		// The test is the cockpit's waiter, so Close does not wait again.
		ui.waitOnce.Do(func() { cockpitErr = <-ui.cockpitDone })

		if cockpitErr == nil {
			t.Fatal("expected the owning cockpit to refuse a buffer")
		}
		assertNamesTheBufferNotTheWrapper(t, cockpitErr.Error())
	})
}

func TestRunBrowserEndWithoutASelectionPrintsOneLine(t *testing.T) {
	t.Parallel()
	const closedWithQ = "roundfix: Run Browser closed (q); no Run selected.\n"
	for _, tt := range []struct {
		name       string
		args       []string
		seed       bool
		outcome    roundtui.BrowserOutcome
		wantStderr string
	}{
		{name: "attach", args: []string{"attach"}, seed: true, outcome: roundtui.BrowserOutcome{Cancelled: true, Reason: "q"}, wantStderr: closedWithQ},
		{name: "runs", args: []string{"runs"}, seed: true, outcome: roundtui.BrowserOutcome{Cancelled: true, Reason: "q"}, wantStderr: closedWithQ},
		{name: "runs without a Run Database", args: []string{"runs"}, outcome: roundtui.BrowserOutcome{Cancelled: true, Reason: "q"}, wantStderr: closedWithQ},
		{name: "no reason", args: []string{"attach"}, seed: true, outcome: roundtui.BrowserOutcome{Cancelled: true}, wantStderr: "roundfix: Run Browser closed; no Run selected.\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			homeDir, repoDir := withCLIWorkspace(t)
			if tt.seed {
				seedOneRunForBrowser(t, homeDir, repoDir)
			}
			withAttachInteractiveInput(t, true)
			withRunsInteractiveInput(t, true)
			setCommandEnvForTest(t, "ROUNDFIX_TUI", "always")
			cockpit := withBrowserAttachCockpit(t, nil)
			sessions := withRunBrowserSession(t, tt.outcome)
			var stdout, stderr bytes.Buffer

			code := runCLI(t, tt.args, &stdout, &stderr)

			if code != exitOK {
				t.Fatalf("exit=%d stderr=%q; want 0", code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected empty stdout, got %q", stdout.String())
			}
			if stderr.String() != tt.wantStderr {
				t.Fatalf("stderr=%q; want %q", stderr.String(), tt.wantStderr)
			}
			if len(*sessions) != 1 || len(*cockpit) != 0 {
				t.Fatalf("expected one browser pass and no cockpit: sessions=%d cockpit=%+v", len(*sessions), *cockpit)
			}
		})
	}
}
