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
	"sync"
	"testing"

	"roundfix/internal/daemon"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestImplementGatedGraphRunsRepositoryVerificationAtSettlement(t *testing.T) {
	t.Parallel()
	testImplementSettlementRepositoryCommands(t, "", true)
}

func TestImplementRepositoryAtSettlementOffRunsOnlyDeclaredCommands(t *testing.T) {
	t.Parallel()
	testImplementSettlementRepositoryCommands(t, "verification:\n  repository_at_settlement: false\n", false)
}

func testImplementSettlementRepositoryCommands(t *testing.T, configKey string, wantRepository bool) {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
		implementQAGateSeed("", "task_01"),
	})
	const repositoryCommand = "repository check"
	writeUserConfig(t, home, "defaults:\n  verification: "+repositoryCommand+"\n"+configKey)
	runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}, qaReport: implementQAReport("pass")}
	committer, verifier, _, _ := withImplementCollaborators(t, runner)
	recording := &implementSettlementVerifier{fakeVerifier: verifier}
	withVerifier(t, recording)
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("implement exit = %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	// The real marker command fails in the pre-work probe, then the fake Agent
	// records its work and the command passes. The wrapper records only the
	// non-QA Task's settlement commands.
	want := []string{implementFixtureVerificationCommand("task_01")}
	if wantRepository {
		want = append(want, repositoryCommand)
	}
	if !reflect.DeepEqual(recording.taskCommands, want) {
		t.Fatalf("non-QA Task commands = %v, want %v", recording.taskCommands, want)
	}
	repositoryCalls := 0
	for _, command := range verifier.commands {
		if command == repositoryCommand {
			repositoryCalls++
		}
	}
	wantRepositoryCalls := 1 // QA always runs repository Verification.
	if wantRepository {
		wantRepositoryCalls++
	}
	if repositoryCalls != wantRepositoryCalls {
		t.Fatalf("repository calls = %d, want %d", repositoryCalls, wantRepositoryCalls)
	}
	if committer.calls != 2 || runner.calls != 2 {
		t.Fatalf("Task commits = %d, Agent calls = %d", committer.calls, runner.calls)
	}
}

func TestImplementUngatedGraphRunsRepositoryVerificationAtSettlement(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
	})
	const repositoryCommand = "repository check"
	writeUserConfig(t, home, "defaults:\n  verification: "+repositoryCommand+"\n")
	runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}}
	committer, verifier, _, _ := withImplementCollaborators(t, runner)
	recording := &implementSettlementVerifier{fakeVerifier: verifier}
	withVerifier(t, recording)
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("implement exit = %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	// The real marker command fails in the pre-work probe, then the fake Agent
	// records its work and the command passes. The wrapper records only the
	// non-QA Task's settlement commands.
	want := []string{implementFixtureVerificationCommand("task_01"), repositoryCommand}
	if !reflect.DeepEqual(recording.taskCommands, want) {
		t.Fatalf("non-QA Task commands = %v, want %v", recording.taskCommands, want)
	}
	repositoryCalls := 0
	for _, command := range verifier.commands {
		if command == repositoryCommand {
			repositoryCalls++
		}
	}
	wantRepositoryCalls := 1 // Only the non-QA Task settlement runs it.
	if repositoryCalls != wantRepositoryCalls {
		t.Fatalf("repository calls = %d, want %d", repositoryCalls, wantRepositoryCalls)
	}
	if committer.calls != 1 || runner.calls != 1 {
		t.Fatalf("Task commits = %d, Agent calls = %d", committer.calls, runner.calls)
	}
}

// Keep QA's own repository Verification out of the non-QA assertion.
type implementSettlementVerifier struct {
	*fakeVerifier
	taskCommands []string
}

func (v *implementSettlementVerifier) Verify(ctx context.Context, req daemon.VerifyRequest) (daemon.VerifyResult, error) {
	if strings.Contains(filepath.Base(req.OutputPath), "batch-001-attempt-") {
		v.taskCommands = append(v.taskCommands, req.Command)
	}
	return v.fakeVerifier.Verify(ctx, req)
}

func TestImplementUngatedRunEndsUnresolvedWhenTheIntegratedTreeIsRed(t *testing.T) {
	t.Parallel()
	testImplementIntegratedRepository(t, true)
}

func TestImplementUngatedRunEndsCleanWhenTheIntegratedTreeIsGreen(t *testing.T) {
	t.Parallel()
	testImplementIntegratedRepository(t, false)
}

type implementIntegratedVerifier struct {
	mu         sync.Mutex
	base       fakeVerifier
	fail       bool
	integrated []daemon.VerifyRequest
}

func (v *implementIntegratedVerifier) Verify(ctx context.Context, req daemon.VerifyRequest) (daemon.VerifyResult, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if filepath.Base(req.OutputPath) == "integrated-repository.log" {
		v.integrated = append(v.integrated, req)
		if v.fail {
			if err := os.MkdirAll(filepath.Dir(req.OutputPath), 0o755); err != nil {
				return daemon.VerifyResult{}, err
			}
			if err := os.WriteFile(req.OutputPath, []byte("integrated tree is red\n"), 0o644); err != nil {
				return daemon.VerifyResult{}, err
			}
			return daemon.VerifyResult{OutputPath: req.OutputPath}, &daemon.VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: errors.New("exit status 1")}
		}
	}
	return v.base.Verify(ctx, req)
}

func testImplementIntegratedRepository(t *testing.T, red bool) {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the first slice"},
		{id: "task_02", title: "Build the second slice"},
	})
	writeUserConfig(t, home, "defaults:\n  verification: repository check\nworktree:\n  concurrency: 2\n")
	runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted, "task_02": spec.StatusCompleted}}
	withAgentRunner(t, runner)
	verifier := &implementIntegratedVerifier{fail: red}
	withVerifier(t, verifier)
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)
	wantCode, wantState := exitOK, store.StateClean
	if red {
		wantCode, wantState = exitRunFailed, store.StateUnresolved
	}
	if code != wantCode {
		t.Fatalf("exit=%d, want %d; stdout=%s stderr=%s", code, wantCode, stdout.String(), stderr.String())
	}
	run := implementRunFromStore(t, home, implementRunIDFromStderr(t, stderr.String()))
	if run.State != wantState {
		t.Fatalf("Run state=%s, want %s; stderr=%s", run.State, wantState, stderr.String())
	}
	if len(verifier.integrated) != 1 || verifier.integrated[0].WorkDir != run.WorkDir {
		t.Fatalf("integrated calls=%+v; Run=%+v", verifier.integrated, run)
	}
	if runner.calls != 2 {
		t.Fatalf("Agent calls=%d, want 2", runner.calls)
	}
	if !red {
		assertRunWorktreeRemoved(t, run.WorkDir)
		return
	}
	assertRunWorktreeExists(t, run.WorkDir)
	reason := fmt.Sprintf("Integrated Repository Verification failed: repository check; diagnostics: %s", verifier.integrated[0].OutputPath)
	transcript := stderr.String()
	reasonAt := strings.Index(transcript, reason+"\n")
	terminalAt := strings.Index(transcript, "Implement Run "+run.ID+" reached Unresolved.")
	if reasonAt < 0 || terminalAt <= reasonAt {
		t.Fatalf("missing or out-of-order failure reason: %s", transcript)
	}
	found := false
	for _, entry := range runEventsForRun(t, home, run.ID) {
		if entry.Event.Kind != runevent.KindDaemonOutcome {
			continue
		}
		var payload struct {
			Reason string `json:"reason"`
		}
		if err := json.Unmarshal(entry.Event.Payload, &payload); err != nil {
			t.Fatal(err)
		}
		if payload.Reason == reason {
			found = true
		}
	}
	if !found {
		t.Fatalf("terminal event missing reason %q", reason)
	}
	for _, task := range []string{"task_01", "task_02"} {
		content, err := os.ReadFile(filepath.Join(run.WorkDir, "docs", "specs", implementTestSlug, task+".md"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), "status: completed") {
			t.Fatalf("%s not completed: %s", task, content)
		}
		userContent, err := os.ReadFile(filepath.Join(repo, "docs", "specs", implementTestSlug, task+".md"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(userContent), "status: pending") {
			t.Fatalf("red Run integrated into user branch: %s", userContent)
		}
	}
}
