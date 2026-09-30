package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/daemon"
	"roundfix/internal/spec"
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
