// Suite: Run Window command
// Invariant: repository-scoped cutoff state changes only through explicit valid commands.
// Boundary IN: CLI parsing, Git-root resolution, and the real Run Database.
// Boundary OUT: Implement Command preflight enforcement, owned by task_03.

package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/store"
)

func TestWindowSetResolvesNextOccurrenceFromNestedWorktreePath(t *testing.T) {
	location := time.FixedZone("BRT", -3*60*60)
	now := time.Date(2026, time.August, 26, 23, 0, 0, 0, location)
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	nestedDir := filepath.Join(repoDir, "nested", "work")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("create nested work directory: %v", err)
	}
	setCommandEnvironmentForTest(t, homeDir, nestedDir)
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.currentRunWindowTime = func() time.Time { return now }
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"window", "set", "07:00"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("window set exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("window set stderr = %q, want empty", stderr.String())
	}
	if output := stdout.String(); !strings.Contains(output, "Run Window set") || !strings.Contains(output, "2026-08-27 07:00 BRT") {
		t.Fatalf("window set output = %q, want effective next-day cutoff", output)
	}

	window, found := readWindowForTest(t, homeDir, repoDir)
	if !found {
		t.Fatal("window set stored no Run Window")
	}
	want := time.Date(2026, time.August, 27, 7, 0, 0, 0, location)
	if !window.CutoffAt.Equal(want) {
		t.Fatalf("stored cutoff = %s, want %s", window.CutoffAt, want)
	}
}

func TestWindowSetAcceptsSameDayAndFutureAbsoluteCutoffs(t *testing.T) {
	location := time.FixedZone("BRT", -3*60*60)
	tests := []struct {
		name   string
		now    time.Time
		input  string
		cutoff time.Time
	}{
		{
			name:   "same-day wall clock still ahead",
			now:    time.Date(2026, time.August, 26, 6, 0, 0, 0, location),
			input:  "07:00",
			cutoff: time.Date(2026, time.August, 26, 7, 0, 0, 0, location),
		},
		{
			name:   "future absolute local instant",
			now:    time.Date(2026, time.August, 26, 23, 0, 0, 0, location),
			input:  "2026-08-28T09:15",
			cutoff: time.Date(2026, time.August, 28, 9, 15, 0, 0, location),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
				dependencies.currentRunWindowTime = func() time.Time { return tt.now }
			})

			runWindowCommandForTest(t, []string{"window", "set", tt.input}, exitOK)
			window, found := readWindowForTest(t, homeDir, repoDir)
			if !found {
				t.Fatal("window set stored no Run Window")
			}
			if !window.CutoffAt.Equal(tt.cutoff) {
				t.Fatalf("stored cutoff = %s, want %s", window.CutoffAt, tt.cutoff)
			}
		})
	}
}

func TestWindowSetRejectsPastAbsoluteCutoffWithoutStoring(t *testing.T) {
	location := time.FixedZone("BRT", -3*60*60)
	now := time.Date(2026, time.August, 26, 23, 0, 0, 0, location)
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.currentRunWindowTime = func() time.Time { return now }
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"window", "set", "2026-08-26T22:59"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("past window set exit = %d, want %d", code, exitPreflight)
	}
	if stdout.Len() != 0 {
		t.Fatalf("past window set stdout = %q, want empty", stdout.String())
	}
	if output := stderr.String(); !strings.Contains(output, "must be in the future") || !strings.Contains(output, "2026-08-26 22:59 BRT") {
		t.Fatalf("past window set stderr = %q, want literal cutoff refusal", output)
	}
	if _, found := readWindowForTest(t, homeDir, repoDir); found {
		t.Fatal("past window set stored a Run Window")
	}
}

func TestWindowSetPreservesExistingWindowUnlessForced(t *testing.T) {
	location := time.FixedZone("BRT", -3*60*60)
	now := time.Date(2026, time.August, 26, 23, 0, 0, 0, location)
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.currentRunWindowTime = func() time.Time { return now }
	})

	runWindowCommandForTest(t, []string{"window", "set", "07:00"}, exitOK)
	standing, found := readWindowForTest(t, homeDir, repoDir)
	if !found {
		t.Fatal("initial window set stored no Run Window")
	}

	stdout := runWindowCommandForTest(t, []string{"window", "set", "08:00"}, exitOK)
	if !strings.Contains(stdout, "unchanged without --force") || !strings.Contains(stdout, "2026-08-27 07:00 BRT") {
		t.Fatalf("non-forced set output = %q, want standing cutoff", stdout)
	}
	preserved, found := readWindowForTest(t, homeDir, repoDir)
	if !found || preserved.CutoffAt != standing.CutoffAt || preserved.CreatedAt != standing.CreatedAt {
		t.Fatalf("non-forced set changed window: got %#v, want %#v", preserved, standing)
	}

	stdout = runWindowCommandForTest(t, []string{"window", "set", "08:00", "--force"}, exitOK)
	if !strings.Contains(stdout, "Run Window replaced") || !strings.Contains(stdout, "2026-08-27 08:00 BRT") {
		t.Fatalf("forced set output = %q, want replacement cutoff", stdout)
	}
	replaced, found := readWindowForTest(t, homeDir, repoDir)
	if !found {
		t.Fatal("forced set removed Run Window")
	}
	want := time.Date(2026, time.August, 27, 8, 0, 0, 0, location)
	if !replaced.CutoffAt.Equal(want) {
		t.Fatalf("forced cutoff = %s, want %s", replaced.CutoffAt, want)
	}
}

func TestWindowShowReportsSetAndAbsentStates(t *testing.T) {
	location := time.FixedZone("BRT", -3*60*60)
	now := time.Date(2026, time.August, 26, 23, 30, 0, 0, location)
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.currentRunWindowTime = func() time.Time { return now }
	})

	stdout := runWindowCommandForTest(t, []string{"window", "show"}, exitOK)
	if !strings.Contains(stdout, "No Run Window is set") || !strings.Contains(stdout, runWindowBoundsExplanation) {
		t.Fatalf("absent window show output = %q, want absence and bounds", stdout)
	}

	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	cutoff := time.Date(2026, time.August, 27, 8, 0, 0, 0, location)
	if _, written, err := runStore.SetRunWindow(context.Background(), repoDir, cutoff, false); err != nil {
		_ = runStore.Close()
		t.Fatalf("seed Run Window: %v", err)
	} else if !written {
		_ = runStore.Close()
		t.Fatal("seed Run Window reported no write")
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}

	stdout = runWindowCommandForTest(t, []string{"window", "show"}, exitOK)
	for _, want := range []string{
		"Cutoff: 2026-08-27 08:00 BRT",
		"Current time: 2026-08-26 23:30 BRT",
		"Remaining: 8h30m0s",
		runWindowBoundsExplanation,
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("set window show output = %q, want %q", stdout, want)
		}
	}
}

func TestWindowClearReportsWhetherWindowExisted(t *testing.T) {
	location := time.FixedZone("BRT", -3*60*60)
	now := time.Date(2026, time.August, 26, 23, 0, 0, 0, location)
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.currentRunWindowTime = func() time.Time { return now }
	})

	runWindowCommandForTest(t, []string{"window", "set", "07:00"}, exitOK)
	stdout := runWindowCommandForTest(t, []string{"window", "clear"}, exitOK)
	if !strings.Contains(stdout, "Run Window cleared") {
		t.Fatalf("first clear output = %q, want cleared report", stdout)
	}
	if _, found := readWindowForTest(t, homeDir, repoDir); found {
		t.Fatal("window clear left the Run Window stored")
	}

	stdout = runWindowCommandForTest(t, []string{"window", "clear"}, exitOK)
	if !strings.Contains(stdout, "No Run Window was set") {
		t.Fatalf("second clear output = %q, want absent report", stdout)
	}
}

func TestWindowHelpExplainsStartAndRunBounds(t *testing.T) {
	_, _ = newImplementWorkspace(t, []implementSeed{{id: "task_01"}})

	stdout := runWindowCommandForTest(t, []string{"window", "--help"}, exitOK)
	for _, want := range []string{"roundfix window <set|show|clear>", "when a Run may start", "budget.max_run_duration", "how long one may run"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("window help output = %q, want %q", stdout, want)
		}
	}
	stdout = runWindowCommandForTest(t, []string{"--help"}, exitOK)
	if !strings.Contains(stdout, "roundfix window <set|show|clear>") || !strings.Contains(stdout, "window     Set, show, or clear this repository's Run Window") {
		t.Fatalf("root help output does not advertise window command: %q", stdout)
	}
}

func TestWindowRejectsMalformedInputAndUnknownSubcommand(t *testing.T) {
	location := time.FixedZone("BRT", -3*60*60)
	now := time.Date(2026, time.August, 26, 23, 0, 0, 0, location)
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.currentRunWindowTime = func() time.Time { return now }
	})

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "missing cutoff", args: []string{"window", "set"}, want: "requires <HH:MM|YYYY-MM-DDTHH:MM>"},
		{name: "malformed cutoff", args: []string{"window", "set", "25:00"}, want: "must match HH:MM or YYYY-MM-DDTHH:MM"},
		{name: "unknown subcommand", args: []string{"window", "move"}, want: "unknown window command"},
		{name: "show argument", args: []string{"window", "show", "extra"}, want: "unexpected argument"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := runCLIContext(t, context.Background(), tt.args, &stdout, &stderr)
			if code != exitPreflight {
				t.Fatalf("exit = %d, want %d", code, exitPreflight)
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("stderr = %q, want %q", stderr.String(), tt.want)
			}
		})
	}
	if _, found := readWindowForTest(t, homeDir, repoDir); found {
		t.Fatal("invalid window commands stored a Run Window")
	}
}

func runWindowCommandForTest(t *testing.T, args []string, wantCode int) string {
	t.Helper()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), args, &stdout, &stderr)
	if code != wantCode {
		t.Fatalf("%v exit = %d, want %d; stderr=%q", args, code, wantCode, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("%v stderr = %q, want empty", args, stderr.String())
	}
	return stdout.String()
}

func readWindowForTest(t *testing.T, homeDir, gitRoot string) (store.RunWindow, bool) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close Run Database: %v", err)
		}
	}()
	window, found, err := runStore.RunWindowFor(context.Background(), gitRoot)
	if err != nil {
		t.Fatalf("read Run Window: %v", err)
	}
	return window, found
}
