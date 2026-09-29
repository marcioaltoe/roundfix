// Suite: Delivery Queue limits and Pending Question command surface.
// Invariant: the public CLI validates, persists, and presents queue limits and exactly one operator question.
// Boundary IN: public CLI dispatch and the real SQLite Delivery Queue store.
// Boundary OUT: detached owner launch, injected as the operating-system boundary.
package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/store"
)

func TestDeliverStartRecordsAndPrintsItsLimits(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	setImplementFixtureAuthorizationOperations(t, repoDir, allDeliveryOperations...)
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	before := time.Now().UTC()

	code := runCLI(t, []string{"deliver", "start", implementTestSlug, "--max-duration", "2h", "--max-retries", "2"}, &stdout, &stderr)
	after := time.Now().UTC()

	if code != exitOK || stderr.Len() != 0 || started != 1 {
		t.Fatalf("deliver start exit=%d starts=%d stderr=%q stdout=%q", code, started, stderr.String(), stdout.String())
	}
	queue := openDeliveryQueueForCLI(t, homeDir, repoDir)
	wantEarliest := before.Add(2 * time.Hour).Truncate(time.Second)
	wantLatest := after.Add(2 * time.Hour).Truncate(time.Second)
	if queue.Limits.Deadline.Before(wantEarliest) || queue.Limits.Deadline.After(wantLatest) || queue.Limits.MaxRetries != 2 {
		t.Fatalf("recorded limits = %+v, want deadline in [%s, %s] and retries 2", queue.Limits, wantEarliest, wantLatest)
	}
	want := "Limits: deadline " + queue.Limits.Deadline.Format(time.RFC3339) + ", retries per item 2, concurrency 1, spend not measured\n"
	if stdout.String() != want {
		t.Fatalf("deliver start stdout = %q, want %q", stdout.String(), want)
	}
}

func TestDeliverStartRecordsNoneForOmittedLimits(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	setImplementFixtureAuthorizationOperations(t, repoDir, allDeliveryOperations...)
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver start exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	queue := openDeliveryQueueForCLI(t, homeDir, repoDir)
	if !queue.Limits.Deadline.IsZero() || queue.Limits.MaxRetries != 0 {
		t.Fatalf("recorded omitted limits = %+v, want zero values", queue.Limits)
	}
	want := "Limits: deadline none, retries per item none, concurrency 1, spend not measured\n"
	if stdout.String() != want {
		t.Fatalf("deliver start stdout = %q, want %q", stdout.String(), want)
	}
}

func TestDeliverStartRefusesANonPositiveMaxDuration(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{name: "zero", value: "0s"},
		{name: "negative", value: "-1s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertDeliverStartLimitsRefused(t, []string{"--max-duration", tt.value}, "max-duration must be greater than zero")
		})
	}
}

func TestDeliverStartRefusesAMalformedMaxDuration(t *testing.T) {
	assertDeliverStartLimitsRefused(t, []string{"--max-duration", "not-a-duration"}, "invalid value")
}

func TestDeliverStartRefusesAMaxRetriesBelowOne(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{name: "zero", value: "0"},
		{name: "negative", value: "-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertDeliverStartLimitsRefused(t, []string{"--max-retries", tt.value}, "max-retries must be at least 1")
		})
	}
}

func TestDeliverStatusPrintsOnePendingQuestion(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	queue, err := runStore.CreateDeliveryQueue(t.Context(), repoDir, []string{"first-parked", "second-parked"})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	first := queue.Items[0]
	if _, _, _, err := runStore.RecordDeliveryQueueItemWorktree(
		t.Context(),
		repoDir,
		first.SpecSlug,
		"roundfix/deliver-first-parked",
		"/worktrees/first-parked",
	); err != nil {
		t.Fatalf("record first Delivery Queue item worktree: %v", err)
	}
	first.Stage = store.DeliveryStageParked
	first.Blocker = delivery.BlockerRevalidationFailed + ": SC-REF-UNRESOLVED"
	first.Branch = "roundfix/deliver-first-parked"
	first.Worktree = "/worktrees/first-parked"
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), repoDir, first); err != nil {
		t.Fatalf("park first Delivery Queue item: %v", err)
	}
	second := queue.Items[1]
	second.Stage = store.DeliveryStageParked
	second.Blocker = delivery.BlockerReviewStale
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), repoDir, second); err != nil {
		t.Fatalf("park second Delivery Queue item: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver status exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if strings.Count(stdout.String(), "Pending question:") != 1 {
		t.Fatalf("deliver status Pending Question count = %d, want 1:\n%s", strings.Count(stdout.String(), "Pending question:"), stdout.String())
	}
	for _, want := range []string{
		"Pending question: first-parked parked revalidation-failed: SC-REF-UNRESOLVED\n",
		"Answer: amend the Spec on its item branch in /worktrees/first-parked, then run roundfix deliver retry first-parked\n",
		"Waiting behind it: 1 parked item(s)\n",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("deliver status does not contain %q:\n%s", want, stdout.String())
		}
	}
}

func TestDeliverStatusPrintsNoQuestionWithoutAParkedItem(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	if _, err := runStore.CreateDeliveryQueue(t.Context(), repoDir, []string{implementTestSlug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver status exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if strings.Contains(stdout.String(), "Pending question:") || strings.Contains(stdout.String(), "Answer:") {
		t.Fatalf("deliver status printed a question without a parked item:\n%s", stdout.String())
	}
}

func TestDeliverHelpNamesTheLimitFlags(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"deliver", "--help"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver help exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	for _, flag := range []string{"--max-duration <duration>", "--max-retries <n>"} {
		if !strings.Contains(stdout.String(), flag) {
			t.Fatalf("deliver help does not name %q:\n%s", flag, stdout.String())
		}
	}
}

func assertDeliverStartLimitsRefused(t *testing.T, flags []string, diagnostic string) {
	t.Helper()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	args := append([]string{"deliver", "start", implementTestSlug}, flags...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, args, &stdout, &stderr)

	if code != exitPreflight || stdout.Len() != 0 || started != 0 {
		t.Fatalf("invalid limits %v exit=%d starts=%d stdout=%q stderr=%q", flags, code, started, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), diagnostic) {
		t.Fatalf("invalid limits diagnostic = %q, want %q", stderr.String(), diagnostic)
	}
	if _, err := os.Stat(store.DatabasePath(homeDir)); !os.IsNotExist(err) {
		t.Fatalf("invalid limits %v created Run Database: err=%v repo=%q", flags, err, repoDir)
	}
}
